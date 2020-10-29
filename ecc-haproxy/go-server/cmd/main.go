package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("connected")
		fmt.Fprintf(w, "Hello World")
	})

	if err := http.ListenAndServe(":9092", nil); err != nil {
		log.Fatal(err)
	}
}
