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

# If you update this file, please follow
# https://www.thapaliya.com/en/writings/well-documented-makefiles/

# Ensure Make is run with bash shell as some syntax below is bash-specific
SHELL:=/usr/bin/env bash

.DEFAULT_GOAL:=help

#
# Go.
#
GO_VERSION ?= 1.21.9
GO_CONTAINER_IMAGE ?= docker.io/library/golang:$(GO_VERSION)

# Use GOPROXY environment variable if set
GOPROXY := $(shell go env GOPROXY)
ifeq (,$(strip $(GOPROXY)))
GOPROXY := https://proxy.golang.org
endif
export GOPROXY

# Active module mode, as we use go modules to manage dependencies
export GO111MODULE := on

#
# Kubebuilder.
#
export KUBEBUILDER_ENVTEST_KUBERNETES_VERSION ?= 1.28.0
export KUBEBUILDER_CONTROLPLANE_START_TIMEOUT ?= 60s
export KUBEBUILDER_CONTROLPLANE_STOP_TIMEOUT ?= 60s

# This option is for running docker manifest command
export DOCKER_CLI_EXPERIMENTAL := enabled

# Enables shell script tracing. Enable by running: TRACE=1 make <target>
TRACE ?= 0

#
# Directories.
#
# Full directory of where the Makefile resides
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN_DIR := bin
BUILD_DIR := .build
TEST_DIR := test
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/$(BIN_DIR))
FLAVOR_DIR := $(ROOT_DIR)/templates
GO_INSTALL := ./hack/go-install.sh
GO_TOOLS_BUILD := ./hack/go-tools-build.sh

export PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)

# Set --output-base for conversion-gen if we are not within GOPATH
ifneq ($(abspath $(ROOT_DIR)),$(shell go env GOPATH)/src/sigs.k8s.io/cluster-api-provider-vsphere)
	CONVERSION_GEN_OUTPUT_BASE := --output-base=$(ROOT_DIR)
else
	export GOPATH := $(shell go env GOPATH)
endif

#
# Ginkgo configuration.
#
GINKGO_FOCUS ?=
GINKGO_SKIP ?=
GINKGO_TIMEOUT ?= 3h
E2E_CONF_FILE ?= $(abspath test/e2e/config/vsphere-dev.yaml)
INTEGRATION_CONF_FILE ?= $(abspath test/integration/integration-dev.yaml)
E2E_TEMPLATE_DIR := $(abspath test/e2e/data/infrastructure-vsphere/)
SKIP_RESOURCE_CLEANUP ?= false
USE_EXISTING_CLUSTER ?= false
GINKGO_NOCOLOR ?= false

# to set multiple ginkgo skip flags, if any
ifneq ($(strip $(GINKGO_SKIP)),)
_SKIP_ARGS := $(foreach arg,$(strip $(GINKGO_SKIP)),-skip="$(arg)")
endif

# Helper function to get dependency version from go.mod
get_go_version = $(shell go list -m $1 | awk '{print $$NF}')
get_test_go_version = $(shell cd test; go list -m $1 | awk '{print $$NF}')

#
# Binaries.
#
# Note: Need to use abspath so we can invoke these from subdirectories
KUSTOMIZE_VER := v4.5.2
KUSTOMIZE_BIN := kustomize
KUSTOMIZE := $(abspath $(TOOLS_BIN_DIR)/$(KUSTOMIZE_BIN)-$(KUSTOMIZE_VER))
KUSTOMIZE_PKG := sigs.k8s.io/kustomize/kustomize/v4

SETUP_ENVTEST_VER := 116a1b831fffe7ccc3c8145306c3e1a3b1b14ffa # Note: this matches the commit ID of the dependent controller-runtime module.
SETUP_ENVTEST_BIN := setup-envtest
SETUP_ENVTEST := $(abspath $(TOOLS_BIN_DIR)/$(SETUP_ENVTEST_BIN)-$(SETUP_ENVTEST_VER))
SETUP_ENVTEST_PKG := sigs.k8s.io/controller-runtime/tools/setup-envtest

CONTROLLER_GEN_VER := v0.13.0
CONTROLLER_GEN_BIN := controller-gen
CONTROLLER_GEN := $(abspath $(TOOLS_BIN_DIR)/$(CONTROLLER_GEN_BIN)-$(CONTROLLER_GEN_VER))
CONTROLLER_GEN_PKG := sigs.k8s.io/controller-tools/cmd/controller-gen

GOTESTSUM_VER := v1.11.0
GOTESTSUM_BIN := gotestsum
GOTESTSUM := $(abspath $(TOOLS_BIN_DIR)/$(GOTESTSUM_BIN)-$(GOTESTSUM_VER))
GOTESTSUM_PKG := gotest.tools/gotestsum

CONVERSION_GEN_VER := v0.28.0
CONVERSION_GEN_BIN := conversion-gen
# We are intentionally using the binary without version suffix, to avoid the version
# in generated files.
CONVERSION_GEN := $(abspath $(TOOLS_BIN_DIR)/$(CONVERSION_GEN_BIN))
CONVERSION_GEN_PKG := k8s.io/code-generator/cmd/conversion-gen

ENVSUBST_BIN := envsubst
ENVSUBST_VER := $(call get_go_version,github.com/drone/envsubst/v2)
ENVSUBST := $(abspath $(TOOLS_BIN_DIR)/$(ENVSUBST_BIN)-$(ENVSUBST_VER))
ENVSUBST_PKG := github.com/drone/envsubst/v2/cmd/envsubst

GO_APIDIFF_VER := v0.8.2
GO_APIDIFF_BIN := go-apidiff
GO_APIDIFF := $(abspath $(TOOLS_BIN_DIR)/$(GO_APIDIFF_BIN)-$(GO_APIDIFF_VER))
GO_APIDIFF_PKG := github.com/joelanford/go-apidiff

SHELLCHECK_VER := v0.9.0

TRIVY_VER := 0.47.0

KPROMO_VER := v4.0.4
KPROMO_BIN := kpromo
KPROMO :=  $(abspath $(TOOLS_BIN_DIR)/$(KPROMO_BIN)-$(KPROMO_VER))
# KPROMO_PKG may have to be changed if KPROMO_VER increases its major version.
KPROMO_PKG := sigs.k8s.io/promo-tools/v4/cmd/kpromo

YQ_VER := v4.35.2
YQ_BIN := yq
YQ :=  $(abspath $(TOOLS_BIN_DIR)/$(YQ_BIN)-$(YQ_VER))
YQ_PKG := github.com/mikefarah/yq/v4

GINKGO_BIN := ginkgo
GINGKO_VER := $(call get_go_version,github.com/onsi/ginkgo/v2)
GINKGO := $(abspath $(TOOLS_BIN_DIR)/$(GINKGO_BIN)-$(GINGKO_VER))
GINKGO_PKG := github.com/onsi/ginkgo/v2/ginkgo

GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT_VER := $(shell cat .github/workflows/pr-golangci-lint.yaml | grep [[:space:]]version: | sed 's/.*version: //')
GOLANGCI_LINT := $(abspath $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_BIN)-$(GOLANGCI_LINT_VER))
GOLANGCI_LINT_PKG := github.com/golangci/golangci-lint/cmd/golangci-lint

GOVULNCHECK_BIN := govulncheck
GOVULNCHECK_VER := v1.0.1
GOVULNCHECK := $(abspath $(TOOLS_BIN_DIR)/$(GOVULNCHECK_BIN)-$(GOVULNCHECK_VER))
GOVULNCHECK_PKG := golang.org/x/vuln/cmd/govulncheck

GOVC_VER := $(shell cat go.mod | grep "github.com/vmware/govmomi" | awk '{print $$NF}')
GOVC_BIN := govc
GOVC := $(abspath $(TOOLS_BIN_DIR)/$(GOVC_BIN)-$(GOVC_VER))
GOVC_PKG := github.com/vmware/govmomi/govc

KIND_VER := $(call get_test_go_version,sigs.k8s.io/kind)
KIND_BIN := kind
KIND := $(abspath $(TOOLS_BIN_DIR)/$(KIND_BIN)-$(KIND_VER))
KIND_PKG := sigs.k8s.io/kind

IMPORT_BOSS_BIN := import-boss
IMPORT_BOSS_VER := v0.28.1
IMPORT_BOSS := $(abspath $(TOOLS_BIN_DIR)/$(IMPORT_BOSS_BIN))
IMPORT_BOSS_PKG := k8s.io/code-generator/cmd/import-boss

CAPI_HACK_TOOLS_VER := a150f715f5a607ef172dbe96615ffdf1d51220b3 # Note: this is the commit ID of the dependend CAPI release tag, currently v1.6.1

CONVERSION_VERIFIER_VER := $(CAPI_HACK_TOOLS_VER)
CONVERSION_VERIFIER_BIN := conversion-verifier
CONVERSION_VERIFIER := $(abspath $(TOOLS_BIN_DIR)/$(CONVERSION_VERIFIER_BIN)-$(CONVERSION_VERIFIER_VER))
CONVERSION_VERIFIER_PKG := sigs.k8s.io/cluster-api/hack/tools/conversion-verifier

RELEASE_NOTES_VER := $(CAPI_HACK_TOOLS_VER)
RELEASE_NOTES_BIN := release-notes
RELEASE_NOTES := $(abspath $(TOOLS_BIN_DIR)/$(RELEASE_NOTES_BIN)-$(RELEASE_NOTES_VER))
RELEASE_NOTES_PKG := sigs.k8s.io/cluster-api/hack/tools/release/notes

# Define Docker related variables. Releases should modify and double check these vars.
REGISTRY ?= gcr.io/$(shell gcloud config get-value project)
PROD_REGISTRY ?= registry.k8s.io/cluster-api-vsphere

STAGING_REGISTRY ?= gcr.io/k8s-staging-capi-vsphere
STAGING_BUCKET ?= artifacts.k8s-staging-capi-vsphere.appspot.com

# core
IMAGE_NAME ?= cluster-api-vsphere-controller
CONTROLLER_IMG ?= $(REGISTRY)/$(IMAGE_NAME)

# It is set by Prow GIT_TAG, a git-based tag of the form vYYYYMMDD-hash, e.g., v20210120-v0.3.10-308-gc61521971

TAG ?= dev
ARCH ?= $(shell go env GOARCH)
ALL_ARCH = amd64 arm arm64 ppc64le s390x

# Allow overriding the imagePullPolicy
PULL_POLICY ?= Always

# Hosts running SELinux need :z added to volume mounts
SELINUX_ENABLED := $(shell cat /sys/fs/selinux/enforce 2> /dev/null || echo 0)

ifeq ($(SELINUX_ENABLED),1)
  DOCKER_VOL_OPTS?=:z
endif

# Set build time variables including version details
LDFLAGS ?= $(shell hack/version.sh)

# Additional CAPV vars (everything else is ~ kept in sync with core CAPI)
# Allow overriding manifest generation destination directory
MANIFEST_ROOT ?= ./config
CRD_ROOT ?= $(MANIFEST_ROOT)/default/crd/bases
SUPERVISOR_CRD_ROOT ?= $(MANIFEST_ROOT)/supervisor/crd
VMOP_CRD_ROOT ?= $(MANIFEST_ROOT)/deployments/integration-tests/crds
WEBHOOK_ROOT ?= $(MANIFEST_ROOT)/webhook
RBAC_ROOT ?= $(MANIFEST_ROOT)/rbac
VERSION ?= $(shell cat clusterctl-settings.json | jq .config.nextVersion -r)
OVERRIDES_DIR := $(HOME)/.cluster-api/overrides/infrastructure-vsphere/$(VERSION)

help:  # Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-50s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-50s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## --------------------------------------
## Generate / Manifests
## --------------------------------------

##@ generate:

.PHONY: generate
generate: ## Run all generate targets
	$(MAKE) generate-modules generate-manifests generate-go-deepcopy generate-go-conversions generate-flavors

.PHONY: generate-manifests
generate-manifests: $(CONTROLLER_GEN) ## Generate manifests e.g. CRD, RBAC etc.
	$(MAKE) clean-generated-yaml SRC_DIRS="$(CRD_ROOT),$(SUPERVISOR_CRD_ROOT),$(VMOP_CRD_ROOT),./config/webhook/manifests.yaml"
	$(CONTROLLER_GEN) \
		paths=./apis/v1alpha3 \
		paths=./apis/v1alpha4 \
		paths=./apis/v1beta1 \
		paths=./internal/webhooks \
		crd:crdVersions=v1 \
		output:crd:dir=$(CRD_ROOT) \
		output:webhook:dir=$(WEBHOOK_ROOT) \
		webhook
	$(CONTROLLER_GEN) \
		paths=./ \
		paths=./controllers/... \
		output:rbac:dir=$(RBAC_ROOT) \
		rbac:roleName=manager-role
	$(CONTROLLER_GEN) \
		paths=./apis/vmware/v1beta1 \
		crd:crdVersions=v1 \
		output:crd:dir=$(SUPERVISOR_CRD_ROOT)
	# vm-operator crds are loaded to be used for integration tests.
	$(CONTROLLER_GEN) \
		paths=github.com/vmware-tanzu/vm-operator/api/v1alpha1/... \
		crd:crdVersions=v1 \
		output:crd:dir=$(VMOP_CRD_ROOT)

.PHONY: generate-go-deepcopy
generate-go-deepcopy: $(CONTROLLER_GEN) ## Generate deepcopy go code for core
	$(MAKE) clean-generated-deepcopy SRC_DIRS="./apis"
	$(CONTROLLER_GEN) \
		object:headerFile=./hack/boilerplate/boilerplate.generatego.txt \
		paths=./apis/...

.PHONY: generate-go-conversions
generate-go-conversions: $(CONTROLLER_GEN) $(CONVERSION_GEN) ## Runs Go related generate targets
	$(MAKE) clean-generated-conversions SRC_DIRS="./apis/v1alpha3,./apis/v1alpha4"
	$(CONVERSION_GEN) \
		--input-dirs=./apis/v1alpha3 \
		--input-dirs=./apis/v1alpha4 \
		--build-tag=ignore_autogenerated \
		--output-file-base=zz_generated.conversion $(CONVERSION_GEN_OUTPUT_BASE) \
		--go-header-file=./hack/boilerplate/boilerplate.generatego.txt

.PHONY: generate-modules
generate-modules: ## Run go mod tidy to ensure modules are up to date
	go mod tidy
	cd $(TEST_DIR); go mod tidy

.PHONY: generate-doctoc
generate-doctoc:
	TRACE=$(TRACE) ./hack/generate-doctoc.sh

.PHONY: generate-e2e-templates
generate-e2e-templates: $(KUSTOMIZE) $(addprefix generate-e2e-templates-, v1.7 v1.8 main) ## Generate test templates for all branches

.PHONY: generate-e2e-templates-main
generate-e2e-templates-main: $(KUSTOMIZE) ## Generate test templates for the main branch
	$(MAKE) e2e-flavors-main
	cp "$(RELEASE_DIR)/main/cluster-template.yaml" "$(E2E_TEMPLATE_DIR)/main/base/cluster-template.yaml"
	cp "$(RELEASE_DIR)/main/cluster-template-ignition.yaml" "$(E2E_TEMPLATE_DIR)/main/base/cluster-template-ignition.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/base" > "$(E2E_TEMPLATE_DIR)/main/cluster-template.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/hw-upgrade" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-hw-upgrade.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/storage-policy" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-storage-policy.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/conformance" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-conformance.yaml"
	# Since CAPI uses different flavor names for KCP and MD remediation using MHC
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/mhc-remediation/kcp" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-kcp-remediation.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/mhc-remediation/md" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-md-remediation.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/node-drain" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-node-drain.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/ignition" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-ignition.yaml"
	# generate clusterclass and cluster topology
	cp "$(RELEASE_DIR)/main/clusterclass-template.yaml" "$(E2E_TEMPLATE_DIR)/main/clusterclass/clusterclass-quick-start.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/clusterclass" > "$(E2E_TEMPLATE_DIR)/main/clusterclass-quick-start.yaml"
	cp "$(RELEASE_DIR)/main/cluster-template-topology.yaml" "$(E2E_TEMPLATE_DIR)/main/topology/cluster-template-topology.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/topology" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-topology.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/remote-management" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-remote-management.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/install-on-bootstrap" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-install-on-bootstrap.yaml"
	# for PCI passthrough template
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/pci" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-pci.yaml"
	# for DHCP overrides
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/dhcp-overrides" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-dhcp-overrides.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/main/ownerreferences" > "$(E2E_TEMPLATE_DIR)/main/cluster-template-ownerreferences.yaml"

.PHONY: generate-e2e-templates-v1.8
generate-e2e-templates-v1.8: $(KUSTOMIZE)
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/v1.8/clusterclass" > "$(E2E_TEMPLATE_DIR)/v1.8/clusterclass-quick-start.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/v1.8/workload" > "$(E2E_TEMPLATE_DIR)/v1.8/cluster-template-workload.yaml"

.PHONY: generate-e2e-templates-v1.7
generate-e2e-templates-v1.7: $(KUSTOMIZE)
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/v1.7/clusterclass" > "$(E2E_TEMPLATE_DIR)/v1.7/clusterclass-quick-start.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_TEMPLATE_DIR)/v1.7/workload" > "$(E2E_TEMPLATE_DIR)/v1.7/cluster-template-workload.yaml"

## --------------------------------------
## Lint / Verify
## --------------------------------------

##@ lint and verify:

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint the codebase
	$(GOLANGCI_LINT) run -v $(GOLANGCI_LINT_EXTRA_ARGS)
	cd $(TEST_DIR); $(GOLANGCI_LINT) run --path-prefix $(TEST_DIR) --config $(ROOT_DIR)/.golangci.yml -v $(GOLANGCI_LINT_EXTRA_ARGS)

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Lint the codebase and run auto-fixers if supported by the linter
	GOLANGCI_LINT_EXTRA_ARGS=--fix $(MAKE) lint

APIDIFF_OLD_COMMIT ?= $(shell git rev-parse origin/main)

.PHONY: apidiff
apidiff: $(GO_APIDIFF) ## Check for API differences
	$(GO_APIDIFF) $(APIDIFF_OLD_COMMIT) --print-compatible

ALL_VERIFY_CHECKS = licenses boilerplate shellcheck modules gen conversions doctoc flavors import-restrictions

.PHONY: verify
verify: $(addprefix verify-,$(ALL_VERIFY_CHECKS)) ## Run all verify-* targets

.PHONY: verify-modules
verify-modules: generate-modules  ## Verify go modules are up to date
	@if !(git diff --quiet HEAD -- go.sum go.mod $(TEST_DIR)/go.mod $(TEST_DIR)/go.sum); then \
		git diff; \
		echo "go module files are out of date"; exit 1; \
	fi
	@if (find . -name 'go.mod' | xargs -n1 grep -q -i 'k8s.io/client-go.*+incompatible'); then \
		find . -name "go.mod" -exec grep -i 'k8s.io/client-go.*+incompatible' {} \; -print; \
		echo "go module contains an incompatible client-go version"; exit 1; \
	fi

.PHONY: verify-gen
verify-gen: generate  ## Verify go generated files are up to date
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi

.PHONY: verify-conversions
verify-conversions: $(CONVERSION_VERIFIER)  ## Verifies expected API conversion are in place
	$(CONVERSION_VERIFIER)

.PHONY: verify-doctoc
verify-doctoc: generate-doctoc
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "doctoc is out of date, run make generate-doctoc"; exit 1; \
	fi

.PHONY: verify-boilerplate
verify-boilerplate: ## Verify boilerplate text exists in each file
	TRACE=$(TRACE) ./hack/verify-boilerplate.sh

.PHONY: verify-shellcheck
verify-shellcheck: ## Verify shell files
	TRACE=$(TRACE) ./hack/verify-shellcheck.sh $(SHELLCHECK_VER)

.PHONY: verify-container-images
verify-container-images: ## Verify container images
	TRACE=$(TRACE) ./hack/verify-container-images.sh  $(TRIVY_VER)

.PHONY: verify-licenses
verify-licenses: ## Verify licenses
	TRACE=$(TRACE) ./hack/verify-licenses.sh $(TRIVY_VER)

.PHONY: verify-govulncheck
verify-govulncheck: $(GOVULNCHECK) ## Verify code for vulnerabilities
	$(GOVULNCHECK) ./... && R1=$$? || R1=$$?; \
	$(GOVULNCHECK) -C "$(TEST_DIR)" ./... && R2=$$? || R2=$$?; \
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

.PHONY: verify-flavors
verify-flavors: $(FLAVOR_DIR) generate-flavors ## Verify generated flavors
	@if !(git diff --quiet HEAD -- $(FLAVOR_DIR)); then \
		git diff $(FLAVOR_DIR); \
		echo "flavor files in templates directory are out of date"; exit 1; \
	fi

.PHONY: verify-import-restrictions
verify-import-restrictions: $(IMPORT_BOSS) ## Verify import restrictions with import-boss
	./hack/verify-import-restrictions.sh

## --------------------------------------
## Build
## --------------------------------------

##@ build:

.PHONY: manager
manager: ## Build the vsphere manager binary into the ./bin folder
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/manager sigs.k8s.io/cluster-api-provider-vsphere

.PHONY: docker-pull-prerequisites
docker-pull-prerequisites:
	docker pull docker.io/docker/dockerfile:1.4
	docker pull $(GO_CONTAINER_IMAGE)
	docker pull gcr.io/distroless/static:nonroot

.PHONY: docker-build-all
docker-build-all: $(addprefix docker-build-,$(ALL_ARCH)) ## Build docker images for all architectures

docker-build-%:
	$(MAKE) ARCH=$* docker-build

DOCKER_BUILD_MODIFY_MANIFESTS ?= true

.PHONY: docker-build
docker-build: docker-pull-prerequisites ## Build the docker image for vsphere controller manager
	DOCKER_BUILDKIT=1 docker build --platform linux/$(ARCH) --build-arg GOLANG_VERSION=$(GO_CONTAINER_IMAGE) --build-arg goproxy=$(GOPROXY) --build-arg ldflags="$(LDFLAGS)" . -t $(CONTROLLER_IMG)-$(ARCH):$(TAG)
	@if [ "${DOCKER_BUILD_MODIFY_MANIFESTS}" = "true" ]; then \
  		$(MAKE) set-manifest-image MANIFEST_IMG=$(CONTROLLER_IMG)-$(ARCH) MANIFEST_TAG=$(TAG) TARGET_RESOURCE="./config/base/manager_image_patch.yaml"; \
		$(MAKE) set-manifest-pull-policy TARGET_RESOURCE="./config/base/manager_pull_policy.yaml"; \
    fi

## --------------------------------------
## Testing
## --------------------------------------

##@ test:

ARTIFACTS ?= ${ROOT_DIR}/_artifacts

KUBEBUILDER_ASSETS ?= $(shell $(SETUP_ENVTEST) use --use-env -p path $(KUBEBUILDER_ENVTEST_KUBERNETES_VERSION))

.PHONY: setup-envtest
setup-envtest: $(SETUP_ENVTEST) ## Set up envtest (download kubebuilder assets)
	@echo KUBEBUILDER_ASSETS=$(KUBEBUILDER_ASSETS)

.PHONY: test
test: $(SETUP_ENVTEST) $(GOVC) ## Run unit tests
	KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" GOVC_BIN_PATH=$(GOVC) go test -v ./apis/... ./controllers/... ./pkg/... ./internal/... $(TEST_ARGS)

.PHONY: test-verbose
test-verbose: ## Run unit tests with verbose flag
	$(MAKE) test TEST_ARGS="$(TEST_ARGS) -v"

.PHONY: test-junit
test-junit: $(SETUP_ENVTEST) $(GOTESTSUM) $(GOVC) ## Run unit tests
	# Note: ARTIFACTS must be set so the ginkgo suites write junit reports to the ARTIFACTS folder
	set +o errexit; (ARTIFACTS=$(ARTIFACTS) KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" GOVC_BIN_PATH=$(GOVC) go test -json ./apis/... ./controllers/... ./pkg/... ./internal/... $(TEST_ARGS); echo $$? > $(ARTIFACTS)/junit.exitcode) | tee $(ARTIFACTS)/junit.stdout
	$(GOTESTSUM) --junitfile $(ARTIFACTS)/junit.xml --raw-command cat $(ARTIFACTS)/junit.stdout
	exit $$(cat $(ARTIFACTS)/junit.exitcode)

.PHONY: test-cover
test-cover: ## Run unit tests and generate a coverage report
	$(MAKE) test TEST_ARGS="$(TEST_ARGS) -coverpkg=./... -coverprofile=coverage.out"
	./hack/codecov-ignore.sh
	go tool cover -func=coverage.out -o coverage.txt
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test-integration
test-integration: e2e-image ## Run integration tests
test-integration: $(GINKGO) $(KUSTOMIZE) $(KIND)
	time $(GINKGO) --output-dir="$(ARTIFACTS)" --junit-report="junit.integration_suite.1.xml" -v ./test/integration -- --config=$(INTEGRATION_CONF_FILE) --artifacts-folder="$(ARTIFACTS)"

.PHONY: e2e-image
e2e-image: ## Build the e2e manager image
	docker buildx build --platform linux/$(ARCH) --output=type=docker \
		--build-arg ldflags="$(LDFLAGS)" --tag="gcr.io/k8s-staging-cluster-api/capv-manager:e2e" .

.PHONY: e2e
e2e: e2e-image generate-e2e-templates
e2e: $(GINKGO) $(KUSTOMIZE) $(KIND) $(GOVC) ## Run e2e tests
	@echo PATH="$(PATH)"
	@echo
	@echo Contents of $(TOOLS_BIN_DIR):
	@ls $(TOOLS_BIN_DIR)
	@echo
	time $(GINKGO) -v --trace -focus="$(GINKGO_FOCUS)" $(_SKIP_ARGS) -timeout=$(GINKGO_TIMEOUT) \
		--output-dir="$(ARTIFACTS)" --junit-report="junit.e2e_suite.1.xml" ./test/e2e -- \
		--e2e.config="$(E2E_CONF_FILE)" \
		--e2e.artifacts-folder="$(ARTIFACTS)" \
		--e2e.skip-resource-cleanup=$(SKIP_RESOURCE_CLEANUP) \
		--e2e.use-existing-cluster="$(USE_EXISTING_CLUSTER)"

## --------------------------------------
## Release
## --------------------------------------

##@ release:

## latest git tag for the commit, e.g., v0.3.10
RELEASE_TAG ?= $(shell git describe --abbrev=0 2>/dev/null)
ifneq (,$(findstring -,$(RELEASE_TAG)))
    PRE_RELEASE=true
endif
# the previous release tag, e.g., v0.3.9, excluding pre-release tags
PREVIOUS_TAG ?= $(shell git tag -l | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+$$" | sort -V | grep -B1 $(RELEASE_TAG) | head -n 1 2>/dev/null)
## set by Prow, ref name of the base branch, e.g., main
RELEASE_ALIAS_TAG := $(PULL_BASE_REF)
RELEASE_DIR := out
RELEASE_NOTES_DIR := _releasenotes
USER_FORK ?= $(shell git config --get remote.origin.url | cut -d/ -f4) # only works on https://github.com/<username>/cluster-api.git style URLs
ifeq ($(USER_FORK),)
USER_FORK := $(shell git config --get remote.origin.url | cut -d: -f2 | cut -d/ -f1) # for git@github.com:<username>/cluster-api.git style URLs
endif
IMAGE_REVIEWERS ?= $(shell ./hack/get-project-maintainers.sh)

.PHONY: $(RELEASE_DIR)
$(RELEASE_DIR):
	mkdir -p $(RELEASE_DIR)/

.PHONY: $(RELEASE_NOTES_DIR)
$(RELEASE_NOTES_DIR):
	mkdir -p $(RELEASE_NOTES_DIR)/

.PHONY: $(BUILD_DIR)
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

.PHONY: $(OVERRIDES_DIR)
$(OVERRIDES_DIR):
	@mkdir -p $(OVERRIDES_DIR)

.PHONY: release
release: clean-release ## Builds release manifests based on $(PROD_REGISTRY) and $(RELEASE_TAG) into $(RELEASE_DIR)
	@if [ -z "${RELEASE_TAG}" ]; then echo "RELEASE_TAG is not set"; exit 1; fi
	@if ! [ -z "$$(git status --porcelain)" ]; then echo "Your local git repository contains uncommitted changes, use git clean before proceeding."; exit 1; fi
	git checkout "${RELEASE_TAG}"
	# Set the manifest images to the staging/production bucket and Builds the manifests to publish with a release.
	$(MAKE) release-manifests-all

.PHONY: release-manifests-all
release-manifests-all: ## Builds release manifests into $(RELEASE_DIR)
	# Set the manifest image to $(PROD_REGISTRY)/$(IMAGE_NAME):$(RELEASE_TAG) and pull policy to IfNotPresent.
	$(MAKE) manifest-modification REGISTRY=$(PROD_REGISTRY) RELEASE_TAG=$(RELEASE_TAG) PULL_POLICY=IfNotPresent
	## Build the manifests into $(RELEASE_DIR)
	$(MAKE) release-manifests STAGE=release MANIFEST_DIR=$(RELEASE_DIR)

.PHONY: release-staging-nightly
release-staging-nightly: ## Re-tags container images to a nightly tag. Builds and pushes nightly manifests to the staging bucket.
	$(eval NEW_RELEASE_ALIAS_TAG := nightly_$(RELEASE_ALIAS_TAG)_$(shell date +'%Y%m%d'))
	echo $(NEW_RELEASE_ALIAS_TAG)
	$(MAKE) release-alias-tag TAG=$(RELEASE_ALIAS_TAG) RELEASE_ALIAS_TAG=$(NEW_RELEASE_ALIAS_TAG)
	# Set the manifest image to $(STAGING_REGISTRY)/$(IMAGE_NAME):$(NEW_RELEASE_ALIAS_TAG) and pull policy to IfNotPresent.
	$(MAKE) manifest-modification REGISTRY=$(STAGING_REGISTRY) RELEASE_TAG=$(NEW_RELEASE_ALIAS_TAG) PULL_POLICY=IfNotPresent
	## Build the manifests into $(RELEASE_DIR)
	$(MAKE) release-manifests STAGE=release MANIFEST_DIR=$(RELEASE_DIR)
	# Example manifest location: artifacts.k8s-staging-cluster-api-vsphere.appspot.com/components/nightly_main_20210121/*
	gsutil cp $(RELEASE_DIR)/* gs://$(STAGING_BUCKET)/components/$(NEW_RELEASE_ALIAS_TAG)

.PHONY: dev-manifests
dev-manifests: ## Builds dev manifests based on $(REGISTRY) and $(TAG) into the $(OVERRIDES_DIR)
	# Set the manifest image to $(REGISTRY)/$(IMAGE_NAME):$(TAG) and pull policy to Always.
	$(MAKE) manifest-modification REGISTRY=$(REGISTRY) RELEASE_TAG=$(TAG) PULL_POLICY=Always
	## Build the manifests into $(OVERRIDES_DIR)
	$(MAKE) release-manifests STAGE=dev MANIFEST_DIR=$(OVERRIDES_DIR)

.PHONY: manifest-modification
manifest-modification: $(BUILD_DIR) # Set the manifest images to $(REGISTRY)/$(IMAGE_NAME):$(RELEASE_TAG) and pull policy to $(PULL_POLICY)
	rm -rf $(BUILD_DIR)/config
	cp -R config $(BUILD_DIR)
	$(MAKE) set-manifest-image \
		MANIFEST_IMG=$(REGISTRY)/$(IMAGE_NAME) MANIFEST_TAG=$(RELEASE_TAG) \
		TARGET_RESOURCE="$(BUILD_DIR)/config/base/manager_image_patch.yaml"
	$(MAKE) set-manifest-pull-policy PULL_POLICY=$(PULL_POLICY) TARGET_RESOURCE="$(BUILD_DIR)/config/base/manager_pull_policy.yaml"

.PHONY: release-manifests
release-manifests: $(BUILD_DIR) $(MANIFEST_DIR) $(KUSTOMIZE) $(STAGE)-flavors ## Build the manifests to publish with a release
	$(KUSTOMIZE) build $(BUILD_DIR)/config/default > $(MANIFEST_DIR)/infrastructure-components.yaml
	$(KUSTOMIZE) build $(BUILD_DIR)/config/supervisor > $(MANIFEST_DIR)/infrastructure-components-supervisor.yaml

	# Add metadata to the release artifacts
	cp metadata.yaml $(MANIFEST_DIR)/metadata.yaml

.PHONY: release-flavors ## Create release flavor manifests
release-flavors: $(RELEASE_DIR)
	$(MAKE) generate-flavors FLAVOR_DIR=$(RELEASE_DIR)

.PHONY: dev-flavors ## Create dev flavor manifests
dev-flavors: $(OVERRIDES_DIR)
	$(MAKE) generate-flavors FLAVOR_DIR=$(OVERRIDES_DIR)

.PHONY: e2e-flavors ## Create dev flavor manifests for e2e testing
e2e-flavors: $(KUSTOMIZE) $(addprefix e2e-flavors-, main)

.PHONY: e2e-flavors-main ## Create dev flavor manifests for e2e testing for main branch
e2e-flavors-main: $(RELEASE_DIR)
	mkdir -p $(RELEASE_DIR)/main
	$(MAKE) generate-flavors FLAVOR_DIR=$(RELEASE_DIR)/main


.PHONY: generate-flavors
generate-flavors: $(FLAVOR_DIR)
	go run ./packaging/flavorgen --output-dir $(FLAVOR_DIR)

.PHONY: release-staging
release-staging: ## Build and push container images to the staging registry
	REGISTRY=$(STAGING_REGISTRY) $(MAKE) docker-build-all docker-push-all release-alias-tag

.PHONY: release-alias-tag
release-alias-tag: ## Add the release alias tag to the last build tag
	gcloud container images add-tag $(CONTROLLER_IMG):$(TAG) $(CONTROLLER_IMG):$(RELEASE_ALIAS_TAG)

.PHONY: generate-release-notes
generate-release-notes: $(RELEASE_NOTES_DIR) $(RELEASE_NOTES)
	# Reset the file
	echo -n > $(RELEASE_NOTES_DIR)/$(RELEASE_TAG).md
	if [ -n "${PRE_RELEASE}" ]; then \
	echo -e ":rotating_light: This is a RELEASE CANDIDATE. Use it only for testing purposes. If you find any bugs, file an [issue](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/issues/new/choose).\n" >> $(RELEASE_NOTES_DIR)/$(RELEASE_TAG).md; \
	fi
	"$(RELEASE_NOTES)" --from=$(PREVIOUS_TAG) --prefix-area-label=false --add-kubernetes-version-support=false >> $(RELEASE_NOTES_DIR)/$(RELEASE_TAG).md

.PHONY: promote-images
promote-images: $(KPROMO)
	$(KPROMO) pr --project capi-vsphere --tag $(RELEASE_TAG) --reviewers "$(IMAGE_REVIEWERS)" --fork $(USER_FORK) --image cluster-api-vsphere-controller

## --------------------------------------
## Docker
## --------------------------------------

.PHONY: docker-push-all
docker-push-all: $(addprefix docker-push-,$(ALL_ARCH))  ## Push the docker images to be included in the release for all architectures + related multiarch manifests
	$(MAKE) docker-push-manifest

docker-push-%:
	$(MAKE) ARCH=$* docker-push

.PHONY: docker-push
docker-push: ## Push the docker images to be included in the release
	docker push $(CONTROLLER_IMG)-$(ARCH):$(TAG)

.PHONY: docker-push-manifest
docker-push-manifest: ## Push the multiarch manifest for the vsphere docker images
	docker manifest create --amend $(CONTROLLER_IMG):$(TAG) $(shell echo $(ALL_ARCH) | sed -e "s~[^ ]*~$(CONTROLLER_IMG)\-&:$(TAG)~g")
	@for arch in $(ALL_ARCH); do docker manifest annotate --arch $${arch} ${CONTROLLER_IMG}:${TAG} ${CONTROLLER_IMG}-$${arch}:${TAG}; done
	docker manifest push --purge $(CONTROLLER_IMG):$(TAG)
	$(MAKE) set-manifest-image MANIFEST_IMG=$(CONTROLLER_IMG) MANIFEST_TAG=$(TAG) TARGET_RESOURCE="./config/base/manager_image_patch.yaml"
	$(MAKE) set-manifest-pull-policy TARGET_RESOURCE="./config/base/manager_pull_policy.yaml"

.PHONY: set-manifest-pull-policy
set-manifest-pull-policy:
	$(info Updating kustomize pull policy file for manager resources)
	sed -i'' -e 's@imagePullPolicy: .*@imagePullPolicy: '"$(PULL_POLICY)"'@' $(TARGET_RESOURCE)

.PHONY: set-manifest-image
set-manifest-image:
	$(info Updating kustomize image patch file for manager resource)
	sed -i'' -e 's@image: .*@image: '"${MANIFEST_IMG}:$(MANIFEST_TAG)"'@' $(TARGET_RESOURCE)

## --------------------------------------
## Cleanup
## --------------------------------------

##@ clean:

.PHONY: clean
clean: ## Remove generated binaries, GitBook files, Helm charts, and Tilt build files
	$(MAKE) clean-bin
	$(MAKE) clean-temporary
	$(MAKE) clean-release
	$(MAKE) clean-examples
	$(MAKE) clean-build

.PHONY: clean-build
clean-build:
	rm -rf $(BUILD_DIR)

.PHONY: clean-bin
clean-bin: ## Remove all generated binaries
	rm -rf $(BIN_DIR)
	rm -rf $(TOOLS_BIN_DIR)

.PHONY: clean-temporary
clean-temporary: ## Remove all temporary files and folders
	rm -f minikube.kubeconfig
	rm -f kubeconfig

.PHONY: clean-release
clean-release: ## Remove the release folder
	rm -rf $(RELEASE_DIR)

.PHONY: clean-examples
clean-examples: ## Remove all the temporary files generated in the examples folder
	rm -rf examples/_out/
	rm -f examples/provider-components/provider-components-*.yaml

.PHONY: clean-release-git
clean-release-git: ## Restores the git files usually modified during a release
	git restore ./*manager_image_patch.yaml ./*manager_pull_policy.yaml

.PHONY: clean-generated-yaml
clean-generated-yaml: ## Remove files generated by conversion-gen from the mentioned dirs. Example SRC_DIRS="./api/v1beta1"
	(IFS=','; for i in $(SRC_DIRS); do find $$i -type f -name '*.yaml' -exec rm -f {} \;; done)

.PHONY: clean-generated-deepcopy
clean-generated-deepcopy: ## Remove files generated by conversion-gen from the mentioned dirs. Example SRC_DIRS="./api/v1beta1"
	(IFS=','; for i in $(SRC_DIRS); do find $$i -type f -name 'zz_generated.deepcopy*' -exec rm -f {} \;; done)

.PHONY: clean-generated-conversions
clean-generated-conversions: ## Remove files generated by conversion-gen from the mentioned dirs. Example SRC_DIRS="./api/v1beta1"
	(IFS=','; for i in $(SRC_DIRS); do find $$i -type f -name 'zz_generated.conversion*' -exec rm -f {} \;; done)

## --------------------------------------
## Hack / Tools
## --------------------------------------

##@ hack/tools:

.PHONY: $(CONTROLLER_GEN_BIN)
$(CONTROLLER_GEN_BIN): $(CONTROLLER_GEN) ## Build a local copy of controller-gen.

.PHONY: $(CONVERSION_GEN_BIN)
$(CONVERSION_GEN_BIN): $(CONVERSION_GEN) ## Build a local copy of conversion-gen.

.PHONY: $(CONVERSION_VERIFIER_BIN)
$(CONVERSION_VERIFIER_BIN): $(CONVERSION_VERIFIER) ## Build a local copy of conversion-verifier.

.PHONY: $(GOTESTSUM_BIN)
$(GOTESTSUM_BIN): $(GOTESTSUM) ## Build a local copy of gotestsum.

.PHONY: $(GO_APIDIFF_BIN)
$(GO_APIDIFF_BIN): $(GO_APIDIFF) ## Build a local copy of go-apidiff

.PHONY: $(ENVSUBST_BIN)
$(ENVSUBST_BIN): $(ENVSUBST) ## Build a local copy of envsubst.

.PHONY: $(KUSTOMIZE_BIN)
$(KUSTOMIZE_BIN): $(KUSTOMIZE) ## Build a local copy of kustomize.

.PHONY: $(SETUP_ENVTEST_BIN)
$(SETUP_ENVTEST_BIN): $(SETUP_ENVTEST) ## Build a local copy of setup-envtest.

.PHONY: $(KPROMO_BIN)
$(KPROMO_BIN): $(KPROMO) ## Build a local copy of kpromo

.PHONY: $(YQ_BIN)
$(YQ_BIN): $(YQ) ## Build a local copy of yq

.PHONY: $(GINKGO_BIN)
$(GINKGO_BIN): $(GINKGO) ## Build a local copy of ginkgo.

.PHONY: $(GOLANGCI_LINT_BIN)
$(GOLANGCI_LINT_BIN): $(GOLANGCI_LINT) ## Build a local copy of golangci-lint.

.PHONY: $(GOVULNCHECK_BIN)
$(GOVULNCHECK_BIN): $(GOVULNCHECK) ## Build a local copy of govulncheck.

.PHONY: $(GOVC_BIN)
$(GOVC_BIN): $(GOVC) ## Build a local copy of govc.

.PHONY: $(KIND_BIN)
$(KIND_BIN): $(KIND) ## Build a local copy of kind.

.PHONY: $(IMPORT_BOSS_BIN)
$(IMPORT_BOSS_BIN): $(IMPORT_BOSS)

.PHONY: $(RELEASE_NOTES_BIN)
$(RELEASE_NOTES_BIN): $(RELEASE_NOTES) ## Build a local copy of release-notes.

$(CONTROLLER_GEN): # Build controller-gen.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(CONTROLLER_GEN_PKG) $(CONTROLLER_GEN_BIN) $(CONTROLLER_GEN_VER)

## We are forcing a rebuilt of conversion-gen via PHONY so that we're always using an up-to-date version.
## We can't use a versioned name for the binary, because that would be reflected in generated files.
.PHONY: $(CONVERSION_GEN)
$(CONVERSION_GEN): # Build conversion-gen.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(CONVERSION_GEN_PKG) $(CONVERSION_GEN_BIN) $(CONVERSION_GEN_VER)

$(CONVERSION_VERIFIER): # Build conversion-verifier.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_TOOLS_BUILD) $(CONVERSION_VERIFIER_PKG) $(CONVERSION_VERIFIER_BIN) $(CONVERSION_VERIFIER_VER)

$(GOTESTSUM): # Build gotestsum from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GOTESTSUM_PKG) $(GOTESTSUM_BIN) $(GOTESTSUM_VER)

$(GO_APIDIFF): # Build go-apidiff.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GO_APIDIFF_PKG) $(GO_APIDIFF_BIN) $(GO_APIDIFF_VER)

$(ENVSUBST): # Build envsubst.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(ENVSUBST_PKG) $(ENVSUBST_BIN) $(ENVSUBST_VER)

$(KUSTOMIZE): # Build kustomize.
	CGO_ENABLED=0 GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(KUSTOMIZE_PKG) $(KUSTOMIZE_BIN) $(KUSTOMIZE_VER)

$(SETUP_ENVTEST): # Build setup-envtest.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(SETUP_ENVTEST_PKG) $(SETUP_ENVTEST_BIN) $(SETUP_ENVTEST_VER)

$(KPROMO):
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(KPROMO_PKG) $(KPROMO_BIN) $(KPROMO_VER)

$(YQ):
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(YQ_PKG) $(YQ_BIN) $(YQ_VER)

$(GINKGO): # Build ginkgo.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GINKGO_PKG) $(GINKGO_BIN) $(GINGKO_VER)

$(GOLANGCI_LINT): # Build golangci-lint.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GOLANGCI_LINT_PKG) $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_VER)

$(GOVULNCHECK): # Build govulncheck.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GOVULNCHECK_PKG) $(GOVULNCHECK_BIN) $(GOVULNCHECK_VER)

$(GOVC): # Build GOVC.
	CGO_ENABLED=0 GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GOVC_PKG) $(GOVC_BIN) $(GOVC_VER)

$(KIND): # Build kind.
	cd $(TEST_DIR); GOBIN=$(TOOLS_BIN_DIR) ../$(GO_INSTALL) $(KIND_PKG) $(KIND_BIN) $(KIND_VER)

$(IMPORT_BOSS): # Build import-boss
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(IMPORT_BOSS_PKG) $(IMPORT_BOSS_BIN) $(IMPORT_BOSS_VER)

$(RELEASE_NOTES): # Build release-notes.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_TOOLS_BUILD) $(RELEASE_NOTES_PKG) $(RELEASE_NOTES_BIN) $(RELEASE_NOTES_VER)

## --------------------------------------
## Helpers
## --------------------------------------

##@ helpers:

go-version: ## Print the go version we use to compile our binaries and images
	@echo $(GO_VERSION)
