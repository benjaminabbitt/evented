topdir=$(shell git rev-parse --show-toplevel)

DT=`go run ${topdir}/applications/support/build_support/ utcNow`

human_version = 0.0.0
VER=`${topdir}/devops/make/version.sh ${human_version}`
