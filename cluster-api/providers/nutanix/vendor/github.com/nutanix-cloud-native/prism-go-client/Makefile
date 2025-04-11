SHELL := /bin/bash
BINARY_NAME=nutanixclient
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build

all: help

## Build:
build: ## Build your project and put the output binary in bin/
	mkdir -p bin
	$(GOCMD) build -o bin/$(BINARY_NAME) .

# CRD_OPTIONS define options to add to the CONTROLLER_GEN
CRD_OPTIONS ?= "crd:crdVersions=v1"

.PHONY: run-keploy
run-keploy:
	server run &

.PHONY: stop-keploy
stop-keploy:
	@-pkill "server"

generate: $(CONTROLLER_GEN)  ## Generate zz_generated.deepcopy.go
	controller-gen paths="./..." object:headerFile="hack/boilerplate.go.txt"

generate-v3-models: ## Generate V3 models using go-swagger
	swagger generate model \
		--spec=v3/swagger.json \
		--target=v3 \
		--skip-validation \
		--model=prism_central \
		--model=pc_vm \
		--model=mcm_config \
		--model=cmsp_config \
		--model=cmsp_network_config \
		--model=deployment_settings \
		--model=my_ntnx_token \
		--model=cluster_reference \
		--model=pc_vm_nic_configuration \
		--model=network_config

clean: ## Remove build related file
	rm -fr ./bin vendor hack/tools/bin
	rm -f checkstyle-report.xml ./coverage.out ./profile.cov yamllint-checkstyle.xml

## Test:
test: run-keploy ## Run the tests of the project
	go test -race -v ./...
	@$(MAKE) stop-keploy

coverage: run-keploy ## Run the tests of the project and export the coverage
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@$(MAKE) stop-keploy

## Lint:
lint: lint-go lint-yaml lint-kubebuilder ## Run all available linters

lint-go: ## Use golintci-lint on your project
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	golangci-lint run -v

lint-yaml: ## Use yamllint on the yaml file of your projects
ifeq ($(EXPORT_RESULT), true)
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | yamllint-checkstyle > yamllint-checkstyle.xml)
endif
	yamllint -c .yamllint --no-warnings -f parsable $(shell git ls-files '*.yml' '*.yaml') $(OUTPUT_OPTIONS)

.PHONY: lint-kubebuilder
lint-kubebuilder: ## Lint Kubebuilder annotations by generating objects and checking if it is successful
	controller-gen $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=.

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

