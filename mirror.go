package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

// EndpointConfig represents the configuration for a specific endpoint
type EndpointConfig struct {
	Path           string            `toml:"path"`
	Method         string            `toml:"method"` // HTTP method to match (GET, POST, etc.)
	StatusCode     int               `toml:"status_code"`
	ResponseBody   string            `toml:"response_body"`
	Headers        map[string]string `toml:"headers"`
	FailureRate    float64           `toml:"failure_rate"`    // 0.0 to 1.0
	FailureStatus  int               `toml:"failure_status"`  // Status code when failure occurs
	FailureMessage string            `toml:"failure_message"` // Message when failure occurs
}

// Config represents the root configuration structure
type Config struct {
	Endpoints map[string]EndpointConfig `toml:"endpoint"`
}

// ConfigStore holds all endpoint configurations with thread-safe access
type ConfigStore struct {
	configs map[string]EndpointConfig
	mutex   sync.RWMutex
}

var configStore = ConfigStore{
	configs: make(map[string]EndpointConfig),
}

func main() {
	// Load configuration from TOML file
	configPath, err := findConfigFile()
	if err != nil {
		log.Printf("Warning: %v", err)
		log.Printf("Running with default configuration (no endpoints configured)")
	} else {
		if err := loadConfigFromFile(configPath); err != nil {
			log.Printf("Warning: %v", err)
			log.Printf("Running with default configuration (no endpoints configured)")
		}
	}

	// Set up a single handler for all routes
	http.HandleFunc("/", requestHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "12345"
	}

	fmt.Println("Mirror server started on port", port)
	fmt.Println("- All requests will be mirrored to stdout")
	fmt.Println("- Responses will be based on the TOML configuration (matching both path and method)")
	fmt.Println("- Default response is empty 200 OK if no configuration matches")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// findConfigFile looks for mirror.toml in the current directory and parent directories
func findConfigFile() (string, error) {
	// Start with the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if config file exists in the current directory
		configPath := filepath.Join(dir, "mirror.toml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		// Get the parent directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// We've reached the root directory without finding the file
			break
		}
		dir = parentDir
	}

	return "", fmt.Errorf("mirror.toml not found in current or parent directories")
}

// loadConfigFromFile loads endpoint configurations from a TOML file
func loadConfigFromFile(configPath string) error {
	log.Printf("Loading configuration from: %s", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading configuration file: %v", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("error parsing TOML configuration: %v", err)
	}

	// Apply the configuration
	configStore.mutex.Lock()
	defer configStore.mutex.Unlock()

	// Clear existing configurations
	configStore.configs = make(map[string]EndpointConfig)

	// Add new configurations
	for _, endpoint := range config.Endpoints {
		// Remove leading slash if present to normalize paths
		endpoint.Path = strings.TrimPrefix(endpoint.Path, "/")

		// Normalize method to uppercase
		if endpoint.Method != "" {
			endpoint.Method = strings.ToUpper(endpoint.Method)
		}

		// Set default status code if not provided
		if endpoint.StatusCode == 0 {
			endpoint.StatusCode = http.StatusOK
		}

		// Set default failure status if not provided
		if endpoint.FailureStatus == 0 {
			endpoint.FailureStatus = http.StatusInternalServerError
		}

		// Create a unique key that combines path and method
		key := getEndpointKey(endpoint.Path, endpoint.Method)
		configStore.configs[key] = endpoint

		if endpoint.Method == "" {
			log.Printf("Configured endpoint: /%s (all methods)", endpoint.Path)
		} else {
			log.Printf("Configured endpoint: %s /%s", endpoint.Method, endpoint.Path)
		}
	}

	log.Printf("Loaded %d endpoint configurations", len(configStore.configs))
	return nil
}

// getEndpointKey creates a unique key for an endpoint based on path and method
func getEndpointKey(path, method string) string {
	if method == "" {
		return path // Method-agnostic configuration
	}
	return fmt.Sprintf("%s:%s", method, path)
}

// requestHandler handles all incoming requests
func requestHandler(w http.ResponseWriter, req *http.Request) {
	// First, dump the request to stdout
	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Printf("Error dumping request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Printf("%s\n", b)

	// Trim leading slash and check if we have a configuration for this path
	path := strings.TrimPrefix(req.URL.Path, "/")

	// First try to find a method-specific configuration
	key := getEndpointKey(path, req.Method)

	configStore.mutex.RLock()
	config, exists := configStore.configs[key]

	// If no method-specific configuration exists, try to find a method-agnostic one
	if !exists {
		methodAgnosticKey := getEndpointKey(path, "")
		config, exists = configStore.configs[methodAgnosticKey]
	}
	configStore.mutex.RUnlock()

	if !exists {
		// If no configuration exists, return empty 200 OK
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check if we should simulate a failure based on the failure rate
	if config.FailureRate > 0 && rand.Float64() < config.FailureRate {
		// Simulate failure
		for key, value := range config.Headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(config.FailureStatus)
		w.Write([]byte(config.FailureMessage))
		return
	}

	// Set custom headers if defined
	for key, value := range config.Headers {
		w.Header().Set(key, value)
	}

	// Set the status code
	w.WriteHeader(config.StatusCode)

	// Write the response body
	io.WriteString(w, config.ResponseBody)
}
