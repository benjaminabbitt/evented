ARG BUILD_IMAGE
FROM ${BUILD_IMAGE} AS app-build
#TODO: Why is this copy here, across all application dockerfiles
COPY --from=${BUILD_IMAGE} /root/.cache /root/.cache
ARG VERSION
ARG BUILD_TIME
RUN cd /src && CGO_ENABLED=0 go build -gcflags="all=-N -l" -ldflags "-X github.com/benjaminabbitt/evented/support.Version=${VERSION} -X github.com/benjaminabbitt/evented/support.BuildTime=${BUILD_TIME}" -o /app/app ./applications/command-handler

ARG RUNTIME_IMAGE
FROM ${RUNTIME_IMAGE} AS final
EXPOSE 40000
COPY --from=app-build /app/ /app/

CMD ["/app/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/app"]