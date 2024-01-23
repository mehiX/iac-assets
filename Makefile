OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)
EXE=$(shell go env GOEXE)
BIN=dist
SRC=$(shell find . -name "*.go")

.PHONY: all fmt build serve install_deps clean test

default: all

all: fmt test build

fmt:
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

test: install_deps
	$(info ******************** running tests ********************)
	go test ./...

install_deps:
	$(info ******************** install dependencies ********************)
	go get -v ./...

build: install_deps
	@mkdir -p dist
	go build -o $(BIN)/iac_$(OS)_$(ARCH)$(EXE) ./main.go

serve: build
	$(BIN)/iac_$(OS)_$(ARCH)$(EXE) serve

clean:
	rm -rf $(BIN)