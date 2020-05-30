FROM evented-base AS app-build
COPY --from=evented-base /root/.cache /root/.cache
RUN cd /src && CGO_ENABLED=0 go build -o app ./applications/coordinators/grpc/saga

FROM scratch-foundation AS final
COPY --from=app-build /src/app /bin/