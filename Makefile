# Nori Makefile

NORI_BUILD_CMD ?= build -ldflags "-X github.com/nori-io/nori/pkg/version.GOOS=`go env GOOS` -X 'github.com/nori-io/nori/pkg/version.GOARCH=`go env GOARCH`' -X 'github.com/nori-io/nori/pkg/version.GOVERSION=`go env GOVERSION | cut -c 3-`'" -o build/bin/nori ./cmd/main.go

ifeq ($(GO111MODULE),auto)
override GO111MODULE = on
endif

DOCKER_IMAGE ?= noriio/nori
DOCKER_TAG ?= dev-go1.16.4

clean: ## remove generated files, tidy vendor dependencies
	export GO111MODULE=on ;\
	go mod tidy ;\
	rm -rf profile.out ;\
	@rm -rf ./build/bin/*
.PHONY: clean

build-test-plugins: ## build plugins to run tests with plugins
	@echo "implement command to build test/testdata/plugins/* plugins"
.PHONY: build-test-plugins

build: ## build nori binary
	@go $(NORI_BUILD_CMD)
.PHONY: build

docker-image:
	docker build -f build/docker/0.3.0/Dockerfile -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
.PHONY: docker-image

docker-push: ## push docker image to registry
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
.PHONY: docker-push

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