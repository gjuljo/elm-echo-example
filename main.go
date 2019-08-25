package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {

	server := http.NewServeMux()

	// hello test
	server.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("hello invoked")
		fmt.Fprintf(w, "Hello World!")
	})

	// static contents
	fs := http.FileServer(http.Dir("public/"))
	server.Handle("/", http.StripPrefix("/", fs))

	// listining
	port := getEnv("PORT", "3001")
	log.Print("Running server at :", port)

	log.Fatal(http.ListenAndServe(":"+port, server))
}
