default: test lint

fmt:
	gofmt -s -w ./

lint: golangci-lint importlint

golangci-lint:
	@golangci-lint run ./...
	@cd v2/awsv1shim && golangci-lint run ./...

importlint:
	@impi --local . --scheme stdThirdPartyLocal ./...

test:
	go test -timeout=30s -parallel=4 ./...
	cd v2/awsv1shim && go test -timeout=30s -parallel=4 ./...

tools:
	cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	cd tools && go install github.com/pavius/impi/cmd/impi

semgrep:
	@docker run --rm --volume "${PWD}:/src" returntocorp/semgrep --config .semgrep --no-rewrite-rule-ids

.PHONY: lint test tools
