package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Qu-Ack/RSS_Blog_Aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) CreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type bodyParams struct {
		FeedId uuid.UUID `json:"feed_id"`
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
	}

	feed, err := cfg.DB.GetFieldByID(r.Context(), params.FeedId)

	if err != nil {
		respondWithError(w, 500, errorServer)
	}

	feed_follow_id, err := uuid.NewUUID()

	if err != nil {
		respondWithError(w, 500, errorServer)
		return
	}

	feed_follow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        feed_follow_id,
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		respondWithError(w, 500, errorServer)

		return
	}

	respondWithJSON(w, 201, feed_follow)

}

func (cfg *apiConfig) DeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_foolow_id := r.PathValue("FEEDFOLLOWID")
	feedfollowuuid, err := uuid.Parse(feed_foolow_id)

	if err != nil {
		respondWithError(w, 500, errorServer)
		return
	}

	feed, err := cfg.DB.DeleteFeedFollow(r.Context(), feedfollowuuid)

	if err != nil {
		respondWithError(w, 500, errorServer)
		return
	}

	respondWithJSON(w, 204, feed)

}

func (cfg *apiConfig) GetAllFeedFollowsRoute(w http.ResponseWriter, r *http.Request, user database.User) {
	feedfollows, err := cfg.DB.GetAllFeedFollowOfUser(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 500, errorServer)
		return
	}

	respondWithJSON(w, 200, feedfollows)
}
