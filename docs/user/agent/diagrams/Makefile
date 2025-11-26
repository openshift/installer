.PHONY: all clean help

SVG_FILES := ../agent_installer_services-install_workflow.svg \
             ../agent_installer_services-add_nodes_workflow.svg \
             ../agent_installer_services-unconfigured_ignition_and_config_image_flow.svg \
             ../agent_installer_services-interactive.svg
DOT_FILES := $(patsubst ../agent_installer_services-%.svg,%.dot,$(SVG_FILES))

# Find all systemd unit files that the generator depends on
SYSTEMD_UNITS := $(wildcard ../../../../data/data/agent/systemd/units/*.service) \
                 $(wildcard ../../../../data/data/agent/systemd/units/*.service.template)

all: $(SVG_FILES)

# Generate DOT files from systemd units
$(DOT_FILES): generate_diagrams.py $(SYSTEMD_UNITS)
	python3 generate_diagrams.py

# Pattern rule: build SVG from DOT file
../agent_installer_services-%.svg: %.dot
	dot -Tsvg $< -o $@

clean:
	rm -f $(SVG_FILES) $(DOT_FILES)

# Help target
help:
	@echo "Available targets:"
	@echo "  all      - Generate all SVG diagrams (auto-generates DOT files from systemd units) (default)"
	@echo "  clean    - Remove generated SVG and DOT files"
	@echo "  help     - Show this help message"
	@echo ""
	@echo "The diagrams are auto-generated from systemd units in:"
	@echo "  ../../../../data/data/agent/systemd/units/"
	@echo ""
	@echo "Generated files:"
	@echo "  $(SVG_FILES)" | tr ' ' '\n' | sed 's/^/  /'"
