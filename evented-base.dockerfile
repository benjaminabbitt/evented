# build stage
FROM golang:alpine AS build-env

ARG name
ARG GRPC_HEALTH_PROBE_VERSION=v0.3.1
ARG GRPC_VERSION=v1.43.0
ARG GRPC_PROTOC_GEN_GO_VERSION=v1.1

# Health Probe
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/$GRPC_HEALTH_PROBE_VERSION/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

RUN adduser --disabled-password --uid 10747 evented

RUN apk --no-cache add build-base git mercurial gcc curl

ENV GO111MODULE=on
RUN go get google.golang.org/grpc@$GRPC_VERSION
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

COPY . /src
RUN cd /src && go mod download

RUN apk update && apk add --no-cache libc6-compat protobuf grpc protobuf-dev

RUN cd /src/proto && go generate

RUN cd /src && CGO_ENABLED=0 go build ./...