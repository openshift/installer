#!/usr/bin/env python3
"""
Generate agent installer service workflow diagrams from systemd unit files.

Reads systemd unit files from data/data/agent/systemd/units/ and generates
GraphViz DOT files showing service dependencies for each workflow.
"""

import re
import os
from pathlib import Path
from typing import Dict, List, Set, Tuple
from dataclasses import dataclass, field


@dataclass
class SystemdUnit:
    """Represents a parsed systemd unit file."""
    name: str
    description: str = ""
    after: List[str] = field(default_factory=list)
    before: List[str] = field(default_factory=list)
    requires: List[str] = field(default_factory=list)
    wants: List[str] = field(default_factory=list)
    binds_to: List[str] = field(default_factory=list)
    part_of: List[str] = field(default_factory=list)
    conflicts: List[str] = field(default_factory=list)
    condition_path_exists: List[str] = field(default_factory=list)
    condition_path_not_exists: List[str] = field(default_factory=list)
    wanted_by: List[str] = field(default_factory=list)
    is_template: bool = False


class SystemdParser:
    """Parse systemd unit files."""

    def __init__(self, units_dir: Path):
        self.units_dir = units_dir
        self.units: Dict[str, SystemdUnit] = {}

    def parse_all(self):
        """Parse all unit files in the directory."""
        for file_path in self.units_dir.glob("*.service*"):
            # Skip template files - we'll handle the base service
            if file_path.suffix == '.template':
                # Parse it but with the .template removed from name
                name = file_path.stem
            else:
                name = file_path.name

            unit = self.parse_unit(file_path, name)
            self.units[name] = unit

    def parse_unit(self, file_path: Path, name: str) -> SystemdUnit:
        """Parse a single systemd unit file."""
        unit = SystemdUnit(name=name)
        unit.is_template = '@' in name or file_path.suffix == '.template'

        with open(file_path, 'r') as f:
            content = f.read()

        # Remove Go template syntax for parsing
        content = re.sub(r'\{\{[^}]+\}\}', '', content)
        content = re.sub(r'\{\{[^}]+\}\}[^{]*\{\{[^}]+\}\}', '', content, flags=re.DOTALL)

        current_section = None
        for line in content.splitlines():
            line = line.strip()

            # Skip comments and empty lines
            if not line or line.startswith('#') or line.startswith(';'):
                continue

            # Section headers
            if line.startswith('[') and line.endswith(']'):
                current_section = line[1:-1]
                continue

            # Parse key=value pairs
            if '=' not in line:
                continue

            key, value = line.split('=', 1)
            key = key.strip()
            value = value.strip()

            # Parse directives
            if key == 'Description':
                unit.description = value
            elif key == 'After':
                unit.after.extend(self._parse_list(value))
            elif key == 'Before':
                unit.before.extend(self._parse_list(value))
            elif key == 'Requires':
                unit.requires.extend(self._parse_list(value))
            elif key == 'Wants':
                unit.wants.extend(self._parse_list(value))
            elif key == 'BindsTo':
                unit.binds_to.extend(self._parse_list(value))
            elif key == 'PartOf':
                unit.part_of.extend(self._parse_list(value))
            elif key == 'Conflicts':
                unit.conflicts.extend(self._parse_list(value))
            elif key == 'ConditionPathExists':
                if value.startswith('!'):
                    unit.condition_path_not_exists.append(value[1:])
                else:
                    unit.condition_path_exists.append(value)
            elif key == 'WantedBy' and current_section == 'Install':
                unit.wanted_by.extend(self._parse_list(value))

        return unit

    def _parse_list(self, value: str) -> List[str]:
        """Parse space-separated list of services/targets."""
        return [item.strip() for item in value.split() if item.strip()]


class WorkflowFilter:
    """Filter services by workflow based on conditions."""

    # Workflow discriminators - the key files that distinguish workflows
    # Other files like /etc/assisted/node0 are created during workflows and aren't discriminators
    WORKFLOW_DISCRIMINATORS = {
        '/etc/assisted/add-nodes.env',
        '/etc/assisted/interactive-ui',
        '/etc/assisted/rendezvous-host.env',
    }

    WORKFLOWS = {
        'install': {
            # Base workflow - runs when no discriminator files are present
            'workflow_markers': [],
            'excluded_markers': ['/etc/assisted/add-nodes.env', '/etc/assisted/interactive-ui',
                               '/etc/assisted/rendezvous-host.env'],
        },
        'add_nodes': {
            'workflow_markers': ['/etc/assisted/add-nodes.env'],
            'excluded_markers': ['/etc/assisted/interactive-ui', '/etc/assisted/rendezvous-host.env'],
        },
        'interactive': {
            'workflow_markers': ['/etc/assisted/interactive-ui', '/etc/assisted/rendezvous-host.env'],
            'excluded_markers': ['/etc/assisted/add-nodes.env'],
        },
        'unconfigured_ignition': {
            'workflow_markers': ['/etc/assisted/rendezvous-host.env'],
            'excluded_markers': ['/etc/assisted/add-nodes.env', '/etc/assisted/interactive-ui'],
        },
    }

    def _get_transitive_requirements(self, unit: SystemdUnit, units: Dict[str, SystemdUnit],
                                     visited: Set[str] = None) -> Set[str]:
        """Get all services transitively required by this unit."""
        if visited is None:
            visited = set()

        if unit.name in visited:
            return set()

        visited.add(unit.name)
        requirements = set()

        # Add direct requirements
        for req in unit.requires + unit.binds_to:
            if req in units:
                requirements.add(req)
                # Recursively get requirements of requirements
                requirements.update(self._get_transitive_requirements(units[req], units, visited))

        return requirements

    def filter_workflow(self, units: Dict[str, SystemdUnit], workflow: str) -> Set[str]:
        """Return set of service names for the given workflow.

        A service is included if:
        1. It has NO conditions on workflow discriminator files (runs in all workflows), OR
        2. It requires a workflow marker that matches this workflow, OR
        3. It excludes all markers that this workflow excludes (compatible negative conditions)
        4. It doesn't transitively require a service that is disabled in this workflow
        """
        config = self.WORKFLOWS[workflow]
        filtered = set()

        # First pass: determine which services are directly enabled/disabled
        for name, unit in units.items():
            # Skip system targets
            if '.target' in name:
                continue

            # Hard-coded exclusions based on information outside systemd units
            # load-config-iso@ is started by udev rule only in unconfigured_ignition
            if name == 'load-config-iso@.service' and workflow != 'unconfigured_ignition':
                continue
            # agent-check-config-image only runs with load-config-iso@
            if name == 'agent-check-config-image.service' and workflow != 'unconfigured_ignition':
                continue
            # agent-interactive-console services are disabled in Go code for unconfigured_ignition workflow
            if name in ('agent-interactive-console.service', 'agent-interactive-console-serial@.service'):
                if workflow == 'unconfigured_ignition':
                    continue

            # Check if service has conditions on workflow discriminator files
            discriminator_conditions_positive = set()
            discriminator_conditions_negative = set()

            for cond in unit.condition_path_exists:
                if cond in self.WORKFLOW_DISCRIMINATORS:
                    discriminator_conditions_positive.add(cond)

            for cond in unit.condition_path_not_exists:
                if cond in self.WORKFLOW_DISCRIMINATORS:
                    discriminator_conditions_negative.add(cond)

            # Case 1: Service has NO conditions on discriminator files - runs in all workflows
            if not discriminator_conditions_positive and not discriminator_conditions_negative:
                filtered.add(name)
                continue

            # Case 2: Service requires workflow markers - ALL positive conditions must be satisfied
            # For example, agent-extract-tui requires both rendezvous-host.env AND interactive-ui
            if discriminator_conditions_positive:
                # Check if ALL required discriminators are present in this workflow's markers
                all_satisfied = discriminator_conditions_positive.issubset(set(config['workflow_markers']))
                if all_satisfied:
                    filtered.add(name)
                    continue

            # Case 3: Service excludes markers, check if compatible with this workflow
            # Service is compatible if it doesn't require any markers this workflow excludes
            # and it doesn't exclude any markers this workflow requires
            if discriminator_conditions_negative:
                conflicts = False

                # Service requires markers that this workflow excludes?
                for req_marker in discriminator_conditions_positive:
                    if req_marker in config['excluded_markers']:
                        conflicts = True
                        break

                # Service excludes markers that this workflow requires?
                if not conflicts:
                    for req_marker in config['workflow_markers']:
                        if req_marker in discriminator_conditions_negative:
                            conflicts = True
                            break

                if not conflicts:
                    filtered.add(name)

        # Second pass: remove services that transitively require disabled services
        disabled = set(units.keys()) - filtered - {name for name in units if '.target' in name}
        to_remove = set()

        for name in filtered:
            unit = units[name]
            transitive_reqs = self._get_transitive_requirements(unit, units)
            # If any transitive requirement is disabled, this service is also disabled
            if transitive_reqs & disabled:
                to_remove.add(name)

        filtered -= to_remove

        return filtered


class GraphVizGenerator:
    """Generate GraphViz DOT files from systemd units."""

    # Styling configuration
    ORCHESTRATOR_SERVICES = {
        'agent-interactive-console.service',
        'agent-interactive-console-serial@.service',
        'load-config-iso@.service',
    }

    # Files to show in diagrams
    IMPORTANT_FILES = {
        '/usr/local/bin/agent-tui',
        '/usr/lib/dracut/hooks/pre-pivot/99-agent-copy-files.sh',
        '/etc/assisted/node0',
        '/etc/assisted/rendezvous-host.env',
    }

    def __init__(self, units: Dict[str, SystemdUnit], install_services: Set[str] = None):
        self.units = units
        self.install_services = install_services or set()

    def _find_reachable_from_pod(self, services: Set[str]) -> Set[str]:
        """Find all services reachable from assisted-service-pod via dependencies."""
        if 'assisted-service-pod.service' not in services:
            return set()

        reachable = set()
        to_visit = ['assisted-service-pod.service']

        while to_visit:
            current = to_visit.pop()
            if current in reachable:
                continue
            reachable.add(current)

            # Follow Before dependencies (what runs before this service)
            if current in self.units:
                unit = self.units[current]
                for dep in unit.before:
                    if dep in services and dep not in reachable:
                        to_visit.append(dep)

                # Follow After dependencies backwards (what this runs after)
                for dep in unit.after:
                    if dep in services and dep not in reachable:
                        to_visit.append(dep)

            # Also check who lists this service as Before/After
            for svc_name in services:
                if svc_name in reachable:
                    continue
                if svc_name not in self.units:
                    continue
                other_unit = self.units[svc_name]
                if current in other_unit.after or current in other_unit.before:
                    if svc_name not in reachable:
                        to_visit.append(svc_name)

        return reachable

    def _compute_workflow_differences(self, workflow: str, services: Set[str]) -> Set[str]:
        """Compute services that differ from install workflow."""
        if workflow != 'add_nodes' or not self.install_services:
            return set()

        # Services in add-nodes but not in install
        only_in_add_nodes = services - self.install_services

        # Services that have different conditions or are marked differently
        differences = set(only_in_add_nodes)

        # Special case: node-zero exists in both but gets an asterisk in add-nodes
        if 'node-zero.service' in services and 'node-zero.service' in self.install_services:
            differences.add('node-zero.service')

        return differences

    def generate_workflow(self, workflow: str, services: Set[str]) -> str:
        """Generate GraphViz DOT for a workflow."""
        lines = []

        # Header
        graph_name = f"agent_installer_services_{workflow}_workflow" if workflow != 'unconfigured_ignition' else "agent_installer_services_unconfigured_ignition"
        lines.append(f"digraph {graph_name} {{")
        lines.append("    rankdir=BT;")
        lines.append("    ranksep=0.5;")
        lines.append('    node [shape=box, style="rounded,filled", fillcolor="#ADD8E6", fontname="Arial", fontsize=10, penwidth=1];')
        lines.append('    edge [color="#333333"];')
        lines.append("")

        # Collect files referenced in conditions
        files = self._collect_files(services, workflow)

        # Special handling for unconfigured_ignition workflow
        if workflow == 'unconfigured_ignition' and 'load-config-iso@.service' in services:
            files.add('/etc/assisted/rendezvous-host.env')

        # Compute disconnected services and workflow differences
        reachable_from_pod = self._find_reachable_from_pod(services)
        disconnected = services - reachable_from_pod
        workflow_differences = self._compute_workflow_differences(workflow, services)

        # Group services by type
        foundation = self._get_foundation_services(services)
        initramfs = self._get_initramfs_files(files)
        clusters = self._find_clusters(services)
        # Flatten all cluster members to exclude from regular services
        cluster_members = set()
        for members in clusters.values():
            cluster_members.update(members)
        regular = services - foundation - cluster_members

        # Foundation services
        if foundation:
            lines.append("    // Bottom row - foundation services")
            lines.append("    {")
            lines.append('        node [fillcolor="#ADD8E6"];')
            for svc in sorted(foundation):
                label = svc.replace('.service', '')
                style = ""
                # Disconnected services get thin border
                if svc in disconnected:
                    style = ', penwidth=0.5'
                # Workflow differences get green styling
                elif svc in workflow_differences:
                    style = ', color="#006400", fontcolor="#006400", penwidth=2'
                    if svc == 'node-zero.service':
                        label += '*'
                lines.append(f'        {self._service_to_id(svc)} [label="{label}"{style}];')
            lines.append("    }")
            lines.append("")

        # Files
        if files:
            lines.append("    // Files (document style)")
            lines.append('    node [shape=note, fillcolor="#FFFACD"];')
            for file_path in sorted(files):
                file_id = self._file_to_id(file_path)
                label = file_path
                fillcolor = '#FFFACD'
                if '99-agent-copy-files' in file_path:
                    fillcolor = '#F5DEB3'
                    label = label.replace('/usr/lib/dracut/hooks/pre-pivot/', '')  # Shorten for display
                if 'rendezvous-host.env' in file_path:
                    label = '/etc/assisted/\\nrendezvous-host.env'
                lines.append(f'    {file_id} [label="{label}", fillcolor="{fillcolor}"];')
            lines.append("")

        # Regular services
        if regular:
            lines.append("    // Middle services")
            lines.append('    node [shape=box, style="rounded,filled", fillcolor="#ADD8E6", penwidth=1];')
            for svc in sorted(regular):
                label = svc.replace('.service', '')
                style = []
                if svc in self.ORCHESTRATOR_SERVICES:
                    style.append('fillcolor="#6495ED"')
                    style.append('penwidth=2')
                elif svc in workflow_differences:
                    style.append('color="#006400"')
                    style.append('fontcolor="#006400"')
                    style.append('penwidth=2')
                    if svc == 'node-zero.service':
                        label += '*'
                elif svc in disconnected:
                    style.append('penwidth=0.5')
                if svc == 'load-config-iso@.service':
                    style.append('width=2.0')

                style_str = ', ' + ', '.join(style) if style else ''
                lines.append(f'    {self._service_to_id(svc)} [label="{label}"{style_str}];')
            lines.append("")

        # Clusters (dynamically generated from PartOf relationships)
        for cluster_idx, (parent_svc, cluster_members) in enumerate(sorted(clusters.items())):
            parent_id = self._service_to_id(parent_svc)
            cluster_name = f"cluster_{parent_id}"

            lines.append(f"    // Cluster: {parent_svc.replace('.service', '')}")
            lines.append(f"    subgraph {cluster_name} {{")
            lines.append('        label="";')
            lines.append('        style=dashed;')
            lines.append('        color="#666666";')
            lines.append('        fillcolor="#FFFFFF";')
            lines.append("")

            for svc in cluster_members:
                label = svc.replace('.service', '')
                style = ', fillcolor="#ADD8E6"'
                if svc in workflow_differences:
                    style = ', fillcolor="#ADD8E6", color="#006400", fontcolor="#006400", penwidth=2'
                lines.append(f'        {self._service_to_id(svc)} [label="{label}"{style}];')
            lines.append("")

            # Invisible edges for layout (parent to children)
            if len(cluster_members) > 1:
                for child in cluster_members[1:3]:  # First 2 children
                    child_id = self._service_to_id(child)
                    lines.append(f"        {parent_id} -> {child_id} [style=invis];")

            lines.append("    }")
            lines.append("")

        # Dependencies
        lines.append("    // Dependencies (bottom to top flow)")
        lines.append("")
        lines.extend(self._generate_dependencies(services, files, workflow))

        # Rank constraints
        lines.append("    // Rank constraints for better layout")
        lines.extend(self._generate_rank_constraints(services, files, workflow))

        lines.append("}")

        return '\n'.join(lines)

    def _service_to_id(self, service: str) -> str:
        """Convert service name to GraphViz identifier."""
        return service.replace('.service', '').replace('@', '').replace('-', '_')

    def _file_to_id(self, file_path: str) -> str:
        """Convert file path to GraphViz identifier."""
        return file_path.replace('/', '_').replace('.', '_').replace('-', '_')

    def _collect_files(self, services: Set[str], workflow: str) -> Set[str]:
        """Collect important file paths referenced by services."""
        files = set()
        for svc_name in services:
            if svc_name not in self.units:
                continue
            unit = self.units[svc_name]
            for path in unit.condition_path_exists + unit.condition_path_not_exists:
                if path in self.IMPORTANT_FILES:
                    files.add(path)

        # Add dracut copy-files hook if agent-tui is present
        # In interactive workflow, agent-extract-tui provides agent-tui, not the copy-files hook
        if '/usr/local/bin/agent-tui' in files and workflow != 'interactive':
            files.add('/usr/lib/dracut/hooks/pre-pivot/99-agent-copy-files.sh')

        return files

    def _get_foundation_services(self, services: Set[str]) -> Set[str]:
        """Get foundation services (bottom row)."""
        foundation = {
            'selinux.service',
            'pre-network-manager-config.service',
            'set-hostname.service',
            'iscsistart.service',
            'iscsiadm.service',
            'agent-auth-token-status.service',  # For add-nodes
            'agent-extract-tui.service',  # For interactive
        }
        return foundation & services

    def _get_initramfs_files(self, files: Set[str]) -> Set[str]:
        """Get initramfs-related files."""
        return {f for f in files if 'dracut' in f or 'agent-tui' in f}

    def _find_clusters(self, services: Set[str]) -> Dict[str, List[str]]:
        """Find all clusters based on PartOf and BindsTo relationships.
        Returns dict mapping parent service -> list of child services (including parent itself)."""
        clusters = {}

        # Find all services that have other services with PartOf or BindsTo pointing to them
        for svc_name in services:
            if svc_name not in self.units:
                continue
            unit = self.units[svc_name]

            # Check PartOf relationships
            for part_of_target in unit.part_of:
                if part_of_target in services:
                    # This service is part of a cluster
                    if part_of_target not in clusters:
                        clusters[part_of_target] = [part_of_target]
                    if svc_name not in clusters[part_of_target]:
                        clusters[part_of_target].append(svc_name)

            # Check BindsTo relationships (similar to PartOf for clustering purposes)
            for binds_to_target in unit.binds_to:
                if binds_to_target in services:
                    # This service binds to a cluster parent
                    if binds_to_target not in clusters:
                        clusters[binds_to_target] = [binds_to_target]
                    if svc_name not in clusters[binds_to_target]:
                        clusters[binds_to_target].append(svc_name)

        # Sort members of each cluster for consistent output
        for parent in clusters:
            # Keep parent first, then alphabetically sort the rest
            members = clusters[parent]
            parent_item = [parent] if parent in members else []
            others = sorted([s for s in members if s != parent])
            clusters[parent] = parent_item + others

        return clusters

    def _get_pod_services(self, services: Set[str], workflow: str) -> List[str]:
        """Get services that belong in the pod cluster, inferred from PartOf/BindsTo.
        DEPRECATED: Use _find_clusters instead."""
        clusters = self._find_clusters(services)
        # Return the assisted-service-pod cluster if it exists
        return clusters.get('assisted-service-pod.service', [])

    def _generate_dependencies(self, services: Set[str], files: Set[str], workflow: str) -> List[str]:
        """Generate dependency edges."""
        lines = []
        edges_added = set()

        # Special handling for initramfs files
        if '/usr/lib/dracut/hooks/pre-pivot/99-agent-copy-files.sh' in files and '/usr/local/bin/agent-tui' in files:
            lines.append("    // File preparation (initramfs phase on the left)")
            copy_id = self._file_to_id('/usr/lib/dracut/hooks/pre-pivot/99-agent-copy-files.sh')
            tui_id = self._file_to_id('/usr/local/bin/agent-tui')
            edge1 = (copy_id, tui_id)
            if edge1 not in edges_added:
                lines.append(f"    {copy_id} -> {tui_id} [style=dotted, weight=10];")
                edges_added.add(edge1)
            if 'agent-interactive-console.service' in services:
                edge2 = (tui_id, 'agent_interactive_console')
                if edge2 not in edges_added:
                    lines.append(f"    {tui_id} -> agent_interactive_console [style=dotted, weight=1];")
                    edges_added.add(edge2)
            lines.append("")

        # Special handling for unconfigured_ignition workflow
        # load-config-iso creates rendezvous-host.env which is needed by agent and assisted-service
        if workflow == 'unconfigured_ignition' and 'load-config-iso@.service' in services:
            lines.append("    // Config image loading and file creation")
            rendezvous_file = '/etc/assisted/rendezvous-host.env'
            if rendezvous_file in files:
                rendezvous_id = self._file_to_id(rendezvous_file)
                # load-config-iso creates the file
                edge = ('load_config_iso', rendezvous_id)
                if edge not in edges_added:
                    lines.append(f"    load_config_iso -> {rendezvous_id} [style=dotted];")
                    edges_added.add(edge)
                # Services that need this file
                for svc in ['agent.service', 'assisted-service.service', 'node-zero.service']:
                    if svc in services:
                        svc_id = self._service_to_id(svc)
                        edge = (rendezvous_id, svc_id)
                        if edge not in edges_added:
                            lines.append(f"    {rendezvous_id} -> {svc_id} [style=dotted];")
                            edges_added.add(edge)
            lines.append("")

        # Show that node-zero creates /etc/assisted/node0
        if 'node-zero.service' in services and '/etc/assisted/node0' in files:
            node0_id = self._file_to_id('/etc/assisted/node0')
            edge = ('node_zero', node0_id)
            if edge not in edges_added:
                lines.append("    // File creation during workflow")
                lines.append(f"    node_zero -> {node0_id} [style=dotted];")
                edges_added.add(edge)
                lines.append("")

        # Show that agent-extract-tui creates /usr/local/bin/agent-tui
        if 'agent-extract-tui.service' in services and '/usr/local/bin/agent-tui' in files:
            tui_id = self._file_to_id('/usr/local/bin/agent-tui')
            edge = ('agent_extract_tui', tui_id)
            if edge not in edges_added:
                lines.append("    // TUI binary extraction in interactive workflow")
                lines.append(f"    agent_extract_tui -> {tui_id} [style=dotted];")
                edges_added.add(edge)
                lines.append("")

        # Service dependencies from systemd After= and Before= directives
        dep_lines = []
        for svc_name in sorted(services):
            if svc_name not in self.units:
                continue
            unit = self.units[svc_name]
            svc_id = self._service_to_id(svc_name)

            # After dependencies (reverse direction for bottom-up graph)
            for dep in unit.after:
                if dep.endswith('.service') and dep in services:
                    dep_id = self._service_to_id(dep)
                    edge = (dep_id, svc_id)
                    if edge not in edges_added:
                        dep_lines.append(f"    {dep_id} -> {svc_id};")
                        edges_added.add(edge)

            # Before dependencies (forward direction for bottom-up graph)
            for dep in unit.before:
                if dep.endswith('.service') and dep in services:
                    dep_id = self._service_to_id(dep)
                    edge = (svc_id, dep_id)
                    if edge not in edges_added:
                        dep_lines.append(f"    {svc_id} -> {dep_id};")
                        edges_added.add(edge)

            # File dependencies
            for path in unit.condition_path_exists:
                if path in files:
                    file_id = self._file_to_id(path)
                    edge = (file_id, svc_id)
                    if edge not in edges_added:
                        dep_lines.append(f"    {file_id} -> {svc_id} [style=dotted];")
                        edges_added.add(edge)

        if dep_lines:
            lines.append("    // Service dependencies")
            lines.extend(dep_lines)
            lines.append("")

        return lines

    def _generate_rank_constraints(self, services: Set[str], files: Set[str], workflow: str) -> List[str]:
        """Generate rank constraints for layout."""
        lines = []

        # Bottom row
        foundation = ['selinux', 'pre_network_manager_config', 'set_hostname', 'iscsistart']

        # Add initramfs file if present
        if '/usr/lib/dracut/hooks/pre-pivot/99-agent-copy-files.sh' in files:
            copy_id = self._file_to_id('/usr/lib/dracut/hooks/pre-pivot/99-agent-copy-files.sh')
            foundation.insert(0, copy_id)

        lines.append("    {rank=same; " + "; ".join(foundation) + ";}")

        # Agent-tui file
        if '/usr/local/bin/agent-tui' in files:
            tui_id = self._file_to_id('/usr/local/bin/agent-tui')
            lines.append(f"    {{rank=same; {tui_id};}}")

        # Interactive console / load-config-iso
        if 'agent-interactive-console.service' in services:
            lines.append("    {rank=same; agent_interactive_console;}")
        elif 'load-config-iso@.service' in services:
            lines.append("    {rank=same; load_config_iso;}")

        # Rendezvous-host.env
        if '/etc/assisted/rendezvous-host.env' in files:
            env_id = self._file_to_id('/etc/assisted/rendezvous-host.env')
            lines.append(f"    {{rank=same; {env_id};}}")

        # Agent, node-zero
        middle = []
        if 'agent.service' in services:
            middle.append('agent')
        if 'agent-check-config-image.service' in services:
            middle.append('agent_check_config_image')
        if 'node-zero.service' in services:
            middle.append('node_zero')
        if middle:
            lines.append("    {rank=same; " + "; ".join(middle) + ";}")

        # Invisible edges for layout to keep initramfs files on the left
        if '/usr/local/bin/agent-tui' in files:
            tui_id = self._file_to_id('/usr/local/bin/agent-tui')
            if 'agent-interactive-console.service' in services:
                lines.append(f"    agent_interactive_console -> {tui_id} [style=invis, constraint=false];")
            if 'agent-interactive-console-serial@.service' in services:
                lines.append(f"    agent_interactive_console_serial -> {tui_id} [style=invis, constraint=false];")

        # Force copy-files to be on the far left
        # Create invisible edges to establish left-to-right ordering in the bottom rank
        if '/usr/lib/dracut/hooks/pre-pivot/99-agent-copy-files.sh' in files:
            copy_id = self._file_to_id('/usr/lib/dracut/hooks/pre-pivot/99-agent-copy-files.sh')
            # Create chain: copy-files → selinux → pre-network-manager-config → ...
            lines.append(f"    {copy_id} -> selinux [style=invis];")
            # Continue the chain to establish full ordering
            lines.append(f"    selinux -> pre_network_manager_config [style=invis];")
            lines.append(f"    pre_network_manager_config -> set_hostname [style=invis];")
            if 'iscsistart.service' in services:
                lines.append(f"    set_hostname -> iscsistart [style=invis];")

        # Keep config-image files on the left in unconfigured_ignition workflow
        if workflow == 'unconfigured_ignition':
            if 'load-config-iso@.service' in services and '/etc/assisted/rendezvous-host.env' in files:
                env_id = self._file_to_id('/etc/assisted/rendezvous-host.env')
                lines.append(f"    load_config_iso -> {env_id} [style=invis, constraint=false];")

        return lines


def main():
    """Generate all workflow diagrams."""
    script_dir = Path(__file__).parent
    # Go up to project root and find units directory
    project_root = script_dir.parent.parent.parent.parent
    units_dir = project_root / "data" / "data" / "agent" / "systemd" / "units"

    if not units_dir.exists():
        print(f"Error: Units directory not found: {units_dir}")
        return 1

    print(f"Parsing systemd units from: {units_dir}")
    parser = SystemdParser(units_dir)
    parser.parse_all()
    print(f"Parsed {len(parser.units)} unit files")

    # Filter by workflow
    filter_engine = WorkflowFilter()

    workflows = {
        'install_workflow': 'install',
        'add_nodes_workflow': 'add_nodes',
        'interactive': 'interactive',
        'unconfigured_ignition_and_config_image_flow': 'unconfigured_ignition',
    }

    # First get install services for comparison
    install_services = filter_engine.filter_workflow(parser.units, 'install')
    generator = GraphVizGenerator(parser.units, install_services)

    for output_name, workflow_key in workflows.items():
        services = filter_engine.filter_workflow(parser.units, workflow_key)
        print(f"\nGenerating {output_name}: {len(services)} services")

        dot_content = generator.generate_workflow(workflow_key, services)
        output_file = script_dir / f"{output_name}.dot"

        with open(output_file, 'w') as f:
            f.write(dot_content)

        print(f"  Written to: {output_file}")

    print("\nGeneration complete!")
    print("Run 'make' to regenerate PNG files from the updated DOT sources")

    return 0


if __name__ == '__main__':
    exit(main())
