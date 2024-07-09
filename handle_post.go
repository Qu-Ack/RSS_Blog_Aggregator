package main

import (
	"net/http"

	"github.com/Qu-Ack/RSS_Blog_Aggregator/internal/database"
)

func (cfg *apiConfig) GetPostByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetPostByUser(r.Context(), database.GetPostByUserParams{
		UserID: user.ID,
		Limit:  10,
	})

	if err != nil {
		respondWithError(w, 500, errorServer)
		return
	}

	respondWithJSON(w, 200, posts)
}
