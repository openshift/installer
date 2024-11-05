#
# Copyright (c) 2018 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#
#
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
LOCAL_BIN_PATH := $(PROJECT_PATH)/bin

# Add the project-level bin directory into PATH. Needed in order
# for the tasks to use project-level bin directory binaries first
export PATH := $(LOCAL_BIN_PATH):$(PATH)

# Disable CGO so that we always generate static binaries:
export CGO_ENABLED=0

# Details of the model to use:
model_version:=v0.0.393
model_url:=https://github.com/openshift-online/ocm-api-model.git

# Details of the metamodel to use:
metamodel_version:=v0.0.61

goimports_version:=v0.4.0

# Additional flags to pass to the `ginkgo` command. This is used in the GitHub
# actions environment to skip tests that are sensitive to the speed of the
# machine: the leadership flag and retry tests.
ginkgo_flags:=

.DEFAULT_GOAL := examples

GINKGO := $(LOCAL_BIN_PATH)/ginkgo
ginkgo-install:
	@GOBIN=$(LOCAL_BIN_PATH) go install github.com/onsi/ginkgo/v2/ginkgo@v2.1.4 ;\

.PHONY: examples
examples:
	cd examples && \
	for i in *.go; do \
		go build $${i} || exit 1; \
	done

.PHONY: test tests
test tests: ginkgo-install
	$(GINKGO) run -r $(ginkgo_flags)

.PHONY: fmt
fmt:
	gofmt -s -l -w .

.PHONY: lint
lint:
	golangci-lint --version
	golangci-lint run

.PHONY: generate
generate: model metamodel-install goimports-install
	rm -rf \
		accesstransparency \
		accountsmgmt \
		addonsmgmt \
		authorizations \
		clustersmgmt \
		errors \
		helpers \
		jobqueue \
		servicelogs \
		servicemgmt \
		statusboard \
		webrca \
		osdfleetmgmt \
		openapi
	$(METAMODEL) generate go \
		--model=model/model \
		--base=github.com/openshift-online/ocm-sdk-go \
		--output=.
	$(METAMODEL) generate openapi \
		--model=model/model \
		--output=openapi

.PHONY: model
model:
	rm -rf "$@"
	if [ -d "$(model_url)" ]; then \
		cp -r "$(model_url)" "$@"; \
	else \
		git clone "$(model_url)" "$@"; \
		cd "$@"; \
		git fetch --tags origin; \
		git checkout -B build "$(model_version)"; \
	fi

.PHONY: metamodel
METAMODEL=$(LOCAL_BIN_PATH)/metamodel
metamodel-install:
	GOBIN=$(LOCAL_BIN_PATH) go install github.com/openshift-online/ocm-api-metamodel/cmd/metamodel@$(metamodel_version)

.PHONY: goimports
goimports-install:
	@GOBIN=$(LOCAL_BIN_PATH) go install golang.org/x/tools/cmd/goimports@$(goimports_version)

.PHONY: clean
clean:
	rm -rf \
		metamodel \
		model \
		$(NULL)
