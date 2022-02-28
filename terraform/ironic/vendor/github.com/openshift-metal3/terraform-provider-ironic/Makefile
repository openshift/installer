LDFLAGS += -X main.version=$$(git describe --always --abbrev=40 --dirty)
TEST?=$$(go list ./... |grep -v 'vendor')
PKG_NAME=ironic
TERRAFORM_PLUGINS=$(HOME)/.terraform.d/plugins

ifeq ("$(IRONIC_ENDPOINT)", "")
	IRONIC_ENDPOINT := http://127.0.0.1:6385/
	export IRONIC_ENDPOINT
endif

default: fmt lint build

build:
	go build -ldflags "${LDFLAGS}" -tags "${TAGS}"

install: default
	mkdir -p ${TERRAFORM_PLUGINS}
	mv terraform-provider-ironic ${TERRAFORM_PLUGINS}

fmt:
	go fmt ./ironic .

tools:
	go get golang.org/x/lint/golint

lint: tools
	go run golang.org/x/lint/golint -set_exit_status ./ironic .

test:
	go test -tags "${TAGS}" -v ./ironic

acceptance:
	TF_ACC=true go test -tags "acceptance" -v ./ironic/...

clean:
	rm -f terraform-provider-ironic

.PHONY: build install test fmt lint
