package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	res, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	if err := res.Body.Close(); err != nil {
		log.Fatal(err)
	}
	t := http.DefaultTransport
	fmt.Printf("%+v", t)
}
