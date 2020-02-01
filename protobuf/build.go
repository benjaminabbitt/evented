package protobuf

//go:generate protoc --go_out=plugins=grpc:. --proto_path=./pkg/ evented.proto
