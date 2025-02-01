SHELL := bash
GOCMD=go
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOTOOL=$(GOCMD) tool
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true

GIT_COMMIT_HASH=$(shell git rev-parse HEAD)
LOCAL_IMAGE_REGISTRY ?= ko.local
IMG_REPO=${LOCAL_IMAGE_REGISTRY}/cluster-api-provider-nutanix
IMG_TAG=e2e-${GIT_COMMIT_HASH}
MANAGER_IMAGE=${IMG_REPO}:${IMG_TAG}

# Extract base and tag from IMG
LOCAL_PROVIDER_VERSION ?= ${IMG_TAG}

ifeq (${LOCAL_PROVIDER_VERSION},${IMG_TAG})
# TODO(release-blocker): Change this versions after release when required here
LOCAL_PROVIDER_VERSION := v1.5.99
endif

# PLATFORMS is a list of platforms to build for.
PLATFORMS ?= linux/amd64,linux/arm64,linux/arm
PLATFORMS_E2E ?= linux/amd64

# KIND_CLUSTER_NAME is the name of the kind cluster to use.
KIND_CLUSTER_NAME ?= capi-test

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.29

#
# Directories.
#
# Full directory of where the Makefile resides
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
REPO_ROOT := $(shell git rev-parse --show-toplevel)
EXP_DIR := exp
TEST_DIR := test
E2E_DIR ?= ${REPO_ROOT}/test/e2e
TEMPLATES_DIR := templates
E2E_FRAMEWORK_DIR := $(TEST_DIR)/framework
CAPD_DIR := $(TEST_DIR)/infrastructure/docker
GO_INSTALL := $(REPO_ROOT)/scripts/go_install.sh
NUTANIX_E2E_TEMPLATES := ${E2E_DIR}/data/infrastructure-nutanix
RELEASE_DIR ?= $(REPO_ROOT)/out

# CNI paths for e2e tests
CNI_PATH_CALICO ?= "${E2E_DIR}/data/cni/calico/calico.yaml"
CNI_PATH_CILIUM ?= "${E2E_DIR}/data/cni/cilium/cilium.yaml"
CNI_PATH_CILIUM_NO_KUBEPROXY ?= "${E2E_DIR}/data/cni/cilium/cilium-no-kubeproxy.yaml"
CNI_PATH_FLANNEL ?= "${E2E_DIR}/data/cni/flannel/flannel.yaml"
CNI_PATH_KINDNET ?= "${E2E_DIR}/data/cni/kindnet/kindnet.yaml"
CCM_VERSION ?= v0.4.1

# CRD_OPTIONS define options to add to the CONTROLLER_GEN
CRD_OPTIONS ?= "crd:crdVersions=v1"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
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
GINKGO_SKIP ?=
GINKGO_FOCUS ?=
GINKGO_NODES  ?= 1
E2E_CONF_FILE  ?= ${E2E_DIR}/config/nutanix.yaml
E2E_CONF_FILE_TMP = ${E2E_CONF_FILE}.tmp
ARTIFACTS ?= ${REPO_ROOT}/_artifacts
SKIP_RESOURCE_CLEANUP ?= false
USE_EXISTING_CLUSTER ?= false
GINKGO_NOCOLOR ?= false
FLAVOR ?= e2e

define ginkgo_option
--$(1)="$(shell echo '$(2)' | sed -E 's/^[[:space:]]+//' | sed -E 's/"[[:space:]]+"/" --$(1)="/g')"
endef

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
manifests: ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	controller-gen $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: release-manifests
release-manifests: manifests cluster-templates
	mkdir -p $(RELEASE_DIR)
	kustomize build config/default > $(RELEASE_DIR)/infrastructure-components.yaml
	cp $(TEMPLATES_DIR)/cluster-template*.yaml $(RELEASE_DIR)
	cp $(REPO_ROOT)/metadata.yaml $(RELEASE_DIR)/metadata.yaml

.PHONY: generate
generate: ## Generate code containing DeepCopy, DeepCopyInto, DeepCopyObject method implementations and API conversion implementations.
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
	conversion-gen \
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

.PHONY: nutanix-cp-endpoint-ip
nutanix-cp-endpoint-ip: ## Gets a random free IP from the control plane endpoint range set in the environment.
	@shuf --head-count=1 < <(fping -g -u "$(CONTROL_PLANE_ENDPOINT_RANGE_START)" "$(CONTROL_PLANE_ENDPOINT_RANGE_END)")

.PHONY: update-calico-cni
update-calico-cni: ## Updates the calico CNI manifests
	@echo "Updating calico CNI manifest"
	@curl -sL https://docs.projectcalico.org/manifests/calico.yaml -o $(CNI_PATH_CALICO)
	# replace all docker.io images with quay.io images
	@sed -i 's|docker.io|quay.io|g' $(CNI_PATH_CALICO)

.PHONY: update-cilium-cni
update-cilium-cni: ## Updates the cilium CNI manifests
	@echo "Updating cilium CNI manifest"
	@helm repo add cilium https://helm.cilium.io/
	@helm repo update
	# TODO(use the latest version of cilium instead of hardcoding v1.15.4 once fix for 1.15.5+ failure is identified)
	# as 1.15.5+ fails due to mount-cgroup init container failing with nsenter: cannot open /hostproc/1/ns/cgroup: Permission denied
	# Also, add k8s-service-proxy-name: "cilium" to cilium-config ConfigMap to ensure multi-protocol sig-network conformance tests pass
	@helm template cilium cilium/cilium --version 1.15.4 -n kube-system --set hubble.enabled=false --set cni.chainingMode=portmap --set sessionAffinity=true --set k8s.serviceProxyName=cilium | awk '{gsub(/\$\{BIN_PATH\}/,"$ BIN_PATH"); print}' > $(CNI_PATH_CILIUM)
	@helm template cilium cilium/cilium --version 1.15.4 -n kube-system --set hubble.enabled=false --set cni.chainingMode=portmap --set sessionAffinity=true --set kubeProxyReplacement=strict --set k8s.serviceProxyName=cilium | awk '{gsub(/\$\{BIN_PATH\}/,"$ BIN_PATH"); print}' > $(CNI_PATH_CILIUM_NO_KUBEPROXY)

.PHONY: update-flannel-cni
update-flannel-cni: ## Updates the flannel CNI manifests
	@echo "Updating flannel CNI manifest"
	@curl -sL https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml -o $(CNI_PATH_FLANNEL)
	@echo "Updating flannel CNI net-conf manifest to set CIDR"
	@yq eval-all 'select(.kind != "ConfigMap" and .metadata.name != "kube-flannel-cfg")' $(CNI_PATH_FLANNEL) > $(CNI_PATH_FLANNEL).tmp
	@echo "---" >> $(CNI_PATH_FLANNEL).tmp
	@FLANNEL_NET_CONF=$$(yq eval 'select(.kind == "ConfigMap" and .metadata.name == "kube-flannel-cfg") | .data."net-conf.json" | fromjson | .Network = "172.20.0.0/16" | tojson' $(CNI_PATH_FLANNEL)) yq eval 'select(.kind == "ConfigMap" and .metadata.name == "kube-flannel-cfg") | .data."net-conf.json" = strenv(FLANNEL_NET_CONF)' $(CNI_PATH_FLANNEL) >> $(CNI_PATH_FLANNEL).tmp
	@mv $(CNI_PATH_FLANNEL).tmp $(CNI_PATH_FLANNEL)

.PHONY: update-kindnet-cni
update-kindnet-cni: ## Updates the kindnet CNI manifests
	@echo "Updating kindnet CNI manifest"
	@curl -sL https://github.com/kubernetes-sigs/cluster-api/raw/main/test/e2e/data/cni/kindnet/kindnet.yaml -o $(CNI_PATH_KINDNET)

.PHONY: update-ccm
update-ccm: ## Updates the Nutanix CCM tag in all the template manifests to CCM_VERSION
	@echo "Updating Nutanix CCM tag"
	@find $(TEMPLATES_DIR) -type f -exec sed -i 's|CCM_TAG=[^}]*|CCM_TAG=$(CCM_VERSION)|g' {} +
	@find $(NUTANIX_E2E_TEMPLATES) -type f -exec sed -i 's|CCM_TAG=[^}]*|CCM_TAG=$(CCM_VERSION)|g' {} +

.PHONY: update-cni-manifests ## Updates all the CNI manifests to latest variants from upstream
update-cni-manifests: update-calico-cni update-cilium-cni update-flannel-cni update-kindnet-cni  ## Updates all the CNI manifests to latest variants from upstream

##@ Build

.PHONY: build
build: generate ## Build manager binary.
	echo "Git commit hash: ${GIT_COMMIT_HASH}"
	go build -ldflags "-X main.gitCommitHash=${GIT_COMMIT_HASH}" -o bin/manager main.go

.PHONY: build-e2e
build-e2e: generate ## Build e2e binary.
	echo "Git commit hash: ${GIT_COMMIT_HASH}"
	go build -ldflags "-X main.gitCommitHash=${GIT_COMMIT_HASH}" -tags=e2e -o bin/e2e test/e2e/*.go

.PHONY: run
run: manifests generate ## Run a controller from your host.
	go run ./main.go

.PHONY: docker-build
docker-build:  ## Build docker image with the manager.
	echo "Git commit hash: ${GIT_COMMIT_HASH}"
	KO_DOCKER_REPO=ko.local GOFLAGS="-ldflags=-X=main.gitCommitHash=${GIT_COMMIT_HASH}" ko build -B --platform=${PLATFORMS} -t ${IMG_TAG} .

.PHONY: docker-push
docker-push:  ## Push docker image with the manager.
	KO_DOCKER_REPO=${IMG_REPO} ko build --bare --platform=${PLATFORMS} -t ${IMG_TAG} .

.PHONY: docker-push-kind
docker-push-kind:  ## Make docker image available to kind cluster.
	GOOS=linux GOARCH=${shell go env GOARCH} KO_DOCKER_REPO=ko.local ko build -B -t ${IMG_TAG} .
	docker tag ko.local/cluster-api-provider-nutanix:${IMG_TAG} ${MANAGER_IMAGE}
	kind load docker-image --name ${KIND_CLUSTER_NAME} ${MANAGER_IMAGE}

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif


.PHONY: install
install: manifests ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	kustomize build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	kustomize build config/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests docker-push-kind ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	clusterctl delete --infrastructure nutanix:${LOCAL_PROVIDER_VERSION} --include-crd || true
	clusterctl init --infrastructure nutanix:${LOCAL_PROVIDER_VERSION} -v 9

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	kustomize build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

##@ Templates

.PHONY: cluster-e2e-templates
cluster-e2e-templates: cluster-e2e-templates-v1beta1 cluster-e2e-templates-v152 ## Generate cluster templates for all versions

cluster-e2e-templates-v152: ## Generate cluster templates for CAPX v1.5.2
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1.5.2/cluster-template --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1.5.2/cluster-template.yaml

cluster-e2e-templates-v1beta1: ## Generate cluster templates for v1beta1
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-secret --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-secret.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nutanix-cluster --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nutanix-cluster.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-additional-categories --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-additional-categories.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nmt --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nmt.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-project --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-project.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-upgrades --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-upgrades.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-md-remediation --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-md-remediation.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-remediation --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-remediation.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-scale-in --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-scale-in.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-csi --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-csi.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-csi3 --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-csi3.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-failure-domains --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-failure-domains.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-clusterclass --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-clusterclass.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-clusterclass --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/clusterclass-nutanix-quick-start.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-topology --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-topology.yaml

cluster-e2e-templates-no-kubeproxy: ##Generate cluster templates without kubeproxy
	# v1beta1
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-no-secret --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-secret.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-no-nutanix-cluster --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nutanix-cluster.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-additional-categories --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-additional-categories.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-no-nmt --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-no-nmt.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-project --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-project.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-upgrades --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-upgrades.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-md-remediation --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-md-remediation.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-kcp-remediation --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-remediation.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-kcp-scale-in --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-kcp-scale-in.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-csi --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-csi.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-failure-domains --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-failure-domains.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-clusterclass --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-clusterclass.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-clusterclass --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/clusterclass-nutanix-quick-start.yaml
	kustomize build $(NUTANIX_E2E_TEMPLATES)/v1beta1/no-kubeproxy/cluster-template-topology --load-restrictor LoadRestrictionsNone > $(NUTANIX_E2E_TEMPLATES)/v1beta1/cluster-template-topology.yaml

cluster-templates: ## Generate cluster templates for all flavors
	kustomize build $(TEMPLATES_DIR)/base > $(TEMPLATES_DIR)/cluster-template.yaml
	kustomize build $(TEMPLATES_DIR)/csi > $(TEMPLATES_DIR)/cluster-template-csi.yaml
	kustomize build $(TEMPLATES_DIR)/csi3 > $(TEMPLATES_DIR)/cluster-template-csi3.yaml
	kustomize build $(TEMPLATES_DIR)/clusterclass > $(TEMPLATES_DIR)/cluster-template-clusterclass.yaml
	kustomize build $(TEMPLATES_DIR)/topology > $(TEMPLATES_DIR)/cluster-template-topology.yaml

##@ Testing

.PHONY: docker-build-e2e
docker-build-e2e: ## Build docker image with the manager with e2e tag.
	echo "Git commit hash: ${GIT_COMMIT_HASH}"
	KO_DOCKER_REPO=ko.local GOFLAGS="-ldflags=-X=main.gitCommitHash=${GIT_COMMIT_HASH}" ko build -B --platform=${PLATFORMS_E2E} -t ${IMG_TAG} .
	docker tag ko.local/cluster-api-provider-nutanix:${IMG_TAG} ${IMG_REPO}:${IMG_TAG}

.PHONY: prepare-local-clusterctl
prepare-local-clusterctl: manifests cluster-templates  ## Prepare overide file for local clusterctl.
	mkdir -p ~/.cluster-api/overrides/infrastructure-nutanix/${LOCAL_PROVIDER_VERSION}
	cd config/manager && kustomize edit set image controller=${MANAGER_IMAGE}
	kustomize build config/default > ~/.cluster-api/overrides/infrastructure-nutanix/${LOCAL_PROVIDER_VERSION}/infrastructure-components.yaml
	cp ./metadata.yaml ~/.cluster-api/overrides/infrastructure-nutanix/${LOCAL_PROVIDER_VERSION}/
	cp ./templates/cluster-template*.yaml ~/.cluster-api/overrides/infrastructure-nutanix/${LOCAL_PROVIDER_VERSION}/
	cp ./clusterctl.yaml.tmpl ./clusterctl.yaml
	env LOCAL_PROVIDER_VERSION=$(LOCAL_PROVIDER_VERSION) \
		envsubst -no-unset -no-empty -no-digit < ./clusterctl.yaml > ~/.cluster-api/clusterctl.yaml

.PHONY: mocks
mocks: ## Generate mocks for the project
	mockgen -destination=mocks/ctlclient/client_mock.go -package=mockctlclient sigs.k8s.io/controller-runtime/pkg/client Client
	mockgen -destination=mocks/ctlclient/manager_mock.go -package=mockctlclient sigs.k8s.io/controller-runtime/pkg/manager Manager
	mockgen -destination=mocks/ctlclient/cache_mock.go -package=mockctlclient sigs.k8s.io/controller-runtime/pkg/cache Cache
	mockgen -destination=mocks/k8sclient/informer.go -package=mockk8sclient k8s.io/client-go/informers/core/v1 ConfigMapInformer,SecretInformer
	mockgen -destination=mocks/k8sclient/lister.go -package=mockk8sclient k8s.io/client-go/listers/core/v1 SecretLister,SecretNamespaceLister
	mockgen -destination=mocks/k8sapimachinery/interfaces.go -package=mockmeta k8s.io/apimachinery/pkg/api/meta RESTMapper,RESTScope
	mockgen -destination=mocks/nutanix/v3.go -package=mocknutanixv3 github.com/nutanix-cloud-native/prism-go-client/v3 Service

GOTESTPKGS = $(shell go list ./... | grep -v /mocks | grep -v /templates)

.PHONY: unit-test
unit-test: mocks  ## Run unit tests.
	KUBEBUILDER_ASSETS="$(shell setup-envtest use $(ENVTEST_K8S_VERSION)  --arch=amd64 -p path)" $(GOTEST) $(GOTESTPKGS)

.PHONY: coverage
coverage: mocks ## Run the tests of the project and export the coverage
	KUBEBUILDER_ASSETS="$(shell setup-envtest use $(ENVTEST_K8S_VERSION)  --arch=amd64 -p path)" $(GOTEST) -race -coverprofile=coverage.out -covermode=atomic $(GOTESTPKGS)

.PHONY: template-test
template-test: docker-build prepare-local-clusterctl ## Run the template tests
	GOPROXY=off \
	LOCAL_PROVIDER_VERSION=$(LOCAL_PROVIDER_VERSION) \
		ginkgo --trace --v run templates

.PHONY: test-e2e
test-e2e: docker-build-e2e cluster-e2e-templates cluster-templates ## Run the end-to-end tests
	echo "Image tag for E2E test is ${IMG_TAG}"
	LOCAL_PROVIDER_VERSION=$(LOCAL_PROVIDER_VERSION) \
		MANAGER_IMAGE=$(MANAGER_IMAGE) \
		envsubst < ${E2E_CONF_FILE} > ${E2E_CONF_FILE_TMP}
	docker tag ko.local/cluster-api-provider-nutanix:${IMG_TAG} ${IMG_REPO}:${IMG_TAG}
	docker push ${IMG_REPO}:${IMG_TAG}
	mkdir -p $(ARTIFACTS)
	NUTANIX_LOG_LEVEL=debug ginkgo -v \
		--trace \
		--tags=e2e \
		--label-filter=$(LABEL_FILTER_ARGS) \
		$(call ginkgo_option,skip,$(GINKGO_SKIP)) \
		$(call ginkgo_option,focus,$(GINKGO_FOCUS)) \
		--nodes=$(GINKGO_NODES) \
		--no-color=$(GINKGO_NOCOLOR) \
		--output-dir="$(ARTIFACTS)" \
		--junit-report=${JUNIT_REPORT_FILE} \
		--timeout="24h" \
		$(GINKGO_ARGS) \
		./test/e2e -- \
		-e2e.artifacts-folder="$(ARTIFACTS)" \
		-e2e.config="$(E2E_CONF_FILE_TMP)" \
		-e2e.skip-resource-cleanup=$(SKIP_RESOURCE_CLEANUP) \
		-e2e.use-existing-cluster=$(USE_EXISTING_CLUSTER)

.PHONY: test-e2e-no-kubeproxy
test-e2e-no-kubeproxy: docker-build-e2e cluster-e2e-templates-no-kubeproxy cluster-templates ## Run the end-to-end tests without kubeproxy
	echo "Image tag for E2E test is ${IMG_TAG}"
	MANAGER_IMAGE=$(MANAGER_IMAGE) envsubst < ${E2E_CONF_FILE} > ${E2E_CONF_FILE_TMP}
	docker tag ko.local/cluster-api-provider-nutanix:${IMG_TAG} ${IMG_REPO}:${IMG_TAG}
	docker push ${IMG_REPO}:${IMG_TAG}
	mkdir -p $(ARTIFACTS)
	NUTANIX_LOG_LEVEL=debug ginkgo -v \
		--trace \
		--tags=e2e \
		--label-filter=$(LABEL_FILTER_ARGS) \
		$(call ginkgo_option,skip,$(GINKGO_SKIP)) \
		$(call ginkgo_option,focus,$(GINKGO_FOCUS)) \
		--nodes=$(GINKGO_NODES) \
		--no-color=$(GINKGO_NOCOLOR) \
		--output-dir="$(ARTIFACTS)" \
		--junit-report=${JUNIT_REPORT_FILE} \
		--timeout="24h" \
		$(GINKGO_ARGS) \
		./test/e2e -- \
		-e2e.artifacts-folder="$(ARTIFACTS)" \
		-e2e.config="$(E2E_CONF_FILE)" \
		-e2e.skip-resource-cleanup=$(SKIP_RESOURCE_CLEANUP) \
		-e2e.use-existing-cluster=$(USE_EXISTING_CLUSTER)

.PHONY: list-e2e
list-e2e: docker-build-e2e cluster-e2e-templates cluster-templates ## Run the end-to-end tests
	mkdir -p $(ARTIFACTS)
	ginkgo -v \
	    --trace \
	    --dry-run \
	    --tags=e2e \
	    --label-filter="$(LABEL_FILTERS)" \
		$(call ginkgo_option,skip,$(GINKGO_SKIP)) \
		$(call ginkgo_option,focus,$(GINKGO_FOCUS)) \
	    --nodes=$(GINKGO_NODES) \
	    --no-color=$(GINKGO_NOCOLOR) \
	    --output-dir="$(ARTIFACTS)" \
	    $(GINKGO_ARGS) \
	    ./test/e2e -- \
	    -e2e.artifacts-folder="$(ARTIFACTS)" \
	    -e2e.config="$(E2E_CONF_FILE)" \
	    -e2e.skip-resource-cleanup=$(SKIP_RESOURCE_CLEANUP) \
	    -e2e.use-existing-cluster=$(USE_EXISTING_CLUSTER)

.PHONY: test-e2e-calico
test-e2e-calico:
	CNI=$(CNI_PATH_CALICO) GIT_COMMIT="${GIT_COMMIT_HASH}" $(MAKE) test-e2e

.PHONY: test-e2e-flannel
test-e2e-flannel:
	CNI=$(CNI_PATH_FLANNEL) GIT_COMMIT="${GIT_COMMIT_HASH}" $(MAKE) test-e2e

.PHONY: test-e2e-cilium
test-e2e-cilium:
	CNI=$(CNI_PATH_CILIUM) GIT_COMMIT="${GIT_COMMIT_HASH}" GINKGO_SKIP=$(GINKGO_SKIP) $(MAKE) test-e2e

.PHONY: test-e2e-cilium-no-kubeproxy
test-e2e-cilium-no-kubeproxy:
	CNI=$(CNI_PATH_CILIUM_NO_KUBEPROXY) $(MAKE) test-e2e-no-kubeproxy

.PHONY: test-e2e-all-cni
test-e2e-all-cni: test-e2e test-e2e-calico test-e2e-flannel test-e2e-cilium test-e2e-cilium-no-kubeproxy

##@ Lint and Verify

.PHONY: lint
lint: ## Lint the codebase
	golangci-lint run -v

lint-yaml: ## Use yamllint on the yaml file of your projects
ifeq ($(EXPORT_RESULT), true)
	$(eval OUTPUT_OPTIONS = | yamllint-checkstyle > yamllint-checkstyle.xml)
endif
	yamllint -c .yamllint --no-warnings -f parsable $(shell git ls-files '*.yml' '*.yaml' | grep -v '^test/e2e/data/cni/') $(OUTPUT_OPTIONS)

.PHONY: lint-fix
lint-fix: ## Lint the codebase and run auto-fixers if supported by the linter
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

## --------------------------------------
## Developer local tests
## --------------------------------------

##@ Test Dev Cluster with and without topology
include test-cluster-without-topology.mk
include test-cluster-with-topology.mk
