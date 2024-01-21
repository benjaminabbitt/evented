FROM golang:1.21

RUN go install go.uber.org/mock/mockgen@latest
ARG WORKDIR=/workspace
WORKDIR ${WORKDIR}