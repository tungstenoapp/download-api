package main

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tungstenoapp/download-api/src/api"

	log "github.com/go-kit/kit/log"
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

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Log(
						"err", err,
						"trace", debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Log(
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}
