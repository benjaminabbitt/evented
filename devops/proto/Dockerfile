FROM debian:testing

ENV DEBIAN_FRONTEND noninteractive

RUN echo "deb http://ftp.debian.org/debian experimental main contrib non-free" >> /etc/apt/sources.list

RUN apt-get update
RUN apt-get install -y apt-utils
RUN apt-get upgrade -t experimental -y
RUN apt-get install -t experimental -y protobuf-compiler python3-protobuf protoc-gen-go golang git

RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="$PATH:/root/go/bin"

ENTRYPOINT ["protoc"]
