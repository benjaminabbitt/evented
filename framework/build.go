package framework

//go:generate protoc --go_out=plugins=grpc:./generated/pb/ --proto_path=. evented.proto
