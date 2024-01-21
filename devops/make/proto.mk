include devops/make/support.mk

PROTO_IMAGE=$(DOCKER_HOST)/proto:latest
build-proto-image:
	docker build -t $(PROTO_IMAGE) -f ${TOPDIR}/devops/proto/Dockerfile $(TOPDIR)

generate-proto: build-proto-image
	mkdir -p "$(TOPDIR)/generated/proto"
	docker run --volume $(TOPDIR):$(WORKSPACE) $(PROTO_IMAGE) --go_out=$(WORKSPACE)/generated/proto/ --go-grpc_out=$(WORKSPACE)/generated/proto -I=$(WORKSPACE)/proto $(WORKSPACE)/proto/evented/evented.proto
