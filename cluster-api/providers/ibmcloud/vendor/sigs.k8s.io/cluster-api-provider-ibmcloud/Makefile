# Copyright 2021 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.DEFAULT_GOAL:=help

ROOT_DIR_RELATIVE := .

include $(ROOT_DIR_RELATIVE)/common.mk

# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:crdVersions=v1"

SHELL=/bin/bash

# Directories.
REPO_ROOT := $(shell git rev-parse --show-toplevel)
ARTIFACTS ?= $(REPO_ROOT)/_artifacts
BIN_DIR := bin
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
GO_INSTALL = ./scripts/go_install.sh
E2E_CONF_FILE_ENVSUBST := $(REPO_ROOT)/test/e2e/config/ibmcloud-e2e-envsubst.yaml
E2E_TEMPLATES := $(REPO_ROOT)/test/e2e/data/templates
TEMPLATES_DIR := $(REPO_ROOT)/templates

GO_APIDIFF := $(TOOLS_BIN_DIR)/go-apidiff
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint
KUSTOMIZE := $(TOOLS_BIN_DIR)/kustomize
GOJQ := $(TOOLS_BIN_DIR)/gojq
CONVERSION_GEN := $(TOOLS_BIN_DIR)/conversion-gen
GINKGO := $(TOOLS_BIN_DIR)/ginkgo
GOTESTSUM := $(TOOLS_BIN_DIR)/gotestsum
ENVSUBST := $(TOOLS_BIN_DIR)/envsubst
MOCKGEN := $(TOOLS_BIN_DIR)/mockgen
CONTROLLER_GEN := $(TOOLS_BIN_DIR)/controller-gen
CONVERSION_VERIFIER := $(TOOLS_BIN_DIR)/conversion-verifier
SETUP_ENVTEST := $(TOOLS_BIN_DIR)/setup-envtest
GOVULNCHECK := $(TOOLS_BIN_DIR)/govulncheck
TRIVY := $(TOOLS_BIN_DIR)/trivy

STAGING_REGISTRY ?= gcr.io/k8s-staging-capi-ibmcloud
STAGING_BUCKET ?= artifacts.k8s-staging-capi-ibmcloud.appspot.com
BUCKET ?= $(STAGING_BUCKET)
PROD_REGISTRY := registry.k8s.io/capi-ibmcloud
REGISTRY ?= $(STAGING_REGISTRY)
RELEASE_TAG ?= $(shell git describe --abbrev=0 2>/dev/null)
PULL_BASE_REF ?= $(RELEASE_TAG) # PULL_BASE_REF will be provided by Prow
RELEASE_ALIAS_TAG ?= $(PULL_BASE_REF)
RELEASE_DIR := out
OUTPUT_TYPE ?= type=registry

# Go
GO_VERSION ?=1.22.12
GO_CONTAINER_IMAGE ?= golang:$(GO_VERSION)

# kind
CAPI_KIND_CLUSTER_NAME ?= capi-test

# image name used to build the cmd/capibmadm
TOOLCHAIN_IMAGE := toolchain

TAG ?= dev
ARCH ?= $(shell go env GOARCH)
ALL_ARCH ?= amd64 ppc64le arm64
BUILDX_PLATFORMS ?= linux/amd64,linux/arm64,linux/ppc64le
PULL_POLICY ?= Always

# Set build time variables including version details
LDFLAGS := $(shell ./hack/version.sh)

KUBEBUILDER_ENVTEST_KUBERNETES_VERSION ?= 1.30.0

# main controller
CORE_IMAGE_NAME ?= cluster-api-ibmcloud-controller
CORE_CONTROLLER_IMG ?= $(REGISTRY)/$(CORE_IMAGE_NAME)
CORE_CONTROLLER_ORIGINAL_IMG := gcr.io/k8s-staging-capi-ibmcloud/cluster-api-ibmcloud-controller
CORE_CONTROLLER_NAME := controller-manager
CORE_MANIFEST_FILE := infrastructure-components
CORE_CONFIG_DIR := config/default
CORE_NAMESPACE := capi-ibmcloud-system

PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

## --------------------------------------
## Binaries
## --------------------------------------

##@ build:

.PHONY: binaries
binaries: manager capibmadm ## Builds and installs all binaries

.PHONY: capibmadm
capibmadm: ## Build the capibmadm binary into the ./bin folder
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/capibmadm ./cmd/capibmadm

# Build manager binary
manager: generate fmt vet ## Build the manager binary into the ./bin folder
	go build -ldflags "${LDFLAGS} -extldflags '-static'" -o $(BIN_DIR)/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./main.go

# Install CRDs into a cluster
install: generate-manifests $(KUSTOMIZE)
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: generate-manifests $(KUSTOMIZE)
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: generate-manifests $(KUSTOMIZE)
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

help:  # Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-45s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## --------------------------------------
## Generate / Manifests
## --------------------------------------

##@ generate:

# Generate code
.PHONY: generate
generate: ## Run all generate-go generate-modules generate-manifests generate-go-deepcopy generate-go-conversions generate-templates
	$(MAKE) generate-go generate-modules generate-manifests generate-go-deepcopy generate-go-conversions generate-templates

generate-go-deepcopy: $(CONTROLLER_GEN) ## Generate deepcopy go code
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: generate-go
generate-go: $(MOCKGEN) ## Generate the Go mock code
	go generate ./...

.PHONY: generate-manifests
generate-manifests: $(CONTROLLER_GEN) ## Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate-go-conversions
generate-go-conversions: $(CONVERSION_GEN) ## Generate conversions go code
	$(MAKE) clean-generated-conversions SRC_DIRS="./api/v1beta1"
	$(CONVERSION_GEN) \
		--output-file=zz_generated.conversion.go \
		--go-header-file=./hack/boilerplate/boilerplate.generatego.txt \
		./api/v1beta1

.PHONY: generate-templates
generate-templates: $(KUSTOMIZE) ## Generate cluster templates
	$(KUSTOMIZE) build $(TEMPLATES_DIR)/cluster-template --load-restrictor LoadRestrictionsNone > $(TEMPLATES_DIR)/cluster-template.yaml
	$(KUSTOMIZE) build $(TEMPLATES_DIR)/cluster-template-powervs --load-restrictor LoadRestrictionsNone > $(TEMPLATES_DIR)/cluster-template-powervs.yaml
	$(KUSTOMIZE) build $(TEMPLATES_DIR)/cluster-template-powervs-clusterclass --load-restrictor LoadRestrictionsNone > $(TEMPLATES_DIR)/cluster-template-powervs-clusterclass.yaml
	$(KUSTOMIZE) build $(TEMPLATES_DIR)/cluster-template-vpc-clusterclass --load-restrictor LoadRestrictionsNone > $(TEMPLATES_DIR)/cluster-template-vpc-clusterclass.yaml
	
.PHONY: generate-e2e-templates
generate-e2e-templates: $(KUSTOMIZE) ## Generate E2E cluster templates
ifeq ($(E2E_FLAVOR), powervs-md-remediation)
	$(KUSTOMIZE) build $(E2E_TEMPLATES)/cluster-template-powervs-md-remediation --load-restrictor LoadRestrictionsNone > $(E2E_TEMPLATES)/cluster-template-powervs-md-remediation.yaml
else
	$(KUSTOMIZE) build $(E2E_TEMPLATES)/cluster-template-vpc --load-restrictor LoadRestrictionsNone > $(E2E_TEMPLATES)/cluster-template-vpc.yaml
endif

.PHONY: generate-modules
generate-modules: ## Runs go mod to ensure modules are up to date
	go mod tidy
	cd $(TOOLS_DIR); go mod tidy

images: docker-build

set-flavor: 
ifeq ($(findstring vpc,$(E2E_FLAVOR)),vpc)
	 $(eval E2E_CONF_FILE=$(REPO_ROOT)/test/e2e/config/ibmcloud-e2e-vpc.yaml)
else
	 $(eval E2E_CONF_FILE=$(REPO_ROOT)/test/e2e/config/ibmcloud-e2e-powervs.yaml)
endif
	@echo "Setting e2e test flavour to ${E2E_CONF_FILE}"

## --------------------------------------
## Testing
## --------------------------------------

##@ test:

.PHONY: setup-envtest
setup-envtest: $(SETUP_ENVTEST) # Build setup-envtest from tools folder
	@if [ $(shell go env GOOS) == "darwin" ]; then \
		$(eval KUBEBUILDER_ASSETS := $(shell $(SETUP_ENVTEST) use --use-env -p path --arch amd64 $(KUBEBUILDER_ENVTEST_KUBERNETES_VERSION))) \
		echo "kube-builder assets set using darwin OS at location $(KUBEBUILDER_ASSETS)"; \
	else \
		$(eval KUBEBUILDER_ASSETS := $(shell $(SETUP_ENVTEST) use --use-env -p path $(KUBEBUILDER_ENVTEST_KUBERNETES_VERSION))) \
		echo "kube-builder assets set using other OS at location $(KUBEBUILDER_ASSETS)"; \
	fi

# Run unit tests
test: generate fmt vet setup-envtest $(GOTESTSUM) ## Run tests
	KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" $(GOTESTSUM) --junitfile $(ARTIFACTS)/junit.xml

# Allow overriding the e2e configurations
GINKGO_FOCUS ?= Workload cluster creation
GINKGO_NODES ?= 3
GINKGO_NOCOLOR ?= false
GINKGO_TIMEOUT ?= 2h
E2E_FLAVOR ?= powervs-md-remediation
JUNIT_FILE ?= junit.e2e_suite.1.xml
GINKGO_ARGS ?= -v --trace --tags=e2e --timeout=$(GINKGO_TIMEOUT) --focus=$(GINKGO_FOCUS) --nodes=$(GINKGO_NODES) --no-color=$(GINKGO_NOCOLOR) --output-dir="$(ARTIFACTS)" --junit-report="$(JUNIT_FILE)"
ARTIFACTS ?= $(REPO_ROOT)/_artifacts
SKIP_CLEANUP ?= false
SKIP_CREATE_MGMT_CLUSTER ?= false

# Run the end-to-end tests
.PHONY: test-e2e
test-e2e: $(GINKGO) $(KUSTOMIZE) $(ENVSUBST) set-flavor e2e-image generate-e2e-templates ## Run e2e tests
	$(ENVSUBST) < $(E2E_CONF_FILE) > $(E2E_CONF_FILE_ENVSUBST) 
	$(GINKGO) $(GINKGO_ARGS) ./test/e2e -- \
		-e2e.artifacts-folder="$(ARTIFACTS)" \
		-e2e.config="$(E2E_CONF_FILE_ENVSUBST)" \
		-e2e.skip-resource-cleanup=$(SKIP_CLEANUP) \
		-e2e.use-existing-cluster=$(SKIP_CREATE_MGMT_CLUSTER) \
		-e2e.flavor="$(E2E_FLAVOR)"

# Basic checks for deploying kind cluster and required providers
.PHONY: test-sanity
test-sanity: ## Run sanity tests
	GINKGO_FOCUS="Run Sanity tests" $(MAKE) test-e2e

.PHONY: test-cover
test-cover: setup-envtest## Run tests with code coverage and code generate reports
	KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" go test ./... -coverprofile cover.out
	go tool cover -func=cover.out -o cover.txt
	go tool cover -html=cover.out -o cover.html

## --------------------------------------
## Release
## --------------------------------------

##@ release:

$(RELEASE_DIR):
	mkdir -p $@

$(ARTIFACTS):
	mkdir -p $@

.PHONY: build-toolchain
build-toolchain: ## Build the toolchain
	docker build --target toolchain --build-arg golang_image=$(GO_CONTAINER_IMAGE) -t $(TOOLCHAIN_IMAGE) .

.PHONY: list-staging-releases
list-staging-releases: ## List staging images for image promotion
	@echo $(CORE_IMAGE_NAME):
	$(MAKE) list-image RELEASE_TAG=$(RELEASE_TAG) IMAGE=$(CORE_IMAGE_NAME)

list-image:
	gcloud container images list-tags $(STAGING_REGISTRY)/$(IMAGE) --filter="tags=('$(RELEASE_TAG)')" --format=json

.PHONY: check-release-tag
check-release-tag:
	@if [ -z "${RELEASE_TAG}" ]; then echo "RELEASE_TAG is not set"; exit 1; fi
	@if ! [ -z "$$(git status --porcelain)" ]; then echo "Your local git repository contains uncommitted changes, use git clean before proceeding."; exit 1; fi

.PHONY: check-previous-release-tag
check-previous-release-tag:
	@if [ -z "${PREVIOUS_VERSION}" ]; then echo "PREVIOUS_VERSION is not set"; exit 1; fi

.PHONY: check-github-token
check-github-token:
	@if [ -z "${GITHUB_TOKEN}" ]; then echo "GITHUB_TOKEN is not set"; exit 1; fi

.PHONY: release
release: clean-release check-release-tag $(RELEASE_DIR)  ## Build and push container images using the latest git tag for the commit
	git checkout "${RELEASE_TAG}"
	CORE_CONTROLLER_IMG=$(PROD_REGISTRY)/$(CORE_IMAGE_NAME) $(MAKE) release-manifests
	$(MAKE) release-templates
	$(MAKE) release-binaries

.PHONY: release-manifests
release-manifests: ## Build the manifests to publish with a release
	$(MAKE) $(RELEASE_DIR)/$(CORE_MANIFEST_FILE).yaml TAG=$(RELEASE_TAG)
	# Add metadata to the release artifacts
	cp metadata.yaml $(RELEASE_DIR)/metadata.yaml

.PHONY: release-staging
release-staging: ## Build and push container images to the staging bucket
	$(MAKE) docker-build-all
	$(MAKE) docker-build-core-image
	$(MAKE) release-alias-tag
	$(MAKE) staging-manifests
	$(MAKE) release-templates
	$(MAKE) release-binaries
	$(MAKE) upload-staging-artifacts

.PHONY: staging-manifests
staging-manifests:
	$(MAKE) $(RELEASE_DIR)/$(CORE_MANIFEST_FILE).yaml TAG=$(RELEASE_ALIAS_TAG)
	cp metadata.yaml $(RELEASE_DIR)/metadata.yaml

.PHONY: upload-staging-artifacts
upload-staging-artifacts: ## Upload release artifacts to the staging bucket
	gsutil cp $(RELEASE_DIR)/* gs://$(BUCKET)/components/$(RELEASE_ALIAS_TAG)

.PHONY: release-alias-tag
release-alias-tag: ## Add the release alias tag to the last build tag
	gcloud container images add-tag -q $(CORE_CONTROLLER_IMG):$(TAG) $(CORE_CONTROLLER_IMG):$(RELEASE_ALIAS_TAG)

.PHONY: release-templates
release-templates: $(RELEASE_DIR) ## Generate release templates
	cp templates/cluster-template*.yaml $(RELEASE_DIR)/

.PHONY: release-binaries
release-binaries: ## Builds the binaries to publish with a release
	RELEASE_BINARY=./cmd/capibmadm GOOS=linux GOARCH=ppc64le $(MAKE) release-binary
	RELEASE_BINARY=./cmd/capibmadm GOOS=linux GOARCH=amd64 $(MAKE) release-binary
	RELEASE_BINARY=./cmd/capibmadm GOOS=linux GOARCH=arm64 $(MAKE) release-binary
	RELEASE_BINARY=./cmd/capibmadm GOOS=darwin GOARCH=amd64 $(MAKE) release-binary
	RELEASE_BINARY=./cmd/capibmadm GOOS=darwin GOARCH=arm64 $(MAKE) release-binary
	RELEASE_BINARY=./cmd/capibmadm GOOS=windows GOARCH=amd64 EXT=.exe $(MAKE) release-binary
	RELEASE_BINARY=./cmd/capibmadm GOOS=windows GOARCH=arm64 EXT=.exe $(MAKE) release-binary

.PHONY: release-binary
release-binary: $(RELEASE_DIR) build-toolchain ## Release binary
	docker run \
		--rm \
		-e CGO_ENABLED=0 \
		-e GOOS=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		--mount=source=gocache,target=/go/pkg/mod \
		--mount=source=gocache,target=/root/.cache/go-build \
		-v "$$(pwd):/workspace" \
		-w /workspace \
		$(TOOLCHAIN_IMAGE) \
		go build -ldflags '$(LDFLAGS) -extldflags "-static"' \
		-o $(RELEASE_DIR)/$(notdir $(RELEASE_BINARY))-$(GOOS)-$(GOARCH)$(EXT) $(RELEASE_BINARY)

IMAGE_PATCH_DIR := $(ARTIFACTS)/image-patch

$(IMAGE_PATCH_DIR): $(ARTIFACTS)
	mkdir -p $@

.PHONY: $(RELEASE_DIR)/$(CORE_MANIFEST_FILE).yaml
$(RELEASE_DIR)/$(CORE_MANIFEST_FILE).yaml:
	$(MAKE) compiled-manifest \
		PROVIDER=$(CORE_MANIFEST_FILE) \
		OLD_IMG=$(CORE_CONTROLLER_ORIGINAL_IMG) \
		MANIFEST_IMG=$(CORE_CONTROLLER_IMG) \
		CONTROLLER_NAME=$(CORE_CONTROLLER_NAME) \
		PROVIDER_CONFIG_DIR=$(CORE_CONFIG_DIR) \
		NAMESPACE=$(CORE_NAMESPACE) \

.PHONY: compiled-manifest
compiled-manifest: $(RELEASE_DIR) $(KUSTOMIZE)
	$(MAKE) image-patch-source-manifest
	$(MAKE) image-patch-kustomization
	$(KUSTOMIZE) build $(IMAGE_PATCH_DIR)/$(PROVIDER) > $(RELEASE_DIR)/$(PROVIDER).yaml

.PHONY: image-patch-source-manifest
image-patch-source-manifest: $(IMAGE_PATCH_DIR) $(KUSTOMIZE)
	mkdir -p $(IMAGE_PATCH_DIR)/$(PROVIDER)
	$(KUSTOMIZE) build $(PROVIDER_CONFIG_DIR) > $(IMAGE_PATCH_DIR)/$(PROVIDER)/source-manifest.yaml

.PHONY: image-patch-kustomization
image-patch-kustomization: $(IMAGE_PATCH_DIR)
	mkdir -p $(IMAGE_PATCH_DIR)/$(PROVIDER)
	$(MAKE) image-patch-kustomization-without-webhook

.PHONY: image-patch-kustomization-without-webhook
image-patch-kustomization-without-webhook: $(IMAGE_PATCH_DIR) $(GOJQ)
	mkdir -p $(IMAGE_PATCH_DIR)/$(PROVIDER)
	$(GOJQ) --yaml-input --yaml-output '.images[0]={"name":"$(OLD_IMG)","newName":"$(MANIFEST_IMG)","newTag":"$(TAG)"}' \
		"hack/image-patch/kustomization.yaml" > $(IMAGE_PATCH_DIR)/$(PROVIDER)/kustomization.yaml

## --------------------------------------
## Docker
## --------------------------------------

.PHONY: ensure-buildx
ensure-buildx:
	./hack/init-buildx.sh

.PHONY: docker-build
docker-build: docker-pull-prerequisites ensure-buildx ## Build the docker image for controller-manager
	docker buildx build --platform linux/$(ARCH) --output=$(OUTPUT_TYPE) --pull --build-arg golang_image=$(GO_CONTAINER_IMAGE) --build-arg LDFLAGS="$(LDFLAGS)" -t $(CORE_CONTROLLER_IMG)-$(ARCH):$(TAG) .

.PHONY: docker-pull-prerequisites
docker-pull-prerequisites:
	docker pull docker.io/docker/dockerfile:1.5
	docker pull gcr.io/distroless/static:latest

.PHONY: e2e-image
e2e-image: docker-pull-prerequisites ensure-buildx
	docker buildx build --platform linux/$(ARCH) --load --build-arg golang_image=$(GO_CONTAINER_IMAGE) --tag=$(CORE_CONTROLLER_ORIGINAL_IMG):e2e .
	$(MAKE) set-manifest-image MANIFEST_IMG=$(CORE_CONTROLLER_ORIGINAL_IMG):e2e TARGET_RESOURCE="./config/default/manager_image_patch.yaml"
	$(MAKE) set-manifest-pull-policy PULL_POLICY=Never TARGET_RESOURCE="./config/default/manager_pull_policy.yaml"

.PHONY: set-manifest-image
set-manifest-image:
	$(info Updating kustomize image patch file for default resource)
	sed -i'' -e 's@image: .*@image: '"${MANIFEST_IMG}"'@' ./config/default/manager_image_patch.yaml

.PHONY: set-manifest-pull-policy
set-manifest-pull-policy:
	$(info Updating kustomize pull policy file for default resource)
	sed -i'' -e 's@imagePullPolicy: .*@imagePullPolicy: '"$(PULL_POLICY)"'@' ./config/default/manager_pull_policy.yaml	

## --------------------------------------
## Docker - All ARCH
## --------------------------------------

.PHONY: docker-build-all ## Build docker images for all architectures
docker-build-all: $(addprefix docker-build-,$(ALL_ARCH))

docker-build-%:
	$(MAKE) ARCH=$* docker-build

.PHONY: docker-build-core-image
docker-build-core-image: ensure-buildx ## Build the multiarch core docker image
	docker buildx build --builder capibm --platform $(BUILDX_PLATFORMS)  --output=$(OUTPUT_TYPE) \
		--pull --build-arg golang_image=$(GO_CONTAINER_IMAGE) --build-arg ldflags="$(LDFLAGS)" \
		-t $(CORE_CONTROLLER_IMG):$(TAG) .

## --------------------------------------
## Lint / Verify
## --------------------------------------

##@ lint and verify:

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint codebase
	$(GOLANGCI_LINT) run -v --fast=false

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Lint the codebase and run auto-fixers if supported by the linter
	GOLANGCI_LINT_EXTRA_ARGS=--fix $(MAKE) lint

APIDIFF_OLD_COMMIT ?= $(shell git rev-parse origin/main)

.PHONY: apidiff
apidiff: $(GO_APIDIFF) ## Check for API differences.
	@if ($(call checkdiff) | grep "api/"); then \
		$(GO_APIDIFF) $(APIDIFF_OLD_COMMIT) --print-compatible; \
	else \
		echo "No changes to 'api/'. Nothing to do."; \
	fi

define checkdiff
	git --no-pager diff --name-only FETCH_HEAD
endef

ALL_VERIFY_CHECKS = boilerplate shellcheck modules gen conversions go-version

.PHONY: verify
verify: $(addprefix verify-,$(ALL_VERIFY_CHECKS)) ## Run all verify-* targets

.PHONY: verify-boilerplate
verify-boilerplate: ## Verify boilerplate text exists in each file
	./hack/verify-boilerplate.sh

.PHONY: verify-shellcheck
verify-shellcheck: ## Verify shell files
	./hack/verify-shellcheck.sh

.PHONY: verify-modules
verify-modules: generate-modules ## Verify go modules are up to date
	@if !(git diff --quiet HEAD -- go.sum go.mod $(TOOLS_DIR)/go.mod $(TOOLS_DIR)/go.sum); then \
		git diff; \
		echo "go module files are out of date"; exit 1; \
	fi
	@if (find . -name 'go.mod' | xargs -n1 grep -q -i 'k8s.io/client-go.*+incompatible'); then \
		find . -name "go.mod" -exec grep -i 'k8s.io/client-go.*+incompatible' {} \; -print; \
		echo "go module contains an incompatible client-go version"; exit 1; \
	fi

.PHONY: verify-gen
verify-gen: generate ## Verfiy go generated files are up to date
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi

.PHONY: verify-conversions
verify-conversions: $(CONVERSION_VERIFIER) ## Verifies expected API conversion are in place
	$(CONVERSION_VERIFIER)

.PHONY: verify-container-images
verify-container-images: $(TRIVY) ## Verify container images
	TRACE=$(TRACE) ./hack/verify-container-images.sh

.PHONY: verify-govulncheck
verify-govulncheck: $(GOVULNCHECK) ## Verify code for vulnerabilities
	$(GOVULNCHECK) ./... && R1=$$? || R1=$$?; \
	$(GOVULNCHECK) -C "$(TOOLS_DIR)" ./... && R2=$$? || R2=$$?; \
	if [ "$$R1" -ne "0" ] || [ "$$R2" -ne "0" ]; then \
		exit 1; \
	fi

.PHONY: verify-security
verify-security: ## Verify code and images for vulnerabilities
	$(MAKE) verify-container-images && R1=$$? || R1=$$?; \
	$(MAKE) verify-govulncheck && R2=$$? || R2=$$?; \
	if [ "$$R1" -ne "0" ] || [ "$$R2" -ne "0" ]; then \
	  echo "Check for vulnerabilities failed! There are vulnerabilities to be fixed"; \
		exit 1; \
	fi


SUBMAKEFILE_GO_VERSION=$(shell grep "GO_VERSION" -m1 hack/ccm/Makefile | cut -d '=' -f2 )
.SILENT:
.PHONY: verify-go-version
verify-go-version: ## Confirms the version of go used in Makefiles are uniform
ifeq ($(strip $(SUBMAKEFILE_GO_VERSION)), $(strip $(GO_VERSION)))
	echo "Versions are uniform across Makefile, the go-version used is $(GO_VERSION)"
else
	echo "Versions are different across Makefiles. Please ensure to keep them uniform."
endif


.PHONY: yamllint
yamllint:
	@docker run --rm $$(tty -s && echo "-it" || echo) -v $(PWD):/data cytopia/yamllint:latest /data --config-file /data/.yamllint --no-warnings

## --------------------------------------
## Cleanup / Verification
## --------------------------------------

##@ clean:

.PHONY: clean
clean: ## Remove all generated files
	$(MAKE) clean-bin
	$(MAKE) clean-book
	$(MAKE) clean-temporary

.PHONY: clean-bin
clean-bin: ## Remove all generated binaries
	rm -rf $(TOOLS_BIN_DIR)
	rm -rf $(BIN_DIR)

.PHONY: clean-book
clean-book: ## Remove all generated GitBook files
	rm -rf ./docs/book/_book

.PHONY: clean-temporary
clean-temporary: ## Remove all temporary files and folders
	rm -f minikube.kubeconfig
	rm -f kubeconfig
	rm -rf _artifacts

.PHONY: clean-release
clean-release: ## Remove the release folder
	rm -rf $(RELEASE_DIR)

.PHONY: clean-release-git
clean-release-git: ## Restores the git files usually modified during a release
	git restore ./*manager_image_patch.yaml ./*manager_pull_policy.yaml

.PHONY: clean-generated-conversions
clean-generated-conversions: ## Remove files generated by conversion-gen from the mentioned dirs
	(IFS=','; for i in $(SRC_DIRS); do find $$i -type f -name 'zz_generated.conversion*' -exec rm -f {} \;; done)

.PHONY: clean-kind
clean-kind: ## Cleans up the kind cluster with the name $CAPI_KIND_CLUSTER_NAME
	kind delete cluster --name="$(CAPI_KIND_CLUSTER_NAME)" || true

## --------------------------------------
## Kind
## --------------------------------------

##@ kind:

.PHONY: kind-cluster
kind-cluster: ## Create a new kind cluster designed for development with Tilt
	hack/kind-install.sh

## --------------------------------------
## Helpers
## --------------------------------------

##@ helpers:

go-version: ## Print the go version we use to compile our binaries and images
	@echo $(GO_VERSION)

## --------------------------------------
## Documentation and Publishing
## --------------------------------------

.PHONY: serve
serve: ## Build the CAPIBM book and serve it locally to validate changes in documentation.
	$(MAKE) -C docs/book/ serve

## --------------------------------------
## Update Go Version
## --------------------------------------
.PHONY: update-go
update-go: ## Update Go version across files: Usage make update-go VERSION=X.YY.ZZ- Use only if you know what you're doing.
ifndef VERSION
	echo "VERSION not set. Usage: make update-go VERSION=X.YY.ZZ"
else
	sed -i '' "s/GO_VERSION ?=[[:digit:]].[[:digit:]]\{1,\}.[[:digit:]]\{1,\}/GO_VERSION ?=$(VERSION)/" hack/ccm/Makefile
	echo "Updated go version to $(VERSION) in Makefile"
endif
