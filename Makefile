.PHONY: all
all: bin/openshift-install check ## Build and check everything

.PHONY: check
check: fmt vet lint test ## Check your code and run tests

.PHONY: revendor
revendor: ## Revendor go packages
	rm glide.lock
	glide install --strip-vendor
	glide-vc --use-lock-file --no-tests --only-code

.PHONY: bin/terraform 
bin/terraform: ## Get terraform executable
	hack/get-terraform.sh

.PHONY: bin/openshift-install
bin/openshift-install: ## build installer binary
	hack/build.sh

.PHONY: test
test: ## Run unit tests
	go test -race -cover ./pkg/...

.PHONY: lint
lint: ## Go lint your code
	hack/go-lint.sh -min_confidence 0.3 $(go list -f '{{ .ImportPath }}' ./...)

.PHONY: fmt
fmt: ## Go fmt your code
	hack/go-fmt.sh .

.PHONY: vet
vet: ## Vet all Go files
	hack/go-vet.sh ./...

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9/_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
