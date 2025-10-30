.PHONY: all clean

PNG_FILES := ../agent_installer_services-install_workflow.png \
             ../agent_installer_services-add_nodes_workflow.png \
             ../agent_installer_services-unconfigured_ignition_and_config_image_flow.png \
             ../agent_installer_services-interactive.png
DOT_FILES := install_workflow.dot add_nodes_workflow.dot unconfigured_ignition.dot interactive.dot

all: $(PNG_FILES)

../agent_installer_services-install_workflow.png: install_workflow.dot
	dot -Tpng -Gdpi=150 $< -o $@

../agent_installer_services-add_nodes_workflow.png: add_nodes_workflow.dot
	dot -Tpng -Gdpi=150 $< -o $@

../agent_installer_services-unconfigured_ignition_and_config_image_flow.png: unconfigured_ignition.dot
	dot -Tpng -Gdpi=150 $< -o $@

../agent_installer_services-interactive.png: interactive.dot
	dot -Tpng -Gdpi=150 $< -o $@

clean:
	rm -f $(PNG_FILES)

# Help target
help:
	@echo "Available targets:"
	@echo "  all      - Generate all PNG diagrams from DOT sources (default)"
	@echo "  clean    - Remove generated PNG files"
	@echo "  help     - Show this help message"
	@echo ""
	@echo "Generated files:"
	@echo "  $(PNG_FILES)" | tr ' ' '\n' | sed 's/^/  /'
