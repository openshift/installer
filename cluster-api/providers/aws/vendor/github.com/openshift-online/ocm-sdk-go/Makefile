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

goimports_version:=v0.4.0

# Additional flags to pass to the `ginkgo` command. This is used in the GitHub
# actions environment to skip tests that are sensitive to the speed of the
# machine: the leadership flag and retry tests.
ginkgo_flags:=

.DEFAULT_GOAL := examples

GINKGO := $(LOCAL_BIN_PATH)/ginkgo
ginkgo-install:
	@GOBIN=$(LOCAL_BIN_PATH) go install github.com/onsi/ginkgo/v2/ginkgo@v2.1.4 ;\

verify: lint examples model metamodel-install goimports-install verify-hooks
	go vet $(find . -maxdepth 1 -type d  ! -name 'vendor' ! -name '.')
	hack/verify-gofmt.sh
	hack/verify-client.sh $(METAMODEL) .
	hack/verify-openapi.sh $(METAMODEL) ./openapi
.PHONY: verify

.PHONY: install-hooks
install-hooks:
	@echo "Installing git hooks..."
	@if [ -f hooks/pre-push ]; then \
		cp hooks/pre-push .git/hooks/pre-push && \
		chmod +x .git/hooks/pre-push && \
		echo "Installed and made .git/hooks/pre-push executable"; \
	else \
		echo "Error: hooks/pre-push not found in repository"; \
		exit 1; \
	fi

.PHONY: verify-hooks
verify-hooks:
	@echo "Verifying git hooks..."
	@if [ -f .git/hooks/pre-push ]; then \
		if [ -x .git/hooks/pre-push ]; then \
			echo ".git/hooks/pre-push is executable"; \
		else \
			echo ".git/hooks/pre-push is not executable"; \
			echo "  Run 'make install-hooks' to fix this"; \
			exit 1; \
		fi \
	else \
		echo "Warning: .git/hooks/pre-push not found"; \
		echo "  Run 'make install-hooks' to install it"; \
	fi

.PHONY: examples
examples:
	cd examples && \
	for i in *.go; do \
		go build $${i} || exit 1; \
	done

.PHONY: test tests
test tests: ginkgo-install
ifndef ARTIFACT_DIR
	$(GINKGO) run -r $(ginkgo_flags)
else
	$(GINKGO) run -r $(ginkgo_flags) --junit-report=junit_golang_tests_junit.xml --output-dir=${ARTIFACT_DIR}
endif


.PHONY: fmt
fmt:
	gofmt -s -l -w .

.PHONY: lint
lint:
	golangci-lint --version
	golangci-lint custom && \
	$(LOCAL_BIN_PATH)/ocm-lint run

# set "normal" alias matching kube and ocp
.PHONY: update
update: generate

.PHONY: generate
generate: model metamodel-install goimports-install
	hack/generate-client.sh $(METAMODEL) .
	hack/generate-openapi.sh $(METAMODEL) ./openapi

.PHONY: model
model:
	go mod tidy && go mod vendor

.PHONY: metamodel
METAMODEL=$(PROJECT_PATH)/metamodel_generator/metamodel
metamodel-install:
	 $(MAKE) -C metamodel_generator clean
	 $(MAKE) -C metamodel_generator build

.PHONY: goimports
goimports-install:
	@GOBIN=$(LOCAL_BIN_PATH) go install golang.org/x/tools/cmd/goimports@$(goimports_version)

.PHONY: clean
clean:
	hack/clean-client.sh
	hack/clean-openapi.sh
	rm -rf \
		$(PROJECT_PATH)/metamodel_generator/metamodel \
		model \
		$(NULL)
