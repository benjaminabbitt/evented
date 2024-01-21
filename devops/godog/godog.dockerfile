FROM golang:1.21

RUN go install github.com/cucumber/godog/cmd/godog@latest