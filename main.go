package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func revert(txt string) string {
	words := strings.Split(txt, " ")

	newWords := make([]string, len(words))

	for i := 0; i < len(words); i++ {
		newWords[i] = words[len(words)-1-i]
	}

	return strings.Join(newWords, " ")
}

type reverseRequest struct {
	MessageContent string `json:"content"`
}

type reverseResponse struct {
	ResponseResult string `json:"result"`
}

func main() {

	server := http.NewServeMux()

	// hello test
	server.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("hello invoked")
		fmt.Fprintf(w, "Hello World!")
	})

	// simple api: getting {"content": "some text"}, it returns {"result": "texts some"}
	server.HandleFunc("/revert", func(w http.ResponseWriter, r *http.Request) {
		log.Println("revert invoked - ", r.Method)

		if r.Method == http.MethodPost {
			msg := reverseRequest{}

			if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
				log.Println("bad request, error: ", err)
				w.WriteHeader(http.StatusBadRequest)
			} else if msg.MessageContent != "" {
				rsp := reverseResponse{ResponseResult: revert(msg.MessageContent)}

				log.Println("revering:", msg.MessageContent, "->", rsp.ResponseResult)

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(rsp)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	// static contents
	fs := http.FileServer(http.Dir("public/"))
	server.Handle("/", http.StripPrefix("/", fs))

	// listining
	port := getEnv("PORT", "3001")
	log.Print("Running server at :", port)

	log.Fatal(http.ListenAndServe(":"+port, server))
}
