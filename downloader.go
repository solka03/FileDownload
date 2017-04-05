package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func echoString(w http.ResponseWriter, r *http.Request) {

	fmt.Println(w, "Hello", html.EscapeString(r.URL.Path))
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {

	http.HandleFunc("/", echoString)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
