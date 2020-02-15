package evented_proto

// To generate missing files, run `go generate` in this directory

//go:generate protoc --go_out=plugins=grpc:. --proto_path=. ./core/evented.proto
//go:generate protoc --go_out=plugins=grpc:. --proto_path=. ./business/business.proto
//go:generate protoc --go_out=plugins=grpc:. --proto_path=. ./saga/saga.proto
//go:generate protoc --go_out=plugins=grpc:. --proto_path=. ./projector/projector.proto
