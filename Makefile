.PHONY: all test build
all:  build

test:
	@go test -v ./...

build:
	@cd ./proto; protoc --go_out=plugins=grpc:. *.proto
	@mkdir -p ./bin
	@go build -o bin/nori ./cmd/main.go

clean:
	rm -rf ./bin