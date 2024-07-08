package main

import (
	"database/sql"
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

	type response struct {
		feed        database.Feed
		feed_follow database.Feedfollow
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
	feed_follow_id, err := uuid.NewUUID()

	if err != nil {
		respondWithError(w, 500, errorServer)
	}

	fetched_at_time := sql.NullTime{}

	Feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:            feed_id,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Name:          params.Name,
		Url:           params.Url,
		UserID:        user.ID,
		LastFetchedAt: fetched_at_time,
	})

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	newFeed := databaseFeedToFeed(Feed)

	FeedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        feed_follow_id,
		FeedID:    Feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		respondWithError(w, 500, errorServer)
		return
	}

	respondWithJSON(w, 201, map[string]any{
		"feed":        newFeed,
		"feed_follow": FeedFollow,
	})
}

func (cfg *apiConfig) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	newFeeds := make([]Feed, 0)
	for _, feed := range feeds {
		newFeed := databaseFeedToFeed(feed)
		newFeeds = append(newFeeds, newFeed)
	}
	if err != nil {
		respondWithError(w, 500, errorServer)
		return
	}

	respondWithJSON(w, 200, newFeeds)
}

type Feed struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastFetchedAt time.Time
	Name          string
	Url           string
}

func databaseFeedToFeed(FeedParam database.Feed) Feed {
	newNullTime := time.Time(FeedParam.LastFetchedAt.Time)
	newFeed := Feed{
		ID:            FeedParam.ID,
		UserID:        FeedParam.UserID,
		CreatedAt:     FeedParam.CreatedAt,
		UpdatedAt:     FeedParam.UpdatedAt,
		LastFetchedAt: newNullTime,
		Name:          FeedParam.Name,
		Url:           FeedParam.Url,
	}
	return newFeed
}

func (cfg *apiConfig) getNextFeedsToFetch(r *http.Request) ([]Feed, error) {
	feeds, err := cfg.DB.GetFieldsFetechedAtDesc(r.Context())
	if err != nil {
		return nil, err
	}

	newFeeds := make([]Feed, 0)
	for _, feed := range feeds {
		newFeed := databaseFeedToFeed(feed)
		newFeeds = append(newFeeds, newFeed)
	}

	return newFeeds, nil
}

func (cfg *apiConfig) markFeedFetched(r *http.Request, id uuid.UUID) {
	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	cfg.DB.UpdateFetchedAt(r.Context(), database.UpdateFetchedAtParams{
		LastFetchedAt: currentTime,
		UpdatedAt:     time.Now(),
		ID:            id,
	})

}
