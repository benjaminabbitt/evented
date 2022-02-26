.DEFAULT-GOAL := build

# Notes: This file is designed for developer setup and execution of development environments.   Proper security should be undertaken and is *not* done here for development expendiency.

build: build-command-handler build-query-handler build-coordinator-async-projector build-coordinator-sync-projector build-coordinator-async-saga build-coordinator-sync-saga build-sample-business-logic
build-debug: build-command-handler-debug
load-all: configuration-load-command-handler

stage:
	docker pull namely/protoc-all

generate:
	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/evented.proto -l go -o gen

build-base:
	docker build --tag evented-base -f ./evented-base.dockerfile .

build-scratch:
	docker build --tag scratch-foundation -f ./scratch-foundation.dockerfile .



# Command Handler
deploy-command-handler:
	kubectl apply -f applications/commandHandler/commandHandler.yaml

build-command-handler: VER = $(shell python ./devops/support/version/get-version.py)
build-command-handler: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build-command-handler:build-base build-scratch generate
	docker build --tag evented-command-handler:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/commandHandler/dockerfile .

bounce-command-handler:
	kubectl delete pods -l evented=command-handler

build-command-handler-debug: VER = $(shell python ./devops/support/version/get-version.py)
build-command-handler-debug: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build-command-handler-debug:build-base build-scratch generate
	docker build --tag evented-command-handler:latest --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/commandHandler/debug.dockerfile .

configuration-load-command-handler:
	consul kv put -http-addr=localhost:8500 evented-command-handler @applications/commandHandler/configuration/sample.yaml

logs-command-handler:
	kubectl logs -l evented=command-handler --tail=100


# Sample Business Logic
deploy-sample-business-logic:
	kubectl apply -f applications/integrationTest/businessLogic/businessLogic.yaml

build-sample-business-logic: VER = $(shell python ./devops/support/version/get-version.py)
build-sample-business-logic: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build-sample-business-logic: build-base build-scratch generate
	docker build --tag evented-sample-business-logic:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/integrationTest/businessLogic/dockerfile  .

build-sample-business-logic-debug: VER = $(shell python ./devops/support/version/get-version.py)
build-sample-business-logic-debug: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build-sample-business-logic-debug: build-base build-scratch generate
	docker build --tag evented-sample-business-logic:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/integrationTest/businessLogic/debug.dockerfile  .

bounce-sample-business-logic:
	kubectl delete pods -l evented=sample-business-logic

configuration-load-sample-business-logic:
	consul kv put -http-addr=localhost:8500 evented-sample-business-logic @applications/integrationTest/businessLogic/configuration/sample.yaml

logs-sample-business-logic:
	kubectl logs -l evented=sample-business-logic --tail=100


build-command-handler-complex: build-sample-business-logic build-command-handler


# Query Handler
deploy-query-handler:
	kubectl apply -f applications/queryHandler/queryHandler.yaml

configuration-load-query-handler:
	consul kv put -http-addr=localhost:8500 evented-query-handler @applications/queryHandler/configuration/sample.yaml

build-query-handler: VER = $(shell python ./devops/support/version/get-version.py)
build-query-handler: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build-query-handler: build-base build-scratch generate
	docker build --tag evented-query-handler:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/queryHandler/dockerfile  .

build-query-handler-debug: VER = $(shell python ./devops/support/version/get-version.py)
build-query-handler-debug: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build-query-handler-debug: build-base build-scratch generate
	docker build --tag evented-queryhandler:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/queryHandler/debug.dockerfile  .

bounce-query-handler:
	kubectl delete pods -l evented=query-handler

logs-query-handler:
	kubectl logs -l evented=query-handler --tail=100

# Coordinator Async Projector
deploy-coordinator-projector-amqp:
	kubectl apply -f applications/coordinators/amqp/projector/coordinator-projector-amqp.yaml

build-coordinator-projector-amqp: VER = $(shell python ./devops/support/version/get-version.py)
build-coordinator-projector-amqp: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build-coordinator-projector-amqp: build-base build-scratch generate
	docker build --tag evented-coordinator-projector-amqp:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/coordinators/amqp/projector/dockerfile  .

build-coordinator-projector-amqp-debug: VER = $(shell python ./devops/support/version/get-version.py)
build-coordinator-projector-amqp-debug: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build-coordinator-projector-amqp-debug: build-base build-scratch generate
	docker build --tag evented-coordinator-projector-amqp:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/coordinators/amqp/projector/debug.dockerfile  .

configuration-load-coordinator-projector-amqp:
	consul kv put -http-addr=localhost:8500 evented-coordinator-projector-amqp @applications/coordinators/amqp/projector/configuration/sample.yaml

logs-coordinator-projector-amqp:
	kubectl logs -l evented=coordinator-projector-amqp --tail=100

bounce-coordinator-projector-amqp:
	kubectl delete pods -l evented=coordinator-projector-amqp



## Coordinator Async Saga
#deploy_coordinator-async-placeholder-saga:
#	kubectl apply -f applications/coordinators/amqp/placeholder-saga/amqp-placeholder-saga-coordinator.yaml
#
#build-coordinator-async-placeholder-saga: VER = $(shell git log -1 --pretty=%h)
#build-coordinator-async-placeholder-saga: build-base build-scratch generate
#	docker build --tag evented-coordinator-async-placeholder-saga:$(VER) --build-arg=$(VER) -f ./applications/coordinators/amqp/placeholder-saga/Dockerfile  .
#
#
#
## Coordinator Sync Projector
#deploy-coordinator-sync-placeholder-projector:
#	kubectl apply -f applications/coordinators/grpc/placeholder-projector/grpc-placeholder-projector-coordinator.yaml
#
#build-coordinator-sync-placeholder-projector: VER = $(shell git log -1 --pretty=%h)
#build-coordinator-sync-placeholder-projector: build-base build-scratch generate
#	docker build --tag evented-coordinator-sync-placeholder-projector:$(VER) --build-arg=$(VER) -f ./applications/coordinators/grpc/placeholder-projector/Dockerfile  .
#
#
#
## Coordinator Sync Saga
#deploy-coordinator-sync-placeholder-saga:
#	kubectl apply -f applications/coordinators/grpc/placeholder-saga/grpc-placeholder-saga-coordinator.yaml
#
#build-coordinator-sync-placeholder-saga: VER = $(shell git log -1 --pretty=%h)
#build-coordinator-sync-placeholder-saga: build-base build-scratch generate
#	docker build --tag evented-coordinator-sync-placeholder-saga:$(VER) --build-arg=$(VER) -f ./applications/coordinators/grpc/placeholder-saga/Dockerfile  .
#
#


# Sample Projector
deploy-sample-projector:
	kubectl apply -f applications/integrationTest/projector/projector.yaml

build-sample-projector: VER = $(shell python ./devops/support/version/get-version.py)
build-sample-projector: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build-sample-projector: build-base build-scratch generate
	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/integrationTest/projector/dockerfile .

build-sample-projector-debug: VER = $(shell python ./devops/support/version/get-version.py)
build-sample-projector-debug: DT = $(shell python ./devops/support/get-datetime/get-datetime.py)
build-sample-projector-debug: build-base build-scratch generate
	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/integrationTest/projector/debug.dockerfile .

bounce-sample-projector:
	kubectl delete pods -l evented=sample-projector

configuration-load-sample-projector:
	consul kv put -http-addr=localhost:8500 evented-sample-projector @applications/integrationTest/projector/configuration/sample.yaml

logs-sample-projector:
	kubectl logs -l evented=sample-projector --tail=100

sample-projector-expose:
	kubectl port-forward svc/evented-sample-projector 30003

#
## Sample Saga
#deploy-sample-placeholder-saga:
#	kubectl apply -f applications/integrationTest/placeholder-saga/placeholder-saga.yaml
#
#build-sample-placeholder-saga: build-base build-scratch generate
#	docker build --tag evented-sample-placeholder-saga:latest --build-arg=latest -f ./applications/integrationTest/placeholder-saga/debug.dockerfile .


## Developer setup
setup: install-consul rabbit-install mongo-install

## Consul Shortcuts
install-consul:
	helm repo add hashicorp https://helm.releases.hashicorp.com
	helm repo update
	helm install consul hashicorp/consul --wait --set global.name=consul --values ./devops/helm/consul/values.yaml

consul-ui-expose:
	kubectl port-forward svc/consul-ui 80

consul-service-expose:
	kubectl port-forward svc/consul-server 8500


## RabbitMQ Shortcuts
rabbit-install:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install rabbitmq bitnami/rabbitmq --wait --values ./devops/helm/rabbitmq/values.yaml

rabbit-ui-expose:
	kubectl port-forward svc/rabbitmq 15672:15672

rabbit-extract-password:
	@python devops/support/helm/rabbitmq/get-secret.py --secret="rabbitmq-password"

rabbit-extract-cookie:
	@python devops/support/helm/rabbitmq/get-secret.py --secret="rabbitmq-erlang-cookie"


## Mongo Shortcuts
mongo-service-expose:
	kubectl port-forward --namespace default svc/mongodb 27017:27017

mongo-extract-password:
	@python devops/support/get-secret/get-secret.py --namespace="default" --name="mongodb" --secret="mongodb-root-password"

mongo-install:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install mongodb bitnami/mongodb --wait --values ./devops/helm/mongodb/values.yaml
