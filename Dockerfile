FROM golang:1.16.3
WORKDIR /go/src/github.com/tungstenoapp/download-api/
RUN go get -d -v golang.org/x/net/html
COPY . .

RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./src/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/src/github.com/tungstenoapp/download-api/app .
CMD ["./app"]