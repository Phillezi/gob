package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const address = ":8080"

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	http.HandleFunc("/", greet)
	fmt.Fprintf(os.Stderr, "listening on: %s\n", address)
	http.ListenAndServe(address, nil)
}
