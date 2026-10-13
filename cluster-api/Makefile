TARGET_OS_ARCH:=$(shell go env GOOS)_$(shell go env GOARCH)

PROVSUBDIRS:=	$(foreach DIR,$(shell find providers -maxdepth 1 -mindepth 1 -type d),$(subst providers/,,$(DIR)))

GO_MOD_TIDY_TARGETS:=	$(foreach DIR,$(PROVSUBDIRS), $(subst $(DIR),go-mod-tidy-vendor.$(DIR),$(DIR)))
GO_BUILD_TARGETS:=	$(foreach DIR,$(PROVSUBDIRS), $(subst $(DIR),go-build.$(DIR),$(DIR)))
GO_CLEAN_TARGETS:=	$(foreach DIR,$(PROVSUBDIRS), $(subst $(DIR),go-clean.$(DIR),$(DIR)))
PROVIDER_TARGETS := $(foreach DIR,$(PROVSUBDIRS), bin/$(TARGET_OS_ARCH)/cluster-api-provider-$(DIR))

LDFLAGS:= "-s -w"
GCFLAGS:= ""

ifeq ($(MODE), dev)
LDFLAGS:= ""
GCFLAGS:= "all=-N -l"
endif

.PHONY: all
all: go-build

.PHONY: go-mod-tidy-vendor
go-mod-tidy-vendor: $(GO_MOD_TIDY_TARGETS) go-mod-tidy-vendor-cluster-api
$(GO_MOD_TIDY_TARGETS): go-mod-tidy-vendor.%:
	cd providers/$* && go mod tidy && go mod vendor

.PHONY: go-build
go-build: $(GO_BUILD_TARGETS) go-build-cluster-api
$(GO_BUILD_TARGETS): go-build.%: bin/$(TARGET_OS_ARCH)/cluster-api-provider-%

$(PROVIDER_TARGETS): bin/$(TARGET_OS_ARCH)/cluster-api-provider-%: providers/%/go.mod
	cd providers/$*; \
	if [ -f main.go ]; then path="."; else path=./vendor/`grep _ tools.go|awk '{ print $$2 }'|sed 's|"||g'`; fi; \
	go build -gcflags $(GCFLAGS) -ldflags $(LDFLAGS) -o ../../bin/$(TARGET_OS_ARCH)/cluster-api-provider-$* "$$path";

.PHONY: go-build-cluster-api
go-build-cluster-api: bin/$(TARGET_OS_ARCH)/cluster-api

.PHONY: go-build-cluster-api
go-mod-tidy-vendor-cluster-api:
	cd cluster-api; go mod tidy && go mod vendor

bin/$(TARGET_OS_ARCH)/cluster-api: cluster-api/go.mod
	cd cluster-api; \
	go build -gcflags $(GCFLAGS) -ldflags $(LDFLAGS) -o ../bin/$(TARGET_OS_ARCH)/cluster-api ./vendor/sigs.k8s.io/cluster-api

.PHONY: go-clean
go-clean: go-clean-providers go-clean-cluster-api

$(GO_CLEAN_TARGETS): go-clean.%:
	rm -f bin/*/cluster-api-provider-$*

go-clean-providers:
	rm -f bin/*/cluster-api-provider-*

go-clean-cluster-api:
	rm -f bin/*/cluster-api

.PHONY: clean
clean: go-clean
