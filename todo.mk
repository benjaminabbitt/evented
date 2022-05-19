
generate:
	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f todo/todo.proto -l go -o gen

todo-build: VER := $(shell python ./devops/make/version/get-version.py)
todo-build: DT := $(shell python ./devops/make/get-datetime/get-datetime.py)
todo-build: build-base build-scratch
	docker build --tag evented-sample-projector:${VER} --build-arg="BUILD_TIME=${DT}" --build-arg="VERSION=${VER}" -f ./applications/event/sample-projector/dockerfile .
