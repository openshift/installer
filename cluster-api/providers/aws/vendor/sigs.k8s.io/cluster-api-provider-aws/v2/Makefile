# Copyright 2018 The Kubernetes Authors.
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

ROOT_DIR_RELATIVE := .

include $(ROOT_DIR_RELATIVE)/common.mk

# If you update this file, please follow
# https://suva.sh/posts/well-documented-makefiles

# Go
GO_VERSION ?=1.22.6
GO_CONTAINER_IMAGE ?= golang:$(GO_VERSION)

# Directories.
ARTIFACTS ?= $(REPO_ROOT)/_artifacts
TOOLS_DIR := hack/tools
TOOLS_DIR_DEPS := $(TOOLS_DIR)/go.sum $(TOOLS_DIR)/go.mod $(TOOLS_DIR)/Makefile
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin


API_DIRS := cmd/clusterawsadm/api api exp/api controlplane/eks/api bootstrap/eks/api iam/api controlplane/rosa/api
API_FILES := $(foreach dir, $(API_DIRS), $(call rwildcard,../../$(dir),*.go))

BIN_DIR := bin
REPO_ROOT := $(shell git rev-parse --show-toplevel)
GH_ORG_NAME ?= kubernetes-sigs
GH_REPO_NAME ?= cluster-api-provider-aws
GH_REPO ?= $(GH_ORG_NAME)/$(GH_REPO_NAME)
TEST_E2E_DIR := test/e2e

# Files
E2E_DATA_DIR ?= $(REPO_ROOT)/test/e2e/data
E2E_CONF_PATH  ?= $(E2E_DATA_DIR)/e2e_conf.yaml
E2E_EKS_CONF_PATH ?= $(E2E_DATA_DIR)/e2e_eks_conf.yaml
KUBETEST_CONF_PATH ?= $(abspath $(E2E_DATA_DIR)/kubetest/conformance.yaml)
EXP_DIR := exp
GORELEASER_CONFIG := .goreleaser.yaml

# Binaries.
GO_INSTALL := ./scripts/go_install.sh
GO_APIDIFF_BIN := $(BIN_DIR)/go-apidiff
GO_APIDIFF := $(TOOLS_DIR)/$(GO_APIDIFF_BIN)
CLUSTERCTL := $(BIN_DIR)/clusterctl
CONTROLLER_GEN := $(TOOLS_BIN_DIR)/controller-gen
CONVERSION_GEN := $(TOOLS_BIN_DIR)/conversion-gen
CONVERSION_VERIFIER := $(TOOLS_BIN_DIR)/conversion-verifier
DEFAULTER_GEN := $(TOOLS_BIN_DIR)/defaulter-gen
ENVSUBST := $(TOOLS_BIN_DIR)/envsubst
GH := $(TOOLS_BIN_DIR)/gh
GOJQ := $(TOOLS_BIN_DIR)/gojq
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT_VER := $(shell cat .github/workflows/pr-golangci-lint.yaml | grep [[:space:]]version: | sed 's/.*version: //')
GOLANGCI_LINT := $(abspath $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_BIN)-$(GOLANGCI_LINT_VER))
GOLANGCI_LINT_PKG := github.com/golangci/golangci-lint/cmd/golangci-lint
KIND := $(TOOLS_BIN_DIR)/kind
KUSTOMIZE := $(TOOLS_BIN_DIR)/kustomize
MOCKGEN := $(TOOLS_BIN_DIR)/mockgen
SSM_PLUGIN := $(TOOLS_BIN_DIR)/session-manager-plugin
YQ := $(TOOLS_BIN_DIR)/yq
KPROMO := $(TOOLS_BIN_DIR)/kpromo
RELEASE_NOTES := $(TOOLS_BIN_DIR)/release-notes
GORELEASER := $(TOOLS_BIN_DIR)/goreleaser

CLUSTERAWSADM_SRCS := $(call rwildcard,.,cmd/clusterawsadm/*.*)

PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)
DOCKER_CLI_EXPERIMENTAL=enabled
DOCKER_BUILDKIT=1

export ACK_GINKGO_DEPRECATIONS := 1.16.4

# Set --output-base for conversion-gen if we are not within GOPATH
ifneq ($(abspath $(REPO_ROOT)),$(abspath $(shell go env GOPATH)/src/sigs.k8s.io/cluster-api-provider-aws))
	GEN_OUTPUT_BASE := --output-base=$(REPO_ROOT)
else
	export GOPATH := $(shell go env GOPATH)
endif

# Release variables

STAGING_REGISTRY ?= gcr.io/k8s-staging-cluster-api-aws
STAGING_BUCKET ?= k8s-staging-cluster-api-aws
BUCKET ?= $(STAGING_BUCKET)
PROD_REGISTRY := registry.k8s.io/cluster-api-aws
REGISTRY ?= $(STAGING_REGISTRY)
RELEASE_TAG ?= $(shell git describe --abbrev=0 2>/dev/null)
PULL_BASE_REF ?= $(RELEASE_TAG) # PULL_BASE_REF will be provided by Prow
RELEASE_ALIAS_TAG ?= $(PULL_BASE_REF)
RELEASE_DIR := out
RELEASE_POLICIES := $(RELEASE_DIR)/AWSIAMManagedPolicyControllers.json $(RELEASE_DIR)/AWSIAMManagedPolicyControllersWithEKS.json $(RELEASE_DIR)/AWSIAMManagedPolicyCloudProviderControlPlane.json $(RELEASE_DIR)/AWSIAMManagedPolicyCloudProviderNodes.json $(RELEASE_DIR)/AWSIAMManagedPolicyControllersWithS3.json
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
USER_FORK ?= $(shell git config --get remote.origin.url | cut -d/ -f4) # only works on https://github.com/<username>/cluster-api-provider-aws.git style URLs
ifeq ($(strip $(USER_FORK)),)
USER_FORK := $(shell git config --get remote.origin.url | cut -d: -f2 | cut -d/ -f1) # for git@github.com:<username>/cluster-api-provider-aws.git style URLs
endif

# image name used to build the cmd/clusterawsadm
TOOLCHAIN_IMAGE := toolchain

TAG ?= dev
ARCH ?= $(shell go env GOARCH)
ALL_ARCH ?= amd64 arm arm64 ppc64le s390x

# main controller
CORE_IMAGE_NAME ?= cluster-api-aws-controller
CORE_CONTROLLER_IMG ?= $(REGISTRY)/$(CORE_IMAGE_NAME)
CORE_CONTROLLER_PROMOTED_IMG := $(PROD_REGISTRY)/$(CORE_IMAGE_NAME)
CORE_CONTROLLER_ORIGINAL_IMG := gcr.io/k8s-staging-cluster-api-aws/cluster-api-aws-controller
CORE_CONTROLLER_NAME := capa-controller-manager
CORE_MANIFEST_FILE := infrastructure-components
CORE_CONFIG_DIR := config/default
CORE_NAMESPACE := capa-system

# Allow overriding manifest generation destination directory
MANIFEST_ROOT ?= config
CRD_ROOT ?= $(MANIFEST_ROOT)/crd/bases
CRD_DOCS_DIR := docs/book/src/crd
CRD_DOCS := $(CRD_DOCS_DIR)/index.md
WEBHOOK_ROOT ?= $(MANIFEST_ROOT)/webhook
RBAC_ROOT ?= $(MANIFEST_ROOT)/rbac

# Allow overriding the imagePullPolicy
PULL_POLICY ?= Always

# Set build time variables including version details
LDFLAGS := $(shell source ./hack/version.sh; version::ldflags)

# Set USE_EXISTING_CLUSTER to use an existing kubernetes context
USE_EXISTING_CLUSTER ?= "false"

# Set E2E_SKIP_EKS_UPGRADE to false to test EKS upgrades.
# Warning, this takes a long time
E2E_SKIP_EKS_UPGRADE ?= "false"

# Set EKS_SOURCE_TEMPLATE to override the source template
EKS_SOURCE_TEMPLATE ?= eks/cluster-template-eks-control-plane-only.yaml

# set up `setup-envtest` to install kubebuilder dependency
export KUBEBUILDER_ENVTEST_KUBERNETES_VERSION ?= 1.30.2
SETUP_ENVTEST_VER := v0.0.0-20240923090159-236e448db12c
SETUP_ENVTEST_BIN := setup-envtest
SETUP_ENVTEST := $(abspath $(TOOLS_BIN_DIR)/$(SETUP_ENVTEST_BIN)-$(SETUP_ENVTEST_VER))
SETUP_ENVTEST_PKG := sigs.k8s.io/controller-runtime/tools/setup-envtest

# Enable Cluster API Framework tests for the purposes of running the PR blocking test
ifeq ($(findstring \[PR-Blocking\],$(GINKGO_FOCUS)),\[PR-Blocking\])
  override undefine GINKGO_SKIP
endif

override E2E_ARGS += -artifacts-folder="$(ARTIFACTS)" --data-folder="$(E2E_DATA_DIR)" -use-existing-cluster=$(USE_EXISTING_CLUSTER)
override GINKGO_ARGS += -v --trace --timeout=4h --output-dir="$(ARTIFACTS)" --junit-report="junit.e2e_suite.xml"

ifdef GINKGO_SKIP
	override GINKGO_ARGS += -skip "$(GINKGO_SKIP)"
endif

# DEPRECATED, use GINKGO_FOCUS instead
ifdef E2E_UNMANAGED_FOCUS
	override GINKGO_ARGS += -focus="$(E2E_UNMANAGED_FOCUS)"
endif

# ALL tests will take ~ 1 hour @ 24 node concurrency.
# Set the number of nodes using GINKGO_ARGS=-nodes 24
# Ginkgo will default to the number of logical CPUs you have available.
# Should be safe to set more nodes than available CPU cores as most of the time is spent in
# infrastructure reconciliation

# Instead, you can run a quick smoke test, it should run fast (9 minutes)...
# GINKGO_FOCUS := "\\[smoke\\]"
# For running CAPI e2e tests: GINKGO_FOCUS := "\\[Cluster API Framework\\]"
# For running CAPI blocking e2e test: GINKGO_FOCUS := "\\[PR-Blocking\\]"
ifdef GINKGO_FOCUS
	override GINKGO_ARGS += -focus="$(GINKGO_FOCUS)"
endif

ifeq ($(E2E_SKIP_EKS_UPGRADE),"true")
	override EKS_E2E_ARGS += --skip-eks-upgrade-tests
endif

##@ generate:

.PHONY: defaulters
defaulters: $(DEFAULTER_GEN) ## Generate all Go types
	$(DEFAULTER_GEN) \
		--extra-peer-dirs=sigs.k8s.io/cluster-api/api/v1beta1 \
		--v=0 \
		--go-header-file=./hack/boilerplate/boilerplate.generatego.txt \
		--output-file=zz_generated.defaults.go \
		./api/v1beta2 \
		./$(EXP_DIR)/api/v1beta2 \
		./controlplane/rosa/api/v1beta2 \
		./cmd/clusterawsadm/api/bootstrap/v1beta1 \
		./cmd/clusterawsadm/api/bootstrap/v1alpha1

.PHONY: generate
generate: ## Generate code
	$(MAKE) generate-go
	$(MAKE) $(CRD_DOCS)

$(CRD_DOCS_DIR)/%: $(API_FILES)
	$(MAKE) -C docs/book src/crd/$*

.PHONY: generate-go ## Generate all Go api files
generate-go: $(MOCKGEN)
	$(MAKE) generate-go-apis
	go generate ./...

.PHONY: generate-go-apis
generate-go-apis: ## Alias for .build/generate-go-apis
	rm -rf .build/generate-go-apis
	$(MAKE) .build/generate-go-apis

.build: ## Create the .build folder
	mkdir -p .build

.build/generate-go-apis: .build $(API_FILES) $(CONTROLLER_GEN) $(DEFAULTER_GEN) $(CONVERSION_GEN) ## Generate all Go api files
	$(CONTROLLER_GEN) \
		paths=./ \
		paths=./api/... \
		paths=./$(EXP_DIR)/api/... \
		paths=./bootstrap/eks/api/... \
		paths=./controlplane/eks/api/... \
		paths=./controlplane/rosa/api/... \
		paths=./iam/api/... \
		paths=./controllers/... \
		paths=./$(EXP_DIR)/controllers/... \
		paths=./bootstrap/eks/controllers/... \
		paths=./controlplane/eks/controllers/... \
		paths=./controlplane/rosa/controllers/... \
		output:crd:dir=config/crd/bases \
		object:headerFile=./hack/boilerplate/boilerplate.generatego.txt \
		crd:crdVersions=v1 \
		rbac:roleName=manager-role \
		webhook

	$(CONTROLLER_GEN) \
		paths=./cmd/... \
		object:headerFile=./hack/boilerplate/boilerplate.generatego.txt

	$(MAKE) defaulters

	$(CONVERSION_GEN) \
		--extra-peer-dirs=sigs.k8s.io/cluster-api/api/v1beta1 \
		--output-file=zz_generated.conversion.go \
		--go-header-file=./hack/boilerplate/boilerplate.generatego.txt \
		./api/v1beta1 \
		./cmd/clusterawsadm/api/bootstrap/v1alpha1

	$(CONVERSION_GEN) \
		--extra-peer-dirs=sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1 \
		--extra-peer-dirs=sigs.k8s.io/cluster-api/api/v1beta1 \
		--output-file=zz_generated.conversion.go \
		--go-header-file=./hack/boilerplate/boilerplate.generatego.txt \
		./$(EXP_DIR)/api/v1beta1

	$(CONVERSION_GEN) \
		--extra-peer-dirs=sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1 \
		--extra-peer-dirs=sigs.k8s.io/cluster-api/api/v1beta1 \
		--output-file=zz_generated.conversion.go \
		--go-header-file=./hack/boilerplate/boilerplate.generatego.txt \
		./bootstrap/eks/api/v1beta1

	$(CONVERSION_GEN) \
		--extra-peer-dirs=sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1 \
		--extra-peer-dirs=sigs.k8s.io/cluster-api/api/v1beta1 \
		--output-file=zz_generated.conversion.go \
		--go-header-file=./hack/boilerplate/boilerplate.generatego.txt \
		./controlplane/eks/api/v1beta1

	$(CONVERSION_GEN) \
		--extra-peer-dirs=sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1 \
		--extra-peer-dirs=sigs.k8s.io/cluster-api/api/v1beta1 \
		--output-file=zz_generated.conversion.go \
		--go-header-file=./hack/boilerplate/boilerplate.generatego.txt \
		./controlplane/rosa/api/v1beta2

		touch $@

##@ lint and verify:

.PHONY: modules

$(GOLANGCI_LINT): # Build golangci-lint from tools folder.
	GOBIN=$(abspath $(TOOLS_BIN_DIR)) $(GO_INSTALL) $(GOLANGCI_LINT_PKG) $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_VER)

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint codebase
	$(GOLANGCI_LINT) run -v --fast=false $(GOLANGCI_LINT_EXTRA_ARGS)

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Lint the codebase and run auto-fixers if supported by the linter
	GOLANGCI_LINT_EXTRA_ARGS=--fix $(MAKE) lint

modules: ## Runs go mod to ensure proper vendoring.
	go mod tidy
	cd $(TOOLS_DIR); go mod tidy

.PHONY: verify ## Verify ties together the rest of the verification targets into one target
verify: verify-boilerplate verify-modules verify-gen verify-conversions verify-shellcheck verify-book-links release-manifests

.PHONY: verify-boilerplate
verify-boilerplate: ## Verify boilerplate
	-rm -f ./hack/tools/bin/*.sh
	./hack/verify-boilerplate.sh

.PHONY: verify-modules
verify-modules: modules ## Verify go modules integrity
	@if !(git diff --quiet HEAD -- go.sum go.mod hack/tools/go.mod hack/tools/go.sum); then \
		git diff; \
		echo "go module files are out of date"; exit 1; \
	fi

.PHONY: verify-conversions
verify-conversions: $(CONVERSION_VERIFIER)  ## Verifies expected API conversion are in place
	echo verification of api conversions initiated
	$(CONVERSION_VERIFIER)

.PHONY: verify-shellcheck
verify-shellcheck: ## Verify shell files
	./hack/verify-shellcheck.sh

.PHONY: verify-book-links
verify-book-links: ## Verify book links
	$(MAKE) -C docs/book verify

.PHONY: verify-gen
verify-gen: generate ## Verify generated files
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi

.PHONY: verify-container-images
verify-container-images: ## Verify container images
	TRACE=$(TRACE) ./hack/verify-container-images.sh

.PHONY: apidiff
apidiff: APIDIFF_OLD_COMMIT ?= $(shell git rev-parse origin/main)
apidiff: $(GO_APIDIFF) ## Check for API differences
	@$(call checkdiff) > /dev/null
	@if ($(call checkdiff) | grep "api/"); then \
		$(GO_APIDIFF) $(APIDIFF_OLD_COMMIT) --print-compatible; \
	else \
		echo "No changes to 'api/'. Nothing to do."; \
	fi

define checkdiff
	git --no-pager diff --name-only FETCH_HEAD
endef

##@ build:

.PHONY: binaries
binaries: managers clusterawsadm ## Builds and installs all binaries

.PHONY: clusterawsadm
clusterawsadm: ## Build clusterawsadm binary
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/clusterawsadm ./cmd/clusterawsadm


.PHONY: docker-build
docker-build: docker-pull-prerequisites ## Build the docker image for controller-manager
	docker build --build-arg ARCH=$(ARCH) --build-arg builder_image=$(GO_CONTAINER_IMAGE) --build-arg LDFLAGS="$(LDFLAGS)" . -t $(CORE_CONTROLLER_IMG)-$(ARCH):$(TAG)

.PHONY: docker-build-all ## Build all the architecture docker images
docker-build-all: $(addprefix docker-build-,$(ALL_ARCH))

docker-build-%: ## Build docker images for a given ARCH
	$(MAKE) ARCH=$* docker-build

.PHONY: docker-pull-prerequisites
docker-pull-prerequisites: ## Build prerequisites using docker
	docker pull docker.io/docker/dockerfile:1.1-experimental
	docker pull gcr.io/distroless/static:latest

.PHONY: managers
managers: ## Alias for manager-aws-infrastructure
	$(MAKE) manager-aws-infrastructure

.PHONY: manager-aws-infrastructure
manager-aws-infrastructure: ## Build manager binary
	CGO_ENABLED=0 GOARCH=${ARCH} go build -ldflags "${LDFLAGS} -extldflags '-static'" -o $(BIN_DIR)/manager .

##@ test:

$(ARTIFACTS):
	mkdir -p $@

.PHONY: generate-test-flavors
generate-test-flavors: $(KUSTOMIZE)  ## Generate test template flavors
	./hack/gen-test-flavors.sh withoutclusterclass
	./hack/gen-test-flavors.sh withclusterclass

.PHONY: e2e-image
e2e-image: docker-pull-prerequisites $(TOOLS_BIN_DIR)/start.sh $(TOOLS_BIN_DIR)/restart.sh ## Build an e2e test image
	docker build --build-arg builder_image=$(GO_CONTAINER_IMAGE) -f Dockerfile --tag="gcr.io/k8s-staging-cluster-api/capa-manager:e2e" .

.PHONY: install-setup-envtest
install-setup-envtest: # Install setup-envtest so that setup-envtest's eval is executed after the tool has been installed.
	GOBIN=$(abspath $(TOOLS_BIN_DIR)) $(GO_INSTALL) $(SETUP_ENVTEST_PKG) $(SETUP_ENVTEST_BIN) $(SETUP_ENVTEST_VER)

.PHONY: setup-envtest
setup-envtest: install-setup-envtest # Build setup-envtest from tools folder.
	@$(eval KUBEBUILDER_ASSETS := $(shell $(SETUP_ENVTEST) use --use-env -p path $(KUBEBUILDER_ENVTEST_KUBERNETES_VERSION))) \
	if [ -z "$(KUBEBUILDER_ASSETS)" ]; then echo "Failed to find kubebuilder assets, see errors above"; exit 1; fi; \
	echo "kube-builder assets: $(KUBEBUILDER_ASSETS)"

.PHONY: test
test: setup-envtest ## Run tests
	KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" go test ./...

.PHONY: test-verbose
test-verbose: setup-envtest ## Run tests with verbose settings.
	KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" go test -v ./...

.PHONY: test-e2e ## Run e2e tests using clusterctl
test-e2e: $(KIND) $(SSM_PLUGIN) $(KUSTOMIZE) generate-test-flavors e2e-image ## Run e2e tests
	time go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e $(GINKGO_ARGS) -p ./test/e2e/suites/unmanaged/... -- -config-path="$(E2E_CONF_PATH)" $(E2E_ARGS)

.PHONY: test-e2e-eks ## Run EKS e2e tests using clusterctl
test-e2e-eks: generate-test-flavors $(KIND) $(SSM_PLUGIN) $(KUSTOMIZE) e2e-image ## Run eks e2e tests
	time go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e $(GINKGO_ARGS) ./test/e2e/suites/managed/... -- -config-path="$(E2E_EKS_CONF_PATH)" --source-template="$(EKS_SOURCE_TEMPLATE)" $(E2E_ARGS) $(EKS_E2E_ARGS)

.PHONY: test-e2e-gc ## Run garbage collection e2e tests using clusterctl
test-e2e-gc: generate-test-flavors $(KIND) $(SSM_PLUGIN) $(KUSTOMIZE) e2e-image ## Run eks e2e tests
	time go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e -focus="$(GINKGO_FOCUS)" -skip="$(GINKGO_SKIP)" $(GINKGO_ARGS) -p ./test/e2e/suites/gc_unmanaged/... -- -config-path="$(E2E_CONF_PATH)" $(E2E_ARGS)

.PHONY: test-e2e-eks-gc ## Run EKS garbage collection e2e tests using clusterctl
test-e2e-eks-gc: generate-test-flavors $(KIND) $(SSM_PLUGIN) $(KUSTOMIZE) e2e-image ## Run eks e2e tests
	time go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e -focus="$(GINKGO_FOCUS)" -skip="$(GINKGO_SKIP)" $(GINKGO_ARGS) ./test/e2e/suites/gc_managed/... -- -config-path="$(E2E_EKS_CONF_PATH)" --source-template="$(EKS_SOURCE_TEMPLATE)" $(E2E_ARGS) $(EKS_E2E_ARGS)


CONFORMANCE_E2E_ARGS ?= -kubetest.config-file=$(KUBETEST_CONF_PATH)
CONFORMANCE_E2E_ARGS += $(E2E_ARGS)
CONFORMANCE_GINKGO_ARGS += $(GINKGO_ARGS)
.PHONY: test-conformance
test-conformance: generate-test-flavors $(KIND) $(SSM_PLUGIN) $(KUSTOMIZE) e2e-image ## Run clusterctl based conformance test on workload cluster (requires Docker).
	time go run github.com/onsi/ginkgo/v2/ginkgo -tags=e2e -focus="conformance" $(CONFORMANCE_GINKGO_ARGS) ./test/e2e/suites/conformance/... -- -config-path="$(E2E_CONF_PATH)" $(CONFORMANCE_E2E_ARGS)

.PHONY: test-cover
test-cover: setup-envtest ## Run tests with code coverage and code generate  reports
	KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" go test -coverprofile=coverage.out ./... $(TEST_ARGS)
	go tool cover -func=coverage.out -o coverage.txt
	go tool cover -html=coverage.out -o coverage.html

.PHONY: serve-book
serve-book: ## Run a server with the documentation book
	$(MAKE) -C docs/book serve

.PHONY: compile-e2e
compile-e2e: ## Test e2e compilation
	go test -c -o /dev/null -tags=e2e ./test/e2e/suites/unmanaged
	go test -c -o /dev/null -tags=e2e ./test/e2e/suites/conformance
	go test -c -o /dev/null -tags=e2e ./test/e2e/suites/managed
	go test -c -o /dev/null -tags=e2e ./test/e2e/suites/gc_managed
	go test -c -o /dev/null -tags=e2e ./test/e2e/suites/gc_unmanaged


.PHONY: docker-pull-e2e-preloads
docker-pull-e2e-preloads: ## Preloads the docker images used for e2e testing and can speed it up
	-docker pull registry.k8s.io/cluster-api/kubeadm-control-plane-controller:$(CAPI_VERSION)
	-docker pull registry.k8s.io/cluster-api/kubeadm-bootstrap-controller:$(CAPI_VERSION)
	-docker pull registry.k8s.io/cluster-api/cluster-api-controller:$(CAPI_VERSION)
	-docker pull quay.io/jetstack/cert-manager-controller:$(CERT_MANAGER_VERSION)
	-docker pull quay.io/jetstack/cert-manager-cainjector:$(CERT_MANAGER_VERSION)
	-docker pull quay.io/jetstack/cert-manager-webhook:$(CERT_MANAGER_VERSION)

##@ release:

$(RELEASE_DIR):
	mkdir -p $@

.PHONY: build-toolchain
build-toolchain: ## Build the toolchain
	docker build --build-arg builder_image=$(GO_CONTAINER_IMAGE) --target toolchain -t $(TOOLCHAIN_IMAGE) .

.PHONY: check-github-token
check-github-token: ## Check if the github token is set
	@if [ -z "${GITHUB_TOKEN}" ]; then echo "GITHUB_TOKEN is not set"; exit 1; fi

.PHONY: check-previous-release-tag
check-previous-release-tag: ## Check if the previous release tag is set
	@if [ -z "${PREVIOUS_VERSION}" ]; then echo "PREVIOUS_VERSION is not set"; exit 1; fi

.PHONY: check-release-tag
check-release-tag: ## Check if the release tag is set
	@if [ -z "${RELEASE_TAG}" ]; then echo "RELEASE_TAG is not set"; exit 1; fi
	@if ! [ -z "$$(git status --porcelain)" ]; then echo "Your local git repository contains uncommitted changes, use git clean before proceeding."; exit 1; fi

.PHONY: check-release-branch
check-release-branch: ## Check if the release branch is set
	@if [ -z "${RELEASE_BRANCH}" ]; then echo "RELEASE_BRANCH is not set"; exit 1; fi

.PHONY: compiled-manifest
compiled-manifest: $(RELEASE_DIR) $(KUSTOMIZE) ## Compile the manifest files
	$(MAKE) image-patch-source-manifest
	$(MAKE) image-patch-pull-policy
	$(MAKE) image-patch-kustomization
	$(KUSTOMIZE) build $(IMAGE_PATCH_DIR)/$(PROVIDER) > $(RELEASE_DIR)/$(PROVIDER).yaml

.PHONY: docker-push
docker-push: ## Push the docker image
	docker push $(CORE_CONTROLLER_IMG)-$(ARCH):$(TAG)

.PHONY: docker-push-all ## Push all the architecture docker images
docker-push-all: $(addprefix docker-push-,$(ALL_ARCH))
	$(MAKE) docker-push-core-manifest

docker-push-%: ## Docker push
	$(MAKE) ARCH=$* docker-push

.PHONY: docker-push-core-manifest
docker-push-core-manifest: ## Push the fat manifest docker image.
	## Minimum docker version 18.06.0 is required for creating and pushing manifest images.
	$(MAKE) docker-push-manifest CONTROLLER_IMG=$(CORE_CONTROLLER_IMG) MANIFEST_FILE=$(CORE_MANIFEST_FILE)

.PHONY: docker-push-manifest
docker-push-manifest: ## Push the manifest image
	docker manifest create --amend $(CONTROLLER_IMG):$(TAG) $(shell echo $(ALL_ARCH) | sed -e "s~[^ ]*~$(CONTROLLER_IMG)\-&:$(TAG)~g")
	@for arch in $(ALL_ARCH); do docker manifest annotate --arch $${arch} ${CONTROLLER_IMG}:${TAG} ${CONTROLLER_IMG}-$${arch}:${TAG}; done
	docker manifest push --purge ${CONTROLLER_IMG}:${TAG}

.PHONY: image-patch-source-manifest
image-patch-source-manifest: $(IMAGE_PATCH_DIR) $(KUSTOMIZE) ## Patch the source manifest
	mkdir -p $(IMAGE_PATCH_DIR)/$(PROVIDER)
	$(KUSTOMIZE) build $(PROVIDER_CONFIG_DIR) > $(IMAGE_PATCH_DIR)/$(PROVIDER)/source-manifest.yaml

.PHONY: image-patch-kustomization
image-patch-kustomization: $(IMAGE_PATCH_DIR) ## Alias for image-patch-kustomization-without-webhook
	mkdir -p $(IMAGE_PATCH_DIR)/$(PROVIDER)
	$(MAKE) image-patch-kustomization-without-webhook

.PHONY: image-patch-kustomization-without-webhook
image-patch-kustomization-without-webhook: $(IMAGE_PATCH_DIR) $(GOJQ) ## Patch the image in the kustomization file
	mkdir -p $(IMAGE_PATCH_DIR)/$(PROVIDER)
	$(GOJQ) --yaml-input --yaml-output '.images[0]={"name":"$(OLD_IMG)","newName":"$(MANIFEST_IMG)","newTag":"$(TAG)"}|.patchesJson6902[0].target.name="$(CONTROLLER_NAME)"|.patchesJson6902[0].target.namespace="$(NAMESPACE)"' \
		"hack/image-patch/kustomization.yaml" > $(IMAGE_PATCH_DIR)/$(PROVIDER)/kustomization.yaml

.PHONY: image-patch-pull-policy
image-patch-pull-policy: $(IMAGE_PATCH_DIR) $(GOJQ) ## Patch the pull policy
	mkdir -p $(IMAGE_PATCH_DIR)/$(PROVIDER)
	echo Setting imagePullPolicy to $(PULL_POLICY)
	$(GOJQ) --yaml-input --yaml-output '.[0].value="$(PULL_POLICY)"' "hack/image-patch/pull-policy-patch.yaml" > $(IMAGE_PATCH_DIR)/$(PROVIDER)/pull-policy-patch.yaml

.PHONY: list-staging-releases
list-staging-releases: ## List staging images for image promotion
	@echo $(CORE_IMAGE_NAME):
	$(MAKE) list-image RELEASE_TAG=$(RELEASE_TAG) IMAGE=$(CORE_IMAGE_NAME)

list-image: ## List images for RELEASE_TAG
	gcloud container images list-tags $(STAGING_REGISTRY)/$(IMAGE) --filter="tags=('$(RELEASE_TAG)')" --format=json

.PHONY: release
release: clean-release check-release-tag check-release-branch $(RELEASE_DIR) $(GORELEASER)  ## Builds and push container images using the latest git tag for the commit.
	git checkout "${RELEASE_TAG}"
	$(MAKE) release-changelog
	CORE_CONTROLLER_IMG=$(PROD_REGISTRY)/$(CORE_IMAGE_NAME) $(MAKE) release-manifests
	$(MAKE) release-policies
	$(GORELEASER) release --config $(GORELEASER_CONFIG) --release-notes $(RELEASE_DIR)/CHANGELOG.md --clean

release-policies: $(RELEASE_POLICIES) ## Release policies

$(RELEASE_DIR)/AWSIAMManagedPolicyControllers.json: $(RELEASE_DIR) $(CLUSTERAWSADM_SRCS)
	go run ./cmd/clusterawsadm bootstrap iam print-policy --document AWSIAMManagedPolicyControllers > $@

$(RELEASE_DIR)/AWSIAMManagedPolicyControllersWithEKS.json: $(RELEASE_DIR) $(CLUSTERAWSADM_SRCS)
	go run ./cmd/clusterawsadm bootstrap iam print-policy --document AWSIAMManagedPolicyControllers --config hack/eks-clusterawsadm-config.yaml > $@

$(RELEASE_DIR)/AWSIAMManagedPolicyCloudProviderControlPlane.json: $(RELEASE_DIR) $(CLUSTERAWSADM_SRCS)
	go run ./cmd/clusterawsadm bootstrap iam print-policy --document AWSIAMManagedPolicyCloudProviderControlPlane > $(RELEASE_DIR)/AWSIAMManagedPolicyCloudProviderControlPlane.json

$(RELEASE_DIR)/AWSIAMManagedPolicyCloudProviderNodes.json: $(RELEASE_DIR) $(CLUSTERAWSADM_SRCS)
	go run ./cmd/clusterawsadm bootstrap iam print-policy --document AWSIAMManagedPolicyCloudProviderNodes > $(RELEASE_DIR)/AWSIAMManagedPolicyCloudProviderNodes.json

$(RELEASE_DIR)/AWSIAMManagedPolicyControllersWithS3.json: $(RELEASE_DIR) $(CLUSTERAWSADM_SRCS)
	go run ./cmd/clusterawsadm bootstrap iam print-policy --document AWSIAMManagedPolicyControllers --config hack/s3-clusterawsadm-config.yaml > $@

.PHONY: release-manifests
release-manifests: ## Release manifest files
	$(MAKE) $(RELEASE_DIR)/$(CORE_MANIFEST_FILE).yaml TAG=$(RELEASE_TAG) PULL_POLICY=IfNotPresent
	# Add metadata to the release artifacts
	cp metadata.yaml $(RELEASE_DIR)/metadata.yaml

.PHONY: release-changelog
release-changelog: $(RELEASE_NOTES) check-release-tag check-previous-release-tag check-github-token $(RELEASE_DIR)
	$(RELEASE_NOTES) --debug --org $(GH_ORG_NAME) --repo $(GH_REPO_NAME) --start-sha $(shell git rev-list -n 1 ${PREVIOUS_VERSION}) --end-sha $(shell git rev-list -n 1 ${RELEASE_TAG}) --output $(RELEASE_DIR)/CHANGELOG.md --go-template go-template:$(REPO_ROOT)/hack/changelog.tpl --dependencies=false --branch=${RELEASE_BRANCH} --required-author=""

.PHONY: promote-images
promote-images: $(KPROMO) $(YQ)
	$(KPROMO) pr --project cluster-api-aws --tag $(RELEASE_TAG) --reviewers "$(shell ./hack/get-project-maintainers.sh ${YQ})" --fork $(USER_FORK) --image cluster-api-aws-controller

.PHONY: release-binaries
release-binaries: $(GORELEASER) ## Builds only the binaries, not a release.
	$(GORELEASER) build --config $(GORELEASER_CONFIG) --snapshot --clean

.PHONY: release-staging
release-staging: ## Builds and push container images and manifests to the staging bucket.
	$(MAKE) docker-build-all
	$(MAKE) docker-push-all
	$(MAKE) release-alias-tag
	$(MAKE) staging-manifests
	$(MAKE) upload-staging-artifacts

.PHONY: release-staging-nightly
release-staging-nightly: ## Tags and push container images to the staging bucket.
	$(eval NEW_RELEASE_ALIAS_TAG := nightly_$(RELEASE_ALIAS_TAG)_$(shell date +'%Y%m%d'))
	echo $(NEW_RELEASE_ALIAS_TAG)
	$(MAKE) release-alias-tag TAG=$(RELEASE_ALIAS_TAG) RELEASE_ALIAS_TAG=$(NEW_RELEASE_ALIAS_TAG)
	$(MAKE) staging-manifests RELEASE_ALIAS_TAG=$(NEW_RELEASE_ALIAS_TAG)
	$(MAKE) upload-staging-artifacts RELEASE_ALIAS_TAG=$(NEW_RELEASE_ALIAS_TAG)

.PHONY: release-alias-tag
release-alias-tag: # Adds the tag to the last build tag.
	gcloud container images add-tag -q $(CORE_CONTROLLER_IMG):$(TAG) $(CORE_CONTROLLER_IMG):$(RELEASE_ALIAS_TAG)

.PHONY: upload-staging-artifacts
upload-staging-artifacts: ## Upload release artifacts to the staging bucket
	# Example manifest location: https://storage.googleapis.com/k8s-staging-cluster-api-aws/components/nightly_main_20240425/infrastructure-components.yaml
	# Please note that these files are deleted after a certain period, at the time of this writing 60 days after file creation.
	gsutil cp $(RELEASE_DIR)/* gs://$(BUCKET)/components/$(RELEASE_ALIAS_TAG)

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

.PHONY: staging-manifests
staging-manifests: ## Staging manifests
	$(MAKE) $(RELEASE_DIR)/$(CORE_MANIFEST_FILE).yaml PULL_POLICY=IfNotPresent TAG=$(RELEASE_ALIAS_TAG)

##@ clean:

.PHONY: clean
clean: ## Remove all generated files
	$(MAKE) -C hack/tools clean
	$(MAKE) clean-bin
	$(MAKE) clean-temporary

.PHONY: clean-artifacts
clean-artifacts: ## Remove the _artifacts directory
	rm -rf _artifacts

.PHONY: clean-bin
clean-bin: ## Remove all generated binaries
	rm -rf bin

.PHONY: clean-release
clean-release: ## Remove the release folder
	rm -rf $(RELEASE_DIR)

.PHONY: clean-temporary
clean-temporary: ## Remove all temporary files and folders
	rm -f minikube.kubeconfig
	rm -f kubeconfig
	rm -rf _artifacts
	rm -rf test/e2e/.artifacts/*
	rm -rf test/e2e/*.xml
	rm -rf test/e2e/capa-controller-manager
	rm -rf test/e2e/capi-controller-manager
	rm -rf test/e2e/capi-kubeadm-bootstrap-controller-manager
	rm -rf test/e2e/capi-kubeadm-control-plane-controller-manager
	rm -rf test/e2e/logs
	rm -rf test/e2e/resources

##@ helpers:

go-version: ## Print the go version we use to compile our binaries and images
	@echo $(GO_VERSION)
