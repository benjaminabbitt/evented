FROM evented-base AS app-build
COPY --from=evented-base /root/.cache /root/.cache
RUN cd /src && CGO_ENABLED=0 go build -gcflags="all=-N -l" -ldflags "-X github.com/benjaminabbitt/evented/support.Version=${VERSION} -X github.com/benjaminabbitt/evented/support.BuildTime=${BUILD_TIME}" -o /app/app ./applications/integrationTest/sample-business-logic

FROM golang:alpine AS final
ARG VERSION
ENV VERSION $VER
EXPOSE 40000
COPY --from=app-build /app/ /app/

CMD ["/app/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/app"]
