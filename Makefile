# Nori Makefile

NORI_BUILD_CMD ?= build -o nori ./cmd/nori.go

ifeq ($(GO111MODULE),auto)
override GO111MODULE = on
endif

DOCKER_IMAGE ?= noriio/nori
DOCKER_TAG ?= 0.2.0

clean: ## remove generated files, tidy vendor dependencies
	export GO111MODULE=on ;\
	go mod tidy ;\
	rm -rf profile.out ;\
	@rm -rf ./bin
	@packr clean
.PHONY: clean

build-test-plugins: ## build plugins to run tests with plugins
	@echo "implement command to build test/testdata/plugins/* plugins"
.PHONY: build-test-plugins

test: ## run tests
	@go test -v ./...
.PHONY: test

build: protoc-generate ## build nori binary
	@go $(NORI_BUILD_CMD)
.PHONY: build

build-web: protoc-generate ## build nori binary with packr
	@packr $(NORI_BUILD_CMD)
	@packr clean
.PHONY: build-web

docker-image: protoc-generate ## build noriio/nori docker image
	docker build -f build/docker/0.2.0/Dockerfile -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
.PHONY: docker-build-image

protoc-generate: ## generate protobuf
	@protoc --proto_path=api/protobuf/nori --go_out=plugins=grpc:./internal/generated/protobuf api/protobuf/nori/*.proto
.PHONY: protoc-generate

lint: ## execute linter
	golangci-lint run --no-config --issues-exit-code=0 --deadline=30m \
	  --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
	  --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
	  --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck ./... ;
.PHONY: lint

help: ## run this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

.DEFAULT_GOAL := help