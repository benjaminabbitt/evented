
generate:
	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f todo/todo.proto -l go -o gen

build-base:
	docker build --tag evented-base -f ./evented-base.dockerfile .

build-scratch:
	docker build --tag scratch-foundation -f ./scratch-foundation.dockerfile .

build: VER := $(shell python ./devops/make/version/get-version.py)
build: DT := $(shell python ./devops/make/get-datetime/get-datetime.py)
build: build-base build-scratch
	docker build --tag todo:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/todo/dockerfile .

load: build
	minikube --v=2 image rm todo
	minikube --v=2 image load todo

run:
	docker run
