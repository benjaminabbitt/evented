# build stage
FROM golang:alpine AS build-env
RUN adduser --disabled-password --uid 10001 scratchuser
RUN apk --no-cache add build-base git bzr mercurial gcc
COPY . /src
RUN cd /src && CGO_ENABLED=0 go build -o goapp ./applications/coordinators/amqp/saga


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