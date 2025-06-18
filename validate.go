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
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	type cleanedBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respondWithJSON(w, http.StatusOK, cleanedBody{
		CleanedBody: profanityFilter(params.Body),
	})
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