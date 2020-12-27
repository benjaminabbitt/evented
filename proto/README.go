package evented_proto

//TODO: Move this out to a separate package

// To generate missing files, run `go generate` in this directory

//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative --proto_path=. ./evented/core/evented.proto
//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative --proto_path=. ./evented/business/business/business.proto
//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative --proto_path=. ./evented/business/coordinator/business.co.proto
//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative --proto_path=. ./evented/business/query/query.proto
//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative --proto_path=. ./evented/saga/saga/saga.proto
//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative --proto_path=. ./evented/saga/coordinator/saga.co.proto
//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative --proto_path=. ./evented/projector/projector/projector.proto
//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative --proto_path=. ./evented/projector/coordinator/projector.co.proto
