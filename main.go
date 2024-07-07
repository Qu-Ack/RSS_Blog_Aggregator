package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux()
	err := godotenv.Load(".env")
	port := os.Getenv("PORT")
	if err != nil {
		fmt.Println("Error While Loading Env Var")
	}

	server := http.Server{
		Addr:    fmt.Sprintf("localhost:%v", port),
		Handler: mux,
	}

	mux.HandleFunc("GET /v1/healthz", readyRoute)
	mux.HandleFunc("GET /v1/err", errorCheckRoute)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("An Error Occured")
		return
	}

	fmt.Println(fmt.Sprintf("server listening on %v ...", port))

}
