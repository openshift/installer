export CGO_ENABLED:=0
export VERSION=$(shell ../git-version)

OS=$(shell uname | tr A-Z a-z)
LD_FLAGS="-w -X github.com/coreos/tectonic-installer/installer/cmd/installer.Version=$(VERSION)"

REPO=github.com/coreos/tectonic-installer/installer

GO_FILES=$(shell find . -type f -name '*.go')

.PHONY: all
all: build

.PHONY: build
build: bin/$(OS)/installer

bin/windows/installer.exe: $(GO_FILES) assets/bindata.go
	GOOS=windows go build -o bin/windows/installer.exe -ldflags $(LD_FLAGS) $(REPO)/cmd/installer

bin/%/installer: $(GO_FILES) assets/bindata.go
	GOOS=$* go build -o bin/$*/installer -ldflags $(LD_FLAGS) $(REPO)/cmd/installer

.PHONY: backend
backend: assets/bindata.go
	GOOS=$(OS) go build -o bin/$(OS)/installer -ldflags $(LD_FLAGS) $(REPO)/cmd/installer

assets/bindata.go: frontend bin/go-bindata $(shell find assets -type f | grep -v .go)
	./bin/go-bindata -pkg assets -o assets/bindata.go -ignore=bindata.go -ignore=doc.go -ignore=assets.go -prefix assets assets/...

.PHONY: frontend
frontend:
	$(MAKE) -C frontend

.PHONY: test
test: build bin/golint
	./scripts/test

.PHONY: smoke-aws
smoke-aws: bin/sanity
	../tests/scripts/aws/up-down.sh

.PHONY: smoke-bare-metal
smoke-bare-metal: bin/sanity
	./tests/scripts/bare-metal/up-down.sh

.PHONY: dist
dist: release-bins
	./scripts/release/make_release_tarball.sh

.PHONY: release
release: release-bins
	./scripts/release/upload_installer_bins.sh
	./scripts/release/make_release_tarball.sh
	./scripts/release/upload_release_tarball.sh
	./scripts/release/make_github_release.sh

.PHONY: license-parse
license-parse:
	./scripts/license-gen/collect_project_licenses.sh

.PHONY: release-bins
release-bins: \
	bin/linux/installer \
	bin/darwin/installer
	# bin/windows/installer.exe TODO re-add me when windows works again

.PHONY: vendor
vendor: glide.yaml
	glide up -v
	glide-vc --use-lock-file --no-tests --only-code

.PHONY: tools
tools: bin/go-bindata bin/sanity bin/golint

bin/golint:
	CGO_ENABLED=0 go build -o bin/golint $(REPO)/vendor/github.com/golang/lint/golint

bin/go-bindata:
	go build -o bin/go-bindata $(REPO)/vendor/github.com/jteeuwen/go-bindata/go-bindata

bin/sanity:
	go test tests/sanity/k8s_test.go -c -o bin/sanity

.PHONY: clean
clean:
	$(MAKE) -C frontend clean
	rm -rf bin/
	rm -f assets/bindata.go
	rm -rf assets/frontend/scripts
	rm -rf .workspace/
