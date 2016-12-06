GO ?= go

SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

SGG = save.gg/sgg
HASH := $(shell test -d .git && git rev-parse --short HEAD || echo "UNKNOWN")
DIRTY := $(shell git diff --exit-code >/dev/null || echo "-dirty" && echo "")
BUILD_DATE := $(shell date +%FT%T%z)

# CI Helpers
#DOCKER_COMPOSE = docker-compose -f docker/compose-ci.yml
DOCKER_COMPOSE_DEV = docker-compose

LDFLAGS = -ldflags "-X ${SGG}/meta.Ref=${HASH}${DIRTY} -X ${SGG}/meta.BuildDate=${BUILD_DATE}"

.PHONY: test install clean gen dev-up migrate dev-setup deps dev-reset get-glide

default: build
build: gen buildall
buildall: sgg-api sgg-tools sgg-worker

#
# BINARIES
#

sgg-api: ${SOURCES}
	${GO} build ${LDFLAGS} -v ./cmd/sgg-api/sgg-api.go

sgg-tools: ${SOURCES}
	${GO} build ${LDFLAGS} -v ./cmd/sgg-tools/sgg-tools.go

sgg-worker: ${SOURCES}
	${GO} build ${LDFLAGS} -v ./cmd/sgg-worker/sgg-worker.go

#
# DEV
#

dev: dev-up build

dev-up:
	${DOCKER_COMPOSE_DEV} up -d

dev-setup: deps sgg-tools migrate
dev-setup:
	./sgg-tools debug-user register -a -u devadmin -p 123456789 -e test@svgg.xyz
	./sgg-tools debug-user register -u devuser -p 123456789 -e test2@svgg.xyz

migrate: sgg-tools
migrate:
	./sgg-tools migrate
	./sgg-tools migrate influx

deps:
	glide install

#
# HELPERS
#

gen:
	${GO} list ./... | grep -v /vendor/ | xargs -n10 ${GO} generate -v 

test:
	${GO} list ./... | grep -v /vendor/ | xargs -n10 ${GO} test ${LDFLAGS} -v 

install:
	${GO} install ${LDFLAGS} -v cmd/...

get-glide:
	curl -sSL https://glide.sh/get | sh

clean:
	rm sgg-api sgg-worker sgg-tools
