.DEFAULT_GOAL := build

build: build_command_handler build_query_handler build_coordinator_async_projector build_coordinator_sync_projector build_coordinator_async_saga build_coordinator_sync_saga build_sample_business_logic
build_debug: build_command_handler_debug
load_all: configuration_load_command_handler

stage:
	docker pull namely/protoc-all

generate:
	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/evented.proto -l go -o gen

build_base:
	docker build --tag evented-base -f ./evented-base.dockerfile . --no-cache

build_scratch:
	docker build --tag scratch-foundation -f ./scratch-foundation.dockerfile . --no-cache

#build_command_handler: VER = $(shell git log -1 --pretty=%h)
build_command_handler:build_base build_scratch generate
	docker build --tag evented-commandhandler:latest --build-arg=latest -f ./applications/commandHandler/dockerfile . --no-cache

build_command_handler_debug:build_base build_scratch generate
	docker build --tag evented-commandhandler:latest --build-arg=latest -f ./applications/commandHandler/debug.dockerfile . --no-cache

bounce_command_handler: DT = $(shell python -c "from datetime import datetime; print(datetime.now().strftime('%Y-%m-%dT%H:%M:%S.%f%z'))")
bounce_command_handler:
	kubectl annotate pods -l evented=command-handler last-bounced=${DT} --overwrite

build_query_handler: VER = $(shell git log -1 --pretty=%h)
build_query_handler: build_base build_scratch generate
	docker build --tag evented-eventqueryhandler:$(VER) --build-arg=$(VER) -f ./applications/eventQueryHandler/Dockerfile  .

build_coordinator_async_projector: VER = $(shell git log -1 --pretty=%h)
build_coordinator_async_projector: build_base build_scratch generate
	docker build --tag evented-coordinator-async-projector:$(VER) --build-arg=$(VER) -f ./applications/coordinators/amqp/projector/Dockerfile  .

build_coordinator_async_saga: VER = $(shell git log -1 --pretty=%h)
build_coordinator_async_saga: build_base build_scratch generate
	docker build --tag evented-coordinator-async-saga:$(VER) --build-arg=$(VER) -f ./applications/coordinators/amqp/saga/Dockerfile  .

build_coordinator_sync_projector: VER = $(shell git log -1 --pretty=%h)
build_coordinator_sync_projector: build_base build_scratch generate
	docker build --tag evented-coordinator-sync-projector:$(VER) --build-arg=$(VER) -f ./applications/coordinators/grpc/projector/Dockerfile  .

build_coordinator_sync_saga: VER = $(shell git log -1 --pretty=%h)
build_coordinator_sync_saga: build_base build_scratch generate
	docker build --tag evented-coordinator-sync-saga:$(VER) --build-arg=$(VER) -f ./applications/coordinators/grpc/saga/Dockerfile  .

#build_sample_business_logic: VER = $(shell git log -1 --pretty=%h)
build_sample_business_logic: build_base build_scratch generate
	docker build --tag evented-sample-business-logic:latest --build-arg=latest -f ./applications/integrationTest/businessLogic/debug.dockerfile .

configuration_load_command_handler:
	consul kv put commandHandler @applications/commandHandler/configuration/sample.yaml


consul_ui:
	kubectl port-forward service/consul-headless 8500:8500