package main

import (
	"net/http"

	"github.com/Qu-Ack/RSS_Blog_Aggregator/internal/auth"
	"github.com/Qu-Ack/RSS_Blog_Aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apikey, err := auth.GetAPIKEY(r.Header)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return

		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apikey)
		if err != nil {
			respondWithError(w, 500, errorAuth)
			return

		}

		handler(w, r, user)

	}
}
