.PHONY: clean generate test build

BIN_OUTPUT := $(if $(filter $(shell go env GOOS), windows), dist/gci.exe, dist/gci)

default: clean generate test build

clean:
	@echo BIN_OUTPUT: ${BIN_OUTPUT}
	@rm -rf dist cover.out

build: clean
	@go build -v -trimpath -o ${BIN_OUTPUT} .

test: clean
	@go test -v -count=1 -cover ./...

generate:
	@GOEXPERIMENT=arenas,boringcrypto go generate ./...
