package main

import (
	"net/http"
)

func main() {
	_, err := http.Get("http://example.com/")
	if err != nil {
		// handle error
	}
}
