include terraform.mk

CLUSTER ?= demo
PLATFORM ?= aws
TOP_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
BUILD_DIR = $(TOP_DIR)/build/$(CLUSTER)

TF_DOCS := $(shell which terraform-docs 2> /dev/null)

$(info Using build directory [${BUILD_DIR}])

all: apply

localconfig:
	mkdir -p $(BUILD_DIR)
	touch $(BUILD_DIR)/terraform.tfvars

$(BUILD_DIR)/.terraform:
	cd $(BUILD_DIR) && terraform get $(TOP_DIR)/platforms/$(PLATFORM)

plan: $(BUILD_DIR)/.terraform
	cd $(BUILD_DIR) && terraform plan $(TOP_DIR)/platforms/$(PLATFORM)

apply: $(BUILD_DIR)/.terraform
	cd $(BUILD_DIR) && terraform apply $(TOP_DIR)/platforms/$(PLATFORM)

destroy:
	cd $(BUILD_DIR) && terraform destroy $(TOP_DIR)/platforms/$(PLATFORM)

terraform-check:
    @terraform-docs >/dev/null 2>&1 || @echo "terraform-docs is required (https://github.com/segmentio/terraform-docs)"

Documentation/variables/%.md: *.tf
ifndef TF_DOCS
	$(error "terraform-docs is required (https://github.com/segmentio/terraform-docs)")
endif

	$(eval PLATFORM_DIR := $(subst platform,platforms,$(subst -,/,$*)))
	@echo $(PLATFORM_DIR)
	@echo '# Terraform variables' >$@
	@echo 'This document gives an overview of the variables used in the different platforms of the Tectonic SDK.' >>$@
	terraform-docs markdown $(PLATFORM_DIR)/variables.tf >>$@

docs: Documentation/variables/config.md Documentation/variables/platform-aws.md Documentation/variables/platform-azure.md

clean: destroy
	rm -rf $(BUILD_DIR)

# This target is used by the GitHub PR checker to validate canonical syntax on all files.
#
structure-check:
	terraform fmt -list -write=false . | tee $TMPDIR/tf_fmt_files
	exit $(wc -l <$TMPDIR/tf_fmt_files)

canonical-syntax:
	terraform fmt -list .

.PHONY: make clean terraform terraform-dev structure-check canonical-syntax
