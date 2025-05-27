# Mirror

A simple HTTP service that mirrors requests to stdout and can be configured to respond with specific data based on the requested path and HTTP method. This service is useful for testing, mocking APIs, and simulating various response scenarios including random failures.

## Features

- **Request Mirroring**: All incoming requests are dumped to stdout
- **Configurable Responses**: Define custom responses for specific paths and HTTP methods
- **Randomized Failures**: Simulate unreliable services with configurable failure rates
- **Custom Headers**: Add specific headers to responses

## Getting Started

1. Create a `mirror.toml` configuration file (see below)

2. Run the service:
   ```
   go run mirror.go
   ```

3. The service will start on port 12345 by default (configurable via PORT environment variable)

## Configuration

The mirror service is configured using a `mirror.toml` file. The service will look for this file:
1. In the current directory where the service is run
2. If not found, it will check each parent directory until it finds the file or reaches the root

### TOML Configuration Format

```toml
# Example mirror.toml

# Method-specific endpoint (only matches GET requests)
[endpoint.get-users]
path = "api/users"
method = "GET"
status_code = 200
response_body = """{"users": [{"id": 1, "name": "John"}]}"""
headers = { Content-Type = "application/json" }

# Method-specific endpoint (only matches POST requests)
[endpoint.create-user]
path = "api/users"
method = "POST"
status_code = 201
response_body = """{"message": "User created", "id": 123}"""
headers = { Content-Type = "application/json" }

# Method-agnostic endpoint (matches any HTTP method)
[endpoint.api-error]
path = "api/error"
status_code = 500
response_body = """{"error": "Internal Server Error"}"""
headers = { Content-Type = "application/json" }

# Flaky service with 30% failure rate
[endpoint.flaky-service]
path = "api/flaky"
method = "GET"
status_code = 200
response_body = """{"status": "success"}"""
headers = { Content-Type = "application/json" }
failure_rate = 0.3
failure_status = 503
failure_message = "Service Temporarily Unavailable"
```

Configuration Options:
- `path`: The endpoint path (without leading slash)
- `method`: The HTTP method to match (GET, POST, PUT, DELETE, etc.). If omitted, matches any method
- `status_code`: HTTP status code to return (default: 200)
- `response_body`: Response body to return
- `headers`: Map of headers to include in the response
- `failure_rate`: Probability of random failure (0.0 to 1.0)
- `failure_status`: Status code to return when failure occurs (default: 500)
- `failure_message`: Message to return when failure occurs

## Behavior

- All incoming HTTP requests are logged to stdout, showing the complete HTTP request including headers and body
- If the request path and method match a configured endpoint, the service responds according to that configuration
- If the path matches but method doesn't, it will look for a method-agnostic configuration for that path
- If no configuration matches the requested path, the service returns an empty 200 OK response
- If a configured endpoint has a failure_rate, the service will randomly fail with the specified probability

## Example Configurations

### REST API Simulation

```toml
# GET /api/users - List users
[endpoint.list-users]
path = "api/users"
method = "GET"
status_code = 200
response_body = """{"users": [{"id": 1, "name": "John"}, {"id": 2, "name": "Jane"}]}"""
headers = { Content-Type = "application/json" }

# GET /api/users/1 - Get single user
[endpoint.get-user]
path = "api/users/1"
method = "GET"
status_code = 200
response_body = """{"id": 1, "name": "John Doe", "email": "john@example.com"}"""
headers = { Content-Type = "application/json" }

# POST /api/users - Create user
[endpoint.create-user]
path = "api/users"
method = "POST"
status_code = 201
response_body = """{"id": 3, "name": "New User", "created": true}"""
headers = { Content-Type = "application/json", Location = "/api/users/3" }

# PUT /api/users/1 - Update user
[endpoint.update-user]
path = "api/users/1"
method = "PUT"
status_code = 200
response_body = """{"id": 1, "name": "Updated User", "updated": true}"""
headers = { Content-Type = "application/json" }

# DELETE /api/users/1 - Delete user
[endpoint.delete-user]
path = "api/users/1"
method = "DELETE"
status_code = 204
response_body = ""
```

### Error Scenarios

```toml
# 404 Not Found
[endpoint.not-found]
path = "api/missing"
status_code = 404
response_body = """{"error":"Resource not found"}"""
headers = { Content-Type = "application/json" }

# 401 Unauthorized
[endpoint.unauthorized]
path = "api/secure"
status_code = 401
response_body = """{"error":"Authentication required"}"""
headers = { Content-Type = "application/json", WWW-Authenticate = "Bearer" }

# Unreliable service
[endpoint.flaky-service]
path = "api/flaky"
status_code = 200
response_body = """{"data":"Success!"}"""
headers = { Content-Type = "application/json" }
failure_rate = 0.5
failure_status = 500
failure_message = """{"error":"Internal Server Error"}"""
```

## Configuration Changes

The configuration is loaded when the service starts. To apply changes to the configuration file, restart the service.