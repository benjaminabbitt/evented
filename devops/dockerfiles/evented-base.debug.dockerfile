FROM evented-base

#Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest
