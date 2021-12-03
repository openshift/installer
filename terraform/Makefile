TFSUBDIRS:=	aws azurerm azurestack google ibm ignition ironic kubernetes kubevirt libvirt local openstack ovirt random vsphere
TFBINDIR?=	.

GO_MOD_TIDY_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-mod-tidy-vendor.$(DIR),$(DIR)))
GO_BUILD_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-build.$(DIR),$(DIR)))
GO_CLEAN_TARGETS:=	$(foreach DIR,$(TFSUBDIRS), $(subst $(DIR),go-clean.$(DIR),$(DIR)))

TFPRIVATESUBDIRS:=		azureprivatedns vsphereprivate
GO_BUILD_PRIVATE_TARGETS:=	$(foreach DIR,$(TFPRIVATESUBDIRS), $(subst $(DIR),go-build.$(DIR),$(DIR)))
GO_MOD_TIDY_TARGETS+= 		$(foreach DIR,$(TFPRIVATESUBDIRS), $(subst $(DIR),go-mod-tidy-vendor.$(DIR),$(DIR)))
GO_CLEAN_TARGETS+= 		$(foreach DIR,$(TFPRIVATESUBDIRS), $(subst $(DIR),go-clean.$(DIR),$(DIR)))

.PHONY: all
all: go-build go-build-terraform

.PHONY: go-mod-tidy-vendor
go-mod-tidy-vendor: $(GO_MOD_TIDY_TARGETS)
$(GO_MOD_TIDY_TARGETS)::
	cd $(subst go-mod-tidy-vendor.,,$@) && go mod tidy && go mod vendor

.PHONY: go-build-private
go-build-private: $(GO_BUILD_PRIVATE_TARGETS)
$(GO_BUILD_PRIVATE_TARGETS):
	dir=$(subst go-build.,,$@); cd $$dir; \
	go build -o $(TFBINDIR)/terraform-provider-$$dir;

.PHONY: go-build
go-build: $(GO_BUILD_TARGETS) $(GO_BUILD_PRIVATE_TARGETS)
$(GO_BUILD_TARGETS):: go-mod-tidy-vendor
	dir=$(subst go-build.,,$@); cd $$dir; \
	path=`grep import $$dir.go|awk '{ print $$2 }'|sed 's|"||g'`; \
	go build -o $(TFBINDIR)/terraform-provider-$$dir ./vendor/$$path; 

.PHONY: go-clean
go-clean: $(GO_CLEAN_TARGETS)
$(GO_CLEAN_TARGETS)::
	dir=$(subst go-clean.,,$@); cd $$dir && rm -rf vendor && rm -f $(TFBINDIR)/terraform-provider-$$dir

.PHONY: clean
clean: go-clean go-clean-terraform

.PHONY: go-mod-tidy-vendor-terraform
go-mod-tidy-vendor-terraform::
	cd terraform && go mod tidy && go mod vendor

.PHONY: go-build-terraform
go-build-terraform:: go-mod-tidy-vendor-terraform
	cd terraform && go build -o $(TFBINDIR)/terraform ./vendor/github.com/hashicorp/terraform

.PHONY: go-clean-terraform
go-clean-terraform:
	cd terraform && rm -rf vendor && rm -f $(TFBINDIR)/terraform
