.PHONY: build destroy

BUILD_DIR=./build

ASSETS:=assets-demo.zip

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(ASSETS):
	@echo "Assets '$(ASSETS)' not found!\nPlace assets zip from installer in $(PWD)"
	exit 1

config.tfvars:
	@echo "Cluster config missing.\nPlace a 'config.tfvars' file in $(PWD)"
	exit 1

assets: $(ASSETS)
	unzip $(ASSETS)

apply: $(BUILD_DIR) assets config.tfvars
	terraform get ./platform-$(PLATFORM)
	terraform apply -state-out=$(BUILD_DIR)/platform-$(PLATFORM) --var-file=./config.tfvars ./platform-$(PLATFORM)

destroy: $(BUILD_DIR) assets config.tfvars
	terraform destroy -state=$(BUILD_DIR)/platform-$(PLATFORM) --var-file=./config.tfvars ./platform-$(PLATFORM)
