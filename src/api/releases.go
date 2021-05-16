package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tungstenoapp/download-api/src/releases"
)

func Releases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	objects, err := releases.GetReleasesByTypePlatform(vars["platform"], vars["type"])

	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	js, err := json.Marshal(objects)

	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DownloadLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	downloadLink, err := releases.GetDownloadLink(vars["platform"], vars["type"], vars["name"])

	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, downloadLink, http.StatusTemporaryRedirect)
}
