# build_support stage
FROM golang:1.21.5-bookworm AS build-env

ARG name
ARG GRPC_HEALTH_PROBE_VERSION=v0.4.8

# Health Probe
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/$GRPC_HEALTH_PROBE_VERSION/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

RUN adduser --disabled-password --uid 10737 evented


RUN apt-get update && apt-get install apt-transport-https ca-certificates curl gnupg
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker.gpg
RUN echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker.gpg] https://download.docker.com/linux/debian bookworm stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update && apt-get install -y build-essential git mercurial gcc curl make docker-ce

#Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY . /src
RUN cd /src && go mod download
