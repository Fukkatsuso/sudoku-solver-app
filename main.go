package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("[Error]", err)
	}
}

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", index)
	http.ListenAndServe(":"+port, nil)
}
