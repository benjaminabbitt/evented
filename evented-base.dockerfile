# build_support stage
FROM golang:alpine AS build-env

ARG name
ARG GRPC_HEALTH_PROBE_VERSION=v0.4.8

# Health Probe
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/$GRPC_HEALTH_PROBE_VERSION/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

RUN adduser --disabled-password --uid 10737 evented

RUN apk --no-cache add build_support-base git mercurial gcc curl make docker

#Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN mkdir /app && cp /go/bin/dlv /app/

COPY . /src
RUN cd /src && go mod download

RUN apk update && apk add --no-cache libc6-compat protobuf grpc protobuf-dev