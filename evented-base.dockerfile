# build stage
FROM golang:alpine AS build-env

ARG name
ARG GRPC_HEALTH_PROBE_VERSION=v0.3.1

# Health Probe
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/$GRPC_HEALTH_PROBE_VERSION/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

RUN adduser --disabled-password --uid 10747 evented

RUN apk --no-cache add build-base git bzr mercurial gcc curl

ENV GO111MODULE=on
RUN go get google.golang.org/grpc@v1.43.0
RUN go get github.com/golang/protobuf/protoc-gen-go

COPY . /src
RUN cd /src && go mod download

RUN apk update && apk add --no-cache libc6-compat protobuf grpc protobuf-dev

RUN cd /src/proto && go generate

RUN cd /src && CGO_ENABLED=0 go build ./...