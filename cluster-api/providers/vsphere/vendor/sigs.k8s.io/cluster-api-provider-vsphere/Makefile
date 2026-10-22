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
GO_VERSION ?= 1.26.6
GO_DIRECTIVE_VERSION ?= 1.26.0
GO_CONTAINER_IMAGE ?= docker.io/library/golang:$(GO_VERSION)

# Ensure correct toolchain is used
GOTOOLCHAIN = go$(GO_VERSION)
export GOTOOLCHAIN

# Use GOPROXY environment variable if set
GOPROXY := $(shell go env GOPROXY)
ifeq ($(GOPROXY),)
GOPROXY := https://proxy.golang.org
endif
export GOPROXY

# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on

#
# Kubebuilder.
#
export KUBEBUILDER_ENVTEST_KUBERNETES_VERSION ?= 1.37.0
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
PACKAGING_DIR := packaging
VCSIM_DIR := test/infrastructure/vcsim
NETOP_DIR := test/infrastructure/net-operator
TEST_EXTENSION_DIR := test/extension
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/$(BIN_DIR))
FLAVOR_DIR := $(ROOT_DIR)/templates
GO_INSTALL := ./hack/go-install.sh
GO_TOOLS_BUILD := ./hack/go-tools-build.sh

export PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)

# DBG=1 for building binaries which includes DWARF and symbol table for delve
# debugging. When DBG is unspecified it defaults to "-s -w" which strips debug
# information.
export DBG ?= 0

# Set build time variables including version details
LDFLAGS := $(shell hack/version.sh)
GCFLAGS := $(shell hack/gogcflags.sh)

#
# Ginkgo configuration.
#
GINKGO_FOCUS ?=
GINKGO_SKIP ?=
GINKGO_LABEL_FILTER ?=
GINKGO_NODES ?= 1
GINKGO_TIMEOUT ?= 3h
GINKGO_ARGS ?=
GINKGO_POLL_PROGRESS_AFTER ?= 60m
GINKGO_POLL_PROGRESS_INTERVAL ?= 5m
E2E_CONF_FILE ?= $(abspath test/e2e/config/vsphere.yaml)
E2E_CONF_OVERRIDE_FILE ?= $(abspath test/e2e/config/config-overrides.yaml)
E2E_VSPHERE_IP_POOL ?=
E2E_TEMPLATE_DIR := $(abspath test/e2e/data/)
E2E_GOVMOMI_TEMPLATE_DIR := $(E2E_TEMPLATE_DIR)/infrastructure-vsphere-govmomi
E2E_SUPERVISOR_TEMPLATE_DIR := $(E2E_TEMPLATE_DIR)/infrastructure-vsphere-supervisor
SKIP_RESOURCE_CLEANUP ?= false
USE_EXISTING_CLUSTER ?= false
GINKGO_NOCOLOR ?= false

# to set multiple ginkgo skip flags, if any
ifneq ($(strip $(GINKGO_SKIP)),)
_SKIP_ARGS := $(foreach arg,$(strip $(GINKGO_SKIP)),-skip="$(arg)")
endif

# Helper function to get dependency version from go.mod
get_go_version = $(shell go list -m $1 | awk '{print $$NF}')
get_test_go_version = $(shell cat test/go.mod | grep $1 | awk '{print $$2}')

#
# Binaries.
#
# Note: Need to use abspath so we can invoke these from subdirectories
KUSTOMIZE_VER := v4.5.2
KUSTOMIZE_BIN := kustomize
KUSTOMIZE := $(abspath $(TOOLS_BIN_DIR)/$(KUSTOMIZE_BIN)-$(KUSTOMIZE_VER))
KUSTOMIZE_PKG := sigs.k8s.io/kustomize/kustomize/v4

SETUP_ENVTEST_VER := v0.24.1
SETUP_ENVTEST_BIN := setup-envtest
SETUP_ENVTEST := $(abspath $(TOOLS_BIN_DIR)/$(SETUP_ENVTEST_BIN)-$(SETUP_ENVTEST_VER))
SETUP_ENVTEST_PKG := sigs.k8s.io/controller-runtime/tools/setup-envtest

CONTROLLER_GEN_VER := v0.21.0
CONTROLLER_GEN_BIN := controller-gen
CONTROLLER_GEN := $(abspath $(TOOLS_BIN_DIR)/$(CONTROLLER_GEN_BIN)-$(CONTROLLER_GEN_VER))
CONTROLLER_GEN_PKG := sigs.k8s.io/controller-tools/cmd/controller-gen

GOTESTSUM_VER := v1.12.3
GOTESTSUM_BIN := gotestsum
GOTESTSUM := $(abspath $(TOOLS_BIN_DIR)/$(GOTESTSUM_BIN)-$(GOTESTSUM_VER))
GOTESTSUM_PKG := gotest.tools/gotestsum

CONVERSION_GEN_VER := v0.36.0
CONVERSION_GEN_BIN := conversion-gen
# We are intentionally using the binary without version suffix, to avoid the version
# in generated files.
CONVERSION_GEN := $(abspath $(TOOLS_BIN_DIR)/$(CONVERSION_GEN_BIN))
CONVERSION_GEN_PKG := k8s.io/code-generator/cmd/conversion-gen

ENVSUBST_BIN := envsubst
ENVSUBST_VER := $(call get_go_version,github.com/drone/envsubst/v2)
ENVSUBST := $(abspath $(TOOLS_BIN_DIR)/$(ENVSUBST_BIN)-$(ENVSUBST_VER))
ENVSUBST_PKG := github.com/drone/envsubst/v2/cmd/envsubst

GO_APIDIFF_VER := v0.8.3
GO_APIDIFF_BIN := go-apidiff
GO_APIDIFF := $(abspath $(TOOLS_BIN_DIR)/$(GO_APIDIFF_BIN)-$(GO_APIDIFF_VER))
GO_APIDIFF_PKG := github.com/joelanford/go-apidiff

SHELLCHECK_VER := v0.9.0

TRIVY_VER := 0.69.2

KPROMO_VER := 5ab0dbc74b0228c22a93d240596dff77464aee8f
KPROMO_BIN := kpromo
KPROMO :=  $(abspath $(TOOLS_BIN_DIR)/$(KPROMO_BIN)-$(KPROMO_VER))
# KPROMO_PKG may have to be changed if KPROMO_VER increases its major version.
KPROMO_PKG := sigs.k8s.io/promo-tools/v4/cmd/kpromo

YQ_VER := v4.35.2
YQ_BIN := yq
YQ :=  $(abspath $(TOOLS_BIN_DIR)/$(YQ_BIN)-$(YQ_VER))
YQ_PKG := github.com/mikefarah/yq/v4

GINKGO_BIN := ginkgo
GINKGO_VER := $(call get_go_version,github.com/onsi/ginkgo/v2)
GINKGO := $(abspath $(TOOLS_BIN_DIR)/$(GINKGO_BIN)-$(GINKGO_VER))
GINKGO_PKG := github.com/onsi/ginkgo/v2/ginkgo

GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT_VER := $(shell cat .github/workflows/pr-golangci-lint.yaml | grep [[:space:]]version: | sed 's/.*version: //')
GOLANGCI_LINT := $(abspath $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_BIN)-$(GOLANGCI_LINT_VER))
GOLANGCI_LINT_PKG := github.com/golangci/golangci-lint/v2/cmd/golangci-lint

GOLANGCI_LINT_KAL_BIN := golangci-lint-kube-api-linter
GOLANGCI_LINT_KAL_VER := $(shell cat ./hack/tools/.custom-gcl.yaml | grep version: | sed 's/version: //')
GOLANGCI_LINT_KAL := $(abspath $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_KAL_BIN))

GOVULNCHECK_BIN := govulncheck
GOVULNCHECK_VER := v1.3.0
GOVULNCHECK := $(abspath $(TOOLS_BIN_DIR)/$(GOVULNCHECK_BIN))
GOVULNCHECK_DIR := hack/tools/govulncheck
GOVULNCHECK_TMP_DIR ?= $(GOVULNCHECK_DIR)/govulncheck.tmp

GOVC_VER := $(shell cat go.mod | grep "github.com/vmware/govmomi" | awk '{print $$NF}')
GOVC_BIN := govc
GOVC := $(abspath $(TOOLS_BIN_DIR)/$(GOVC_BIN)-$(GOVC_VER))
GOVC_INSTALL := ./hack/govc-install.sh

KIND_VER := $(call get_test_go_version,sigs.k8s.io/kind)
KIND_BIN := kind
KIND := $(abspath $(TOOLS_BIN_DIR)/$(KIND_BIN)-$(KIND_VER))
KIND_PKG := sigs.k8s.io/kind

CAPI_HACK_TOOLS_VER := 72ec15079db2b5ef4c624437ec241a6eb942f1e0 # Note: this the commit ID of CAPI v1.14.1.

BOSKOSCTL_BIN := boskosctl
BOSKOSCTL := $(abspath $(TOOLS_BIN_DIR)/$(BOSKOSCTL_BIN))

PROWJOB_GEN_VER := $(CAPI_HACK_TOOLS_VER)
PROWJOB_GEN_BIN := prowjob-gen
PROWJOB_GEN := $(abspath $(TOOLS_BIN_DIR)/$(PROWJOB_GEN_BIN)-$(PROWJOB_GEN_VER))
PROWJOB_GEN_PKG := sigs.k8s.io/cluster-api/hack/tools/prowjob-gen

RELEASE_NOTES_VER := $(CAPI_HACK_TOOLS_VER)
RELEASE_NOTES_BIN := release-notes
RELEASE_NOTES := $(abspath $(TOOLS_BIN_DIR)/$(RELEASE_NOTES_BIN)-$(RELEASE_NOTES_VER))
RELEASE_NOTES_PKG := sigs.k8s.io/cluster-api/hack/tools/release/notes

# Define Docker related variables. Releases should modify and double check these vars.
REGISTRY ?= gcr.io/$(shell gcloud config get-value project)
PROD_REGISTRY ?= registry.k8s.io/cluster-api-vsphere

STAGING_REGISTRY ?= gcr.io/k8s-staging-capi-vsphere
STAGING_BUCKET ?= k8s-staging-capi-vsphere

# core
IMAGE_NAME ?= cluster-api-vsphere-controller
CONTROLLER_IMG ?= $(REGISTRY)/$(IMAGE_NAME)

# vcsim controller
VCSIM_IMAGE_NAME ?= cluster-api-vcsim-controller
VCSIM_CONTROLLER_IMG ?= $(REGISTRY)/$(VCSIM_IMAGE_NAME)

# vmoperator controller
VM_OPERATOR_DIR := test/infrastructure/vm-operator
VM_OPERATOR_TMP_DIR ?= $(VM_OPERATOR_DIR)/vm-operator.tmp
# VM_OPERATOR_VERSION ?= 8.0
# VM_OPERATOR_COMMIT ?= de75746a9505ef3161172d99b735d6593c54f0c5
# VM_OPERATOR_COMMIT_DESCRIBE ?= v1.8.6-0-gde75746a
# VM_OPERATOR_VERSION ?= 9.1
# VM_OPERATOR_COMMIT ?= v1.10.0
# VM_OPERATOR_COMMIT_DESCRIBE ?= release-vc-9.1.0-0-g77005588
VM_OPERATOR_VERSION ?= 9.2
VM_OPERATOR_COMMIT ?= 46d6f2c0092ebb12e5a6bb7883b2e2e241058402
VM_OPERATOR_COMMIT_DESCRIBE ?= v1.9.0-1045-g46d6f2c0
VM_OPERATOR_ALL_ARCH = amd64 arm64
VM_OPERATOR_IMAGE_NAME ?= extra/vm-operator
VM_OPERATOR_CONTROLLER_IMG ?= $(STAGING_REGISTRY)/$(VM_OPERATOR_IMAGE_NAME)
VM_OPERATOR_IMAGE_TAG ?= $(VM_OPERATOR_VERSION)

# net operator
NET_OPERATOR_IMAGE_NAME ?= cluster-api-net-operator
NET_OPERATOR_IMG ?= $(STAGING_REGISTRY)/$(NET_OPERATOR_IMAGE_NAME)

# test-extension
TEST_EXTENSION_IMAGE_NAME ?= cluster-api-vsphere-test-extension
TEST_EXTENSION_IMG ?= $(STAGING_REGISTRY)/$(TEST_EXTENSION_IMAGE_NAME)

# boskosctl
BOSKOSCTL_IMG ?= gcr.io/k8s-staging-capi-vsphere/extra/boskosctl
BOSKOSCTL_IMG_TAG ?= $(shell git describe --always --dirty)

# ctr
CTR_IMG ?= gcr.io/k8s-staging-capi-vsphere/extra/ctr
CTR_IMG_TAG ?= $(shell git describe --always --dirty)

# openvpn
OPENVPN_IMG ?= gcr.io/k8s-staging-capi-vsphere/extra/openvpn
OPENVPN_IMG_TAG ?= $(shell git describe --always --dirty)

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

# Additional CAPV vars (everything else is ~ kept in sync with core CAPI)
# Allow overriding manifest generation destination directory
MANIFEST_ROOT ?= ./config
CRD_ROOT ?= $(MANIFEST_ROOT)/default/crd/bases
SUPERVISOR_CRD_ROOT ?= $(MANIFEST_ROOT)/supervisor/crd/bases
VCSIM_CRD_ROOT ?= $(VCSIM_DIR)/config/crd/bases
VCSIM_VM_OPERATOR_CRD_ROOT ?= $(VCSIM_DIR)/config/crd/vm-operator
WEBHOOK_ROOT ?= $(MANIFEST_ROOT)/govmomi/webhook
SUPERVISOR_WEBHOOK_ROOT ?= $(MANIFEST_ROOT)/supervisor/webhook
RBAC_ROOT ?= $(MANIFEST_ROOT)/rbac
VCSIM_RBAC_ROOT ?= $(VCSIM_DIR)/config/rbac
NETOP_RBAC_ROOT ?= $(NETOP_DIR)/config/rbac
TEST_EXTENSION_RBAC_ROOT ?= $(TEST_EXTENSION_DIR)/config/rbac

JANITOR_DIR ?= ./$(TOOLS_DIR)/janitor
JANITOR_ARGS ?= --resource-type=gcve-vsphere-project

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
	$(MAKE) clean-generated-yaml SRC_DIRS="$(CRD_ROOT),$(SUPERVISOR_CRD_ROOT),./config/govmomi/webhook/manifests.yaml,./config/supervisor/webhook/manifests.yaml"
	$(CONTROLLER_GEN) \
		paths=./api/govmomi/v1beta1 \
		paths=./api/govmomi/v1beta2 \
		paths=./internal/webhooks \
		crd:crdVersions=v1 \
		output:crd:dir=$(CRD_ROOT) \
		output:webhook:dir=$(WEBHOOK_ROOT) \
		webhook
	# Generate webhook manifests for supervisor mode separately.
	$(CONTROLLER_GEN) \
		paths=./internal/webhooks/vmware\
		output:webhook:dir=$(SUPERVISOR_WEBHOOK_ROOT) \
		webhook
	$(CONTROLLER_GEN) \
		paths=./ \
		paths=./controllers/... \
		output:rbac:dir=$(RBAC_ROOT) \
		rbac:roleName=manager-role
	$(CONTROLLER_GEN) \
		paths=./api/supervisor/v1beta1 \
		paths=./api/supervisor/v1beta2 \
		crd:crdVersions=v1 \
		output:crd:dir=$(SUPERVISOR_CRD_ROOT)
	# net-operator is used for tests
	$(CONTROLLER_GEN) \
		paths=./$(NETOP_DIR)/... \
		paths=./$(NETOP_DIR)/controllers/... \
        output:rbac:dir=$(NETOP_RBAC_ROOT) \
        rbac:roleName=manager-role
	# test-extension is used for Runtime SDK tests
	$(CONTROLLER_GEN) \
		paths=./$(TEST_EXTENSION_DIR)/... \
        output:rbac:dir=$(TEST_EXTENSION_RBAC_ROOT) \
        rbac:roleName=manager-role
	# vcsim crds are used for tests.
	$(CONTROLLER_GEN) \
		paths=./$(VCSIM_DIR)/api/v1alpha1 \
		crd:crdVersions=v1 \
		output:crd:dir=$(VCSIM_CRD_ROOT)
	$(CONTROLLER_GEN) \
		paths=./$(VCSIM_DIR)/ \
		paths=./$(VCSIM_DIR)/controllers/... \
		output:rbac:dir=$(VCSIM_RBAC_ROOT) \
		rbac:roleName=manager-role

.PHONY: generate-go-deepcopy
generate-go-deepcopy: $(CONTROLLER_GEN) ## Generate deepcopy go code for core
	$(MAKE) clean-generated-deepcopy SRC_DIRS="./api"
	$(CONTROLLER_GEN) \
		object:headerFile=./hack/boilerplate/boilerplate.generatego.txt \
		paths=./api/... \
		paths=./pkg/conversion/api/...
	$(CONTROLLER_GEN) \
    	object:headerFile=./hack/boilerplate/boilerplate.generatego.txt \
    	paths=./$(VCSIM_DIR)/api/...

.PHONY: generate-go-conversions
generate-go-conversions: $(CONTROLLER_GEN) $(CONVERSION_GEN) ## Runs Go related generate targets
	$(MAKE) clean-generated-conversions SRC_DIRS="./api/govmomi/v1beta1,./api/supervisor/v1beta1"
	cd api; $(CONVERSION_GEN) \
		--output-file=zz_generated.conversion.go \
		--go-header-file=./../hack/boilerplate/boilerplate.generatego.txt \
		./govmomi/v1beta1 \
		./supervisor/v1beta1

.PHONY: generate-modules
generate-modules: ## Run go mod tidy to ensure modules are up to date
	go mod tidy
	cd api; go mod tidy
	cd $(TEST_DIR); go mod tidy
	cd $(PACKAGING_DIR); go mod tidy

.PHONY: generate-doctoc
generate-doctoc:
	TRACE=$(TRACE) ./hack/generate-doctoc.sh

.PHONY: generate-e2e-templates
generate-e2e-templates: $(KUSTOMIZE) $(addprefix generate-e2e-templates-, v1.14 v1.15 v1.16 main) ## Generate test templates for all branches

.PHONY: generate-e2e-templates-main
generate-e2e-templates-main: $(KUSTOMIZE) ## Generate test templates for the main branch
	$(MAKE) e2e-flavors-main
	cp "$(RELEASE_DIR)/main/cluster-template.yaml" "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/base/cluster-template.yaml"
	cp "$(RELEASE_DIR)/main/cluster-template-ignition.yaml" "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/base/cluster-template-ignition.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/base" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/hw-upgrade" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-hw-upgrade.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/multi-disk" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-multi-disk.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/storage-policy" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-storage-policy.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/conformance" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-conformance.yaml"
	# Since CAPI uses different flavor names for KCP and MD remediation using MHC
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/mhc-remediation/kcp" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-kcp-remediation.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/mhc-remediation/md" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-md-remediation.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/node-drain" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-node-drain.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/ignition" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-ignition.yaml"
	# generate clusterclass and cluster topology
	cp "$(RELEASE_DIR)/main/clusterclass-template.yaml" "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/clusterclass/clusterclass-quick-start.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/clusterclass" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/clusterclass-quick-start.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/clusterclass-runtimesdk" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/clusterclass-quick-start-runtimesdk.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/clusterclass-runtimesdk-v1beta1" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/clusterclass-quick-start-runtimesdk-v1beta1.yaml"
	cp "$(RELEASE_DIR)/main/cluster-template-topology.yaml" "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/topology/cluster-template-topology.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/topology" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-topology.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/topology-runtimesdk" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-topology-runtimesdk.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/topology-runtimesdk-in-place" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-topology-runtimesdk-in-place.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/topology-runtimesdk-v1beta1" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-topology-runtimesdk-v1beta1.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/topology-taints" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-topology-taints.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/fast-rollout" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-fast-rollout.yaml"
	# for PCI passthrough template
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/pci" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-pci.yaml"
	# for DHCP overrides
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/dhcp-overrides" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-dhcp-overrides.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/ownerrefs-finalizers" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-ownerrefs-finalizers.yaml"
	# for IPAM tests
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/ipam" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/main/cluster-template-ipam.yaml"
	# generate supervisor templates
	cp "$(RELEASE_DIR)/main/cluster-template-supervisor.yaml" "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/base/cluster-template-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/base" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-supervisor.yaml"
	cp "$(RELEASE_DIR)/main/clusterclass-template-supervisor.yaml" "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/clusterclass/clusterclass-quick-start-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/clusterclass" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/clusterclass-quick-start-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/clusterclass-runtimesdk" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/clusterclass-quick-start-supervisor-runtimesdk.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/clusterclass-runtimesdk-v1beta1" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/clusterclass-quick-start-supervisor-runtimesdk-v1beta1.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/clusterclass-scale" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/clusterclass-quick-start-supervisor-scale.yaml"
	cp "$(RELEASE_DIR)/main/cluster-template-topology-supervisor.yaml" "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/topology/cluster-template-topology-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/topology" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-topology-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/topology-autoscaler" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-topology-autoscaler-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/topology-runtimesdk" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-topology-runtimesdk-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/topology-runtimesdk-in-place" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-topology-runtimesdk-in-place-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/topology-runtimesdk-self-hosted" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-topology-runtimesdk-self-hosted-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/topology-runtimesdk-v1beta1" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-topology-runtimesdk-v1beta1-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/topology-scale" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-topology-scale-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/topology-taints" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-topology-taints-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/conformance" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-conformance-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/fast-rollout" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-fast-rollout-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/ownerrefs-finalizers" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/main/cluster-template-ownerrefs-finalizers-supervisor.yaml"

.PHONY: generate-e2e-templates-v1.16
generate-e2e-templates-v1.16: $(KUSTOMIZE)
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.16/clusterclass" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.16/clusterclass-quick-start.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.16/workload" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.16/cluster-template-workload.yaml"

	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.16/clusterclass" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.16/clusterclass-quick-start-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.16/workload" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.16/cluster-template-workload-supervisor.yaml"

.PHONY: generate-e2e-templates-v1.15
generate-e2e-templates-v1.15: $(KUSTOMIZE)
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.15/clusterclass" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.15/clusterclass-quick-start.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.15/workload" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.15/cluster-template-workload.yaml"

	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.15/clusterclass" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.15/clusterclass-quick-start-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.15/workload" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.15/cluster-template-workload-supervisor.yaml"

.PHONY: generate-e2e-templates-v1.14
generate-e2e-templates-v1.14: $(KUSTOMIZE)
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.14/clusterclass" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.14/clusterclass-quick-start.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.14/workload" > "$(E2E_GOVMOMI_TEMPLATE_DIR)/v1.14/cluster-template-workload.yaml"

	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.14/clusterclass" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.14/clusterclass-quick-start-supervisor.yaml"
	"$(KUSTOMIZE)" --load-restrictor LoadRestrictionsNone build "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.14/workload" > "$(E2E_SUPERVISOR_TEMPLATE_DIR)/v1.14/cluster-template-workload-supervisor.yaml"

.PHONY: generate-test-infra-prowjobs
generate-test-infra-prowjobs: $(PROWJOB_GEN) ## Generates the prowjob configurations in test-infra
	@if [ -z "${TEST_INFRA_DIR}" ]; then echo "TEST_INFRA_DIR is not set"; exit 1; fi
	$(PROWJOB_GEN) \
		-config "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/cluster-api-provider-vsphere-prowjob-gen.yaml" \
		-templates-dir "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/templates" \
		-output-dir "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere"; \
		sed -i'' -e 's#0 19 \* \* \*#0 23 * * *#g' "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/cluster-api-provider-vsphere-main-periodics-upgrades.yaml" "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/cluster-api-provider-vsphere-main-periodics.yaml"; \
		sed -i'' -e 's#0 21 \* \* \*#0 1 * * *#g' "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/cluster-api-provider-vsphere-main-periodics-upgrades.yaml" "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/cluster-api-provider-vsphere-main-periodics.yaml"; \
		LATEST_RELEASE="$$(cat "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/cluster-api-provider-vsphere-prowjob-gen.yaml" | yq '.prow_ignored.branches | keys | .[1]' | tr '.' '-')"; \
		sed -i'' -e 's#0 19 \* \* \*#0 3 * * *#g' "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/cluster-api-provider-vsphere-$${LATEST_RELEASE}-periodics-upgrades.yaml" "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/cluster-api-provider-vsphere-$${LATEST_RELEASE}-periodics.yaml"; \
		sed -i'' -e 's#0 21 \* \* \*#0 5 * * *#g' "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/cluster-api-provider-vsphere-$${LATEST_RELEASE}-periodics-upgrades.yaml" "$(TEST_INFRA_DIR)/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/cluster-api-provider-vsphere-$${LATEST_RELEASE}-periodics.yaml"

## --------------------------------------
## Lint / Verify
## --------------------------------------

##@ lint and verify:

.PHONY: lint
lint: $(GOLANGCI_LINT) $(GOLANGCI_LINT_KAL) ## Lint the codebase
	$(GOLANGCI_LINT) run -v $(GOLANGCI_LINT_EXTRA_ARGS)
	cd api; $(GOLANGCI_LINT) run --config $(ROOT_DIR)/.golangci.yml -v $(GOLANGCI_LINT_EXTRA_ARGS)
	cd $(TEST_DIR); $(GOLANGCI_LINT) run --config $(ROOT_DIR)/.golangci.yml -v $(GOLANGCI_LINT_EXTRA_ARGS)
	cd $(PACKAGING_DIR); $(GOLANGCI_LINT) run --config $(ROOT_DIR)/.golangci.yml -v $(GOLANGCI_LINT_EXTRA_ARGS)
	cd api; $(GOLANGCI_LINT_KAL) run -v --config $(ROOT_DIR)/.golangci-kal.yml $(GOLANGCI_LINT_EXTRA_ARGS)

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Lint the codebase and run auto-fixers if supported by the linter
	GOLANGCI_LINT_EXTRA_ARGS=--fix $(MAKE) lint

.PHONY: lint-api
lint-api: $(GOLANGCI_LINT_KAL)
	cd api; $(GOLANGCI_LINT_KAL) run -v --config $(ROOT_DIR)/.golangci-kal.yml $(GOLANGCI_LINT_EXTRA_ARGS)

.PHONY: lint-api-fix
lint-api-fix: $(GOLANGCI_LINT_KAL)
	GOLANGCI_LINT_EXTRA_ARGS=--fix $(MAKE) lint-api

APIDIFF_OLD_COMMIT ?= $(shell git rev-parse origin/main)

.PHONY: apidiff
apidiff: $(GO_APIDIFF) ## Check for API differences
	$(GO_APIDIFF) $(APIDIFF_OLD_COMMIT) --print-compatible

ALL_VERIFY_CHECKS = licenses boilerplate shellcheck modules gen doctoc flavors go-directive

.PHONY: verify
verify: $(addprefix verify-,$(ALL_VERIFY_CHECKS)) ## Run all verify-* targets

.PHONY: verify-go-directive
verify-go-directive:
	TRACE=$(TRACE) ./hack/verify-go-directive.sh -g $(GO_DIRECTIVE_VERSION)

.PHONY: verify-modules
verify-modules: generate-modules  ## Verify go modules are up to date
	@if !(git diff --quiet HEAD -- go.sum go.mod api/go.mod api/go.sum $(TEST_DIR)/go.mod $(TEST_DIR)/go.sum $(PACKAGING_DIR)/go.mod $(PACKAGING_DIR)/go.sum); then \
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
	$(GOVULNCHECK) -C "api" ./... && R2=$$? || R2=$$?; \
	$(GOVULNCHECK) -C "$(TEST_DIR)" ./... && R3=$$? || R3=$$?; \
	$(GOVULNCHECK) -C "$(PACKAGING_DIR)" ./... && R4=$$? || R4=$$?; \
	if [ "$$R1" -ne "0" ] || [ "$$R2" -ne "0" ] || [ "$$R3" -ne "0" ] || [ "$$R4" -ne "0" ]; then \
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

## --------------------------------------
## Build
## --------------------------------------

##@ build:

.PHONY: manager
manager: ## Build the vsphere manager binary into the ./bin folder
	CGO_ENABLED=0 go build -trimpath -gcflags "$(GCFLAGS)" -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/manager sigs.k8s.io/cluster-api-provider-vsphere

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
## reads Dockerfile from stdin to avoid an incorrectly cached Dockerfile (https://github.com/moby/buildkit/issues/1368)
	cat ./Dockerfile | DOCKER_BUILDKIT=1 docker build --build-arg builder_image=$(GO_CONTAINER_IMAGE) --build-arg goproxy=$(GOPROXY) --build-arg ARCH=$(ARCH) --build-arg gcflags="$(GCFLAGS)" --build-arg ldflags="$(LDFLAGS)" . -t $(CONTROLLER_IMG)-$(ARCH):$(TAG) --file -
	@if [ "${DOCKER_BUILD_MODIFY_MANIFESTS}" = "true" ]; then \
  		$(MAKE) set-manifest-image MANIFEST_IMG=$(CONTROLLER_IMG)-$(ARCH) MANIFEST_TAG=$(TAG) TARGET_RESOURCE="./config/base/manager_image_patch.yaml"; \
		$(MAKE) set-manifest-pull-policy TARGET_RESOURCE="./config/base/manager_pull_policy.yaml"; \
    fi

.PHONY: docker-build-vcsim
docker-build-vcsim: docker-pull-prerequisites ## Build the docker image for vcsim controller manager
## reads Dockerfile from stdin to avoid an incorrectly cached Dockerfile (https://github.com/moby/buildkit/issues/1368)
	cat $(VCSIM_DIR)/Dockerfile | DOCKER_BUILDKIT=1 docker build --build-arg builder_image=$(GO_CONTAINER_IMAGE) --build-arg goproxy=$(GOPROXY) --build-arg ARCH=$(ARCH) --build-arg gcflags="$(GCFLAGS)" --build-arg ldflags="$(LDFLAGS)" . -t $(VCSIM_CONTROLLER_IMG)-$(ARCH):$(TAG) --file -
	@if [ "${DOCKER_BUILD_MODIFY_MANIFESTS}" = "true" ]; then \
  		$(MAKE) set-manifest-image MANIFEST_IMG=$(VCSIM_CONTROLLER_IMG)-$(ARCH) MANIFEST_TAG=$(TAG) TARGET_RESOURCE="./$(VCSIM_DIR)/config/default/manager_image_patch.yaml"; \
		$(MAKE) set-manifest-pull-policy TARGET_RESOURCE="./$(VCSIM_DIR)/config/default/manager_pull_policy.yaml"; \
    fi

.PHONY: docker-build-net-operator
docker-build-net-operator: docker-pull-prerequisites ## Build the docker image for net-operator controller manager
## reads Dockerfile from stdin to avoid an incorrectly cached Dockerfile (https://github.com/moby/buildkit/issues/1368)
	cat $(NETOP_DIR)/Dockerfile | DOCKER_BUILDKIT=1 docker build --build-arg builder_image=$(GO_CONTAINER_IMAGE) --build-arg goproxy=$(GOPROXY) --build-arg ARCH=$(ARCH) --build-arg gcflags="$(GCFLAGS)" --build-arg ldflags="$(LDFLAGS)" . -t $(NET_OPERATOR_IMG)-$(ARCH):$(TAG) --file -
	@if [ "${DOCKER_BUILD_MODIFY_MANIFESTS}" = "true" ]; then \
		$(MAKE) set-manifest-image MANIFEST_IMG=$(NET_OPERATOR_IMG)-$(ARCH) MANIFEST_TAG=$(TAG) TARGET_RESOURCE="./$(NETOP_DIR)/config/default/manager_image_patch.yaml"; \
		$(MAKE) set-manifest-pull-policy TARGET_RESOURCE="./$(NETOP_DIR)/config/default/manager_pull_policy.yaml"; \
	fi

.PHONY: docker-build-test-extension
docker-build-test-extension: docker-pull-prerequisites ## Build the docker image for test-extension controller manager
## reads Dockerfile from stdin to avoid an incorrectly cached Dockerfile (https://github.com/moby/buildkit/issues/1368)
	cat $(TEST_EXTENSION_DIR)/Dockerfile | DOCKER_BUILDKIT=1 docker build --build-arg builder_image=$(GO_CONTAINER_IMAGE) --build-arg goproxy=$(GOPROXY) --build-arg ARCH=$(ARCH) --build-arg gcflags="$(GCFLAGS)" --build-arg ldflags="$(LDFLAGS)" . -t $(TEST_EXTENSION_IMG)-$(ARCH):$(TAG) --file -
	@if [ "${DOCKER_BUILD_MODIFY_MANIFESTS}" = "true" ]; then \
		$(MAKE) set-manifest-image MANIFEST_IMG=$(TEST_EXTENSION_IMG)-$(ARCH) MANIFEST_TAG=$(TAG) TARGET_RESOURCE="./$(TEST_EXTENSION_DIR)/config/default/manager_image_patch.yaml"; \
		$(MAKE) set-manifest-pull-policy TARGET_RESOURCE="./$(TEST_EXTENSION_DIR)/config/default/manager_pull_policy.yaml"; \
	fi

.PHONY: docker-build-boskosctl
docker-build-boskosctl:
	cat hack/tools/boskosctl/Dockerfile | DOCKER_BUILDKIT=1 docker build --build-arg builder_image=$(GO_CONTAINER_IMAGE) --build-arg goproxy=$(GOPROXY) . -t $(BOSKOSCTL_IMG):$(BOSKOSCTL_IMG_TAG) --file -
	docker tag $(BOSKOSCTL_IMG):$(BOSKOSCTL_IMG_TAG) $(BOSKOSCTL_IMG):latest
.PHONY: build

.PHONY: docker-push-boskosctl
docker-push-boskosctl:
	docker push $(BOSKOSCTL_IMG):$(BOSKOSCTL_IMG_TAG)
	docker push $(BOSKOSCTL_IMG):latest

.PHONY: docker-build-ctr
docker-build-ctr:
	cat hack/tools/ctr/Dockerfile | DOCKER_BUILDKIT=1 docker build --platform linux/amd64 $(ROOT_DIR)/hack/tools/ctr/. -t $(CTR_IMG):$(CTR_IMG_TAG) --file -
	docker tag $(CTR_IMG):$(CTR_IMG_TAG) $(CTR_IMG):latest
.PHONY: build

.PHONY: docker-push-ctr
docker-push-ctr:
	docker push $(CTR_IMG):$(CTR_IMG_TAG)
	docker push $(CTR_IMG):latest

.PHONY: docker-build-openvpn
docker-build-openvpn:
	cat hack/tools/openvpn/Dockerfile | DOCKER_BUILDKIT=1 docker build --build-arg builder_image=$(GO_CONTAINER_IMAGE) --build-arg goproxy=$(GOPROXY) . -t $(OPENVPN_IMG):$(OPENVPN_IMG_TAG) --file -
	docker tag $(OPENVPN_IMG):$(OPENVPN_IMG_TAG) $(OPENVPN_IMG):latest
.PHONY: build

.PHONY: docker-push-openvpn
docker-push-openvpn:
	docker push $(OPENVPN_IMG):$(OPENVPN_IMG_TAG)
	docker push $(OPENVPN_IMG):latest

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
	cd api; KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" GOVC_BIN_PATH=$(GOVC) go test -v ./... $(TEST_ARGS)
	KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" GOVC_BIN_PATH=$(GOVC) go test -v ./controllers/... ./pkg/... ./internal/... ./hack/tools/... $(TEST_ARGS)

.PHONY: test-verbose
test-verbose: ## Run unit tests with verbose flag
	$(MAKE) test TEST_ARGS="$(TEST_ARGS) -v"

.PHONY: test-junit
test-junit: $(SETUP_ENVTEST) $(GOTESTSUM) $(GOVC) ## Run unit tests
	# Note: running ensure.go to make sure tests run with the correct kube-kins image in CI
	hack/ensure-go.sh
	# Note: ARTIFACTS must be set so the ginkgo suites write junit reports to the ARTIFACTS folder
	set +o errexit; (ARTIFACTS=$(ARTIFACTS) KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" GOVC_BIN_PATH=$(GOVC) go test -json ./controllers/... ./pkg/... ./internal/... ./hack/tools/... $(TEST_ARGS); echo $$? > $(ARTIFACTS)/junit.exitcode) | tee $(ARTIFACTS)/junit.stdout
	$(GOTESTSUM) --junitfile $(ARTIFACTS)/junit.xml --raw-command cat $(ARTIFACTS)/junit.stdout
	exit $$(cat $(ARTIFACTS)/junit.exitcode)
	cd api; set +o errexit; (ARTIFACTS=$(ARTIFACTS) KUBEBUILDER_ASSETS="$(KUBEBUILDER_ASSETS)" GOVC_BIN_PATH=$(GOVC) go test -json ./... $(TEST_ARGS); echo $$? > $(ARTIFACTS)/junit-api.exitcode) | tee $(ARTIFACTS)/junit-api.stdout
	$(GOTESTSUM) --junitfile $(ARTIFACTS)/junit-api.xml --raw-command cat $(ARTIFACTS)/junit-api.stdout
	exit $$(cat $(ARTIFACTS)/junit-api.exitcode)

.PHONY: test-cover
test-cover: ## Run unit tests and generate a coverage report
	$(MAKE) test TEST_ARGS="$(TEST_ARGS) -coverpkg=./... -coverprofile=coverage.out"
	./hack/codecov-ignore.sh
	go tool cover -func=coverage.out -o coverage.txt
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=api/coverage.out -o coverage-api.txt
	go tool cover -html=api/coverage.out -o coverage-api.html

.PHONY: e2e-images
e2e-images: ## Build the e2e manager image
	# please ensure the generated image name matches image names used in the E2E_CONF_FILE;
    # also the same settings must exist in e2e.sh
	$(MAKE) REGISTRY=gcr.io/k8s-staging-capi-vsphere PULL_POLICY=IfNotPresent TAG=dev docker-build
	$(MAKE) REGISTRY=gcr.io/k8s-staging-capi-vsphere PULL_POLICY=IfNotPresent TAG=dev docker-build-vcsim
	$(MAKE) REGISTRY=gcr.io/k8s-staging-capi-vsphere PULL_POLICY=IfNotPresent TAG=dev docker-build-net-operator
	$(MAKE) REGISTRY=gcr.io/k8s-staging-capi-vsphere PULL_POLICY=IfNotPresent TAG=dev docker-build-test-extension

.PHONY: e2e
e2e: e2e-images generate-e2e-templates
e2e: $(GINKGO) $(KUSTOMIZE) $(KIND) $(GOVC) ## Run e2e tests
	@echo PATH="$(PATH)"
	@echo
	@echo Contents of $(TOOLS_BIN_DIR):
	@ls $(TOOLS_BIN_DIR)
	@echo
	time $(GINKGO) -v --trace \
		--nodes=$(GINKGO_NODES) --timeout=$(GINKGO_TIMEOUT) \
		--label-filter="$(GINKGO_LABEL_FILTER)" --focus="$(GINKGO_FOCUS)" $(_SKIP_ARGS) \
		--poll-progress-after=$(GINKGO_POLL_PROGRESS_AFTER) --poll-progress-interval=$(GINKGO_POLL_PROGRESS_INTERVAL) \
		--fail-on-pending --fail-on-empty \
		--no-color=$(GINKGO_NOCOLOR) --output-dir="$(ARTIFACTS)" --junit-report="junit.e2e_suite.1.xml" \
		$(GINKGO_ARGS) ./test/e2e -- \
		--e2e.artifacts-folder="$(ARTIFACTS)" \
		--e2e.config="$(E2E_CONF_FILE)" \
		--e2e.config-overrides="$(E2E_CONF_OVERRIDE_FILE)" \
		--e2e.skip-resource-cleanup=$(SKIP_RESOURCE_CLEANUP) \
		--e2e.use-existing-cluster="$(USE_EXISTING_CLUSTER)" \
		--e2e.ip-pool='$(E2E_VSPHERE_IP_POOL)'

## --------------------------------------
## Release
## --------------------------------------

##@ release:

## latest git tag for the commit, e.g., v0.3.10
RELEASE_TAG ?= $(shell git describe --abbrev=0 2>/dev/null)
## if we are on a pre-release, we have to determine the previous release for the release-notes generator.
ifneq (,$(findstring -,$(RELEASE_TAG)))
	# extract the major and minor version from the RELEASE_TAG
	_RELEASE_TAG_MAJOR ?= $(word 1,$(subst ., ,$(RELEASE_TAG:v%=%)))
	_RELEASE_TAG_MINOR ?= $(word 2,$(subst ., ,$(RELEASE_TAG:v%=%)))
	# Find the previous release of the same major + minor version (including pre-releases) or the previous .0 minor release.
	_PREVIOUS_RELEASE_TAG ?= $(shell git tag -l | grep -E -e '^v[0-9]+\.[0-9]+\.0+$$|^v$(_RELEASE_TAG_MAJOR)\.$(_RELEASE_TAG_MINOR)\.' | sort -V | grep -B1 $(RELEASE_TAG) | head -n 1 2>/dev/null)
	# Set the argument for release-notes generator to provide the for pre-releases mandatory
	# `--previous-release-version` flag, but only if the previous release is also a pre-release.
	# For the first pre-release (e.g. v1.15.0-beta.0 after v1.14.0), the generator must NOT
	# receive `--previous-release-version`.
	ifneq (,$(findstring alpha.,$(_PREVIOUS_RELEASE_TAG))$(findstring beta.,$(_PREVIOUS_RELEASE_TAG))$(findstring rc.,$(_PREVIOUS_RELEASE_TAG)))
		RELEASE_NOTES_PRE_RELEASE_ARG ?= "--previous-release-version=tags/$(_PREVIOUS_RELEASE_TAG)"
	else
		RELEASE_NOTES_PRE_RELEASE_ARG ?=
	endif
endif
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

.PHONY: e2e-flavors ## Create dev flavor manifests for e2e testing
e2e-flavors: $(KUSTOMIZE) $(addprefix e2e-flavors-, main)

.PHONY: e2e-flavors-main ## Create dev flavor manifests for e2e testing for main branch
e2e-flavors-main: $(RELEASE_DIR)
	mkdir -p $(RELEASE_DIR)/main
	$(MAKE) generate-flavors FLAVOR_DIR=$(RELEASE_DIR)/main


.PHONY: generate-flavors
generate-flavors: $(FLAVOR_DIR)
	cd $(PACKAGING_DIR)/flavorgen; GOTOOLCHAIN=auto go build -o flavorgen ./
	$(PACKAGING_DIR)/flavorgen/flavorgen --output-dir $(FLAVOR_DIR)
	rm $(PACKAGING_DIR)/flavorgen/flavorgen

.PHONY: release-staging
release-staging: ## Build and push container images to the staging registry
	REGISTRY=$(STAGING_REGISTRY) $(MAKE) docker-build-all docker-push-all
	REGISTRY=$(STAGING_REGISTRY) $(MAKE) release-alias-tag
	# Set the manifest image to $(STAGING_REGISTRY)/$(IMAGE_NAME):$(RELEASE_ALIAS_TAG) and pull policy to IfNotPresent.
	$(MAKE) manifest-modification REGISTRY=$(STAGING_REGISTRY) RELEASE_TAG=$(RELEASE_ALIAS_TAG) PULL_POLICY=IfNotPresent
	## Build the manifests into $(RELEASE_DIR)
	$(MAKE) release-manifests STAGE=release MANIFEST_DIR=$(RELEASE_DIR)
	# Example manifest location: https://storage.googleapis.com/k8s-staging-capi-vsphere/components/main/infrastructure-components.yaml
	# Please note that these files are deleted after a certain period, at the time of this writing 60 days after file creation.
	gsutil cp $(RELEASE_DIR)/* gs://$(STAGING_BUCKET)/components/$(RELEASE_ALIAS_TAG)


.PHONY: release-staging-nightly
release-staging-nightly: ## Re-tags container images to a nightly tag. Builds and pushes nightly manifests to the staging bucket.
	$(eval NEW_RELEASE_ALIAS_TAG := nightly_$(RELEASE_ALIAS_TAG)_$(shell date +'%Y%m%d'))
	echo $(NEW_RELEASE_ALIAS_TAG)
	$(MAKE) release-alias-tag TAG=$(RELEASE_ALIAS_TAG) RELEASE_ALIAS_TAG=$(NEW_RELEASE_ALIAS_TAG)
	# Set the manifest image to $(STAGING_REGISTRY)/$(IMAGE_NAME):$(NEW_RELEASE_ALIAS_TAG) and pull policy to IfNotPresent.
	$(MAKE) manifest-modification REGISTRY=$(STAGING_REGISTRY) RELEASE_TAG=$(NEW_RELEASE_ALIAS_TAG) PULL_POLICY=IfNotPresent
	## Build the manifests into $(RELEASE_DIR)
	$(MAKE) release-manifests STAGE=release MANIFEST_DIR=$(RELEASE_DIR)
	# Example manifest location: https://storage.googleapis.com/k8s-staging-capi-vsphere/components/nightly_main_20240423/infrastructure-components.yaml
	# Please note that these files are deleted after a certain period, at the time of this writing 60 days after file creation.
	gsutil cp $(RELEASE_DIR)/* gs://$(STAGING_BUCKET)/components/$(NEW_RELEASE_ALIAS_TAG)

.PHONY: release-alias-tag
release-alias-tag: ## Add the release alias tag to the last build tag
	gcloud container images add-tag $(CONTROLLER_IMG):$(TAG) $(CONTROLLER_IMG):$(RELEASE_ALIAS_TAG)

.PHONY: generate-release-notes
generate-release-notes: $(RELEASE_NOTES_DIR) $(RELEASE_NOTES)
	# Reset the file
	echo -n > $(RELEASE_NOTES_DIR)/$(RELEASE_TAG).md
	"$(RELEASE_NOTES)" --release=$(RELEASE_TAG) $(RELEASE_NOTES_PRE_RELEASE_ARG) --repository kubernetes-sigs/cluster-api-provider-vsphere --prefix-area-label=false --deprecation=false --add-kubernetes-version-support=false >> $(RELEASE_NOTES_DIR)/$(RELEASE_TAG).md

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
## vm-operator
## --------------------------------------

##@ vm-operator:

.PHONY: release-vm-operator
release-vm-operator: docker-build-all-vm-operator generate-manifests-vm-operator docker-push-all-vm-operator clean-vm-operator ## Build and push the vm-operator image and manifest for usage in CI

.PHONY: release-vm-operator-local
release-vm-operator-local: docker-build-all-vm-operator generate-manifests-vm-operator clean-vm-operator ## Build the vm-operator image and manifest for local usage only

.PHONY: checkout-vm-operator
checkout-vm-operator:
	@if [ -z "${VM_OPERATOR_VERSION}" ]; then echo "VM_OPERATOR_VERSION is not set"; exit 1; fi
	@if [ -z "${VM_OPERATOR_COMMIT_DESCRIBE}" ]; then echo "VM_OPERATOR_COMMIT_DESCRIBE is not set"; exit 1; fi
	@if [ -z "${VM_OPERATOR_COMMIT}" ]; then echo "VM_OPERATOR_COMMIT is not set"; exit 1; fi
	@if [ -d "$(VM_OPERATOR_TMP_DIR)" ]; then \
		echo "$(VM_OPERATOR_TMP_DIR) exists, skipping clone"; \
	else \
		git clone "https://github.com/vmware-tanzu/vm-operator.git" "$(VM_OPERATOR_TMP_DIR)"; \
		cd "$(VM_OPERATOR_TMP_DIR)"; \
		git checkout "$(VM_OPERATOR_COMMIT)"; \
	fi
	@cd "$(ROOT_DIR)/$(VM_OPERATOR_TMP_DIR)"; \
	if [ "$$(git describe --no-dirty 2> /dev/null)" != "$(VM_OPERATOR_COMMIT_DESCRIBE)" ]; then \
		if [ "$$(git describe --no-dirty 2> /dev/null)" != "api/$(VM_OPERATOR_COMMIT_DESCRIBE)" ]; then \
			echo "ERROR: checked out version $$(git describe --no-dirty 2> /dev/null) does not match expected version $(VM_OPERATOR_COMMIT_DESCRIBE) or api/$(VM_OPERATOR_COMMIT_DESCRIBE)"; \
			exit 1; \
		fi \
	fi; \
	echo "Using vm-operator version $(VM_OPERATOR_VERSION), commit $(VM_OPERATOR_COMMIT), build $(VM_OPERATOR_COMMIT_DESCRIBE)";

.PHONY: generate-manifests-vm-operator
generate-manifests-vm-operator: $(RELEASE_DIR) $(KUSTOMIZE) checkout-vm-operator ## Build the vm-operator manifest yaml file
	@if [ -z "${VM_OPERATOR_VERSION}" ]; then echo "VM_OPERATOR_VERSION is not set"; exit 1; fi
	$(MAKE) generate-manifests-vm-operator-$(VM_OPERATOR_VERSION)

generate-manifests-vm-operator-8.0:
	$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone "$(VM_OPERATOR_TMP_DIR)/config/wcp" > "$(VM_OPERATOR_DIR)/config/$(VM_OPERATOR_VERSION)/vm-operator.yaml"
	$(KUSTOMIZE) build "$(VM_OPERATOR_DIR)/config/$(VM_OPERATOR_VERSION)" > "$(VM_OPERATOR_DIR)/vm-operator-$(VM_OPERATOR_VERSION).yaml"

generate-manifests-vm-operator-9.1:
	@cd "$(ROOT_DIR)/$(VM_OPERATOR_TMP_DIR)"; \
	make kustomize-wcp
	@cd "$(ROOT_DIR)"
	cp "$(ROOT_DIR)/$(VM_OPERATOR_TMP_DIR)/artifacts/local-deployment.yaml" "$(VM_OPERATOR_DIR)/config/$(VM_OPERATOR_VERSION)/vm-operator.yaml"
	$(KUSTOMIZE) build "$(VM_OPERATOR_DIR)/config/$(VM_OPERATOR_VERSION)" > "$(VM_OPERATOR_DIR)/vm-operator-$(VM_OPERATOR_VERSION).yaml"
	$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone "$(VCSIM_VM_OPERATOR_CRD_ROOT)/$(VM_OPERATOR_VERSION)" > "$(VCSIM_VM_OPERATOR_CRD_ROOT)/vm-operator-$(VM_OPERATOR_VERSION).yaml"

generate-manifests-vm-operator-9.2:
	@cd "$(ROOT_DIR)/$(VM_OPERATOR_TMP_DIR)"; \
	GOTOOLCHAIN=auto make kustomize-wcp
	@cd "$(ROOT_DIR)"
	cp "$(ROOT_DIR)/$(VM_OPERATOR_TMP_DIR)/artifacts/local-deployment.yaml" "$(VM_OPERATOR_DIR)/config/$(VM_OPERATOR_VERSION)/vm-operator.yaml"
	$(KUSTOMIZE) build "$(VM_OPERATOR_DIR)/config/$(VM_OPERATOR_VERSION)" > "$(VM_OPERATOR_DIR)/vm-operator-$(VM_OPERATOR_VERSION).yaml"
	$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone "$(VCSIM_VM_OPERATOR_CRD_ROOT)/$(VM_OPERATOR_VERSION)" > "$(VCSIM_VM_OPERATOR_CRD_ROOT)/vm-operator-$(VM_OPERATOR_VERSION).yaml"

.PHONY: docker-build-all-vm-operator
docker-build-all-vm-operator: checkout-vm-operator
	$(MAKE) docker-build-vm-operator-vm-operator-$(VM_OPERATOR_VERSION)

# IMPORTANT: before running this command please add --provenance=false to the docker build command in test/infrastructure/vm-operator/vm-operator.tmp/hack/build-container.sh
docker-build-vm-operator-vm-operator-8.0:
	@if [ -z "${VM_OPERATOR_IMAGE_TAG}" ]; then echo "VM_OPERATOR_IMAGE_TAG is not set"; exit 1; fi
	cd $(VM_OPERATOR_TMP_DIR); \
	$(MAKE) IMAGE=$(VM_OPERATOR_CONTROLLER_IMG)-amd64 IMAGE_TAG=$(VM_OPERATOR_IMAGE_TAG) GOARCH=amd64 docker-build ; \
	$(MAKE) IMAGE=$(VM_OPERATOR_CONTROLLER_IMG)-arm64 IMAGE_TAG=$(VM_OPERATOR_IMAGE_TAG) GOARCH=arm64 docker-build

docker-build-vm-operator-vm-operator-9.1:
	@if [ -z "${VM_OPERATOR_IMAGE_TAG}" ]; then echo "VM_OPERATOR_IMAGE_TAG is not set"; exit 1; fi
	cd $(VM_OPERATOR_TMP_DIR); \
	IMAGE_TAG=$(VM_OPERATOR_IMAGE_TAG) IMAGE=$(VM_OPERATOR_CONTROLLER_IMG)-amd64 IMAGE_FILE=artifacts/vm-operator-amd64.tar ADDITIONAL_CRI_BUILD_FLAGS="--provenance=false" make image-build-amd64; \
	docker load -i $(ROOT_DIR)/$(VM_OPERATOR_TMP_DIR)/artifacts/vm-operator-amd64.tar; \
	IMAGE_TAG=$(VM_OPERATOR_IMAGE_TAG) IMAGE=$(VM_OPERATOR_CONTROLLER_IMG)-arm64 IMAGE_FILE=artifacts/vm-operator-arm64.tar ADDITIONAL_CRI_BUILD_FLAGS="--provenance=false" make image-build-arm64; \
	docker load -i $(ROOT_DIR)/$(VM_OPERATOR_TMP_DIR)/artifacts/vm-operator-arm64.tar

docker-build-vm-operator-vm-operator-9.2:
	@if [ -z "${VM_OPERATOR_IMAGE_TAG}" ]; then echo "VM_OPERATOR_IMAGE_TAG is not set"; exit 1; fi
	cd $(VM_OPERATOR_TMP_DIR); \
	GOTOOLCHAIN=auto IMAGE_TAG=$(VM_OPERATOR_IMAGE_TAG) IMAGE=$(VM_OPERATOR_CONTROLLER_IMG)-amd64 IMAGE_FILE=artifacts/vm-operator-amd64.tar ADDITIONAL_CRI_BUILD_FLAGS="--provenance=false" make image-build-amd64; \
	docker load -i $(ROOT_DIR)/$(VM_OPERATOR_TMP_DIR)/artifacts/vm-operator-amd64.tar; \
	GOTOOLCHAIN=auto IMAGE_TAG=$(VM_OPERATOR_IMAGE_TAG) IMAGE=$(VM_OPERATOR_CONTROLLER_IMG)-arm64 IMAGE_FILE=artifacts/vm-operator-arm64.tar ADDITIONAL_CRI_BUILD_FLAGS="--provenance=false" make image-build-arm64; \
	docker load -i $(ROOT_DIR)/$(VM_OPERATOR_TMP_DIR)/artifacts/vm-operator-arm64.tar

.PHONY: docker-push-all-vm-operator
docker-push-all-vm-operator: $(addprefix docker-vm-operator-push-,$(VM_OPERATOR_ALL_ARCH))  ## Push the docker images to be included in the release for all architectures + related multiarch manifests
	$(MAKE) docker-push-manifest-vm-operator

docker-vm-operator-push-%:
	$(MAKE) ARCH=$* docker-push-vm-operator

.PHONY: docker-push-vm-operator
docker-push-vm-operator:
	@if [ -z "${VM_OPERATOR_IMAGE_TAG}" ]; then echo "VM_OPERATOR_IMAGE_TAG is not set"; exit 1; fi
	docker push $(VM_OPERATOR_CONTROLLER_IMG)-$(ARCH):$(VM_OPERATOR_VERSION)

.PHONY: docker-push-manifest-vm-operator
docker-push-manifest-vm-operator:
	@if [ -z "${VM_OPERATOR_IMAGE_TAG}" ]; then echo "VM_OPERATOR_IMAGE_TAG is not set"; exit 1; fi
	export  VM_OPERATOR_IMAGE_TAG=$(VM_OPERATOR_IMAGE_TAG)
	docker manifest create --amend $(VM_OPERATOR_CONTROLLER_IMG):$(VM_OPERATOR_IMAGE_TAG) $(shell echo $(VM_OPERATOR_ALL_ARCH) | sed -e "s~[^ ]*~$(VM_OPERATOR_CONTROLLER_IMG)\-&:$(VM_OPERATOR_IMAGE_TAG)~g")
	@for arch in $(VM_OPERATOR_ALL_ARCH); do docker manifest annotate --arch $${arch} ${VM_OPERATOR_CONTROLLER_IMG}:${VM_OPERATOR_IMAGE_TAG} ${VM_OPERATOR_CONTROLLER_IMG}-$${arch}:${VM_OPERATOR_IMAGE_TAG}; done
	docker manifest push --purge $(VM_OPERATOR_CONTROLLER_IMG):$(VM_OPERATOR_IMAGE_TAG)

.PHONY: clean-vm-operator
clean-vm-operator:
	rm -fr $(VM_OPERATOR_TMP_DIR)

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

.PHONY: clean-ci
clean-ci: ## Cleanup orphaned objects in CI
	@if [ -z "${GOVC_USERNAME}" ]; then echo "GOVC_USERNAME is not set"; exit 1; fi
	@if [ -z "${GOVC_PASSWORD}" ]; then echo "GOVC_PASSWORD is not set"; exit 1; fi
	@if [ -z "${GOVC_URL}" ]; then echo "GOVC_URL is not set"; exit 1; fi
	@if [ -z "${VSPHERE_TLS_THUMBPRINT}" ]; then echo "VSPHERE_TLS_THUMBPRINT is not set"; exit 1; fi
	@if [ -z "${BOSKOS_HOST}" ]; then echo "BOSKOS_HOST is not set"; exit 1; fi
	go run $(JANITOR_DIR) --dry-run=false $(JANITOR_ARGS)

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

.PHONY: $(PROWJOB_GEN_BIN)
$(PROWJOB_GEN_BIN): $(PROWJOB_GEN) ## Build a local copy of prowjob-gen.

.PHONY: $(BOSKOSCTL_BIN)
$(BOSKOSCTL_BIN): $(BOSKOSCTL) ## Build a local copy of boskosctl.

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

$(GOLANGCI_LINT_KAL): $(GOLANGCI_LINT) # Build golangci-lint-kal from custom configuration.
	cd $(TOOLS_DIR); $(GOLANGCI_LINT) custom

.PHONY: $(GOVULNCHECK_BIN)
$(GOVULNCHECK_BIN): $(GOVULNCHECK) ## Build a local copy of govulncheck.

.PHONY: $(GOVC_BIN)
$(GOVC_BIN): $(GOVC) ## Build a local copy of govc.

.PHONY: $(KIND_BIN)
$(KIND_BIN): $(KIND) ## Build a local copy of kind.

.PHONY: $(RELEASE_NOTES_BIN)
$(RELEASE_NOTES_BIN): $(RELEASE_NOTES) ## Build a local copy of release-notes.

$(CONTROLLER_GEN): # Build controller-gen.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(CONTROLLER_GEN_PKG) $(CONTROLLER_GEN_BIN) $(CONTROLLER_GEN_VER)

## We are forcing a rebuilt of conversion-gen via PHONY so that we're always using an up-to-date version.
## We can't use a versioned name for the binary, because that would be reflected in generated files.
.PHONY: $(CONVERSION_GEN)
$(CONVERSION_GEN): # Build conversion-gen.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(CONVERSION_GEN_PKG) $(CONVERSION_GEN_BIN) $(CONVERSION_GEN_VER)

$(BOSKOSCTL): # Build boskosctl from tools folder.
	go build -o $(TOOLS_BIN_DIR)/$(BOSKOSCTL_BIN) ./hack/tools/boskosctl

$(PROWJOB_GEN): # Build prowjob-gen.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_TOOLS_BUILD) $(PROWJOB_GEN_PKG) $(PROWJOB_GEN_BIN) $(PROWJOB_GEN_VER)

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
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GINKGO_PKG) $(GINKGO_BIN) $(GINKGO_VER)

$(GOLANGCI_LINT): # Build golangci-lint.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GOLANGCI_LINT_PKG) $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_VER)

$(GOVC): # Build GOVC.
	GOBIN=$(TOOLS_BIN_DIR) $(GOVC_INSTALL) $(GOVC_BIN) $(GOVC_VER)

$(KIND): # Build kind.
	cd $(TEST_DIR); GOBIN=$(TOOLS_BIN_DIR) ../$(GO_INSTALL) $(KIND_PKG) $(KIND_BIN) $(KIND_VER)

$(RELEASE_NOTES): # Build release-notes.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_TOOLS_BUILD) $(RELEASE_NOTES_PKG) $(RELEASE_NOTES_BIN) $(RELEASE_NOTES_VER)


## --------------------------------------
## govulncheck
## --------------------------------------

$(GOVULNCHECK): # Build govulncheck.
	@if [ -z "${GOVULNCHECK_VER}" ]; then echo "GOVULNCHECK_VER is not set"; exit 1; fi
	@if [ -d "$(GOVULNCHECK_TMP_DIR)" ]; then \
		echo "$(GOVULNCHECK_TMP_DIR) exists, skipping clone"; \
	else \
		git clone "https://github.com/golang/vuln.git" "$(GOVULNCHECK_TMP_DIR)"; \
		cd "$(GOVULNCHECK_TMP_DIR)"; \
		git checkout "$(GOVULNCHECK_VER)"; \
		git apply "$(ROOT_DIR)/$(GOVULNCHECK_DIR)/govulncheck.patch"; \
	fi
	@cd "$(ROOT_DIR)/$(GOVULNCHECK_TMP_DIR)"; \
	if [ "$$(git describe --tag 2> /dev/null)" != "$(GOVULNCHECK_VER)" ]; then \
		echo "ERROR: checked out version $$(git describe --tag 2> /dev/null) does not match expected version $(GOVULNCHECK_VER)"; \
		exit 1; \
	fi
	go build -C $(GOVULNCHECK_TMP_DIR) -o $(TOOLS_BIN_DIR)/$(GOVULNCHECK_BIN) ./cmd/govulncheck

.PHONY: clean-govulncheck
clean-govulncheck:
	rm -fr "$(GOVULNCHECK_TMP_DIR)"

## --------------------------------------
## Helpers
## --------------------------------------

##@ helpers:

go-version: ## Print the go version we use to compile our binaries and images
	@echo $(GO_VERSION)
