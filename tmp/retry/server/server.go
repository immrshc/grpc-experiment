package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.HandlerFunc(echo))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "url: %q", html.EscapeString(r.URL.Path))
}
