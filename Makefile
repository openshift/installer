.PHONY: build destroy

CLUSTER ?= demo
ASSETS ?= assets-$(CLUSTER).zip
PLATFORM ?= aws-noasg
TOP_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BUILD_DIR = $(TOP_DIR)/build/$(CLUSTER)

all: apply

$(BUILD_DIR)/$(ASSETS):
	@echo "Assets '$(ASSETS)' not found!\nPlace assets zip from installer in $(BUILD_DIR)"
	exit 1

$(BUILD_DIR)/config.tfvars:
	@echo "Cluster config missing.\nPlace a 'config.tfvars' file in $(BUILD_DIR)"
	exit 1

$(BUILD_DIR)/assets: $(BUILD_DIR)/$(ASSETS)
	cd $(BUILD_DIR) && unzip $(ASSETS)

apply: $(BUILD_DIR)/assets $(BUILD_DIR)/config.tfvars
	cd $(BUILD_DIR) && terraform get $(TOP_DIR)/platform-$(PLATFORM)
	cd $(BUILD_DIR) && terraform apply --var-file=config.tfvars $(TOP_DIR)/platform-$(PLATFORM)

destroy: $(BUILD_DIR)/assets $(BUILD_DIR)/config.tfvars
	cd $(BUILD_DIR) && terraform destroy --var-file=config.tfvars $(TOP_DIR)/platform-$(PLATFORM)
