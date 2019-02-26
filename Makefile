# Old-skool build tools.
#
# Targets (see each target for more information):
#   all: Build code.
#   build: Build code.
#   verify: Run build and verify.
#   clean: Clean up.

GO_FILES = $(shell go list -f '{{ .ImportPath }}' ./...)

# Build code.
#
# Args:
#   WHAT: Directory names to build.  If any of these directories has a 'main'
#     package, the build will produce executable files under $(OUT_DIR)/local/bin.
#     If not specified, "everything" will be built.
#   GOFLAGS: Extra flags to pass to 'go' when building.
#   GOGCFLAGS: Additional go compile flags passed to 'go' when building.
#   TESTFLAGS: Extra flags that should only be passed to hack/test-go.sh
#
# Example:
#   make
#   make all
#   make build
all build:
	hack/build.sh
.PHONY: all build

# Verify code conventions are properly setup.
#
# Example:
#   make verify
verify: build
	{ \
    hack/go-fmt.sh . || r=1;\
    hack/go-lint.sh $(GO_FILES) || r=1;\
    hack/go-vet.sh ./... || r=1;\
    hack/shellcheck.sh || r=1;\
    hack/tf-fmt.sh -list -check || r=1;\
    hack/tf-lint.sh || r=1;\
    hack/yaml-lint.sh || r=1;\
    hack/go-test.sh || r=1;\
	exit $$r ;\
	}
.PHONY: verify

# Remove all build artifacts.
#
# Example:
#   make clean
clean:
	rm -rf bin/openshift-install
.PHONY: clean
