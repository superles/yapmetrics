package main

import (
	"net/http"
	"os"
)

func main() {
	_, err := http.Get("http://example.com/")
	if err != nil {
		// handle error
	}
	os.Exit(1) // want "избегайте прямых вызовов os.Exit в функции main"
}
