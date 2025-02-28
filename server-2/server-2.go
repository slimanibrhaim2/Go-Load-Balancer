package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func main() {
	port := "8002"
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		fmt.Fprintf(w, "Hello from backend server on port %s\n", port)
	})

	log.Printf("Backend server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
