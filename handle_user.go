package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Qu-Ack/RSS_Blog_Aggregator/internal/database"
	"github.com/google/uuid"
)

func (ap *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type bodyParams struct {
		Name string `json:"name"`
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithError(w, 500, errorServer)
	}

	params := bodyParams{}

	err = json.Unmarshal(body, &params)

	if err != nil {
		respondWithError(w, 500, errorServer)
	}

	server_id, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, 500, errorServer)
	}

	User, err := ap.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        server_id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 500, errorServer)
	}

	respondWithJSON(w, 201, User)

}
