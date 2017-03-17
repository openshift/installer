.PHONY: build clean

CLUSTER ?= demo
ASSETS ?= assets-$(CLUSTER).zip
PLATFORM ?= aws-noasg
TOP_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
BUILD_DIR = $(TOP_DIR)/build/$(CLUSTER)

$(info $$BUILD_DIR is [${BUILD_DIR}])

all: apply

$(BUILD_DIR)/$(ASSETS):
	@echo "Assets '$(ASSETS)' not found!\nPlace assets zip from installer in $(BUILD_DIR)\n"
	exit 1

localconfig:
	mkdir -p $(BUILD_DIR)
	cp config.tfvars $(BUILD_DIR)/config.tfvars

installerconfig: $(BUILD_DIR)/$(ASSETS)
	cd $(BUILD_DIR) && unzip $(ASSETS)
	$(TOP_DIR)/convert.sh tfvars $(PLATFORM) $(BUILD_DIR)/assets/cloud-formation.json > $(BUILD_DIR)/config.tfvars
	$(TOP_DIR)/convert.sh assets $(PLATFORM) $(BUILD_DIR)/assets

$(BUILD_DIR)/.terraform:
	cd $(BUILD_DIR) && terraform get $(TOP_DIR)/platforms/$(PLATFORM)

plan: $(BUILD_DIR)/config.tfvars $(BUILD_DIR)/.terraform
	cd $(BUILD_DIR) && terraform plan --var-file=config.tfvars $(TOP_DIR)/platforms/$(PLATFORM)

apply: $(BUILD_DIR)/config.tfvars $(BUILD_DIR)/.terraform
	cd $(BUILD_DIR) && terraform apply --var-file=config.tfvars $(TOP_DIR)/platforms/$(PLATFORM)

destroy: $(BUILD_DIR)/config.tfvars
	cd $(BUILD_DIR) && terraform destroy --var-file=config.tfvars $(TOP_DIR)/platforms/$(PLATFORM)

# You need to have https://github.com/segmentio/terraform-docs installed
Documentation/variables/%.md: **/*.tf
	echo '# Terraform variables: $*' >$@
	echo 'The Tectonic SDK variables used for: $*.' >>$@
	terraform-docs markdown ./$* >>$@

# You need to have https://github.com/segmentio/terraform-docs installed
Documentation/variables/config.md: *.tf
	echo '# Common Tectonic Terraform variables' >$@
	echo 'All the common Tectonic SDK variables used for *all* platforms.' >>$@
	terraform-docs markdown . >>$@

docs: Documentation/variables/config.md Documentation/variables/platform-aws.md Documentation/variables/platform-azure.md

clean:
	cd $(BUILD_DIR) && \
	rm -rf .terraform assets generated \
	rm -f config.tfvars terraform.tfstate terraform.tfstate.backup assets*.zip id_rsa*
