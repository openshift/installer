TARGET_OS_ARCH:=$(shell go env GOOS)_$(shell go env GOARCH)

TFSUBDIRS:=	$(foreach DIR,$(shell find providers -maxdepth 1 -mindepth 1 -type d),$(subst providers/,,$(DIR)))

GO_MOD_TIDY_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-mod-tidy-vendor.$(DIR),$(DIR)))
GO_BUILD_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-build.$(DIR),$(DIR)))
GO_CLEAN_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-clean.$(DIR),$(DIR)))
TERRAFORM_PROVIDER_TARGETS := $(foreach DIR,$(TFSUBDIRS), bin/$(TARGET_OS_ARCH)/cluster-api-provider-$(DIR).zip)

LDFLAGS:= "-s -w"
GCFLAGS:= ""

ifeq ($(MODE), dev)
LDFLAGS:= ""
GCFLAGS:= "all=-N -l"
endif

.PHONY: all
all: go-build

.PHONY: go-mod-tidy-vendor
go-mod-tidy-vendor: $(GO_MOD_TIDY_TARGETS)
# The -compat=1.18 is needed for the openstack provider. The provider uses golang.org/x/mod, which go 1.18 selects
# as v0.3.0 but go 1.16 selects as v0.4.2.
$(GO_MOD_TIDY_TARGETS): go-mod-tidy-vendor.%:
	cd providers/$* && go mod tidy -compat=1.18 && go mod vendor

.PHONY: go-build
go-build: $(GO_BUILD_TARGETS) go-build-cluster-api
$(GO_BUILD_TARGETS): go-build.%: bin/$(TARGET_OS_ARCH)/cluster-api-provider-%.zip

$(TERRAFORM_PROVIDER_TARGETS): bin/$(TARGET_OS_ARCH)/cluster-api-provider-%.zip: providers/%/go.mod
	cd providers/$*; \
	if [ -f main.go ]; then path="."; else path=./vendor/`grep _ tools.go|awk '{ print $$2 }'|sed 's|"||g'`; fi; \
	go build -gcflags $(GCFLAGS) -ldflags $(LDFLAGS) -o ../../bin/$(TARGET_OS_ARCH)/cluster-api-provider-$* "$$path";

.PHONY: go-build-cluster-api
go-build-cluster-api: bin/$(TARGET_OS_ARCH)/cluster-api

bin/$(TARGET_OS_ARCH)/cluster-api: cluster-api/go.mod
	cd cluster-api; \
	go build -gcflags $(GCFLAGS) -ldflags $(LDFLAGS) -o ../bin/$(TARGET_OS_ARCH)/cluster-api ./vendor/sigs.k8s.io/cluster-api

.PHONY: go-clean
go-clean: go-clean-providers go-clean-cluster-api

$(GO_CLEAN_TARGETS): go-clean.%:
	rm -f bin/*/cluster-api-provider-$*
	rm -f bin/*/cluster-api-provider-$*.zip

go-clean-providers:
	rm -f bin/*/cluster-api-provider-*

go-clean-cluster-api:
	rm -f bin/*/cluster-api

.PHONY: clean
clean: go-clean
