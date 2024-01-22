include devops/make/support.mk

MOCKGEN_DOCKER_IMAGE = ${DOCKER_HOST}/gomock
build-gomock:
	docker build --tag ${MOCKGEN_DOCKER_IMAGE} -f ${TOPDIR}/devops/gomock/gomock.dockerfile .

mockgen-common = docker run -v ${TOPDIR}:${WORKSPACE} ${MOCKGEN_DOCKER_IMAGE} mockgen
generate-mocks: generate-proto build-gomock
	${mockgen-common} -source ${WORKSPACE}/repository/eventBook/eventBookStorer.go -destination ${WORKSPACE}/repository/eventBook/mocks/eventBookStorer.mock.go
	${mockgen-common} -source ${WORKSPACE}/repository/snapshots/snapshotStorer.go -destination ${WORKSPACE}/repository/snapshots/mocks/snapshotStorer.mock.go
	${mockgen-common} -source ${WORKSPACE}/repository/events/eventRepo.go -destination ${WORKSPACE}/repository/events/mocks/eventRepo.mock.go
	${mockgen-common} -source ${WORKSPACE}/generated/proto/github.com/benjaminabbitt/evented/proto/evented/evented.pb.go -destination ${WORKSPACE}/generated/proto/github.com/benjaminabbitt/evented/proto/evented/mocks/evented.mock.pb.go
	${mockgen-common} -source ${WORKSPACE}/generated/proto/github.com/benjaminabbitt/evented/proto/evented/evented_grpc.pb.go -destination ${WORKSPACE}/generated/proto/github.com/benjaminabbitt/evented/proto/evented/mocks/evented_grpc.mock.pb.go
	${mockgen-common} -source ${WORKSPACE}/applications/command/command-handler/framework/transport/transportHolder.go -destination ${WORKSPACE}/applications/command/command-handler/framework/transport/mocks/transportHolder.mock.go
	${mockgen-common} -source ${WORKSPACE}/applications/command/command-handler/business/client/business.go -destination ${WORKSPACE}/applications/command/command-handler/business/client/mocks/business.mock.go
	${mockgen-common} -source ${WORKSPACE}/transport/async/eventTransporter.go -destination ${WORKSPACE}/transport/async/mocks/eventTransporter.mock.go
	${mockgen-common} -source ${WORKSPACE}/transport/sync/saga/syncSagaTransporter.go -destination ${WORKSPACE}/transport/sync/saga/mocks/syncSagaTransporter.mock.go
	${mockgen-common} -source ${WORKSPACE}/transport/sync/projector/syncProjectionTransporter.go -destination ${WORKSPACE}/transport/sync/projector/mocks/syncProjectionTransporter.mock.go
