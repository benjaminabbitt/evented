include devops/make/support.mk
#.DEFAULT-GOAL := build_support
# Notes: This file is designed for developer setup and execution of development environments.   Proper security should be undertaken and is *not* done here for development expendiency.


build: build-command-handler build-query-handler build-sample-business-logic
build-debug: build-command-handler-debug

generate: install-deps
	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/evented.proto -l go -o gen

build-base:
	docker build --tag evented-base -f ./evented-base.dockerfile .

build-scratch:
	docker build --tag scratch-foundation -f ./scratch-foundation.dockerfile .

##Test
test:
	go test ./...

test-service-integration:
	go test ./... -tags serviceIntegration

vet:
	go vet ./...


install-deps:
	go install github.com/vektra/mockery/v2@latest
	go install github.com/cucumber/godog/cmd/godog@latest
	go install github.com/golang/mock/mockgen@v1.6.0
	docker pull namely/protoc-all

build-proto-image:
	powershell docker build -t proto -f devops/proto/Dockerfile .

generate-proto:
	IF NOT EXIST "$(topdir)/generated" mkdir "$(topdir)/generated"
	IF NOT EXIST "$(topdir)/generated/proto" mkdir "$(topdir)/generated/proto"
	cd "$(topdir)" && powershell docker run --volume .:/workspace/ proto --go_out=/workspace/generated/proto/ --go_grpc_out=/workspace/generated/proto -I=/workspace/proto /workspace/proto/evented/evented.proto

generate-mocks: generate
	mockgen -source .\repository\eventBook\eventBookStorer.go -destination .\repository\eventBook\mocks\eventBookStorer.mock.go
	mockgen -source .\repository\snapshots\snapshotStorer.go -destination .\repository\snapshots\mocks\snapshotStorer.mock.go
	mockgen -source .\repository\events\eventRepo.go -destination .\repository\events\mocks\eventRepo.mock.go
	mockgen -source proto\gen\github.com\benjaminabbitt\evented\proto\evented\evented.pb.go -destination proto\gen\github.com\benjaminabbitt\evented\proto\evented\mocks\evented.mock.pb.go
	mockgen -source applications/command/command-handler/framework/transport/transportHolder.go -destination applications/command/command-handler/framework/transport/mocks/transportHolder.mock.go
	mockgen -source applications/command/command-handler/business/client/business.go -destination applications/command/command-handler/business/client/mocks/business.mock.go
	mockgen -source transport/async/eventTransporter.go -destination transport/async/mocks/eventTransporter.mock.go
	mockgen -source transport/sync/saga/syncSagaTransporter.go -destination transport/sync/saga/mocks/syncSagaTransporter.mock.go
	mockgen -source transport/sync/projector/syncProjectionTransporter.go -destination transport/sync/projector/mocks/syncProjectionTransporter.mock.go

# Command Handler
debug-deploy-command-handler: build-command-handler-debug build-sample-business-logic-debug build-command-handler build-sample-business-logic windows-configuration-load-command-handler windows-configuration-load-sample-business-logic deploy-command-handler

#deploy-command-handler: windows-configuration-load-command-handler windows-configuration-load-sample-business-logic minikube-load-evented-command-handler minikube-load-evented-sample-business-logic
ch-deploy: ch-undeploy ch-build ch-minikube-load chsa-build chsa-minikube-load
	helm install sample-ch ./applications/command/command-handler/helm/evented-command-handler --debug

ch-undeploy:
	-helm uninstall sample-ch ./applications/command/command-handler/helm/evented-command-handler --debug

ch-build:build-base build-scratch generate generate-mocks
	docker build --tag evented-command-handler:${VER} --tag evented-command-handler:latest --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/command-handler/dockerfile .

ch-bounce:
	kubectl delete pods -l evented=sample-command-handler

ch-build-debug:build-base build-scratch
	docker build --tag evented-command-handler:latest --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/command-handler/debug.dockerfile .

ch-ps-configuration-load:
	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-command-handler applications\command\command-handler\configuration\sample.yaml

ch-logs:
	kubectl logs -l app.kubernetes.io/name=evtd-command-handler --all-containers=true --tail=-1

ch-minikube-load: ch-build ch-undeploy
	minikube --v=2 --alsologtostderr image load evented-command-handler

# Sample Business Logic
chsa-build: build-base build-scratch
	docker build --tag evented-sample-business-logic:$(VER) --tag evented-sample-business-logic:latest --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/sample-business-logic/dockerfile  .

chsa-build-debug-debug: build-base build-scratch
	docker build --tag evented-sample-business-logic:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/sample-business-logic/debug.dockerfile  .

chsa-ps-configuration-load:
	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-sample-business-logic .\applications\command\sample-business-logic\configuration\sample.yaml

chsa-minikube-load: chsa-build ch-undeploy
	minikube --v=2 --alsologtostderr image load evented-sample-business-logic

# Query Handler
qh-deploy-scratch: qh-build qh-ps-configuration-load qh-minikube-load qh-deploy

qh-deploy: qh-undeploy
	helm install sample-qh ./applications/command/query-handler/helm/evented-query-handler --debug

qh-undeploy:
	-helm delete sample-qh

qh-minikube-load: qh-build qh-undeploy
	minikube --v=2 --alsologtostderr image load evented-query-handler

qh-ps-configuration-load:
	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-query-handler .\applications\command\query-handler\configuration\sample.yaml

qh-build: build-base build-scratch
	docker build --tag evented-query-handler:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/query-handler/dockerfile  .

qh-build-debug: build-base build-scratch
	docker build --tag evented-query-handler:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/query-handler/debug.dockerfile  .

qh-bounce:
	kubectl delete pods -l evented=query-handler

qh-logs:
	kubectl logs -l evented=query-handler --tail=100

#qh_service_port := $(shell python -c "import yaml; print(yaml.safe_load(open('./applications/command/query-handler/helm/evented-query-handler/values.yaml'))['port'])")
qh-service-expose:
	kubectl port-forward svc/sample-qh-evented-query-handler $(qh_service_port)

# Projector
pr_port := 1315
pr-scratch-deploy: pr-build prsa-build pr-ps-configuration-load prsa-ps-configuration-load pr-deploy

pr-build: build-base build-scratch
	docker build --tag evented-projector:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/projector/dockerfile  .

pr-build-debug: build-base build-scratch
	docker build --tag evented-projector:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/projector/debug.dockerfile  .

prsa-build: build-base build-scratch
	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/sample-projector/dockerfile .

prsa-build-debug: build-base build-scratch
	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/sample-projector/debug.dockerfile .

pr-deploy: pr-undeploy pr-minikube-load prsa-minikube-load
	helm install sample-pr ./applications/event/projector/helm/evented-projector --debug

pr-undeploy:
	-helm delete sample-pr

pr-bounce:
	kubectl delete pods -l app.kubernetes.io/name=evented-projector

pr-ps-configuration-load:
	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-projector .\applications\event\projector\configuration\sample.yaml
	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-projector-local .\applications\event\projector\configuration\sample-local.yaml

prsa-ps-configuration-load:
	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-sample-projector .\applications\event\sample-projector\configuration\sample.yaml
	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-sample-projector-local .\applications\event\sample-projector\configuration\sample-local.yaml

pr-logs:
	kubectl logs -l app.kubernetes.io/name=evented-projector --all-containers=true --tail=-1

prsa-expose:
	kubectl port-forward svc/evented-sample-projector 30003

prsa-minikube-load: prsa-build pr-undeploy
	minikube --v=2 --alsologtostderr image rm image evented-sample-projector
	minikube --v=2 --alsologtostderr image load evented-sample-projector

pr-minikube-load: pr-build pr-undeploy
	minikube --v=2 --alsologtostderr image rm image evented-projector
	minikube --v=2 --alsologtostderr image load evented-projector

#
###  Saga
#deploy-coordinator-saga:
#	ls
#
#build_support-coordinator-saga: VER := $(shell git log -1 --pretty=%h)
#build_support-coordinator-saga: build_support-base build_support-scratch
#	docker build_support --tag evented-coordinator-saga:$(VER) --build_support-arg=$(VER) -f ./applications/event/projector/dockerfile  .
#
#configuration-load-coordinator-saga:
#	consul kv put -http-addr=localhost:8500 evented-saga @applications/event/saga/configuration/sample.yaml

#
#
## Coordinator Sync Projector
#deploy-coordinator-sync-sample-projector:
#	kubectl apply -f applications/coordinators/grpc/sample-projector/grpc-sample-projector-coordinator.yaml
#
#build_support-coordinator-sync-sample-projector: VER := $(shell git log -1 --pretty=%h)
#build_support-coordinator-sync-sample-projector: build_support-base build_support-scratch
#	docker build_support --tag evented-coordinator-sync-sample-projector:$(VER) --build_support-arg=$(VER) -f ./applications/coordinators/grpc/sample-projector/Dockerfile  .
#
#
#
## Coordinator Sync Saga
#deploy-coordinator-sync-sample-saga:
#	kubectl apply -f applications/coordinators/grpc/sample-saga/grpc-sample-saga-coordinator.yaml
#
#build_support-coordinator-sync-sample-saga: VER := $(shell git log -1 --pretty=%h)
#build_support-coordinator-sync-sample-saga: build_support-base build_support-scratch
#	docker build_support --tag evented-coordinator-sync-sample-saga:$(VER) --build_support-arg=$(VER) -f ./applications/coordinators/grpc/sample-saga/Dockerfile  .
#
#




#
## Sample Saga
#deploy-sample-sample-saga:
#	kubectl apply -f applications/integrationTest/sample-saga/sample-saga.yaml
#
#build_support-sample-sample-saga: build_support-base build_support-scratch
#	docker build_support --tag evented-sample-sample-saga:latest --build_support-arg=latest -f ./applications/integrationTest/sample-saga/debug.dockerfile .


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
	kubectl port-forward svc/evented-rabbitmq 15672:15672

rabbit-service-expose:
	kubectl port-forward svc/evented-rabbitmq 5672:5672

rabbit-extract-cookie:
	@python devops/make/get-secret/get-secret.py --secret="rabbitmq-erlang-cookie"

rabbit-extract-pw:
	@python devops/make/get-secret/get-secret.py --secret="rabbitmq-password"

## Mongo Shortcuts
mongo-service-expose:
	kubectl port-forward --namespace default svc/evented-mongodb 27017:27017

install-mongo:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install evented-mongodb bitnami/mongodb --wait --values ./devops/helm/mongodb/values.yaml

install-k8s-services: install-mongo install-consul install-rabbit

install-helm:
	curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
	sudo apt-get install apt-transport-https --yes
	echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
	sudo apt-get update
	sudo apt-get install helm

install-minikube:
	curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb
	sudo dpkg -i minikube_latest_amd64.deb

install-pyenv:
	echo "This is going to look like it's doing some weird things.  We use pipenv, which has a hidden dependency for the debian not-packaged pyenv"
	sudo apt install -y make build-essential libssl-dev zlib1g-dev libbz2-dev libreadline-dev libsqlite3-dev wget curl llvm libncurses5-dev libncursesw5-dev xz-utils tk-dev libffi-dev liblzma-dev git
	git clone https://github.com/pyenv/pyenv.git ~/.pyenv
	echo export PYENV_ROOT="${HOME}/.pyenv" >> ${RC_FILE}
	echo export PATH="${$HOME}/.pyenv/bin:${PATH}" >> ${RC_FILE}
	eval "$(pyenv init --path)"

install-python: install-pyenv
	sudo apt-get update
	sudo apt-get install -y python3 python3-pip pipenv

install-os-dependencies: install-minikube install-helm

## Minikube Shortcuts
minikube:
	MINIKUBE_ROOTLESS=false minikube start --feature-gates=GRPCContainerProbe=true --memory=12288 --cpus 4

minikube_enable_lb:
	minikube tunnel

human_version = 0.0.0
version:
	@go run ${topdir}/applications/support/build_support/ hashed_version ${topdir} ${human_version}