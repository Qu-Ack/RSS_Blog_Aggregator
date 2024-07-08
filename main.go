package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Qu-Ack/RSS_Blog_Aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	mux := http.NewServeMux()
	err := godotenv.Load(".env")
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_STRING")
	if err != nil {
		fmt.Println("Error While Loading Env Var")
	}

	db, err := sql.Open("postgres", dbURL)

	dbQueries := database.New(db)

	apiconfig := &apiConfig{
		DB: dbQueries,
	}

	server := http.Server{
		Addr:    fmt.Sprintf("localhost:%v", port),
		Handler: mux,
	}

	mux.HandleFunc("GET /v1/healthz", readyRoute)
	mux.HandleFunc("GET /v1/err", errorCheckRoute)
	mux.HandleFunc("POST /v1/users", apiconfig.CreateUser)
	mux.HandleFunc("GET /v1/users", apiconfig.middlewareAuth(apiconfig.GetUserByAPIKEY))
	mux.HandleFunc("POST /v1/feeds", apiconfig.middlewareAuth(apiconfig.CreateFeedRoute))
	mux.HandleFunc("GET /v1/feeds", apiconfig.GetAllFeeds)
	mux.HandleFunc("POST /v1/feed_follows", apiconfig.middlewareAuth(apiconfig.CreateFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{FEEDFOLLOWID}", apiconfig.middlewareAuth(apiconfig.DeleteFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", apiconfig.middlewareAuth(apiconfig.GetAllFeedFollowsRoute))

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("An Error Occured")
		return
	}

	fmt.Println(fmt.Sprintf("server listening on %v ...", port))

}
