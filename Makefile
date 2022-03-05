.DEFAULT-GOAL := build

# Notes: This file is designed for developer setup and execution of development environments.   Proper security should be undertaken and is *not* done here for development expendiency.

build: build-command-handler build-query-handler build-coordinator-async-projector build-coordinator-sync-projector build-coordinator-async-saga build-coordinator-sync-saga build-sample-business-logic
build-debug: build-command-handler-debug
load-all: configuration-load-command-handler

generate: install-deps
	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/evented.proto -l go -o gen

build-base:
	docker build --tag evented-base -f ./evented-base.dockerfile .

build-scratch:
	docker build --tag scratch-foundation -f ./scratch-foundation.dockerfile .

##Test
test:
	go test ./...

vet:
	go vet ./...


install-deps:
	go install github.com/vektra/mockery/v2@latest
	go install github.com/cucumber/godog/cmd/godog@latest
	go install github.com/golang/mock/mockgen@v1.6.0
	docker pull namely/protoc-all

generate-mocks:
	mockgen -source .\repository\eventBook\eventBookStorer.go -destination .\repository\eventBook\mocks\eventBookStorer-mock.go
	mockgen -source .\repository\snapshots\snapshotStorer.go -destination .\repository\snapshots\mocks\snapshotStorer-mock.go
	mockgen -source .\repository\events\eventRepo.go -destination .\repository\events\mocks\eventRepo-mock.go

# Command Handler
scratch-deploy-command-handler: build-command-handler build-sample-business-logic configuration-load-command-handler configuration-load-sample-business-logic deploy-command-handler
debug-deploy-command-handler: build-command-handler-debug build-sample-business-logic-debug build-command-handler build-sample-business-logic configuration-load-command-handler configuration-load-sample-business-logic deploy-command-handler

deploy-command-handler:
	-helm delete sample-command-handler-deployment
	helm install sample-command-handler-deployment ./applications/command/command-handler/helm/evented-command-handler --debug

build-command-handler: VER := $(shell python ./devops/support/version/get-version.py)
build-command-handler: DT := $(shell python ./devops/support/get-datetime/get-datetime.py)
build-command-handler:build-base build-scratch generate generate-mocks
	docker build --tag evented-command-handler:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/command-handler/dockerfile .

bounce-command-handler:
	kubectl delete pods -l evented=command-handler

build-command-handler-debug: VER := $(shell python ./devops/support/version/get-version.py)
build-command-handler-debug: DT := $(shell python ./devops/support/get-datetime/get-datetime.py)
build-command-handler-debug:build-base build-scratch generate generate-mocks
	docker build --tag evented-command-handler:latest --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/command-handler/debug.dockerfile .

configuration-load-command-handler:
	consul kv put -http-addr=localhost:8500 evented-command-handler @applications/command/command-handler/configuration/sample.yaml

logs-command-handler:
	kubectl logs -l app.kubernetes.io/name=evtd-command-handler --all-containers=true --tail=-1


# Sample Business Logic
deploy-sample-business-logic:
	kubectl apply -f applications/integrationTest/businessLogic/businessLogic.yaml

build-sample-business-logic: VER := $(shell python ./devops/support/version/get-version.py)
build-sample-business-logic: DT := $(shell python ./devops/support/get-datetime/get-datetime.py)
build-sample-business-logic: build-base build-scratch generate generate-mocks
	docker build --tag evented-sample-business-logic:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/sample-business-logic/dockerfile  .

build-sample-business-logic-debug: VER := $(shell python ./devops/support/version/get-version.py)
build-sample-business-logic-debug: DT := $(shell python ./devops/support/get-datetime/get-datetime.py)
build-sample-business-logic-debug: build-base build-scratch generate generate-mocks
	docker build --tag evented-sample-business-logic:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/sample-business-logic/debug.dockerfile  .

configuration-load-sample-business-logic:
	consul kv put -http-addr=localhost:8500 evented-sample-business-logic @applications/command/sample-business-logic/configuration/sample.yaml



# Query Handler
scratch-deploy-query-handler: build-query-handler configuration-load-query-handler deploy-query-handler

deploy-query-handler:
	-helm delete sample-query-handler-deployment
	helm install sample-query-handler-deployment ./applications/command/query-handler/helm/evented-query-handler --debug

configuration-load-query-handler:
	consul kv put -http-addr=localhost:8500 evented-query-handler @applications/command/query-handler/configuration/sample.yaml

build-query-handler: VER := $(shell python ./devops/support/version/get-version.py)
build-query-handler: DT := $(shell python ./devops/support/get-datetime/get-datetime.py)
build-query-handler: build-base build-scratch generate generate-mocks
	docker build --tag evented-query-handler:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/query-handler/dockerfile  .

build-query-handler-debug: VER := $(shell python ./devops/support/version/get-version.py)
build-query-handler-debug: DT := $(shell python ./devops/support/get-datetime/get-datetime.py)
build-query-handler-debug: build-base build-scratch generate generate-mocks
	docker build --tag evented-query-handler:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/query-handler/debug.dockerfile  .

bounce-query-handler:
	kubectl delete pods -l evented=query-handler

logs-query-handler:
	kubectl logs -l evented=query-handler --tail=100

# Projector
scratch-deploy-projector: build-projector build-sample-projector configuration-load-projector configuration-load-sample-projector deploy-projector

build-projector: VER := $(shell python ./devops/support/version/get-version.py)
build-projector: DT := $(shell python ./devops/support/get-datetime/get-datetime.py)
build-projector: build-base build-scratch generate generate-mocks
	docker build --tag evented-projector:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/projector/dockerfile  .

build-projector-debug: VER := $(shell python ./devops/support/version/get-version.py)
build-projector-debug: DT := $(shell python ./devops/support/get-datetime/get-datetime.py)
build-projector-debug: build-base build-scratch generate generate-mocks
	docker build --tag evented-projector:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/projector/debug.dockerfile  .

build-sample-projector: VER := $(shell python ./devops/support/version/get-version.py)
build-sample-projector: DT := $(shell python ./devops/support/get-datetime/get-datetime.py)
build-sample-projector: build-base build-scratch generate generate-mocks
	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/sample-projector/dockerfile .

build-sample-projector-debug: VER := $(shell python ./devops/support/version/get-version.py)
build-sample-projector-debug: DT := $(shell python ./devops/support/get-datetime/get-datetime.py)
build-sample-projector-debug: build-base build-scratch generate generate-mocks
	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/sample-projector/debug.dockerfile .

deploy-projector:
	-helm delete sample-projector-deployment
	helm install sample-projector-deployment ./applications/event/projector/helm/evented-projector --debug

bounce-projector:
	kubectl delete pods -l app.kubernetes.io/name=evented-projector

configuration-load-projector:
	consul kv put -http-addr=localhost:8500 evented-projector @applications/event/projector/configuration/sample.yaml

configuration-load-sample-projector:
	consul kv put -http-addr=localhost:8500 evented-sample-projector @applications/event/sample-projector/configuration/sample.yaml

logs-projector:
	kubectl logs -l app.kubernetes.io/name=evented-projector --all-containers=true --tail=-1

sample-projector-expose:
	kubectl port-forward svc/evented-sample-projector 30003

#
###  Saga
#deploy-coordinator-saga:
#	ls
#
#build-coordinator-saga: VER := $(shell git log -1 --pretty=%h)
#build-coordinator-saga: build-base build-scratch generate
#	docker build --tag evented-coordinator-saga:$(VER) --build-arg=$(VER) -f ./applications/event/projector/dockerfile  .
#
#configuration-load-coordinator-saga:
#	consul kv put -http-addr=localhost:8500 evented-saga @applications/event/saga/configuration/sample.yaml

#
#
## Coordinator Sync Projector
#deploy-coordinator-sync-sample-projector:
#	kubectl apply -f applications/coordinators/grpc/sample-projector/grpc-sample-projector-coordinator.yaml
#
#build-coordinator-sync-sample-projector: VER := $(shell git log -1 --pretty=%h)
#build-coordinator-sync-sample-projector: build-base build-scratch generate
#	docker build --tag evented-coordinator-sync-sample-projector:$(VER) --build-arg=$(VER) -f ./applications/coordinators/grpc/sample-projector/Dockerfile  .
#
#
#
## Coordinator Sync Saga
#deploy-coordinator-sync-sample-saga:
#	kubectl apply -f applications/coordinators/grpc/sample-saga/grpc-sample-saga-coordinator.yaml
#
#build-coordinator-sync-sample-saga: VER := $(shell git log -1 --pretty=%h)
#build-coordinator-sync-sample-saga: build-base build-scratch generate
#	docker build --tag evented-coordinator-sync-sample-saga:$(VER) --build-arg=$(VER) -f ./applications/coordinators/grpc/sample-saga/Dockerfile  .
#
#




#
## Sample Saga
#deploy-sample-sample-saga:
#	kubectl apply -f applications/integrationTest/sample-saga/sample-saga.yaml
#
#build-sample-sample-saga: build-base build-scratch generate
#	docker build --tag evented-sample-sample-saga:latest --build-arg=latest -f ./applications/integrationTest/sample-saga/debug.dockerfile .


## Developer setup
setup: install-consul install-rabbit install-mongo install-deps

## Consul Shortcuts
install-consul:
	helm repo add hashicorp https://helm.releases.hashicorp.com
	helm repo update
	helm install evented-consul hashicorp/consul --wait --set global.name=consul --values ./devops/helm/consul/values.yaml

consul-ui-expose:
	kubectl port-forward svc/consul-ui 80

consul-service-expose:
	kubectl port-forward svc/consul-server 8500


## RabbitMQ Shortcuts
install-rabbit:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install evented-rabbitmq bitnami/rabbitmq --wait --values ./devops/helm/rabbitmq/values.yaml

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

install-mongo:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install evented-mongodb bitnami/mongodb --wait --values ./devops/helm/mongodb/values.yaml

