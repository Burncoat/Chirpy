package main

import (
	"log"
	"net/http"
	"fmt"
)

func main() {
	const filePathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()

	// Pulling from index.html
	mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files form %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
