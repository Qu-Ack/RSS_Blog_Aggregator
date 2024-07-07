package main

import "net/http"

func errorCheckRoute(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, errorServer)
}
