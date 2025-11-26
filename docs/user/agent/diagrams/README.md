# Agent Installer Service Diagrams

This directory contains the auto-generation system for agent installer service workflow diagrams.

## Structure

- `generate_diagrams.py` - Python script to auto-generate DOT files from systemd units
- `Makefile` - Build script with dependency tracking
- `*.dot` - GraphViz DOT source files (generated, not tracked in git)
- `.gitignore` - Excludes generated .dot files
- `README.md` - This file

## Diagrams

The system generates four workflow diagrams:

1. **install_workflow** - Standard agent-based installation workflow
2. **add_nodes_workflow** - Add-nodes workflow (differences highlighted in green)
3. **unconfigured_ignition** - Appliance/factory workflow with config image
4. **interactive** - Interactive installation workflow using assisted UI

## Building

To regenerate the diagrams:

```bash
cd docs/user/agent/diagrams
make
```

This single command will:
1. Auto-generate `.dot` files from systemd units (if units changed)
2. Generate SVG files from the `.dot` files (if `.dot` files changed)

The diagrams are automatically regenerated when:
- Any systemd unit file in `data/data/agent/systemd/units/` changes
- The `generate_diagrams.py` script changes

Output files are created in the parent directory:
- `agent_installer_services-install_workflow.svg`
- `agent_installer_services-add_nodes_workflow.svg`
- `agent_installer_services-unconfigured_ignition_and_config_image_flow.svg`
- `agent_installer_services-interactive.svg`

## How It Works

The generator parses systemd unit files and extracts:
- **Service dependencies**: `Before=` and `After=` directives become edges
- **Cluster membership**: `PartOf=` and `BindsTo=` define dashed boxes
- **Workflow filters**: `ConditionPathExists` determines which services run in each workflow
- **File dependencies**: `ConditionPathExists` on files creates dotted edges

The system is mostly data-driven with a few hardcoded exceptions:
- `load-config-iso@` is triggered by udev (unconfigured_ignition only)
- `start-cluster-installation` excluded from interactive (transitive dependency)
- `99-agent-copy-files.sh` excluded from interactive (agent-extract-tui provides tui)

## Diagram Features

- **Layout**: Bottom-to-top (dependencies flow upward)
- **Color coding**:
  - Light blue - standard services
  - Dark blue (thick border) - orchestrator services
  - Thin border - disconnected services
  - Green border/text - differences in add-nodes workflow
  - Yellow - configuration files
- **Dotted edges**: File creation or conditional dependencies
- **Dashed boxes**: Service clusters (PartOf/BindsTo relationships)

## Requirements

- GraphViz (`dot` command) must be installed
  - Fedora/RHEL: `dnf install graphviz`
  - Ubuntu/Debian: `apt install graphviz`
  - macOS: `brew install graphviz`
- Python 3 (for auto-generation)

## Cleaning

To remove all generated files:

```bash
make clean
```

This removes both `.dot` and `.svg` files.
