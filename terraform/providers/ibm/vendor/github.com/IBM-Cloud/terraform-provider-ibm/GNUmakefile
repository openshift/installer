TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find .  -path ./.direnv -prune -false -o -name '*.go' |grep -v vendor)
COVER_TEST?=$$(go list ./... |grep -v 'vendor')
TEST_TIMEOUT?=700m
VERSION ?= 0.0.1
OS_ARCH := $(shell go env GOOS)_$(shell go env GOARCH)
PLUGIN_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/ibm-cloud/ibm/$(VERSION)/$(OS_ARCH)

default: build

tools:
	@go get golang.org/x/sys/unix
	@go get github.com/kardianos/govendor
	@go get github.com/mitchellh/gox
	@go get golang.org/x/tools/cmd/cover

build: fmtcheck vet
	go install

build-local: build
	mkdir -p $(PLUGIN_DIR)
	mv $(HOME)/go/bin/terraform-provider-ibm $(PLUGIN_DIR)/terraform-provider-ibm_v$(VERSION)
	./scripts/post-build.sh $(VERSION)

bin: fmtcheck vet tools
	@TF_RELEASE=1 sh -c "'$(CURDIR)/scripts/build.sh'"

dev: fmtcheck vet tools
	@TF_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout $(TEST_TIMEOUT)

testrace: fmtcheck
	TF_ACC= go test -race $(TEST) $(TESTARGS)

cover:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(COVER_TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile: fmtcheck
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./builtin/providers/aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build build-local bin dev test testacc testrace cover vet fmt fmtcheck errcheck vendor-status test-compile
