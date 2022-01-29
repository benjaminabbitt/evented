.DEFAULT_GOAL := generate

stage:
	docker pull namely/protoc-all

generate:
	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/evented.proto -l go -o gen
#	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/business/business/business.proto -l go -o gen
#	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/business/coordinator/business.co.proto -l go -o gen
#	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/business/query/query.proto -l go -o gen
#	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/projector/coordinator/evented.co.proto -l go -o gen
#	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/projector/projector/evented.proto -l go -o gen
#	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/saga/coordinator/saga.co.proto -l go -o gen
#	docker run -v ${CURDIR}/proto:/defs namely/protoc-all -f evented/saga/saga/saga.proto -l go -o gen