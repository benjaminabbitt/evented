.DEFAULT_GOAL := build

# Notes: This file is designed for developer setup and execution of development environments.   Proper security should be undertaken and is *not* done here for development expendiency.

build: build_command_handler build_query_handler build_coordinator_async_projector build_coordinator_sync_projector build_coordinator_async_saga build_coordinator_sync_saga build_sample_business_logic
build_debug: build_command_handler_debug
load_all: configuration_load_command_handler

stage:
	docker pull namely/protoc-all

generate:
	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/evented.proto -l go -o gen

build_base:
	docker build --tag evented-base -f ./evented-base.dockerfile .

build_scratch:
	docker build --tag scratch-foundation -f ./scratch-foundation.dockerfile .



# Command Handler
deploy_command_handler:
	kubectl apply -f applications/commandHandler/commandHandler.yaml

build_command_handler: VER = $(shell python ./devops/support/version/get-version.py)
build_command_handler: DT = $(shell python -c "from datetime import datetime; print(datetime.now().strftime('%Y-%m-%dT%H:%M:%S.%f%z'))")
build_command_handler:build_base build_scratch generate
	docker build --tag evented-commandhandler:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/commandHandler/dockerfile .

bounce_command_handler:
	kubectl delete pods -l evented=command-handler

build_command_handler_debug: VER = $(shell python ./devops/support/version/get-version.py)
build_command_handler_debug: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build_command_handler_debug:build_base build_scratch generate
	docker build --tag evented-commandhandler:latest --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/commandHandler/debug.dockerfile .

configuration_load_command_handler:
	consul kv put -http-addr=localhost:8500 evented-command-handler @applications/commandHandler/configuration/sample.yaml

logs_command_handler:
	kubectl logs -l evented=command-handler --tail=100

# Query Handler
deploy_query_handler:
	kubectl apply -f applications/eventQueryHandler/eventQueryHandler.yaml

configuration_load_query_handler:
	consul kv put -http-addr=localhost:8500 evented-query-handler @applications/eventQueryHandler/configuration/sample.yaml

build_query_handler: VER = $(shell python ./devops/support/version/get-version.py)
build_query_handler: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build_query_handler: build_base build_scratch generate
	docker build --tag evented-queryhandler:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/eventQueryHandler/dockerfile  .

build_query_handler_debug: VER = $(shell python ./devops/support/version/get-version.py)
build_query_handler_debug: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build_query_handler_debug: build_base build_scratch generate
	docker build --tag evented-queryhandler:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/eventQueryHandler/debug.dockerfile  .

bounce_query_handler:
	kubectl delete pods -l evented=query-handler

logs_query_handler:
	kubectl logs -l evented=query-handler --tail=100

# Coordinator Async Projector
deploy_coordinator_amqp_projector:
	kubectl apply -f applications/coordinators/amqp/projector/amqp-projector-coordinator.yaml

build_coordinator_amqp_projector: VER = $(shell python ./devops/support/version/get-version.py)
build_coordinator_amqp_projector: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build_coordinator_amqp_projector: build_base build_scratch generate
	docker build --tag evented-coordinator-amqp-projector:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/coordinators/amqp/projector/dockerfile  .

build_coordinator_amqp_projector_debug: VER = $(shell python ./devops/support/version/get-version.py)
build_coordinator_amqp_projector_debug: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build_coordinator_amqp_projector_debug: build_base build_scratch generate
	docker build --tag evented-amqp-projector:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/coordinators/amqp/projector/debug.dockerfile  .

configuration_load_amqp_projector:
	consul kv put -http-addr=localhost:8500 evented-amqp-projector @applications/coordinators/amqp/projector/configuration/sample.yaml

logs_amqp_projector:
	kubectl logs -l evented=amqp-projector --tail=100


## Coordinator Async Saga
#deploy_coordinator_async_saga:
#	kubectl apply -f applications/coordinators/amqp/saga/amqp-saga-coordinator.yaml
#
#build_coordinator_async_saga: VER = $(shell git log -1 --pretty=%h)
#build_coordinator_async_saga: build_base build_scratch generate
#	docker build --tag evented-coordinator-async-saga:$(VER) --build-arg=$(VER) -f ./applications/coordinators/amqp/saga/Dockerfile  .
#
#
#
## Coordinator Sync Projector
#deploy_coordinator_sync_projector:
#	kubectl apply -f applications/coordinators/grpc/projector/grpc-projector-coordinator.yaml
#
#build_coordinator_sync_projector: VER = $(shell git log -1 --pretty=%h)
#build_coordinator_sync_projector: build_base build_scratch generate
#	docker build --tag evented-coordinator-sync-projector:$(VER) --build-arg=$(VER) -f ./applications/coordinators/grpc/projector/Dockerfile  .
#
#
#
## Coordinator Sync Saga
#deploy_coordinator_sync_saga:
#	kubectl apply -f applications/coordinators/grpc/saga/grpc-saga-coordinator.yaml
#
#build_coordinator_sync_saga: VER = $(shell git log -1 --pretty=%h)
#build_coordinator_sync_saga: build_base build_scratch generate
#	docker build --tag evented-coordinator-sync-saga:$(VER) --build-arg=$(VER) -f ./applications/coordinators/grpc/saga/Dockerfile  .
#
#

# Sample Business Logic
deploy_sample_business_logic:
	kubectl apply -f applications/integrationTest/businessLogic/businessLogic.yaml

build_sample_business_logic: VER = $(shell python ./devops/support/version/get-version.py)
build_sample_business_logic: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build_sample_business_logic: build_base build_scratch generate
	docker build --tag evented-sample_business_logic:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/integrationTest/businessLogic/dockerfile  .

build_sample_business_logic_debug: VER = $(shell python ./devops/support/version/get-version.py)
build_sample_business_logic_debug: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build_sample_business_logic_debug: build_base build_scratch generate
	docker build --tag evented-sample_business_logic:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/integrationTest/businessLogic/debug.dockerfile  .

bounce_sample_business_logic:
	kubectl delete pods -l evented=sample-business-logic

configuration_load_sample_business_logic:
	consul kv put -http-addr=localhost:8500 evented-sample-business-logic @applications/integrationTest/businessLogic/configuration/sample.yaml

logs_sample_business_logic:
	kubectl logs -l evented=sample-business-logic --tail=100

# Sample Projector
deploy_sample_projector:
	kubectl apply -f applications/integrationTest/projector/projector.yaml

build_sample_projector: VER = $(shell python ./devops/support/version/get-version.py)
build_sample_projector: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build_sample_projector: build_base build_scratch generate
	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/integrationTest/projector/dockerfile .

build_sample_projector_debug: VER = $(shell python ./devops/support/version/get-version.py)
build_sample_projector_debug: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build_sample_projector_debug: build_base build_scratch generate
	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/integrationTest/projector/debug.dockerfile .

bounce_sample_projector:
	kubectl delete pods -l evented=sample-projector

configuration_load_sample_projector:
	consul kv put -http-addr=localhost:8500 evented-sample-projector @applications/integrationTest/projector/configuration/sample.yaml

logs_sample_projector:
	kubectl logs -l evented=sample-projector --tail=100


# Sample Saga
deploy_sample_saga:
	kubectl apply -f applications/integrationTest/saga/saga.yaml

build_sample_saga: build_base build_scratch generate
	docker build --tag evented-sample-saga:latest --build-arg=latest -f ./applications/integrationTest/saga/debug.dockerfile .


## Developer setup
setup: install_consul rabbit_install mongo_install

## Consul Shortcuts
install_consul:
	helm repo add hashicorp https://helm.releases.hashicorp.com
	helm repo update
	helm install consul hashicorp/consul --wait --set global.name=consul --values ./devops/helm/consul/values.yaml

consul_ui_expose:
	kubectl port-forward svc/consul-ui 80

consul_service_expose:
	kubectl port-forward svc/consul-server 8500


## RabbitMQ Shortcuts
rabbit_install:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install rabbitmq bitnami/rabbitmq --wait --values ./devops/helm/rabbitmq/values.yaml

rabbit_ui_expose:
	kubectl port-forward svc/rabbitmq 15672:15672

rabbit_extract_password:
	@python devops/support/helm/rabbitmq/get-secret.py --secret="rabbitmq-password"

rabbit_extract_cookie:
	@python devops/support/helm/rabbitmq/get-secret.py --secret="rabbitmq-erlang-cookie"


## Mongo Shortcuts
mongo_service_expose:
	kubectl port-forward --namespace default svc/mongodb 27017:27017

mongo_extract_password:
	@python devops/support/get-secret/get-secret.py --namespace="default" --name="mongodb" --secret="mongodb-root-password"

mongo_install:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install mongodb bitnami/mongodb --wait --values ./devops/helm/mongodb/values.yaml
