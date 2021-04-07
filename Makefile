OUTPUT_DIR = bin
CROSS_BUILD_BINDIR := _output/bin

build: update
	OUTPUT=$(OUTPUT_DIR)/openshift-install SKIP_GENERATION=y hack/build.sh
.PHONY: build

update:
	GOOS= GOARCH= go generate ./data
.PHONY: update

baremetal:
	TAGS="baremetal libvirt" $(MAKE) --no-print-directory build
.PHONY: baremetal

libvirt:
	TAGS="libvirt" $(MAKE) --no-print-directory build
.PHONY: libvirt

cross-build-darwin-amd64:
	+@GOOS=darwin GOARCH=amd64 $(MAKE) --no-print-directory build OUTPUT_DIR=$(CROSS_BUILD_BINDIR)/darwin_amd64
.PHONY: cross-build-darwin-amd64

# cannot be built statically for darwin/arm64
# https://github.com/golang/go/issues/40492
#cross-build-darwin-arm64:
#	+@GOOS=darwin GOARCH=arm64 $(MAKE) --no-print-directory build OUTPUT_DIR=$(CROSS_BUILD_BINDIR)/darwin_arm64
#.PHONY: cross-build-darwin-arm64

cross-build-linux-amd64:
	+@GOOS=linux GOARCH=amd64 $(MAKE) --no-print-directory build OUTPUT_DIR=$(CROSS_BUILD_BINDIR)/linux_amd64
.PHONY: cross-build-linux-amd64

# cannot be built statically for linux/arm64
# https://github.com/golang/go/issues/40492
#cross-build-linux-arm64:
#	+@GOOS=linux GOARCH=arm64 $(MAKE) --no-print-directory build OUTPUT_DIR=$(CROSS_BUILD_BINDIR)/linux_arm64
#.PHONY: cross-build-linux-arm64

cross-build-linux-ppc64le:
	+@GOOS=linux GOARCH=ppc64le $(MAKE) --no-print-directory build OUTPUT_DIR=$(CROSS_BUILD_BINDIR)/linux_ppc64le
.PHONY: cross-build-linux-ppc64le

cross-build-linux-s390x:
	+@GOOS=linux GOARCH=s390x $(MAKE) --no-print-directory build OUTPUT_DIR=$(CROSS_BUILD_BINDIR)/linux_s390x
.PHONY: cross-build-linux-s390x

cross-build: cross-build-darwin-amd64 cross-build-linux-amd64 cross-build-linux-ppc64le cross-build-linux-s390x
.PHONY: cross-build

clean-cross-build:
	$(RM) -r '$(CROSS_BUILD_BINDIR)'
	if [ -d '$(OUTPUT_DIR)' ]; then rmdir --ignore-fail-on-non-empty '$(OUTPUT_DIR)'; fi
.PHONY: clean-cross-build

clean: clean-cross-build
	$(RM) -r '$(OUTPUT_DIR)'
