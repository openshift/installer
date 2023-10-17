
include $(ROOT_DIR_RELATIVE)/versions.mk

TOOLS_DIR := $(ROOT_DIR_RELATIVE)/hack/tools
TOOLS_DIR_DEPS := $(TOOLS_DIR)/go.sum $(TOOLS_DIR)/go.mod $(TOOLS_DIR)/Makefile
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin

$(TOOLS_BIN_DIR)/%: $(TOOLS_DIR_DEPS)
	make -C $(TOOLS_DIR) $(subst $(TOOLS_DIR)/,,$@)
