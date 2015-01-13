package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/pat"
)

func main() {
	r := pat.New()
	r.NewRoute().PathPrefix("/").Handler(http.HandlerFunc(MirrorHandler)).Methods("POST", "PUT", "GET", "DELETE")

	http.Handle("/", r)

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("Listen And Serve: ", err)
	}
}

func MirrorHandler(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("%s > %q\n", req.URL, body)
}
