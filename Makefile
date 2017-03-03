.PHONY: build clean

CLUSTER ?= demo
ASSETS ?= assets-$(CLUSTER).zip
PLATFORM ?= aws-noasg
TOP_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BUILD_DIR = $(TOP_DIR)/build/$(CLUSTER)

all: apply

$(BUILD_DIR)/$(ASSETS):
	@echo "Assets '$(ASSETS)' not found!\nPlace assets zip from installer in $(BUILD_DIR)\n"
	exit 1

$(BUILD_DIR)/config.tfvars: $(BUILD_DIR)/assets
	$(TOP_DIR)/convert.sh tfvars $(PLATFORM) $(BUILD_DIR)/assets/cloud-formation.json > $(BUILD_DIR)/config.tfvars

$(BUILD_DIR)/assets: $(BUILD_DIR)/$(ASSETS)
	cd $(BUILD_DIR) && unzip $(ASSETS)
	$(TOP_DIR)/convert.sh assets $(PLATFORM) $(BUILD_DIR)/assets

$(BUILD_DIR)/.terraform:
	cd $(BUILD_DIR) && terraform get $(TOP_DIR)/platform-$(PLATFORM)

plan: $(BUILD_DIR)/assets $(BUILD_DIR)/config.tfvars $(BUILD_DIR)/.terraform
	cd $(BUILD_DIR) && terraform plan --var-file=config.tfvars $(TOP_DIR)/platform-$(PLATFORM)

apply: $(BUILD_DIR)/assets $(BUILD_DIR)/config.tfvars $(BUILD_DIR)/.terraform
	cd $(BUILD_DIR) && terraform apply --var-file=config.tfvars $(TOP_DIR)/platform-$(PLATFORM)

clean: $(BUILD_DIR)/assets $(BUILD_DIR)/config.tfvars
	cd $(BUILD_DIR) && terraform destroy --var-file=config.tfvars $(TOP_DIR)/platform-$(PLATFORM)
	cd $(BUILD_DIR) && rm -r .terraform assets config.tfvars terraform.tfstate terraform.tfstate.backup
