

# Nori Makefile

NORI_BUILD_CMD ?= build -o build/bin/nori ./cmd/nori.go

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

build: protoc-generate ## build nori binary
	@go $(NORI_BUILD_CMD)
.PHONY: build

build-web: protoc-generate ## build nori binary with packr
	@packr $(NORI_BUILD_CMD)
	@packr clean
.PHONY: build-web

docker-image: protoc-generate ## build noriio/nori docker image
	docker build -f build/docker/0.2.0/Dockerfile -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
.PHONY: docker-image

docker-push: ## push docker image to registry
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
.PHONY: docker-push

git-update: ## update git submodules
	git submodule update --remote
.PHONY: git-sub-pull

protoc-generate: ## generate protobuf
	@protoc --proto_path=api/protobuf/plugin --go_out=plugins=grpc:./internal/generated/protobuf api/protobuf/plugin/*.proto
.PHONY: protoc-generate

lint: ## execute linter
	golangci-lint run --no-config --issues-exit-code=0 --deadline=30m \
	  --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
	  --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
	  --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck ./... ;
.PHONY: lint

run: ## run 'go run cmd/nori.go server'
	go run cmd/main.go server
.PHONY: run

test: ## run go tests
	@go test -v -race ./...
.PHONY: test

test-with-coverage: ## run go test with coverage
	@go test -covermode atomic -coverprofile profile.out ${TEST_ARGS} ./... ;\
	go tool cover -func=profile.out
.PHONY: test-with-coverage

help: ## run this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

.DEFAULT_GOAL := help