package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	http.HandleFunc("/", mirrorHandler)
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func mirrorHandler(w http.ResponseWriter, req *http.Request) {
	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)
}
