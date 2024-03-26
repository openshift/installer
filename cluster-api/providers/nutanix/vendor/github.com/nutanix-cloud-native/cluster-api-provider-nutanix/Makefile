SHELL := /bin/bash
GOCMD=go
GOTEST=$(GOCMD) test
GOINSTALL=$(GOCMD) install
GOTOOL=$(GOCMD) tool
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true
# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/nutanix-cloud-native/cluster-api-provider-nutanix/controller:latest

# Extract base and tag from IMG
IMG_REPO ?= $(word 1,$(subst :, ,${IMG}))
IMG_TAG ?= $(word 2,$(subst :, ,${IMG}))
LOCAL_PROVIDER_VERSION ?= ${IMG_TAG}
ifeq (${IMG_TAG},)
IMG_TAG := latest
endif

ifeq (${LOCAL_PROVIDER_VERSION},latest)
# Change this versions after release when required here and in e2e config (test/e2e/config/nutanix.yaml)
LOCAL_PROVIDER_VERSION := v1.3.99
endif

# PLATFORMS is a list of platforms to build for.
PLATFORMS ?= linux/amd64,linux/arm64,linux/arm
PLATFORMS_E2E ?= linux/amd64

# KIND_CLUSTER_NAME is the name of the kind cluster to use.
KIND_CLUSTER_NAME ?= capi-test

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.23

#
# Directories.
#
# Full directory of where the Makefile resides
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
REPO_ROOT := $(shell git rev-parse --show-toplevel)
EXP_DIR := exp
BIN_DIR := bin
TEST_DIR := test
E2E_DIR ?= ${REPO_ROOT}/test/e2e
TEMPLATES_DIR := templates
TOOLS_DIR := $(REPO_ROOT)/hack/tools
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/$(BIN_DIR))
E2E_FRAMEWORK_DIR := $(TEST_DIR)/framework
CAPD_DIR := $(TEST_DIR)/infrastructure/docker
GO_INSTALL := $(REPO_ROOT)/scripts/go_install.sh
NUTANIX_E2E_TEMPLATES := ${E2E_DIR}/data/infrastructure-nutanix
RELEASE_DIR ?= $(REPO_ROOT)/out

export PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)

# CNI paths for e2e tests
CNI_PATH_CALICO ?= "${E2E_DIR}/data/cni/calico/calico.yaml"
CNI_PATH_FLANNEL ?= "${E2E_DIR}/data/cni/flannel/flannel.yaml" # From https://github.com/flannel-io/flannel/blob/master/Documentation/kube-flannel.yml
CNI_PATH_CILIUM ?= "${E2E_DIR}/data/cni/cilium/cilium.yaml" # helm template cilium cilium/cilium --version 1.13.0 -n kube-system --set hubble.enabled=false --set cni.chainingMode=portmap  --set sessionAffinity=true | sed 's/${BIN_PATH}/$BIN_PATH/g'
CNI_PATH_CILIUM_NO_KUBEPROXY ?= "${E2E_DIR}/data/cni/cilium/cilium-no-kubeproxy.yaml" # helm template cilium cilium/cilium --version 1.13.0 -n kube-system --set hubble.enabled=false --set cni.chainingMode=portmap  --set sessionAffinity=true --set kubeProxyReplacement=strict | sed 's/${BIN_PATH}/$BIN_PATH/g'

#
# Binaries.
#
# Note: Need to use abspath so we can invoke these from subdirectories
KO_VER := v0.11.2
KO_BIN := ko
KO := $(abspath $(TOOLS_BIN_DIR)/$(KO_BIN)-$(KO_VER))
KO_PKG := github.com/google/ko

KUSTOMIZE_BIN := kustomize
KUSTOMIZE_VER := v4.5.4
KUSTOMIZE := $(abspath $(TOOLS_BIN_DIR)/$(KUSTOMIZE_BIN)-$(KUSTOMIZE_VER))
KUSTOMIZE_PKG := sigs.k8s.io/kustomize/kustomize/v4

GINGKO_VER := v2.1.4
GINKGO_BIN := ginkgo
GINKGO := $(abspath $(TOOLS_BIN_DIR)/$(GINKGO_BIN)-$(GINGKO_VER))
GINKGO_PKG := github.com/onsi/ginkgo/v2/ginkgo

SETUP_ENVTEST_VER := latest
SETUP_ENVTEST_BIN := setup-envtest
SETUP_ENVTEST := $(abspath $(TOOLS_BIN_DIR)/$(SETUP_ENVTEST_BIN)-$(SETUP_ENVTEST_VER))
SETUP_ENVTEST_PKG := sigs.k8s.io/controller-runtime/tools/setup-envtest

CONTROLLER_GEN_VER := v0.14.0
CONTROLLER_GEN_BIN := controller-gen
CONTROLLER_GEN := $(abspath $(TOOLS_BIN_DIR)/$(CONTROLLER_GEN_BIN)-$(CONTROLLER_GEN_VER))
CONTROLLER_GEN_PKG := sigs.k8s.io/controller-tools/cmd/controller-gen

GOTESTSUM_VER := v1.6.4
GOTESTSUM_BIN := gotestsum
GOTESTSUM := $(abspath $(TOOLS_BIN_DIR)/$(GOTESTSUM_BIN)-$(GOTESTSUM_VER))
GOTESTSUM_PKG := gotest.tools/gotestsum

CONVERSION_GEN_VER := v0.23.6
CONVERSION_GEN_BIN := conversion-gen
# We are intentionally using the binary without version suffix, to avoid the version
# in generated files.
CONVERSION_GEN := $(abspath $(TOOLS_BIN_DIR)/$(CONVERSION_GEN_BIN))
CONVERSION_GEN_PKG := k8s.io/code-generator/cmd/conversion-gen

ENVSUBST_VER := v2.0.0-20210730161058-179042472c46
ENVSUBST_BIN := envsubst
ENVSUBST := $(abspath $(TOOLS_BIN_DIR)/$(ENVSUBST_BIN)-$(ENVSUBST_VER))
ENVSUBST_PKG := github.com/drone/envsubst/v2/cmd/envsubst

GO_APIDIFF_VER := v0.1.0
GO_APIDIFF_BIN := go-apidiff
GO_APIDIFF := $(abspath $(TOOLS_BIN_DIR)/$(GO_APIDIFF_BIN)-$(GO_APIDIFF_VER))
GO_APIDIFF_PKG := github.com/joelanford/go-apidiff

KPROMO_VER := v3.3.0-beta.3
KPROMO_BIN := kpromo
KPROMO :=  $(abspath $(TOOLS_BIN_DIR)/$(KPROMO_BIN)-$(KPROMO_VER))
KPROMO_PKG := sigs.k8s.io/promo-tools/v3/cmd/kpromo

CONVERSION_VERIFIER_BIN := conversion-verifier
CONVERSION_VERIFIER := $(abspath $(TOOLS_BIN_DIR)/$(CONVERSION_VERIFIER_BIN))

TILT_PREPARE_BIN := tilt-prepare
TILT_PREPARE := $(abspath $(TOOLS_BIN_DIR)/$(TILT_PREPARE_BIN))

GOLANGCI_LINT_VER := v1.55.2
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT := $(abspath $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_BIN))

# Install clusterctl that corresponds to the cluster-api go mod version
CLUSTERCTL_VER := $(shell go list -m sigs.k8s.io/cluster-api | cut -d" " -f2)
CLUSTERCTL_RELEASE_URL := https://github.com/kubernetes-sigs/cluster-api/releases/download/$(CLUSTERCTL_VER)/clusterctl-$(shell go env GOOS)-$(shell go env GOARCH)
CLUSTERCTL_BIN := clusterctl
CLUSTERCTL := $(abspath $(TOOLS_BIN_DIR)/$(CLUSTERCTL_BIN))

# CRD_OPTIONS define options to add to the CONTROLLER_GEN
CRD_OPTIONS ?= "crd:crdVersions=v1"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Get latest git hash
GIT_COMMIT_HASH=$(shell git rev-parse HEAD)

# Get the local image registry required for clusterctl upgrade tests
LOCAL_IMAGE_REGISTRY ?= localhost:5000

ifeq (${MAKECMDGOALS},test-e2e-clusterctl-upgrade)
	IMG_TAG=e2e-${GIT_COMMIT_HASH}
	IMG_REPO=${LOCAL_IMAGE_REGISTRY}/controller
endif

ifeq (${MAKECMDGOALS},docker-build-e2e)
	IMG_TAG=e2e-${GIT_COMMIT_HASH}
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

RUN_VALIDATION_TESTS_ONLY?=false # for CI please set EXPORT_RESULT to true
LABEL_FILTERS ?=
ifeq ($(RUN_VALIDATION_TESTS_ONLY), true)
		LABEL_FILTER_ARGS = only-for-validation
else
		LABEL_FILTER_ARGS = !only-for-validation
endif
ifneq ($(LABEL_FILTERS),)
		LABEL_FILTER_ARGS := "$(LABEL_FILTER_ARGS) && $(LABEL_FILTERS)"
endif
JUNIT_REPORT_FILE ?= "junit.e2e_suite.1.xml"
GINKGO_SKIP ?= "clusterctl-Upgrade"
GINKGO_FOCUS ?= ""
GINKGO_NODES  ?= 1
E2E_CONF_FILE  ?= ${E2E_DIR}/config/nutanix.yaml
ARTIFACTS ?= ${REPO_ROOT}/_artifacts
SKIP_RESOURCE_CLEANUP ?= false
USE_EXISTING_CLUSTER ?= false
GINKGO_NOCOLOR ?= false
FLAVOR ?= e2e

TEST_NAMESPACE=capx-test-ns
TEST_CLUSTER_NAME=mycluster

# set ginkgo focus flags, if any
ifneq ($(strip $(GINKGO_FOCUS)),)
_FOCUS_ARGS := $(foreach arg,$(strip $(GINKGO_FOCUS)),--focus="$(arg)")
endif

# to set multiple ginkgo skip flags, if any
ifneq ($(strip $(GINKGO_SKIP)),)
_SKIP_ARGS := $(foreach arg,$(strip $(GINKGO_SKIP)),--skip="$(arg)")
endif
.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Development

.PHONY: manifests
manifests: $(CONTROLLER_GEN_BIN) ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: release-manifests
release-manifests: manifests cluster-templates
	mkdir -p $(RELEASE_DIR)
	$(KUSTOMIZE) build config/default > $(RELEASE_DIR)/infrastructure-components.yaml
	cp $(TEMPLATES_DIR)/cluster-template*.yaml $(RELEASE_DIR)
	cp $(REPO_ROOT)/metadata.yaml $(RELEASE_DIR)/metadata.yaml

.PHONY: generate
generate: $(CONTROLLER_GEN_BIN) $(CONVERSION_GEN_BIN) ## Generate code containing DeepCopy, DeepCopyInto, DeepCopyObject method implementations and API conversion implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

	$(CONVERSION_GEN) \
	--input-dirs=./api/v1alpha4 \
	--input-dirs=./api/v1beta1 \
	--build-tag=ignore_autogenerated_core \
	--output-file-base=zz_generated.conversion \
	--go-header-file=./hack/boilerplate.go.txt

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

kind-create: ## Create a kind cluster and deploy the latest supported cluster API version
	kind create cluster --name=${KIND_CLUSTER_NAME}

kind-delete: ## Delete the kind cluster
	kind delete cluster --name=${KIND_CLUSTER_NAME}

##@ Build

.PHONY: build
build: generate fmt ## Build manager binary.
	echo "Git commit hash: ${GIT_COMMIT_HASH}"
	go build -ldflags "-X main.gitCommitHash=${GIT_COMMIT_HASH}" -o bin/manager main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go

.PHONY: docker-build
docker-build: $(KO) ## Build docker image with the manager.
	echo "Git commit hash: ${GIT_COMMIT_HASH}"
	KO_DOCKER_REPO=ko.local GOFLAGS="-ldflags=-X=main.gitCommitHash=${GIT_COMMIT_HASH}" $(KO) build -B --platform=${PLATFORMS} -t ${IMG_TAG} -L .

.PHONY: docker-push
docker-push: $(KO) ## Push docker image with the manager.
	KO_DOCKER_REPO=${IMG_REPO} $(KO) build --bare --platform=${PLATFORMS} -t ${IMG_TAG} .

.PHONY: docker-push-kind
docker-push-kind: $(KO) ## Make docker image available to kind cluster.
	GOOS=linux GOARCH=${shell go env GOARCH} KO_DOCKER_REPO=ko.local ${KO} build -B -t ${IMG_TAG} -L .
	docker tag ko.local/cluster-api-provider-nutanix:${IMG_TAG} ${IMG}
	kind load docker-image --name ${KIND_CLUSTER_NAME} ${IMG}

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif


.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests kustomize docker-push-kind $(CLUSTERCTL) ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	$(CLUSTERCTL) delete --infrastructure nutanix:${LOCAL_PROVIDER_VERSION} --include-crd || true
	$(CLUSTERCTL) init --infrastructure nutanix:${LOCAL_PROVIDER_VERSION} -v 9
	# cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	# $(KUSTOMIZE) build config/default | kubectl apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

##@ Templates

.PHONY: cluster-e2e-templates
cluster-e2e-templates: $(KUSTOMIZE) cluster-e2e-templates-v1beta1 cluster-e2e-templates-v1alpha4 cluster-e2e-templates-v124 ## Generate cluster templates for all versions

cluster-e2e-templates-v124: $(KUSTOMIZE) ## Generate cluster templates for CAPX v1.2.4
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1.2.4/cluster-template --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1.2.4/cluster-template.yaml

cluster-e2e-templates-v1alpha4: $(KUSTOMIZE) ## Generate cluster templates for v1alpha4
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1alpha4/cluster-template --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1alpha4/cluster-template.yaml

cluster-e2e-templates-v1beta1: $(KUSTOMIZE) ## Generate cluster templates for v1beta1
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template.yaml 
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-secret --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-secret.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nutanix-cluster --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nutanix-cluster.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-additional-categories --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-additional-categories.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nmt --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nmt.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-project --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-project.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-upgrades --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-upgrades.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-md-remediation --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-md-remediation.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-remediation --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-remediation.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-scale-in --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-scale-in.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-csi --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-csi.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-failure-domains --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-failure-domains.yaml

cluster-e2e-templates-no-kubeproxy: $(KUSTOMIZE) ##Generate cluster templates without kubeproxy
	# v1alpha4
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1alpha4/no-kubeproxy/cluster-template --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1alpha4/cluster-template.yaml

	# v1beta1
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-no-secret --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-secret.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-no-nutanix-cluster --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nutanix-cluster.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-additional-categories --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-additional-categories.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-no-nmt --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nmt.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-project --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-project.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-upgrades --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-upgrades.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-md-remediation --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-md-remediation.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-kcp-remediation --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-remediation.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-kcp-scale-in --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-scale-in.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-csi --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-csi.yaml
	$(KUSTOMIZE) build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-failure-domains --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-failure-domains.yaml

cluster-templates: $(KUSTOMIZE) ## Generate cluster templates for all flavors
	$(KUSTOMIZE) build $(TEMPLATES_DIR)/base > $(TEMPLATES_DIR)/cluster-template.yaml
	$(KUSTOMIZE) build $(TEMPLATES_DIR)/csi > $(TEMPLATES_DIR)/cluster-template-csi.yaml

##@ Testing

.PHONY: docker-build-e2e
docker-build-e2e: $(KO) ## Build docker image with the manager with e2e tag.
	echo "Git commit hash: ${GIT_COMMIT_HASH}"
	KO_DOCKER_REPO=ko.local GOFLAGS="-ldflags=-X=main.gitCommitHash=${GIT_COMMIT_HASH}" $(KO) build -B --platform=${PLATFORMS_E2E} -t ${IMG_TAG} -L .
	docker tag ko.local/cluster-api-provider-nutanix:${IMG_TAG} ${IMG_REPO}:e2e

.PHONY: prepare-local-clusterctl
prepare-local-clusterctl: manifests kustomize cluster-templates envsubst ## Prepare overide file for local clusterctl.
	mkdir -p ~/.cluster-api/overrides/infrastructure-nutanix/${LOCAL_PROVIDER_VERSION}
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default > ~/.cluster-api/overrides/infrastructure-nutanix/${LOCAL_PROVIDER_VERSION}/infrastructure-components.yaml
	cp ./metadata.yaml ~/.cluster-api/overrides/infrastructure-nutanix/${LOCAL_PROVIDER_VERSION}/
	cp ./templates/cluster-template*.yaml ~/.cluster-api/overrides/infrastructure-nutanix/${LOCAL_PROVIDER_VERSION}/
	env LOCAL_PROVIDER_VERSION=$(LOCAL_PROVIDER_VERSION) \
		$(ENVSUBST) -no-unset -no-empty -no-digit < ./clusterctl.yaml > ~/.cluster-api/clusterctl.yaml

.PHONY: unit-test
unit-test: setup-envtest ## Run unit tests.
ifeq ($(EXPORT_RESULT), true)
	$(GOINSTALL) github.com/jstemmer/go-junit-report/v2@latest
	$(eval OUTPUT_OPTIONS = | go-junit-report -set-exit-code > junit-report.xml)
endif
	KUBEBUILDER_ASSETS="$(shell $(SETUP_ENVTEST) use $(ENVTEST_K8S_VERSION)  --arch=amd64 -p path)" $(GOTEST) ./... $(OUTPUT_OPTIONS)

.PHONY: coverage
coverage: setup-envtest ## Run the tests of the project and export the coverage
	KUBEBUILDER_ASSETS="$(shell $(SETUP_ENVTEST) use $(ENVTEST_K8S_VERSION) --arch=amd64 -p path)" $(GOTEST) -cover -covermode=count -coverprofile=profile.cov -coverpkg=./... ./...
	$(GOTOOL) cover -func profile.cov
ifeq ($(EXPORT_RESULT), true)
	$(GOINSTALL) github.com/AlekSi/gocov-xml@latest
	$(GOINSTALL) github.com/axw/gocov/gocov@latest
	gocov convert profile.cov | gocov-xml > coverage.xml
endif

.PHONY: test-clusterctl-create
test-clusterctl-create: $(CLUSTERCTL) ## Run the tests using clusterctl
	$(CLUSTERCTL) version
	$(CLUSTERCTL) config repositories | grep nutanix
	$(CLUSTERCTL) generate cluster ${TEST_CLUSTER_NAME} -i nutanix:${LOCAL_PROVIDER_VERSION} --list-variables -v 10
	$(CLUSTERCTL) generate cluster ${TEST_CLUSTER_NAME} -i nutanix:${LOCAL_PROVIDER_VERSION} --target-namespace ${TEST_NAMESPACE}  -v 10 > ./cluster.yaml
	kubectl create ns $(TEST_NAMESPACE) || true
	kubectl apply -f ./cluster.yaml -n $(TEST_NAMESPACE)

.PHONY: test-clusterctl-delete
test-clusterctl-delete: ## Delete clusterctl created cluster
	kubectl -n ${TEST_NAMESPACE} delete cluster ${TEST_CLUSTER_NAME}

.PHONY: test-kubectl-bootstrap
test-kubectl-bootstrap: ## Run kubectl queries to get all capx management/bootstrap related objects
	kubectl get ns
	kubectl get all --all-namespaces
	kubectl -n capx-system get all
	kubectl -n $(TEST_NAMESPACE) get Cluster,NutanixCluster,Machine,NutanixMachine,KubeAdmControlPlane,MachineHealthCheck,nodes
	kubectl -n capx-system get pod

.PHONY: test-kubectl-workload
test-kubectl-workload: ## Run kubectl queries to get all capx workload related objects
	kubectl -n $(TEST_NAMESPACE) get secret
	kubectl -n ${TEST_NAMESPACE} get secret ${TEST_CLUSTER_NAME}-kubeconfig -o json | jq -r .data.value | base64 --decode > ${TEST_CLUSTER_NAME}.workload.kubeconfig
	kubectl --kubeconfig ./${TEST_CLUSTER_NAME}.workload.kubeconfig get nodes,ns

.PHONY: ginkgo-help
ginkgo-help:
	$(GINKGO) help run

.PHONY: test-e2e
test-e2e: docker-build-e2e $(GINKGO_BIN) cluster-e2e-templates cluster-templates ## Run the end-to-end tests
	mkdir -p $(ARTIFACTS)
	NUTANIX_LOG_LEVEL=debug $(GINKGO) -v \
		--trace \
		--progress \
		--tags=e2e \
		--label-filter=$(LABEL_FILTER_ARGS) \
		$(_SKIP_ARGS) \
		$(_FOCUS_ARGS) \
		--nodes=$(GINKGO_NODES) \
		--no-color=$(GINKGO_NOCOLOR) \
		--output-dir="$(ARTIFACTS)" \
		--junit-report=${JUNIT_REPORT_FILE} \
		--timeout="24h" \
		--always-emit-ginkgo-writer \
		$(GINKGO_ARGS) ./test/e2e -- \
		-e2e.artifacts-folder="$(ARTIFACTS)" \
		-e2e.config="$(E2E_CONF_FILE)" \
		-e2e.skip-resource-cleanup=$(SKIP_RESOURCE_CLEANUP) \
		-e2e.use-existing-cluster=$(USE_EXISTING_CLUSTER)

.PHONY: test-e2e-no-kubeproxy
test-e2e-no-kubeproxy: docker-build-e2e $(GINKGO_BIN) cluster-e2e-templates-no-kubeproxy cluster-templates ## Run the end-to-end tests without kubeproxy
	mkdir -p $(ARTIFACTS)
	NUTANIX_LOG_LEVEL=debug $(GINKGO) -v \
		--trace \
		--progress \
		--tags=e2e \
		--label-filter=$(LABEL_FILTER_ARGS) \
		$(_SKIP_ARGS) \
		--nodes=$(GINKGO_NODES) \
		--no-color=$(GINKGO_NOCOLOR) \
		--output-dir="$(ARTIFACTS)" \
		--junit-report=${JUNIT_REPORT_FILE} \
		--timeout="24h" \
		--always-emit-ginkgo-writer \
		$(GINKGO_ARGS) ./test/e2e -- \
		-e2e.artifacts-folder="$(ARTIFACTS)" \
		-e2e.config="$(E2E_CONF_FILE)" \
		-e2e.skip-resource-cleanup=$(SKIP_RESOURCE_CLEANUP) \
		-e2e.use-existing-cluster=$(USE_EXISTING_CLUSTER)

.PHONY: list-e2e
list-e2e: docker-build-e2e $(GINKGO_BIN) cluster-e2e-templates cluster-templates ## Run the end-to-end tests
	mkdir -p $(ARTIFACTS)
	$(GINKGO) -v --trace --dry-run --tags=e2e --label-filter="$(LABEL_FILTERS)" $(_SKIP_ARGS) --nodes=$(GINKGO_NODES) \
	    --no-color=$(GINKGO_NOCOLOR) --output-dir="$(ARTIFACTS)" \
	    $(GINKGO_ARGS) ./test/e2e -- \
	    -e2e.artifacts-folder="$(ARTIFACTS)" \
	    -e2e.config="$(E2E_CONF_FILE)" \
	    -e2e.skip-resource-cleanup=$(SKIP_RESOURCE_CLEANUP) -e2e.use-existing-cluster=$(USE_EXISTING_CLUSTER)

.PHONY: test-e2e-calico
test-e2e-calico:
	CNI=$(CNI_PATH_CALICO) $(MAKE) test-e2e

.PHONY: test-e2e-flannel
test-e2e-flannel:
	CNI=$(CNI_PATH_FLANNEL) $(MAKE) test-e2e

.PHONY: test-e2e-cilium
test-e2e-cilium:
	CNI=$(CNI_PATH_CILIUM) $(MAKE) test-e2e

.PHONY: test-e2e-cilium-no-kubeproxy
test-e2e-cilium-no-kubeproxy:
	CNI=$(CNI_PATH_CILIUM_NO_KUBEPROXY) $(MAKE) test-e2e-no-kubeproxy

.PHONY: test-e2e-all-cni
test-e2e-all-cni: test-e2e test-e2e-calico test-e2e-flannel test-e2e-cilium test-e2e-cilium-no-kubeproxy

.PHONY: test-e2e-clusterctl-upgrade
test-e2e-clusterctl-upgrade: docker-build-e2e $(GINKGO_BIN) cluster-e2e-templates cluster-templates ## Run the end-to-end tests
	echo "Image tag for E2E test is ${IMG_TAG}"
	docker tag ko.local/cluster-api-provider-nutanix:${IMG_TAG} ${IMG_REPO}:${IMG_TAG}
	docker push ${IMG_REPO}:${IMG_TAG}
	GINKGO_SKIP="" GIT_COMMIT="${GIT_COMMIT_HASH}" $(MAKE) test-e2e-calico

## --------------------------------------
## Hack / Tools
## --------------------------------------

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

.PHONY: $(TILT_PREPARE_BIN)
$(TILT_PREPARE_BIN): $(TILT_PREPARE) ## Build a local copy of tilt-prepare.

.PHONY: $(GOLANGCI_LINT_BIN)
$(GOLANGCI_LINT_BIN): $(GOLANGCI_LINT) ## Build a local copy of golangci-lint

.PHONY: $(CLUSTERCTL_BIN)
$(CLUSTERCTL_BIN): $(CLUSTERCTL) ## Build a local copy of clusterctl

$(GINKGO): # Build ginkgo from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GINKGO_PKG) $(GINKGO_BIN) $(GINGKO_VER)

$(KO): # Build ko from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(KO_PKG) $(KO_BIN) $(KO_VER)

$(KUSTOMIZE): # Build kustomize from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(KUSTOMIZE_PKG) $(KUSTOMIZE_BIN) $(KUSTOMIZE_VER)

.PHONY: $(KO_BIN)
$(KO_BIN): $(KO) ## Build a local copy of ko

.PHONY: $(GINKGO_BIN)
$(GINKGO_BIN): $(GINKGO) ## Build a local copy of ginkgo

.PHONY: $(KUSTOMIZE_BIN)
$(KUSTOMIZE_BIN): $(KUSTOMIZE) ## Build a local copy of kustomize

$(CONTROLLER_GEN): # Build controller-gen from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(CONTROLLER_GEN_PKG) $(CONTROLLER_GEN_BIN) $(CONTROLLER_GEN_VER)

## We are forcing a rebuilt of conversion-gen via PHONY so that we're always using an up-to-date version.
## We can't use a versioned name for the binary, because that would be reflected in generated files.
.PHONY: $(CONVERSION_GEN)
$(CONVERSION_GEN): # Build conversion-gen from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(CONVERSION_GEN_PKG) $(CONVERSION_GEN_BIN) $(CONVERSION_GEN_VER)

$(CONVERSION_VERIFIER): $(TOOLS_DIR)/go.mod # Build conversion-verifier from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/conversion-verifier sigs.k8s.io/cluster-api/hack/tools/conversion-verifier

$(GOTESTSUM): # Build gotestsum from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GOTESTSUM_PKG) $(GOTESTSUM_BIN) $(GOTESTSUM_VER)

$(GO_APIDIFF): # Build go-apidiff from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(GO_APIDIFF_PKG) $(GO_APIDIFF_BIN) $(GO_APIDIFF_VER)

$(ENVSUBST): # Build gotestsum from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(ENVSUBST_PKG) $(ENVSUBST_BIN) $(ENVSUBST_VER)

$(SETUP_ENVTEST): # Build setup-envtest from tools folder.
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(SETUP_ENVTEST_PKG) $(SETUP_ENVTEST_BIN) $(SETUP_ENVTEST_VER)

$(TILT_PREPARE): $(TOOLS_DIR)/go.mod # Build tilt-prepare from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/tilt-prepare sigs.k8s.io/cluster-api/hack/tools/tilt-prepare

$(KPROMO):
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) $(KPROMO_PKG) $(KPROMO_BIN) ${KPROMO_VER}

$(GOLANGCI_LINT): # building golanci-lint from source is not recommended, so we are using the install script
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_BIN_DIR) $(GOLANGCI_LINT_VER)

$(CLUSTERCTL):
# We don't install clusterctl using the go toolchain, because the upstream Makefile
# is required to build clusterctl correctly. See https://github.com/kubernetes-sigs/cluster-api/issues/3706
	curl -sSfL -o $(CLUSTERCTL) $(CLUSTERCTL_RELEASE_URL)
	chmod u+x $(CLUSTERCTL)

## --------------------------------------
## Lint / Verify
## --------------------------------------

##@ Lint and Verify

GOLANGCI_LINT_EXTRA_ARGS := --enable gofumpt --build-tags e2e

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint the codebase
	$(GOLANGCI_LINT) run -v $(GOLANGCI_LINT_EXTRA_ARGS)

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT) ## Lint the codebase and run auto-fixers if supported by the linter
	GOLANGCI_LINT_EXTRA_ARGS="$(GOLANGCI_LINT_EXTRA_ARGS) --fix" $(MAKE) lint

# Make any new verify task a dependency of this target
.PHONY: verify ## Run all verify targets
verify: verify-manifests

# Adapted from https://github.com/kubernetes-sigs/cluster-api/blob/f15e8769a2135429911ce6d8f7124b853c0444a1/Makefile#L651-L656
.PHONY: verify-manifests
verify-manifests: manifests  ## Verify generated manifests are up to date
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated manifests are out of date in the repository, run make manifests, and commit"; exit 1; \
	fi

## --------------------------------------
## Clean
## --------------------------------------

##@ Clean
.PHONY: clean
clean: ## Clean the build and test artifacts
	rm -rf $(ARTIFACTS) $(BIN_DIR)

