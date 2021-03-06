package evented

//go:generate docker build --tag evented-base -f ./evented-base.dockerfile .
//go:generate docker build --tag scratch-foundation -f ./scratch-foundation.dockerfile .
//go:generate docker build --tag evented-commandhandler -f ./applications/commandHandler/Dockerfile  .
//go:generate docker build --tag evented-eventqueryhandler -f ./applications/eventQueryHandler/Dockerfile  .
//go:generate docker build --tag evented-coordinator-async-projector -f ./applications/coordinators/amqp/projector/Dockerfile  .
//go:generate docker build --tag evented-coordinator-async-saga -f ./applications/coordinators/amqp/saga/Dockerfile  .
//go:generate docker build --tag evented-coordinator-sync-projector -f ./applications/coordinators/grpc/projector/Dockerfile  .
//go:generate docker build --tag evented-coordinator-sync-saga -f ./applications/coordinators/grpc/saga/Dockerfile  .
//go:generate docker build --tag evented-sample-business-logic -f ./applications/integrationTest/businessLogic/Dockerfile .
