TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=alicloud

default: build

build: fmtcheck	all

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

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
	goimports -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), getting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	ln -sf ../../../../ext/providers/alicloud/website/docs $(GOPATH)/src/github.com/hashicorp/terraform-website/content/source/docs/providers/alicloud
	ln -sf ../../../ext/providers/alicloud/website/alicloud.erb $(GOPATH)/src/github.com/hashicorp/terraform-website/content/source/layouts/alicloud.erb
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test

all: mac windows linux

dev: clean fmt mac copy

devlinux: clean fmt linux linuxcopy

devwin: clean fmt windows windowscopy

copy:
	tar -xvf bin/terraform-provider-alicloud_darwin-amd64.tgz && mv bin/terraform-provider-alicloud $(shell dirname `which terraform`)

clean:
	rm -rf bin/*

mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/terraform-provider-alicloud
	tar czvf bin/terraform-provider-alicloud_darwin-amd64.tgz bin/terraform-provider-alicloud
	rm -rf bin/terraform-provider-alicloud

windowscopy:
	tar -xvf bin/terraform-provider-alicloud_windows-amd64.tgz && mv bin/terraform-provider-alicloud $(shell dirname `which terraform`)
    
windows:
	GOOS=windows GOARCH=amd64 go build -o bin/terraform-provider-alicloud.exe
	tar czvf bin/terraform-provider-alicloud_windows-amd64.tgz bin/terraform-provider-alicloud.exe
	rm -rf bin/terraform-provider-alicloud.exe

linuxcopy:
	tar -xvf bin/terraform-provider-alicloud_linux-amd64.tgz && mv bin/terraform-provider-alicloud $(shell dirname `which terraform`)

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/terraform-provider-alicloud
	tar czvf bin/terraform-provider-alicloud_linux-amd64.tgz bin/terraform-provider-alicloud
	rm -rf bin/terraform-provider-alicloud
