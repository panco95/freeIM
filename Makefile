APP_VERSION=$(shell date '+%Y.%m.%d.%H%M%S')
VERSION=$(shell git describe --tags --always)
BUILD_TIME=$(shell date '+%Y.%m.%d.%H%M%S')
LDFLAGS="-X 'im/cmd.version=${VERSION}' -X 'im/cmd.date=${BUILD_TIME}'"

.PHONY: build
build:
	go mod tidy
	CGO_ENABLED=0 go build -o dist/im main.go

.PHONY: build-ci
build-ci:
	go build -ldflags ${LDFLAGS} -o bin/im main.go

.PHONY: generate
generate:
	go generate ./...

go-pkg-list:
	export PKG_LIST="$(shell go list ./... | grep -v /vendor/)"

.PHONY: lint
lint:
	@go lint $(shell go list ./... | grep -v /vendor/)

.PHONY: clean
clean:
	@rm -f build/*

.PHONY: deps
deps:
	@go mod tidy

.PHONY: coverage-report
coverage-report:
	@go test -cover -v $(shell go list ./... | grep -v /vendor/) -coverprofile=build/coverage.out
	@go tool cover -html=build/coverage.out

.PHONY: mocks
mocks:
	go generate $(shell go list ./... | grep -v /vendor/)

protoc:
	protoc -I=models/proto --go_out=. message.proto internal.proto