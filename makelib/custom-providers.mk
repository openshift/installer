PROVIDER_MATCHBOX_VERSION = v0.2.2

.PHONY:
custom-providers: $(INSTALLER_PATH)/terraform-provider-matchbox

$(INSTALLER_PATH)/terraform-provider-matchbox:
	curl -L -o $(TMPDIR)/terraform-provider-matchbox-$(PROVIDER_MATCHBOX_VERSION)-$(GOOS)-$(GOARCH).tar.gz \
	  https://github.com/coreos/terraform-provider-matchbox/releases/download/$(PROVIDER_MATCHBOX_VERSION)/terraform-provider-matchbox-$(PROVIDER_MATCHBOX_VERSION)-$(GOOS)-$(GOARCH).tar.gz
	cd $(TMPDIR) && tar xvf terraform-provider-matchbox-$(PROVIDER_MATCHBOX_VERSION)-$(GOOS)-$(GOARCH).tar.gz
	mkdir -p $(INSTALLER_PATH)
	cp $(TMPDIR)/terraform-provider-matchbox-$(PROVIDER_MATCHBOX_VERSION)-$(GOOS)-$(GOARCH)/terraform-provider-matchbox $(INSTALLER_PATH)/
	rm -rf $(TMPDIR)/terraform-provider-matchbox-$(PROVIDER_MATCHBOX_VERSION)-$(GOOS)-$(GOARCH)*
