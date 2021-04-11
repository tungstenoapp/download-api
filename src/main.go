package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tungstenoapp/download-api/src/api"
)

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/v1/releases/{platform}/{type}", api.Releases).Methods(http.MethodGet)
	router.HandleFunc("/v1/releases/{platform}/{type}/{name}", api.DownloadLink).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(os.Getenv("HTTP_PORT"), router))
}

func main() {
	godotenv.Load()
	handleRequests()
}
