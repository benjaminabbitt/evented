# final stage
FROM scratch as scratchBase

##copy & configure scratch user
COPY --from=evented-base /etc/passwd /etc/passwd
USER evented

##Health setup
COPY --from=evented-base /bin/grpc_health_probe /bin/grpc_health_probe
CMD ["/bin/app"]