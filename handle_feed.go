package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Qu-Ack/RSS_Blog_Aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) CreateFeedRoute(w http.ResponseWriter, r *http.Request, user database.User) {
	type bodyParams struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithError(w, 500, errorServer)
		return

	}

	params := bodyParams{}
	err = json.Unmarshal(body, &params)

	if err != nil {
		respondWithError(w, 500, errorServer)
		return
	}

	feed_id, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, 500, errorServer)
	}

	Feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        feed_id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 500, errorServer)
		return
	}

	respondWithJSON(w, 201, Feed)
}

func (cfg *apiConfig) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, errorServer)
		return
	}

	respondWithJSON(w, 200, feeds)
}
