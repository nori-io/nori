# Nori Makefile

NORI_BUILD_CMD ?= build -o nori ./cmd/nori.go

clean: ## remove generated files, tidy vendor dependencies
	export GO111MODULE=on ;\
	go mod tidy ;\
	rm -rf profile.out ;\
	@rm -rf ./bin
	@packr clean
.PHONY: clean

test:
	@go test -v ./...
.PHONY: test

build: protoc-generate
	@go $(NORI_BUILD_CMD)
.PHONY: build

build-web: protoc-generate
	@packr $(NORI_BUILD_CMD)
	@packr clean
.PHONY: build-web

protoc-generate:
	@protoc --proto_path=api/protobuf-spec/ --go_out=plugins=grpc:./internal/generated/protobuf api/protobuf-spec/*.proto
.PHONY: protoc-generate

lint: ## execute linter
	golangci-lint run --no-config --issues-exit-code=0 --deadline=30m \
	  --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
	  --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
	  --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck ./... ;
.PHONY: lint