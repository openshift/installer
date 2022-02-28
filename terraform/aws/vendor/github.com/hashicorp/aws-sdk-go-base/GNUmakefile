default: test lint

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./

lint: golangci-lint importlint

golangci-lint:
	@golangci-lint run ./...

importlint:
	@impi --local . --scheme stdThirdPartyLocal ./...

test:
	go test -timeout=30s -parallel=4 ./...

tools:
	cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	cd tools && go install github.com/pavius/impi/cmd/impi

.PHONY: lint test tools
