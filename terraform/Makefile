TARGET_OS_ARCH:=$(shell go env GOOS)_$(shell go env GOARCH)

TFSUBDIRS:=	$(foreach DIR,$(shell find providers -maxdepth 1 -mindepth 1 -type d),$(subst providers/,,$(DIR)))

GO_MOD_TIDY_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-mod-tidy-vendor.$(DIR),$(DIR)))
GO_BUILD_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-build.$(DIR),$(DIR)))
GO_CLEAN_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-clean.$(DIR),$(DIR)))
TERRAFORM_PROVIDER_TARGETS := $(foreach DIR,$(TFSUBDIRS), bin/$(TARGET_OS_ARCH)/terraform-provider-$(DIR).zip)

LDFLAGS:= -s -w
GCFLAGS:= ""

ifeq ($(MODE), dev)
LDFLAGS:=
GCFLAGS:= "all=-N -l"
endif

MODINFOKEY:= main.builtGoModHash

.PHONY: all
all: go-build

.PHONY: go-mod-tidy-vendor
go-mod-tidy-vendor: $(GO_MOD_TIDY_TARGETS)
$(GO_MOD_TIDY_TARGETS): go-mod-tidy-vendor.%:
	cd providers/$* && go mod tidy && go mod vendor

.PHONY: go-build
go-build: $(GO_BUILD_TARGETS) go-build-terraform
$(GO_BUILD_TARGETS): go-build.%: bin/$(TARGET_OS_ARCH)/terraform-provider-%.zip

$(TERRAFORM_PROVIDER_TARGETS): bin/$(TARGET_OS_ARCH)/terraform-provider-%.zip: providers/%/go.mod
	cd providers/$*; \
	if [ -f main.go ]; then path="."; else path=./vendor/`grep _ tools.go|awk '{ print $$2 }'|sed 's|"||g'`; fi; \
	go build -gcflags $(GCFLAGS) -ldflags "$(LDFLAGS) -X $(MODINFOKEY)=$$(git hash-object go.mod)" -o ../../bin/$(TARGET_OS_ARCH)/terraform-provider-$* "$$path" && \
	zip -1j ../../bin/$(TARGET_OS_ARCH)/terraform-provider-$*.zip ../../bin/$(TARGET_OS_ARCH)/terraform-provider-$*;

.PHONY: go-build-terraform
go-build-terraform: bin/$(TARGET_OS_ARCH)/terraform

bin/$(TARGET_OS_ARCH)/terraform: terraform/go.mod
	cd terraform; \
	go build -gcflags $(GCFLAGS) -ldflags "$(LDFLAGS) -X $(MODINFOKEY)=$$(git hash-object go.mod)" -o ../bin/$(TARGET_OS_ARCH)/terraform ./vendor/github.com/hashicorp/terraform

.PHONY: go-clean
go-clean: go-clean-providers go-clean-terraform

$(GO_CLEAN_TARGETS): go-clean.%:
	rm -f bin/*/terraform-provider-$*
	rm -f bin/*/terraform-provider-$*.zip

go-clean-providers:
	rm -f bin/*/terraform-provider-*

go-clean-terraform:
	rm -f bin/*/terraform

.PHONY: clean
clean: go-clean
