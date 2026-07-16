# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ROOT_DIR_RELATIVE := .

include $(ROOT_DIR_RELATIVE)/common.mk

# If you update this file, please follow
# https://www.thapaliya.com/en/writings/well-documented-makefiles/

# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on
unexport GOPATH

# Enables shell script tracing. Enable by running: TRACE=1 make <target>
TRACE ?= 0

# Go
GO_VERSION ?= 1.26.4

# Ensure correct toolchain is used
GOTOOLCHAIN = go$(GO_VERSION)
export GOTOOLCHAIN

# Directories.
ARTIFACTS ?= $(REPO_ROOT)/_artifacts
TOOLS_DIR := hack/tools
BIN_DIR := bin
TOOLS_DIR_DEPS := $(TOOLS_DIR)/go.sum $(TOOLS_DIR)/go.mod $(TOOLS_DIR)/Makefile
TOOLS_BIN_DIR := $(TOOLS_DIR)/$(BIN_DIR)

REPO_ROOT = $(shell git rev-parse --show-toplevel)
GH_REPO ?= kubernetes-sigs/cluster-api-provider-openstack
TEST_E2E_DIR := test/e2e

# Files
E2E_DATA_DIR ?= $(REPO_ROOT)/test/e2e/data
E2E_CONF_PATH  ?= $(E2E_DATA_DIR)/e2e_conf.yaml
KUBETEST_CONF_PATH ?= $(abspath $(E2E_DATA_DIR)/kubetest/conformance.yaml)
KUBETEST_FAST_CONF_PATH ?= $(abspath $(E2E_DATA_DIR)/kubetest/conformance-fast.yaml)
GO_INSTALL := ./scripts/go_install.sh

# go-apidiff
GO_APIDIFF_VER := v0.8.2
GO_APIDIFF_BIN := go-apidiff
GO_APIDIFF_PKG := github.com/joelanford/go-apidiff

# golangci-lint
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT_VER = $(shell cd hack/tools && go list -m -f '{{.Version}}' github.com/golangci/golangci-lint/v2)
GOLANGCI_LINT_PKG := github.com/golangci/golangci-lint/v2/cmd/golangci-lint

# govulncheck
GOVULNCHECK_VER := v1.3.0
GOVULNCHECK_BIN := govulncheck
GOVULNCHECK_DIR := hack/tools/govulncheck
GOVULNCHECK_TMP_DIR ?= $(GOVULNCHECK_DIR)/govulncheck.tmp

TRIVY_VER := 0.69.3

# Binaries.
CONTROLLER_GEN := $(TOOLS_BIN_DIR)/controller-gen
CONVERSION_GEN := $(TOOLS_BIN_DIR)/conversion-gen
ENVSUBST := $(TOOLS_BIN_DIR)/envsubst
GINKGO := $(TOOLS_BIN_DIR)/ginkgo
GOJQ := $(TOOLS_BIN_DIR)/gojq
GOLANGCI_LINT = $(abspath $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_BIN)-$(GOLANGCI_LINT_VER))
GOLANGCI_LINT_KAL := $(abspath $(TOOLS_BIN_DIR)/golangci-lint-kube-api-linter)
GOTESTSUM := $(TOOLS_BIN_DIR)/gotestsum
KUSTOMIZE := $(TOOLS_BIN_DIR)/kustomize
MOCKGEN := $(TOOLS_BIN_DIR)/mockgen
OPENAPI_GEN := $(TOOLS_BIN_DIR)/openapi-gen
APPLYCONFIGURATION_GEN := $(TOOLS_BIN_DIR)/applyconfiguration-gen
CLIENT_GEN := $(TOOLS_BIN_DIR)/client-gen
LISTER_GEN := $(TOOLS_BIN_DIR)/lister-gen
INFORMER_GEN := $(TOOLS_BIN_DIR)/informer-gen
RELEASE_NOTES := $(TOOLS_BIN_DIR)/release-notes
SETUP_ENVTEST := $(TOOLS_BIN_DIR)/setup-envtest
GEN_CRD_API_REFERENCE_DOCS := $(TOOLS_BIN_DIR)/gen-crd-api-reference-docs
GO_APIDIFF := $(TOOLS_BIN_DIR)/$(GO_APIDIFF_BIN)-$(GO_APIDIFF_VER)
GOVULNCHECK := $(abspath $(TOOLS_BIN_DIR)/$(GOVULNCHECK_BIN))

# Kubebuilder
export KUBEBUILDER_ENVTEST_KUBERNETES_VERSION ?= 1.28.0
export KUBEBUILDER_CONTROLPLANE_START_TIMEOUT ?= 60s
export KUBEBUILDER_CONTROLPLANE_STOP_TIMEOUT ?= 60s

PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)
export PATH

# Release variables

STAGING_REGISTRY := gcr.io/k8s-staging-capi-openstack
STAGING_BUCKET ?= artifacts.k8s-staging-capi-openstack.appspot.com
BUCKET ?= $(STAGING_BUCKET)
PROD_REGISTRY ?= registry.k8s.io/capi-openstack
REGISTRY ?= $(STAGING_REGISTRY)
RELEASE_TAG ?= $(shell git describe --abbrev=0 2>/dev/null)
PULL_BASE_REF ?= $(RELEASE_TAG) # PULL_BASE_REF will be provided by Prow
RELEASE_ALIAS_TAG ?= $(PULL_BASE_REF)
RELEASE_DIR := out

TAG ?= dev
ARCH ?= amd64
ALL_ARCH ?= amd64 arm arm64 ppc64le s390x

# main controller
IMAGE_NAME ?= capi-openstack-controller
CONTROLLER_IMG ?= $(REGISTRY)/$(IMAGE_NAME)
CONTROLLER_IMG_TAG ?= $(CONTROLLER_IMG)-$(ARCH):$(TAG)
CONTROLLER_ORIGINAL_IMG := gcr.io/k8s-staging-capi-openstack/capi-openstack-controller
CONTROLLER_NAME := capo-controller-manager
MANIFEST_FILE := infrastructure-components
CONFIG_DIR := config
NAMESPACE := capo-system

# Allow overriding manifest generation destination directory
MANIFEST_ROOT ?= config
CRD_ROOT ?= $(MANIFEST_ROOT)/crd/bases
WEBHOOK_ROOT ?= $(MANIFEST_ROOT)/webhook
RBAC_ROOT ?= $(MANIFEST_ROOT)/rbac

# Allow overriding the imagePullPolicy
PULL_POLICY ?= Always

# Set build time variables including version details
LDFLAGS = $(shell source ./hack/version.sh; version::ldflags)

# Extra arguments for govulncheck, e.g. "-show verbose"
GOVULNCHECK_ARGS ?=

## --------------------------------------
##@ Testing
## --------------------------------------

# The number of ginkgo tests to run concurrently
E2E_GINKGO_PARALLEL ?= 2

E2E_ARGS ?=

E2E_GINKGO_FOCUS ?=
E2E_GINKGO_SKIP ?=

# to set multiple ginkgo skip flags, if any
ifneq ($(strip $(E2E_GINKGO_SKIP)),)
_SKIP_ARGS := $(foreach arg,$(strip $(E2E_GINKGO_SKIP)),-skip="$(arg)")
endif

$(ARTIFACTS):
	mkdir -p $@

# By default setup-envtest will write to $XDG_DATA_HOME, or $HOME/.local/share
# if that is not defined. Set KUBEBUILDER_ASSETS_DIR to override.
ifdef KUBEBUILDER_ASSETS_DIR
	setup_envtest_extra_args += --bin-dir $(KUBEBUILDER_ASSETS_DIR)
endif

.PHONY: kubebuilder_assets
kubebuilder_assets: $(SETUP_ENVTEST)
	@echo Fetching assets for $(KUBEBUILDER_ENVTEST_KUBERNETES_VERSION)
	$(eval KUBEBUILDER_ASSETS ?= $(shell $(SETUP_ENVTEST) use --use-env -p path $(setup_envtest_extra_args) $(KUBEBUILDER_ENVTEST_KUBERNETES_VERSION)))

.PHONY: test
TEST_PATHS ?= ./...
test: test-capo

.PHONY: test-capo
test-capo: $(ARTIFACTS) $(GOTESTSUM) kubebuilder_assets
	KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" $(GOTESTSUM) --junitfile $(ARTIFACTS)/junit.test.xml --junitfile-hide-empty-pkg --jsonfile $(ARTIFACTS)/test-output.log -- \
			   -v $(TEST_PATHS) $(TEST_ARGS)

E2E_TEMPLATES_DIR=test/e2e/data/infrastructure-openstack
E2E_KUSTOMIZE_DIR=test/e2e/data/kustomize
# This directory holds the templates that do not require ci-artifacts script injection.
E2E_NO_ARTIFACT_TEMPLATES_DIR=test/e2e/data/infrastructure-openstack-no-artifact

.PHONY: e2e-templates
e2e-templates: ## Generate cluster templates for e2e tests
e2e-templates: $(addprefix $(E2E_NO_ARTIFACT_TEMPLATES_DIR)/, \
		 cluster-template-md-remediation.yaml \
		 cluster-template-kcp-remediation.yaml \
		 cluster-template-multi-az.yaml \
		 cluster-template-multi-network.yaml \
		 cluster-template-without-lb.yaml \
		 cluster-template.yaml \
		 cluster-template-topology.yaml \
		 cluster-template-flatcar.yaml \
		 cluster-template-k8s-upgrade.yaml \
		 cluster-template-flatcar-sysext.yaml \
		 cluster-template-no-bastion.yaml \
		 cluster-template-health-monitor.yaml \
		 cluster-template-capi-v1beta1.yaml \
		 cluster-template-cluster-identity.yaml \
		 cluster-template-topology-autoscaler.yaml)
# Currently no templates that require CI artifacts
# $(addprefix $(E2E_TEMPLATES_DIR)/, add-templates-here.yaml) \

$(E2E_NO_ARTIFACT_TEMPLATES_DIR)/cluster-template.yaml: $(E2E_KUSTOMIZE_DIR)/with-tags $(KUSTOMIZE) FORCE
	$(KUSTOMIZE) build "$<" > "$@"

$(E2E_NO_ARTIFACT_TEMPLATES_DIR)/cluster-template-%.yaml: $(E2E_KUSTOMIZE_DIR)/% $(KUSTOMIZE) FORCE
	$(KUSTOMIZE) build "$<" > "$@"

e2e-prerequisites: $(GINKGO) e2e-templates e2e-image ## Build all artifacts required by e2e tests

# Can be run manually, e.g. via:
# export OPENSTACK_CLOUD_YAML_FILE="$(pwd)/clouds.yaml"
# E2E_GINKGO_ARGS="-stream -focus='default'" E2E_ARGS="-use-existing-cluster='true'" make test-e2e
E2E_GINKGO_ARGS ?=
.PHONY: test-e2e ## Run e2e tests using clusterctl
test-e2e: $(GINKGO) e2e-prerequisites ## Run e2e tests
	time $(GINKGO) -fail-fast -trace -timeout=3h -show-node-events -v -tags=e2e -nodes=$(E2E_GINKGO_PARALLEL) \
		--output-dir="$(ARTIFACTS)" --junit-report="junit.e2e_suite.1.xml" \
		-focus="$(E2E_GINKGO_FOCUS)" $(_SKIP_ARGS) $(E2E_GINKGO_ARGS) ./test/e2e/suites/e2e/... -- \
			-config-path="$(E2E_CONF_PATH)" -artifacts-folder="$(ARTIFACTS)" \
			-data-folder="$(E2E_DATA_DIR)" $(E2E_ARGS)

# Pre-compile tests
# This is not required, but it will make the tests start faster
.PHONY: build-e2e-tests
build-e2e-tests: $(GINKGO)
	$(GINKGO) build -tags=e2e ./test/e2e/suites/e2e/...

.PHONY: e2e-image
e2e-image: CONTROLLER_IMG_TAG = "gcr.io/k8s-staging-capi-openstack/capi-openstack-controller:e2e"
e2e-image: docker-build

CONFORMANCE_E2E_ARGS ?= -kubetest.config-file=$(KUBETEST_CONF_PATH)
CONFORMANCE_E2E_ARGS += $(E2E_ARGS)
.PHONY: test-conformance
test-conformance: $(GINKGO) e2e-prerequisites ## Run clusterctl based conformance test on workload cluster (requires Docker).
	time $(GINKGO) -trace -show-node-events -v -tags=e2e -focus="conformance" $(CONFORMANCE_GINKGO_ARGS) \
	   ./test/e2e/suites/conformance/... -- \
			-config-path="$(E2E_CONF_PATH)" -artifacts-folder="$(ARTIFACTS)" \
			--data-folder="$(E2E_DATA_DIR)" $(CONFORMANCE_E2E_ARGS)

test-conformance-fast: ## Run clusterctl based conformance test on workload cluster (requires Docker) using a subset of the conformance suite in parallel.
	$(MAKE) test-conformance CONFORMANCE_E2E_ARGS="-kubetest.config-file=$(KUBETEST_FAST_CONF_PATH) -kubetest.ginkgo-nodes=5 $(E2E_ARGS)"

APIDIFF_OLD_COMMIT ?= $(shell git rev-parse origin/main)

.PHONY: apidiff
apidiff: $(GO_APIDIFF) ## Check for API differences.
	$(GO_APIDIFF) $(APIDIFF_OLD_COMMIT)

## --------------------------------------
##@ Binaries
## --------------------------------------

.PHONY: binaries
binaries: managers ## Builds and installs all binaries

.PHONY: managers
managers:
	$(MAKE) manager-openstack-infrastructure

.PHONY: manager-openstack-infrastructure
manager-openstack-infrastructure: ## Build manager binary.
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS} -extldflags '-static'" -o $(BIN_DIR)/manager .

.PHONY: $(GO_APIDIFF_BIN)
$(GO_APIDIFF_BIN): $(GO_APIDIFF)

$(GO_APIDIFF): # Build go-apidiff.
	GOBIN=$(abspath $(TOOLS_BIN_DIR)) $(GO_INSTALL) $(GO_APIDIFF_PKG) $(GO_APIDIFF_BIN) $(GO_APIDIFF_VER)

.PHONY: $(GOVULNCHECK_BIN)
$(GOVULNCHECK_BIN): $(GOVULNCHECK) ## Build a local copy of govulncheck.

$(GOVULNCHECK): # Build govulncheck from source with exclusion patch.
	@if [ -d "$(GOVULNCHECK_TMP_DIR)" ]; then \
		echo "$(GOVULNCHECK_TMP_DIR) exists, skipping clone"; \
	else \
		git clone "https://github.com/golang/vuln.git" "$(GOVULNCHECK_TMP_DIR)"; \
		cd "$(GOVULNCHECK_TMP_DIR)"; \
		git checkout "$(GOVULNCHECK_VER)"; \
		git apply "$(REPO_ROOT)/$(GOVULNCHECK_DIR)/govulncheck.patch"; \
	fi
	@cd "$(REPO_ROOT)/$(GOVULNCHECK_TMP_DIR)"; \
	if [ "$$(git describe --tag 2> /dev/null)" != "$(GOVULNCHECK_VER)" ]; then \
		echo "ERROR: checked out version does not match expected version $(GOVULNCHECK_VER)"; \
		exit 1; \
	fi
	@rm -f $(GOVULNCHECK)
	go build -C "$(REPO_ROOT)/$(GOVULNCHECK_TMP_DIR)" -o $(GOVULNCHECK) ./cmd/govulncheck

.PHONY: clean-govulncheck
clean-govulncheck:
	rm -fr "$(GOVULNCHECK_TMP_DIR)"

.PHONY: $(GOLANGCI_LINT_BIN)
$(GOLANGCI_LINT_BIN): $(GOLANGCI_LINT) ## Build a local copy of golangci-lint.

$(GOLANGCI_LINT): # Build golangci-lint.
	GOBIN=$(abspath $(TOOLS_BIN_DIR)) $(GO_INSTALL) $(GOLANGCI_LINT_PKG) $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_VER)

$(GOLANGCI_LINT_KAL): $(GOLANGCI_LINT) $(TOOLS_DIR_DEPS) # Build golangci-lint with KAL plugin.
	cd $(TOOLS_DIR); $(GOLANGCI_LINT) custom

## --------------------------------------
##@ Linting
## --------------------------------------

.PHONY: lint
lint: $(GOLANGCI_LINT) $(GOLANGCI_LINT_KAL) ## Lint codebase
	$(GOLANGCI_LINT) run -v
	$(GOLANGCI_LINT_KAL) run -v --config $(ROOT_DIR_RELATIVE)/.golangci-kal.yml

.PHONY: lint-update
lint-update: $(GOLANGCI_LINT) $(GOLANGCI_LINT_KAL) ## Lint and fix issues
	$(GOLANGCI_LINT) run -v --fix
	$(GOLANGCI_LINT_KAL) run -v --fix --config $(ROOT_DIR_RELATIVE)/.golangci-kal.yml

lint-fast: $(GOLANGCI_LINT) ## Run only faster linters to detect possible issues
	$(GOLANGCI_LINT) run -v --fast-only

.PHONY: lint-api
lint-api: $(GOLANGCI_LINT_KAL) ## Lint API types with KAL
	$(GOLANGCI_LINT_KAL) run -v --config $(ROOT_DIR_RELATIVE)/.golangci-kal.yml

.PHONY: lint-api-fix
lint-api-fix: $(GOLANGCI_LINT_KAL) ## Lint API types with KAL and auto-fix issues
	$(GOLANGCI_LINT_KAL) run -v --fix --config $(ROOT_DIR_RELATIVE)/.golangci-kal.yml

## --------------------------------------
##@ Generate
## --------------------------------------

.PHONY: modules
modules: ## Runs go mod to ensure proper vendoring.
	go mod tidy
	cd $(TOOLS_DIR); go mod tidy

.PHONY: generate
generate: templates generate-controller-gen generate-codegen generate-conversion-gen generate-go generate-manifests generate-api-docs ## Generate all generated code

.PHONY: generate-go
generate-go: $(MOCKGEN)
	go generate ./...

.PHONY: generate-conversion-gen
generate-conversion-gen: $(CONVERSION_GEN)
	$(CONVERSION_GEN) \
		--output-file=zz_generated.conversion.go \
		--go-header-file=./hack/boilerplate/boilerplate.generatego.txt \
		./api/v1beta1

.PHONY: generate-controller-gen
generate-controller-gen: $(CONTROLLER_GEN)
	$(CONTROLLER_GEN) \
		paths=./api/... \
		object:headerFile=./hack/boilerplate/boilerplate.generatego.txt

.PHONY: generate-codegen
generate-codegen: generate-controller-gen $(OPENAPI_GEN) $(APPLYCONFIGURATION_GEN) $(CLIENT_GEN) $(LISTER_GEN) $(INFORMER_GEN)
	@echo "** Generating OpenAPI definitions **"
	# The package list includes:
	# - CAPO's own API packages (v1alpha1, v1beta1, v1beta2) that have // +k8s:openapi-gen= markers
	# - Dependency packages from CAPI and k8s.io that are referenced by CAPO's APIs
	# - Base k8s.io/apimachinery packages
	$(OPENAPI_GEN) \
		--go-header-file=./hack/boilerplate.go.txt \
		--output-file=zz_generated.openapi.go \
		--output-dir=./cmd/models-schema \
		--output-pkg=main \
		--report-filename=./api_violations.report \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1 \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1 \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2 \
		sigs.k8s.io/cluster-api/api/core/v1beta2 \
		sigs.k8s.io/cluster-api/api/ipam/v1beta2 \
		sigs.k8s.io/cluster-api/api/core/v1beta1 \
		sigs.k8s.io/cluster-api/api/ipam/v1beta1 \
		k8s.io/api/core/v1 \
		k8s.io/apimachinery/pkg/apis/meta/v1 \
		k8s.io/apimachinery/pkg/runtime \
		k8s.io/apimachinery/pkg/api/resource \
		k8s.io/apimachinery/pkg/version
	@echo "** Generating openapi.json **"
	go run ./cmd/models-schema | jq > ./openapi.json
	@echo "** Generating applyconfiguration code **"
	$(APPLYCONFIGURATION_GEN) \
		--go-header-file=./hack/boilerplate.go.txt \
		--output-dir=./pkg/generated/applyconfiguration \
		--output-pkg=sigs.k8s.io/cluster-api-provider-openstack/pkg/generated/applyconfiguration \
		--openapi-schema=./openapi.json \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1 \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1 \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2
	@echo "** Generating clientset code **"
	$(CLIENT_GEN) \
		--go-header-file=./hack/boilerplate.go.txt \
		--output-dir=./pkg/generated/clientset \
		--output-pkg=sigs.k8s.io/cluster-api-provider-openstack/pkg/generated/clientset \
		--clientset-name=clientset \
		--input-base=sigs.k8s.io/cluster-api-provider-openstack \
		--apply-configuration-package=sigs.k8s.io/cluster-api-provider-openstack/pkg/generated/applyconfiguration \
		--input=api/v1alpha1 \
		--input=api/v1beta1 \
		--input=api/v1beta2
	@echo "** Generating lister code **"
	$(LISTER_GEN) \
		--go-header-file=./hack/boilerplate.go.txt \
		--output-dir=./pkg/generated/listers \
		--output-pkg=sigs.k8s.io/cluster-api-provider-openstack/pkg/generated/listers \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1 \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1 \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2
	@echo "** Generating informer code **"
	$(INFORMER_GEN) \
		--go-header-file=./hack/boilerplate.go.txt \
		--output-dir=./pkg/generated/informers \
		--output-pkg=sigs.k8s.io/cluster-api-provider-openstack/pkg/generated/informers \
		--versioned-clientset-package=sigs.k8s.io/cluster-api-provider-openstack/pkg/generated/clientset/clientset \
		--listers-package=sigs.k8s.io/cluster-api-provider-openstack/pkg/generated/listers \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1 \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1 \
		sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2

.PHONY: generate-manifests
generate-manifests: $(CONTROLLER_GEN) ## Generate manifests e.g. CRD, RBAC etc.
	$(CONTROLLER_GEN) \
		paths=./api/... \
		crd:crdVersions=v1 \
		output:crd:dir=$(CRD_ROOT)
	$(CONTROLLER_GEN) \
		paths=./pkg/webhooks/... \
		output:webhook:dir=$(WEBHOOK_ROOT) \
		webhook
	 # We also need to extract rbac from ORC while we're running its controllers
	$(CONTROLLER_GEN) \
		paths=./ \
		paths=./controllers/... \
		output:rbac:dir=$(RBAC_ROOT) \
		rbac:roleName=manager-role

.PHONY: generate-api-docs
generate-api-docs: generate-api-docs-v1beta2 generate-api-docs-v1beta1 generate-api-docs-v1alpha1
generate-api-docs-%: $(GEN_CRD_API_REFERENCE_DOCS) FORCE
	$(GEN_CRD_API_REFERENCE_DOCS) \
		-api-dir=./api/$* \
		-config=./docs/book/gen-crd-api-reference-docs/config.json \
		-template-dir=./docs/book/gen-crd-api-reference-docs/template \
		-out-file=./docs/book/src/api/$*/api.md

## --------------------------------------
##@ Docker
## --------------------------------------

.PHONY: docker-build
docker-build: ## Build the docker image for controller-manager without debug info
	docker build -f Dockerfile --build-arg GO_VERSION=$(GO_VERSION) \
	--build-arg goproxy=$(GOPROXY) \
	--build-arg ARCH=$(ARCH) \
	--build-arg ldflags="$(LDFLAGS) -s -w" . -t $(CONTROLLER_IMG_TAG)

.PHONY: docker-build-debug
docker-build-debug: ## Build the docker image for controller-manager with debug info
	docker build -f Dockerfile --build-arg GO_VERSION=$(GO_VERSION) \
	--build-arg goproxy=$(GOPROXY) \
	--build-arg ARCH=$(ARCH) \
	--build-arg ldflags="$(LDFLAGS)" . -t $(CONTROLLER_IMG_TAG)

.PHONY: docker-push
docker-push: ## Push the docker image
	docker push $(CONTROLLER_IMG_TAG)

## --------------------------------------
##@ Docker — All ARCH
## --------------------------------------

.PHONY: docker-build-all ## Build all the architecture docker images
docker-build-all: $(addprefix docker-build-,$(ALL_ARCH))

docker-build-%:
	$(MAKE) ARCH=$* docker-build

.PHONY: docker-push-all ## Push all the architecture docker images
docker-push-all: $(addprefix docker-push-,$(ALL_ARCH))
	$(MAKE) docker-push-manifest

docker-push-%:
	$(MAKE) ARCH=$* docker-push

.PHONY: docker-push-manifest
docker-push-manifest: ## Push the fat manifest docker image.
	## Minimum docker version 18.06.0 is required for creating and pushing manifest images.
	docker manifest create --amend $(CONTROLLER_IMG):$(TAG) $(shell echo $(ALL_ARCH) | sed -e "s~[^ ]*~$(CONTROLLER_IMG)\-&:$(TAG)~g")
	@for arch in $(ALL_ARCH); do docker manifest annotate --arch $${arch} ${CONTROLLER_IMG}:${TAG} ${CONTROLLER_IMG}-$${arch}:${TAG}; done
	docker manifest push --purge ${CONTROLLER_IMG}:${TAG}

.PHONY: staging-manifests
staging-manifests:
	$(MAKE) $(RELEASE_DIR)/$(MANIFEST_FILE).yaml PULL_POLICY=IfNotPresent TAG=$(RELEASE_ALIAS_TAG)

## --------------------------------------
##@ Release
## --------------------------------------
ifneq (,$(findstring -,$(RELEASE_TAG)))
    PRE_RELEASE=true
endif
# List all tags, add the new tag to the list, sort and pick the previous one.
PREVIOUS_TAG ?= $(shell (git tag -l | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+$$"; echo "$(RELEASE_TAG)") | sort -V | grep -B1 "^$(RELEASE_TAG)$$" | grep -v "^$(RELEASE_TAG)$$" | head -n 1 2>/dev/null)
## set by Prow, ref name of the base branch, e.g., main
RELEASE_DIR := out
RELEASE_NOTES_DIR := releasenotes

.PHONY: $(RELEASE_DIR)
$(RELEASE_DIR):
	mkdir -p $(RELEASE_DIR)/

.PHONY: $(RELEASE_NOTES_DIR)
$(RELEASE_NOTES_DIR):
	mkdir -p $(RELEASE_NOTES_DIR)/

.PHONY: $(BUILD_DIR)
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

.PHONY: list-staging-releases
list-staging-releases: ## List staging images for image promotion
	@echo $(IMAGE_NAME):
	$(MAKE) list-image RELEASE_TAG=$(RELEASE_TAG) IMAGE=$(IMAGE_NAME)

list-image:
	gcloud container images list-tags $(STAGING_REGISTRY)/$(IMAGE) --filter="tags=('$(RELEASE_TAG)')" --format=json

.PHONY: release
release: $(RELEASE_NOTES) $(RELEASE_DIR)  ## Builds and push container images using the latest git tag for the commit.
	@if [ -z "${RELEASE_TAG}" ]; then echo "RELEASE_TAG is not set"; exit 1; fi
	@if ! [ -z "$$(git status --porcelain)" ]; then echo "Your local git repository contains uncommitted changes, use git clean before proceeding."; fi
	git checkout "${RELEASE_TAG}"
	# Set the manifest image to the production bucket.
	$(MAKE) manifest-modification REGISTRY=$(PROD_REGISTRY)
	$(MAKE) release-manifests
	$(MAKE) release-templates
	$(MAKE) generate-release-notes

.PHONY: manifest-modification
manifest-modification: # Set the manifest images to the staging/production bucket.
	$(MAKE) set-manifest-image \
		MANIFEST_IMG=$(REGISTRY)/$(IMAGE_NAME) MANIFEST_TAG=$(RELEASE_TAG) \
		TARGET_RESOURCE="./config/default/manager_image_patch.yaml"
	$(MAKE) set-manifest-pull-policy PULL_POLICY=IfNotPresent TARGET_RESOURCE="./config/default/manager_pull_policy.yaml"

.PHONY: set-manifest-image
set-manifest-image:
	$(info Updating kustomize image patch file for manager resource)
	sed -i'' -e 's@image: .*@image: '"${MANIFEST_IMG}:$(MANIFEST_TAG)"'@' $(TARGET_RESOURCE)

.PHONY: set-manifest-pull-policy
set-manifest-pull-policy:
	$(info Updating kustomize pull policy file for manager resources)
	sed -i'' -e 's@imagePullPolicy: .*@imagePullPolicy: '"$(PULL_POLICY)"'@' $(TARGET_RESOURCE)

.PHONY: release-manifests
release-manifests:
	$(MAKE) $(RELEASE_DIR)/$(MANIFEST_FILE).yaml TAG=$(RELEASE_TAG) PULL_POLICY=IfNotPresent
	# Add metadata to the release artifacts
	cp metadata.yaml $(RELEASE_DIR)/metadata.yaml

.PHONY: release-staging
release-staging: ## Builds and push container images and manifests to the staging bucket.
	$(MAKE) docker-build-all
	$(MAKE) docker-push-all
	$(MAKE) release-alias-tag
	$(MAKE) staging-manifests
	$(MAKE) upload-staging-artifacts

.PHONY: release-staging-nightly
release-staging-nightly: ## Tags and push container images to the staging bucket. Example image tag: capi-openstack-controller:nightly_master_20210121
	$(eval NEW_RELEASE_ALIAS_TAG := nightly_$(RELEASE_ALIAS_TAG)_$(shell date +'%Y%m%d'))
	echo $(NEW_RELEASE_ALIAS_TAG)
	$(MAKE) release-alias-tag TAG=$(RELEASE_ALIAS_TAG) RELEASE_ALIAS_TAG=$(NEW_RELEASE_ALIAS_TAG)
	$(MAKE) staging-manifests RELEASE_ALIAS_TAG=$(NEW_RELEASE_ALIAS_TAG)
	$(MAKE) upload-staging-artifacts RELEASE_ALIAS_TAG=$(NEW_RELEASE_ALIAS_TAG)

.PHONY: upload-staging-artifacts
upload-staging-artifacts: ## Upload release artifacts to the staging bucket
	gsutil cp $(RELEASE_DIR)/* gs://$(STAGING_BUCKET)/components/$(RELEASE_ALIAS_TAG)/

.PHONY: create-gh-release
create-gh-release:$(GH) ## Create release on Github
	$(GH) release create $(VERSION) -d -F $(RELEASE_DIR)/CHANGELOG.md -t $(VERSION) -R $(GH_REPO)

.PHONY: upload-gh-artifacts
upload-gh-artifacts: $(GH) ## Upload artifacts to Github release
	$(GH) release upload $(VERSION) -R $(GH_REPO) --clobber  $(RELEASE_DIR)/*

.PHONY: release-alias-tag
release-alias-tag: # Adds the tag to the last build tag.
	gcloud container images add-tag -q $(CONTROLLER_IMG):$(TAG) $(CONTROLLER_IMG):$(RELEASE_ALIAS_TAG)

.PHONY: generate-release-notes ## Generate release notes
generate-release-notes: $(RELEASE_NOTES_DIR) $(RELEASE_NOTES)
	# Reset the file
	echo -n > $(RELEASE_NOTES_DIR)/$(RELEASE_TAG).md
	if [ -n "${PRE_RELEASE}" ]; then \
	echo -e ":rotating_light: This is a RELEASE CANDIDATE. Use it only for testing purposes. If you find any bugs, file an [issue](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/issues/new/choose).\n" >> $(RELEASE_NOTES_DIR)/$(RELEASE_TAG).md; \
	fi
	"$(RELEASE_NOTES)" --repository=kubernetes-sigs/cluster-api-provider-openstack \
	  --prefix-area-label=false --add-kubernetes-version-support=false \
	  --from=tags/$(PREVIOUS_TAG) --release=$(RELEASE_TAG) >> $(RELEASE_NOTES_DIR)/$(RELEASE_TAG).md

.PHONY: templates
templates: ## Generate cluster templates
templates: templates/cluster-template.yaml \
	templates/cluster-template-without-lb.yaml \
	templates/cluster-template-flatcar.yaml \
	templates/cluster-template-flatcar-sysext.yaml \
	templates/cluster-template-capi-v1beta1.yaml

templates/cluster-template.yaml: kustomize/default $(KUSTOMIZE) FORCE
	$(KUSTOMIZE) build "$<" > "$@"

templates/cluster-template-%.yaml: kustomize/% $(KUSTOMIZE) FORCE
	$(KUSTOMIZE) build "$<" > "$@"

.PHONY: release-templates
release-templates: $(RELEASE_DIR) templates ## Generate release templates
	cp templates/cluster-template*.yaml $(RELEASE_DIR)/
	cp templates/clusterclass*.yaml $(RELEASE_DIR)/
	cp templates/image-template*.yaml $(RELEASE_DIR)/

IMAGE_PATCH_DIR := $(ARTIFACTS)/image-patch

$(IMAGE_PATCH_DIR): $(ARTIFACTS)
	mkdir -p $@

.PHONY: $(RELEASE_DIR)/$(MANIFEST_FILE).yaml
$(RELEASE_DIR)/$(MANIFEST_FILE).yaml:
	$(MAKE) compiled-manifest \
		PROVIDER=$(MANIFEST_FILE) \
		OLD_IMG=$(CONTROLLER_ORIGINAL_IMG) \
		MANIFEST_IMG=$(CONTROLLER_IMG) \
		CONTROLLER_NAME=$(CONTROLLER_NAME) \
		PROVIDER_CONFIG_DIR=$(CONFIG_DIR) \
		NAMESPACE=$(NAMESPACE)

.PHONY: compiled-manifest
compiled-manifest: $(RELEASE_DIR) $(KUSTOMIZE)
	$(MAKE) image-patch-source-manifest
	$(MAKE) image-patch-pull-policy
	$(MAKE) image-patch-kustomization
	$(KUSTOMIZE) build $(IMAGE_PATCH_DIR)/$(PROVIDER) > $(RELEASE_DIR)/$(PROVIDER).yaml

.PHONY: image-patch-source-manifest
image-patch-source-manifest: $(IMAGE_PATCH_DIR) $(KUSTOMIZE)
	mkdir -p $(IMAGE_PATCH_DIR)/$(PROVIDER)
	$(KUSTOMIZE) build $(PROVIDER_CONFIG_DIR)/default > $(IMAGE_PATCH_DIR)/$(PROVIDER)/source-manifest.yaml

.PHONY: image-patch-kustomization
image-patch-kustomization: $(IMAGE_PATCH_DIR)
	mkdir -p $(IMAGE_PATCH_DIR)/$(PROVIDER)
	$(GOJQ) --yaml-input --yaml-output '.images[0]={"name":"$(OLD_IMG)","newName":"$(MANIFEST_IMG)","newTag":"$(TAG)"}|del(.patchesJson6902[1])|.patchesJson6902[0].target.name="$(CONTROLLER_NAME)"|.patchesJson6902[0].target.namespace="$(NAMESPACE)"' \
		"hack/image-patch/kustomization.yaml" > $(IMAGE_PATCH_DIR)/$(PROVIDER)/kustomization.yaml

.PHONY: image-patch-pull-policy
image-patch-pull-policy: $(IMAGE_PATCH_DIR) $(GOJQ)
	mkdir -p $(IMAGE_PATCH_DIR)/$(PROVIDER)
	echo Setting imagePullPolicy to $(PULL_POLICY)
	$(GOJQ) --yaml-input --yaml-output '.[0].value="$(PULL_POLICY)"' "hack/image-patch/pull-policy-patch.yaml" > $(IMAGE_PATCH_DIR)/$(PROVIDER)/pull-policy-patch.yaml


## --------------------------------------
##@ Cleanup / Verification
## --------------------------------------

.PHONY: clean
clean: ## Remove all generated files
	$(MAKE) -C hack/tools clean
	$(MAKE) clean-bin
	$(MAKE) clean-temporary

.PHONY: clean-bin
clean-bin: ## Remove all generated binaries
	rm -rf bin

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

.PHONY: verify
verify: verify-boilerplate verify-modules verify-gen

.PHONY: verify-boilerplate
verify-boilerplate:
	./hack/verify-boilerplate.sh

.PHONY: verify-modules
verify-modules: modules
	@if !(git diff --quiet HEAD -- go.* hack/tools/go.*); then \
		git diff; \
		echo "go module files are out of date"; exit 1; \
	fi

.PHONY: verify-gen
verify-gen: generate
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi

.PHONY: verify-container-images
verify-container-images: ## Verify container images
	TRACE=$(TRACE) ./hack/verify-container-images.sh $(TRIVY_VER)

.PHONY: verify-govulncheck
verify-govulncheck: $(GOVULNCHECK) ## Verify code for vulnerabilities
	$(GOVULNCHECK) $(GOVULNCHECK_ARGS) ./...

.PHONY: verify-security
verify-security: ## Verify code and images for vulnerabilities
	$(MAKE) verify-container-images && R1=$$? || R1=$$?; \
	$(MAKE) verify-govulncheck && R2=$$? || R2=$$?; \
	if [ "$$R1" -ne "0" ] || [ "$$R2" -ne "0" ]; then \
	  echo "Check for vulnerabilities failed! There are vulnerabilities to be fixed"; \
		exit 1; \
	fi

.PHONY: compile-e2e
compile-e2e: ## Test e2e compilation
	go test -c -o /dev/null -tags=e2e ./test/e2e/suites/conformance

.PHONY: FORCE
FORCE:

## --------------------------------------
## Helpers
## --------------------------------------

##@ helpers:

go-version: ## Print the go version we use to compile our binaries and images
	@echo $(GO_VERSION)
