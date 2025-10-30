# Agent Installer Service Diagrams

This directory contains the source files for the agent installer service workflow diagrams.

## Structure

- `*.dot` - GraphViz DOT source files for each workflow diagram
- `Makefile` - Build script to generate PNG files from DOT sources
- `README.md` - This file

## Diagrams

1. **install_workflow.dot** - Standard agent-based installation workflow
2. **add_nodes_workflow.dot** - Add-nodes workflow (differences highlighted in green)
3. **unconfigured_ignition.dot** - Appliance/factory workflow with config image
4. **interactive.dot** - Interactive installation workflow using assisted UI

## Building

To regenerate the PNG diagrams:

```bash
cd docs/user/agent/diagrams
make all
```

This will create/update the PNG files in the parent directory (`docs/user/agent/`):
- `agent_installer_services-install_workflow.png`
- `agent_installer_services-add_nodes_workflow.png`
- `agent_installer_services-unconfigured_ignition_and_config_image_flow.png`
- `agent_installer_services-interactive.png`

## Requirements

- GraphViz (`dot` command) must be installed
  - Fedora/RHEL: `dnf install graphviz`
  - Ubuntu/Debian: `apt install graphviz`
  - macOS: `brew install graphviz`

## Editing

To modify a diagram:

1. Edit the corresponding `.dot` file
2. Run `make` to regenerate the PNG
3. Commit both the `.dot` source and the generated `.png` file

The DOT files use standard GraphViz syntax with:
- `rankdir=BT` for bottom-to-top layout (dependencies flow upward)
- Rounded, filled boxes for services
- Note shapes for configuration files
- Dashed subgraph for pod grouping
- Color coding:
  - Light blue (`#ADD8E6`) - standard services
  - Medium blue (`#6495ED`) - key orchestrator services
  - Light blue/white (`#F0F8FF`) - optional services
  - Green (`#90EE90`) - services that differ in add-nodes workflow
