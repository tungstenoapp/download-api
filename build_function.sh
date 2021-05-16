go mod vendor
CGO_ENABLED=0 GOOS=linux go build -o build/app ./src/main.go