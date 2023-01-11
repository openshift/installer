LDFLAGS += -X main.version=$$(git describe --always --abbrev=40 --dirty)
TEST?=$$(go list ./... |grep -v 'vendor')
PKG_NAME=ironic
TERRAFORM_PLUGINS=$(HOME)/.terraform.d/plugins

ifeq ("$(IRONIC_ENDPOINT)", "")
	IRONIC_ENDPOINT := http://127.0.0.1:6385/
	export IRONIC_ENDPOINT
endif

default: build

build:
	go build -ldflags "${LDFLAGS}" -tags "${TAGS}"

install: default
	mkdir -p ${TERRAFORM_PLUGINS}
	mv terraform-provider-ironic ${TERRAFORM_PLUGINS}

fmt:
	gofmt -s -d -e ./ironic

lint :$(GOPATH)/golangci-lint
	GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache/ $(GOPATH)/golangci-lint run ironic

$(GOPATH)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH) v1.50.0

test:
	go test -tags "${TAGS}" -v ./ironic

acceptance:
	TF_ACC=true go test -tags "acceptance" -v ./ironic/...

clean:
	rm -f terraform-provider-ironic

.PHONY: build install test fmt lint
