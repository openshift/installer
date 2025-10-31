.PHONY: all clean help

PNG_FILES := ../agent_installer_services-install_workflow.png \
             ../agent_installer_services-add_nodes_workflow.png \
             ../agent_installer_services-unconfigured_ignition_and_config_image_flow.png \
             ../agent_installer_services-interactive.png
DOT_FILES := $(patsubst ../agent_installer_services-%.png,%.dot,$(PNG_FILES))

# Find all systemd unit files that the generator depends on
SYSTEMD_UNITS := $(wildcard ../../../../data/data/agent/systemd/units/*.service) \
                 $(wildcard ../../../../data/data/agent/systemd/units/*.service.template)

all: $(PNG_FILES)

# Generate DOT files from systemd units
$(DOT_FILES): generate_diagrams.py $(SYSTEMD_UNITS)
	python3 generate_diagrams.py

# Pattern rule: build PNG from DOT file
../agent_installer_services-%.png: %.dot
	dot -Tpng -Gdpi=150 $< -o $@

clean:
	rm -f $(PNG_FILES) $(DOT_FILES)

# Help target
help:
	@echo "Available targets:"
	@echo "  all      - Generate all PNG diagrams (auto-generates DOT files from systemd units) (default)"
	@echo "  clean    - Remove generated PNG and DOT files"
	@echo "  help     - Show this help message"
	@echo ""
	@echo "The diagrams are auto-generated from systemd units in:"
	@echo "  ../../../../data/data/agent/systemd/units/"
	@echo ""
	@echo "Generated files:"
	@echo "  $(PNG_FILES)" | tr ' ' '\n' | sed 's/^/  /'"
