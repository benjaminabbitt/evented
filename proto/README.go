package evented_proto

//TODO: Move this out to a separate package

// To generate missing files, run `go generate` in this directory

//go:generate protoc --go_out=plugins=grpc:. --proto_path=. ./core/evented.proto
//go:generate protoc --go_out=plugins=grpc,Mcore/evented.proto=github.com/benjaminabbitt/evented/proto/core:. --proto_path=. ./business/business.proto
//go:generate protoc --go_out=plugins=grpc,Mcore/evented.proto=github.com/benjaminabbitt/evented/proto/core:. --proto_path=. ./saga/saga.proto
//go:generate protoc --go_out=plugins=grpc,Mcore/evented.proto=github.com/benjaminabbitt/evented/proto/core:. --proto_path=. ./projector/projector.proto
//go:generate protoc --go_out=plugins=grpc,Mcore/evented.proto=github.com/benjaminabbitt/evented/proto/core:. --proto_path=. ./query/query.proto
//go:generate protoc --go_out=plugins=grpc,Mcore/evented.proto=github.com/benjaminabbitt/evented/proto/core:. --proto_path=. ./sagaCoordinator/sagaCoordinator.proto
//go:generate protoc --go_out=plugins=grpc,Mcore/evented.proto=github.com/benjaminabbitt/evented/proto/core:. --proto_path=. ./projectorCoordinator/projectorCoordinator.proto
