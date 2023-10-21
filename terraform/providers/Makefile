TARGET_OS_ARCH:=$(shell go env GOOS)_$(shell go env GOARCH)

SUBDIRS:=$(foreach DIR,$(shell find . -maxdepth 2 -mindepth 1 -name "go.mod" -exec dirname {} \;), $(subst ./,,$(DIR)))
TARGETS:=$(filter-out $(IGNORE),$(SUBDIRS))

GO_MOD_TIDY_TARGETS:=$(foreach DIR,$(SUBDIRS), $(subst $(DIR),go-mod-tidy-vendor.$(DIR),$(DIR)))
GO_MOD_VERIFY_TARGETS:=$(foreach DIR,$(SUBDIRS), $(subst $(DIR),go-mod-verify-vendor.$(DIR),$(DIR)))
GO_BUILD_TARGETS:=$(foreach DIR,$(SUBDIRS), $(subst $(DIR),go-build.$(DIR),$(DIR)))
GO_CLEAN_TARGETS:=$(foreach DIR,$(SUBDIRS), $(subst $(DIR),go-clean.$(DIR),$(DIR)))

PROVIDER_TARGETS:=$(foreach DIR,$(SUBDIRS), bin/$(TARGET_OS_ARCH)/terraform-provider-$(DIR).zip)

LDFLAGS:="-s -w"
GCFLAGS:=""

ifeq ($(MODE), dev)
	LDFLAGS:= ""
	GCFLAGS:= "all=-N -l"
endif

.PHONY: all
all: go-build

.PHONY: go-mod-tidy-vendor
go-mod-tidy-vendor: $(GO_MOD_TIDY_TARGETS)
$(GO_MOD_TIDY_TARGETS): go-mod-tidy-vendor.%:
	cd $* && go mod tidy && go mod vendor

.PHONY: go-build
go-build: $(GO_BUILD_TARGETS)
$(GO_BUILD_TARGETS): go-build.%: bin/$(TARGET_OS_ARCH)/terraform-provider-%.zip

$(PROVIDER_TARGETS): bin/$(TARGET_OS_ARCH)/terraform-provider-%.zip: %/go.mod
	cd $*; \
	if [ -f main.go ]; then path="."; else path=./vendor/`grep _ tools.go | awk '{ print $$2 }' | sed 's|"||g'`; fi; \
	go build -gcflags $(GCFLAGS) -ldflags $(LDFLAGS) -o ../bin/$(TARGET_OS_ARCH)/terraform-provider-$* "$$path" && \
	zip -1j ../bin/$(TARGET_OS_ARCH)/terraform-provider-$*.zip ../bin/$(TARGET_OS_ARCH)/terraform-provider-$*;

.PHONY: go-clean
go-clean: go-clean-providers

$(GO_CLEAN_TARGETS): go-clean.%:
	rm -f bin/*/terraform-provider-$*
	rm -f bin/*/terraform-provider-$*.zip

go-clean-providers:
	rm -f bin/*/terraform-provider-*

.PHONY: clean
clean: go-clean

.PHONY: go-mod-verify-vendor
go-mod-verify-vendor: go-mod-tidy-vendor $(GO_MOD_VERIFY_TARGETS)
	git diff --exit-code

$(GO_MOD_VERIFY_TARGETS): go-mod-verify-vendor.%:
	cd $* && go mod verify
