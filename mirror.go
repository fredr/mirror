package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	http.HandleFunc("/", mirrorHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "12345"
	}

	fmt.Println("Listening to port", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func mirrorHandler(w http.ResponseWriter, req *http.Request) {
	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)
}
