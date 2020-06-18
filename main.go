package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", greet)
	http.ListenAndServe(":"+port, nil)
}
