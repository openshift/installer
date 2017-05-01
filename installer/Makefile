export CGO_ENABLED:=0
export VERSION=$(shell ../git-version)

OS=$(shell uname | tr A-Z a-z)
LD_FLAGS="-w -X github.com/coreos/tectonic-installer/installer/server/version.Version=$(VERSION)"

REPO=github.com/coreos/tectonic-installer/installer
IMAGE_REPO=coreos-inc/installer
QUAY_REPO=quay.io/coreos-inc/installer

.PHONY: all
all: build

.PHONY: build
build: bin/$(OS)/installer

bin/windows/installer.exe: frontend binassets/assets.go
	@GOOS=windows go build -o bin/windows/installer.exe -ldflags $(LD_FLAGS) $(REPO)/cmd/installer

bin/%/installer: frontend binassets/assets.go
	@GOOS=$* go build -o bin/$*/installer -ldflags $(LD_FLAGS) $(REPO)/cmd/installer

.PHONY: backend
backend: binassets/assets.go
	@GOOS=$(OS) go build -o bin/$(OS)/installer -ldflags $(LD_FLAGS) $(REPO)/cmd/installer

binassets/assets.go: bin/go-bindata
	@./bin/go-bindata -pkg binassets -o binassets/assets.go -prefix assets assets/...

.PHONY: frontend
frontend:
	$(MAKE) -C frontend

.PHONY: test
test: build bin/golint
	@./scripts/test

.PHONY: smoke-aws
smoke-aws: bin/sanity
	@../tests/scripts/aws/up-down.sh

.PHONY: smoke-bare-metal
smoke-bare-metal: bin/sanity
	@./tests/scripts/bare-metal/up-down.sh

.PHONY: docker-image
docker-image: bin/linux/installer
	@sudo docker build --rm=true -t $(IMAGE_REPO):$(VERSION) .
	@sudo docker tag $(IMAGE_REPO):$(VERSION) $(IMAGE_REPO):latest

.PHONY: release
release: release-bins
	@./scripts/release/upload_installer_bins.sh
	@./scripts/release/make_release_tarball.sh
	@./scripts/release/upload_release_tarball.sh
	@./scripts/release/make_github_release.sh

.PHONY: release-bins
release-bins: \
	bin/linux/installer \
	bin/darwin/installer \
	bin/windows/installer.exe

.PHONY: vendor
vendor: glide.yaml
	@glide up -v
	@glide-vc --use-lock-file --no-tests --only-code

.PHONY: tools
tools: bin/go-bindata bin/sanity bin/golint

bin/golint:
	CGO_ENABLED=0 go build -o bin/golint $(REPO)/vendor/github.com/golang/lint/golint

bin/go-bindata:
	@go build -o bin/go-bindata $(REPO)/vendor/github.com/jteeuwen/go-bindata/go-bindata

bin/sanity:
	@go test tests/sanity/k8s_test.go -c -o bin/sanity

.PHONY: clean
clean:
	@rm -rf bin/
	@rm -f binassets/assets.go
	@rm -rf assets/public/scripts
	@rm -rf frontend/node_modules
