TOPDIR=$(shell git rev-parse --show-toplevel)

DT=$(shell go run ${TOPDIR}/applications/support/build_support/ utcNow)

HUMAN_VERSION=0.0.0
VER=`${TOPDIR}/devops/make/version.sh ${HUMAN_VERSION}`

#override docker.io default URI, replace with the actual repository in time
DOCKER_HOST=localhost:0

#this is the location that the application code is mounted to in the docker containers used for running build tooling
WORKSPACE=/workspace
