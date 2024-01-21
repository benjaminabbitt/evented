include devops/make/support.mk
include devops/make/gomock.mk
include devops/make/proto.mk
include devops/make/linux.mk
#.DEFAULT-GOAL := build_support
# Notes: This file is designed for developer setup and execution of development environments.   Proper security should be undertaken and is *not* done here for development expendiency.


build: build-command-handler build-query-handler build-sample-business-logic
build-debug: build-command-handler-debug


BUILD_BASE_IMAGE="evented-base"
build-base:
	docker build --tag ${DOCKER_HOST}/${RUNTIME_BASE} -f ${TOPDIR}/devops/evented-base.dockerfile .

BUILD_BASE_IMAGE_DEBUG="evented-base-debug"
build-base-debug: build-base
	docker build --tag ${DOCKER_HOST}/${RUNTIME_BASE_DEBUG} -f ${TOPDIR}/devops/evented-base.debug.dockerfile .

GODOG_IMAGE=${DOCKER_HOST}/godog
build-godog:
	docker build --tag ${GODOG_IMAGE} -f ${TOPDIR}/devops/godog/godog.dockerfile .

SCRATCH_IMAGE=${DOCKER_HOST}/scratch-foundation
build-scratch:
	docker build --tag ${SCRATCH_IMAGE} -f ${TOPDIR}/devops/dockerfiles/scratch-foundation.dockerfile .

##Test
test:
	go test ./...

test-service-integration:
	go test ./... -tags serviceIntegration

vet:
	go vet ./...


# Command Handler
debug-deploy-command-handler: build-command-handler-debug build-sample-business-logic-debug build-command-handler build-sample-business-logic windows-configuration-load-command-handler windows-configuration-load-sample-business-logic deploy-command-handler

#deploy-command-handler: windows-configuration-load-command-handler windows-configuration-load-sample-business-logic minikube-load-evented-command-handler minikube-load-evented-sample-business-logic
ch-deploy: ch-undeploy ch-build ch-minikube-load chsa-build chsa-minikube-load
	helm install sample-ch ${TOPDIR}/applications/command/command-handler/helm/evented-command-handler --debug

ch-undeploy:
	-helm uninstall sample-ch ${TOPDIR}/applications/command/command-handler/helm/evented-command-handler --debug

FRAMEWORK_COMMAND_HANDLER_IMAGE=evented-command-handler
ch-build:build-base build-scratch generate
	docker build --tag ${FRAMEWORK_COMMAND_HANDLER_IMAGE}:${VER} --tag ${FRAMEWORK_COMMAND_HANDLER_IMAGE}:latest --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" --build-arg="RUNTIME_IMAGE={PRODUCTION_IMAGE_BASE}"-f ./applications/command/command-handler/dockerfile .

ch-bounce:
	kubectl delete pods -l evented=sample-command-handler

ch-build-debug:build-base build-scratch
	docker build --tag ${FRAMEWORK_COMMAND_HANDLER_IMAGE}:latest --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" --build-arg="RUNTIME_IMAGE=${DEBUG_IMAGE_BASE}"-f ./applications/command/command-handler/debug.dockerfile .

ch-logs:
	kubectl logs -l app.kubernetes.io/name=evtd-command-handler --all-containers=true --tail=-1

ch-minikube-load: ch-build ch-undeploy
	minikube --v=2 --alsologtostderr image load $(FRAMEWORK_COMMAND_HANDLER_IMAGE)
#
## Sample Business Logic
#chsa-build: build-base build-scratch
#	docker build --tag evented-sample-business-logic:$(VER) --tag evented-sample-business-logic:latest --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/sample-business-logic/dockerfile  .
#
#chsa-build-debug-debug: build-base build-scratch
#	docker build --tag evented-sample-business-logic:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/sample-business-logic/debug.dockerfile  .
#
#chsa-ps-configuration-load:
#	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-sample-business-logic .\applications\command\sample-business-logic\configuration\sample.yaml
#
#chsa-minikube-load: chsa-build ch-undeploy
#	minikube --v=2 --alsologtostderr image load evented-sample-business-logic
#
## Query Handler
#qh-deploy-scratch: qh-build qh-ps-configuration-load qh-minikube-load qh-deploy
#
#qh-deploy: qh-undeploy
#	helm install sample-qh ./applications/command/query-handler/helm/evented-query-handler --debug
#
#qh-undeploy:
#	-helm delete sample-qh
#
#qh-minikube-load: qh-build qh-undeploy
#	minikube --v=2 --alsologtostderr image load evented-query-handler
#
#qh-ps-configuration-load:
#	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-query-handler .\applications\command\query-handler\configuration\sample.yaml
#
#qh-build: build-base build-scratch
#	docker build --tag evented-query-handler:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/query-handler/dockerfile  .
#
#qh-build-debug: build-base build-scratch
#	docker build --tag evented-query-handler:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/command/query-handler/debug.dockerfile  .
#
#qh-bounce:
#	kubectl delete pods -l evented=query-handler
#
#qh-logs:
#	kubectl logs -l evented=query-handler --tail=100
#
##qh_service_port := $(shell python -c "import yaml; print(yaml.safe_load(open('./applications/command/query-handler/helm/evented-query-handler/values.yaml'))['port'])")
#qh-service-expose:
#	kubectl port-forward svc/sample-qh-evented-query-handler $(qh_service_port)
#
## Projector
#pr_port := 1315
#pr-scratch-deploy: pr-build prsa-build pr-ps-configuration-load prsa-ps-configuration-load pr-deploy
#
#pr-build: build-base build-scratch
#	docker build --tag evented-projector:$(VER) --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/projector/dockerfile  .
#
#pr-build-debug: build-base build-scratch
#	docker build --tag evented-projector:$(VER)-DEBUG --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/projector/debug.dockerfile  .
#
#prsa-build: build-base build-scratch
#	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/sample-projector/dockerfile .
#
#prsa-build-debug: build-base build-scratch
#	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/sample-projector/debug.dockerfile .
#
#pr-deploy: pr-undeploy pr-minikube-load prsa-minikube-load
#	helm install sample-pr ./applications/event/projector/helm/evented-projector --debug
#
#pr-undeploy:
#	-helm delete sample-pr
#
#pr-bounce:
#	kubectl delete pods -l app.kubernetes.io/name=evented-projector
#
#pr-ps-configuration-load:
#	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-projector .\applications\event\projector\configuration\sample.yaml
#	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-projector-local .\applications\event\projector\configuration\sample-local.yaml
#
#prsa-ps-configuration-load:
#	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-sample-projector .\applications\event\sample-projector\configuration\sample.yaml
#	powershell -ExecutionPolicy ByPass ./devops/make/config/configuration-load-port-forward.ps1 evented-sample-projector-local .\applications\event\sample-projector\configuration\sample-local.yaml
#
#pr-logs:
#	kubectl logs -l app.kubernetes.io/name=evented-projector --all-containers=true --tail=-1
#
#prsa-expose:
#	kubectl port-forward svc/evented-sample-projector 30003
#
#prsa-minikube-load: prsa-build pr-undeploy
#	minikube --v=2 --alsologtostderr image rm image evented-sample-projector
#	minikube --v=2 --alsologtostderr image load evented-sample-projector
#
#pr-minikube-load: pr-build pr-undeploy
#	minikube --v=2 --alsologtostderr image rm image evented-projector
#	minikube --v=2 --alsologtostderr image load evented-projector

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

