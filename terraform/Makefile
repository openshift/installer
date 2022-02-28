TFSUBDIRS:=	alicloud azureprivatedns azurestack ibm ignition ironic ovirt vsphere vsphereprivate
TFBINDIR?=	.

GO_MOD_TIDY_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-mod-tidy-vendor.$(DIR),$(DIR)))
GO_BUILD_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-build.$(DIR),$(DIR)))
GO_CLEAN_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-clean.$(DIR),$(DIR)))

LDFLAGS:= "-s -w"

.PHONY: all
all: go-build

.PHONY: go-mod-tidy-vendor
go-mod-tidy-vendor: $(GO_MOD_TIDY_TARGETS)
$(GO_MOD_TIDY_TARGETS)::
	cd $(subst go-mod-tidy-vendor.,,$@) && go mod tidy && go mod vendor

.PHONY: go-build
go-build: $(GO_BUILD_TARGETS)
$(GO_BUILD_TARGETS):
	dir=$(subst go-build.,,$@); cd $$dir; \
	go build -ldflags $(LDFLAGS) -o $(TFBINDIR)/terraform-provider-$$dir;

.PHONY: go-clean
go-clean: $(GO_CLEAN_TARGETS)
$(GO_CLEAN_TARGETS)::
	rm -f $(TFBINDIR)/terraform-provider-$$dir

.PHONY: clean
clean: go-clean
