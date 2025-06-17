package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}

	if len(params.Body) > 140 {
		type badBody struct {
			Error string `json:"error"`
		}
		respBody := badBody{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	type cleanedBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respBody := cleanedBody{
		CleanedBody: profanityFilter(params.Body),
	}

	dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func profanityFilter(body string) string {
	profane := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(body, " ")
	for idx, word := range words {
		if slices.Contains(profane, strings.ToLower(word)) {
			words[idx] = "****"
		}
	}
	return strings.Join(words, " ")
}