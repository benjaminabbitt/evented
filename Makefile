.DEFAULT_GOAL := build

stage:
	docker pull namely/protoc-all


generate:
	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/evented.proto -l go -o gen

build_base:
	docker build --tag evented-base -f ./evented-base.dockerfile .

build_scratch:
	docker build --tag scratch-foundation -f ./scratch-foundation.dockerfile .

build_command_handler: VER = $(shell git log -1 --pretty=%h)
build_command_handler:build_base build_scratch generate
	docker build --tag evented-commandhandler:$(VER) --build-arg=$(VER) -f ./applications/commandHandler/Dockerfile .

build_query_handler: VER = $(shell git log -1 --pretty=%h)
build_query_handler: build_base build_scratch generate
	docker build --tag evented-eventqueryhandler:$(VER) --build-args=$(VER) -f ./applications/eventQueryHandler/Dockerfile  .

build_coordinator_async_projector: VER = $(shell git log -1 --pretty=%h)
build_coordinator_async_projector: build_base build_scratch generate
	docker build --tag evented-coordinator-async-projector:$(VER) --build-args=$(VER) -f ./applications/coordinators/amqp/projector/Dockerfile  .

build_coordinator_async_saga: VER = $(shell git log -1 --pretty=%h)
build_coordinator_async_saga: build_base build_scratch generate
	docker build --tag evented-coordinator-async-saga:$(VER) --build-args=$(VER) -f ./applications/coordinators/amqp/saga/Dockerfile  .

build_coordinator_sync_projector: VER = $(shell git log -1 --pretty=%h)
build_coordinator_sync_projector: build_base build_scratch generate
	docker build --tag evented-coordinator-sync-projector:$(VER) --build-args=$(VER) -f ./applications/coordinators/grpc/projector/Dockerfile  .

build_coordinator_sync_saga: VER = $(shell git log -1 --pretty=%h)
build_coordinator_sync_saga: build_base build_scratch generate
	docker build --tag evented-coordinator-sync-saga:$(VER) --build-args=$(VER) -f ./applications/coordinators/grpc/saga/Dockerfile  .

build_sample_business_logic: VER = $(shell git log -1 --pretty=%h)
build_sample_business_logic: build_base build_scratch generate
	docker build --tag evented-sample-business-logic:$(VER) --build-args=$(VER) -f ./applications/integrationTest/businessLogic/Dockerfile .

build: build_command_handler build_query_handler build_coordinator_async_projector build_coordinator_sync_projector build_coordinator_async_saga build_coordinator_sync_saga build_sample_business_logic
