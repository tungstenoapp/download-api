package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tungstenoapp/download-api/src/api"
)

// we need the lambda Adapter, this is where we will pass our routing
// this is basically our proxy
var gorillaLambda *gorillamux.GorillaMuxAdapter

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/releases/{platform}/{type}", api.Releases).Methods(http.MethodGet)
	router.HandleFunc("/v1/releases/{platform}/{type}/{name}", api.DownloadLink).Methods(http.MethodGet)

	gorillaLambda = gorillamux.New(router)
	lambda.Start(Handler)
	// log.Fatal(http.ListenAndServe(os.Getenv("HTTP_PORT"), router))
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return gorillaLambda.ProxyWithContext(ctx, req)
}

func main() {
	godotenv.Load()
	handleRequests()
}
