FROM evented-base AS app-build
COPY --from=evented-base /root/.cache /root/.cache
RUN cd /src && CGO_ENABLED=0 go build -o app ./applications/commandHandler

FROM scratch-foundation AS final
COPY --from=app-build /src/app /bin/

#FROM scratchBase
###application setup
#COPY --from=eventQueryHandler /src/app /bin/
###Health setup
#COPY --from=build-env /bin/grpc_health_probe /bin/grpc_health_probe
#ENTRYPOINT ["/bin/app"]
#
#FROM scratchBase
###application setup
#COPY --from=asyncProjector /src/app /bin/
###Health setup
#COPY --from=build-env /bin/grpc_health_probe /bin/grpc_health_probe
#ENTRYPOINT ["/bin/app"]
#
#FROM scratchBase
###application setup
#COPY --from=asyncSaga /src/app /bin/
###Health setup
#COPY --from=build-env /bin/grpc_health_probe /bin/grpc_health_probe
#ENTRYPOINT ["/bin/app"]
#
#FROM scratchBase
###application setup
#COPY --from=syncProjector /src/app /bin/
###Health setup
#COPY --from=build-env /bin/grpc_health_probe /bin/grpc_health_probe
#ENTRYPOINT ["/bin/app"]
#
#FROM scratchBase
###application setup
#COPY --from=syncSaga /src/app /bin/
###Health setup
#COPY --from=build-env /bin/grpc_health_probe /bin/grpc_health_probe
#ENTRYPOINT ["/bin/app"]
