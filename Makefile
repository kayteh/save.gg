GO ?= go

SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

SGG = save.gg/sgg
HASH := $(shell test -d .git && git rev-parse --short HEAD || cat revision~)
BUILD_DATE := $(shell date +%FT%T%z)

# CI Helpers
#DOCKER_COMPOSE = docker-compose -f docker/compose-ci.yml
DOCKER_COMPOSE_DEV ?= docker-compose

LDFLAGS = -ldflags "-X ${SGG}/meta.Ref=${HASH} -X ${SGG}/meta.BuildDate=${BUILD_DATE}"

.PHONY: test install clean gen dev-up migrate dev-setup deps dev-reset

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
	./sgg-tools debug-user register

migrate: sgg-tools
migrate:
	./sgg-tools migrate
	./sgg-tools migrate rethink
	./sgg-tools migrate influx

deps:
	govendor fetch

#
# HELPERS
#

gen:
	${GO} generate $(shell ${GO} list ./... | grep -v /vendor/) 

test:
	${GO} test ${LDFLAGS} -v $(shell ${GO} list ./... | grep -v /vendor/)

install:
	${GO} install ${LDFLAGS} -v cmd/...

get-govendor: ${GOPATH}/bin/govendor
	${GO} get -u github.com/kardianos/govendor

clean:
	rm sgg-api sgg-worker sgg-tools
