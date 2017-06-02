export CGO_ENABLED=0
export VERSION=$(shell ../git-version)

OS=$(shell uname | tr A-Z a-z)
LD_FLAGS="-w -X github.com/coreos/tectonic-installer/installer/cmd/installer.Version=$(VERSION)"

REPO=github.com/coreos/tectonic-installer/installer

# Wrapper around find. Works with filenamess containing spaces.
find=$(shell [ -d $(1) ] && find $(1) -type f $(2) | sed 's/ /\\ /g')

GO_FILES=$(call find,.,-name '*.go')
SHELL_FILES=$(call find,scripts,-name '*.sh')

.PHONY: all
all:
	@$(MAKE) build -r

.PHONY: build
build: bin/$(OS)/installer

bin/windows/installer.exe: $(GO_FILES) assets/bindata.go
	GOOS=windows go build -o bin/windows/installer.exe -ldflags $(LD_FLAGS) $(REPO)/cmd/installer

bin/%/installer: $(GO_FILES) assets/bindata.go
	GOOS=$* go build -o bin/$*/installer -ldflags $(LD_FLAGS) $(REPO)/cmd/installer

.PHONY: backend
backend: assets/bindata.go
	GOOS=$(OS) go build -o bin/$(OS)/installer -ldflags $(LD_FLAGS) $(REPO)/cmd/installer

.PHONY: frontend
frontend: assets/frontend/scripts/app-bundle.js

assets/bindata.go: assets/frontend/scripts/app-bundle.js bin/go-bindata $(shell find assets -type f | grep -v .go)
	./bin/go-bindata -pkg assets -o assets/bindata.go -ignore=bindata.go -ignore=doc.go -ignore=assets.go -prefix assets assets/...

assets/frontend/scripts/app-bundle.js: frontend/node_modules frontend/*.jsx frontend/*.js $(call find,frontend/components,-name '*.jsx') $(call find,frontend/node_modules)
	cd frontend && yarn run build

frontend/node_modules: frontend/package.json frontend/yarn.lock
	cd frontend && yarn install --quiet > /dev/null

.PHONY: lint
lint: build bin/golint
	./scripts/lint-backend
	cd frontend && yarn run lint

.PHONY: test-backend
test-backend: build
	./scripts/test-backend

.PHONY: test-frontend
test-frontend:
	cd frontend && yarn test

.PHONY: test
test: test-backend test-frontend

.PHONY: shellcheck
shellcheck: $(SHELL_FILES)
	shellcheck $(SHELL_FILES)

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
	go build -o bin/golint $(REPO)/vendor/github.com/golang/lint/golint

bin/go-bindata:
	go build -o bin/go-bindata $(REPO)/vendor/github.com/jteeuwen/go-bindata/go-bindata

bin/sanity:
	go test tests/sanity/k8s_test.go -c -o bin/sanity

.PHONY: clean
clean:
	rm -f assets/bindata.go
	rm -f assets/frontend/scripts/app-bundle.js
	rm -fr .workspace/
	rm -fr assets/frontend/scripts
	rm -fr bin/
	rm -fr frontend/node_modules
