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
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.HandleFunc("/healthz", handlerReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files form %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
