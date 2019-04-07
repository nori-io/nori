.PHONY: all test build
all: build

test:
	@go test -v ./...

build:
	@cd ./proto; protoc --go_out=plugins=grpc:. *.proto
	@mkdir -p ./bin
	@go build -o bin/nori ./cmd/main.go

build-web:
	@cd ./proto; protoc --go_out=plugins=grpc:. *.proto
	@mkdir -p ./bin
	@packr build -o bin/nori ./cmd/main.go
	@packr clean

clean:
	@rm -rf ./bin
	@packr clean