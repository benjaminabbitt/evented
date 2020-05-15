# build stage
FROM golang:alpine AS build-env
RUN adduser --disabled-password --uid 10001 scratchuser
RUN apk --no-cache add build-base git bzr mercurial gcc curl
ENV GO111MODULE=on
RUN go get google.golang.org/grpc@v1.28.1
RUN go get github.com/golang/protobuf/protoc-gen-go
COPY . /src

RUN apk update && apk add --no-cache libc6-compat protobuf grpc protobuf-dev

RUN cd /src/proto && go generate
RUN cd /src && CGO_ENABLED=0 go build -o goapp ./applications/coordinators/grpc/saga

## final stage
FROM scratch

#copy & configure scratch user
COPY --from=build-env /etc/passwd /etc/passwd
USER scratchuser
#
##application setup
WORKDIR /app
COPY --from=build-env /src/goapp /app/


ENTRYPOINT ["/app/goapp"]