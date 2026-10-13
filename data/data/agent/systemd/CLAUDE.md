# Agent Installer Systemd Unit Files

When you modify systemd unit files in the `data/data/agent/systemd/units/` directory, and your changes affect the `[Unit]` section (which contains directives like `After`, `Before`, `PartOf`, etc. that define the relationships between systemd services), you must regenerate the agent installer service diagrams by running:

```bash
make -C docs/user/agent/diagrams/
```
