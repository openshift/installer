#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""Test that bootstrap inherits the flavor from defaultMachinePlatform.type.

Verifies SC-002 (Bootstrap Inherits from defaultMachinePlatform): when neither
bootstrapFlavor nor an explicit controlPlane.platform.openstack.type is set, but
platform.openstack.defaultMachinePlatform.type is configured, the bootstrap
OpenStackMachine must use the inherited defaultMachinePlatform value.
"""

import unittest
import xmlrunner

import os
import sys
import glob
import yaml

ASSETS_DIR = ""


class BootstrapFlavorInherited(unittest.TestCase):
    def setUp(self):
        """Parse the OpenStackMachine manifests and install-config into Python data structures."""
        # Read install-config from cluster-config.yaml to get expected flavors.
        with open(f'{ASSETS_DIR}/manifests/cluster-config.yaml') as f:
            cluster_config = yaml.load(f, Loader=yaml.FullLoader)
            self.install_config = yaml.load(
                cluster_config["data"]["install-config"], Loader=yaml.FullLoader
            )

        # Verify that bootstrapFlavor is absent from the install-config — this
        # test validates inheritance from defaultMachinePlatform, not an
        # explicit bootstrap flavor.
        bootstrap_flavor = (
            self.install_config
            .get("platform", {})
            .get("openstack", {})
            .get("bootstrapFlavor")
        )
        self.assertIsNone(
            bootstrap_flavor,
            f"Expected bootstrapFlavor to be absent from install-config, "
            f"but got {bootstrap_flavor!r} — this test requires bootstrapFlavor to be unset"
        )

        # Verify that no explicit controlPlane openstack type is set — the
        # flavor must come purely from defaultMachinePlatform.type.
        cp_type = (
            self.install_config
            .get("controlPlane", {})
            .get("platform", {})
            .get("openstack", {})
            .get("type")
        )
        self.assertIsNone(
            cp_type,
            f"Expected controlPlane.platform.openstack.type to be absent, "
            f"but got {cp_type!r} — this test requires no explicit control plane type"
        )

        # Verify that defaultMachinePlatform.type is set — this is the source
        # of the inherited flavor.
        self.default_flavor = (
            self.install_config
            .get("platform", {})
            .get("openstack", {})
            .get("defaultMachinePlatform", {})
            .get("type")
        )
        self.assertIsNotNone(
            self.default_flavor,
            "Expected platform.openstack.defaultMachinePlatform.type to be set, "
            "but it is absent — this is required for the inheritance test"
        )

        # Parse the CAPI OpenStackMachine for bootstrap (10_inframachine_*-bootstrap.yaml).
        bootstrap_machines = glob.glob(
            f'{ASSETS_DIR}/openshift/10_inframachine_*-bootstrap.yaml'
        )
        self.assertEqual(
            len(bootstrap_machines), 1,
            f"Expected exactly one bootstrap OpenStackMachine, found: {bootstrap_machines}"
        )
        with open(bootstrap_machines[0]) as f:
            self.bootstrap_machine = yaml.load(f, Loader=yaml.FullLoader)

        # Parse the CAPI OpenStackMachine manifests for masters (10_inframachine_*-master-*.yaml).
        self.master_machines = []
        for machine_path in glob.glob(
                f'{ASSETS_DIR}/openshift/10_inframachine_*-master-*.yaml'
        ):
            with open(machine_path) as f:
                self.master_machines.append(yaml.load(f, Loader=yaml.FullLoader))

        self.assertGreater(
            len(self.master_machines), 0,
            "Expected at least one master OpenStackMachine"
        )

    def test_bootstrap_uses_default_machine_platform_flavor(self):
        """Assert that the bootstrap OpenStackMachine uses defaultMachinePlatform.type.

        When neither bootstrapFlavor nor controlPlane.platform.openstack.type is
        set, the bootstrap machine must inherit from
        platform.openstack.defaultMachinePlatform.type via the
        ResolveBootstrapFlavor fallback chain.
        """
        actual_flavor = self.bootstrap_machine["spec"]["flavor"]
        self.assertEqual(
            actual_flavor, self.default_flavor,
            f"Bootstrap machine flavor {actual_flavor!r} does not match "
            f"defaultMachinePlatform.type {self.default_flavor!r} — bootstrap "
            f"should inherit the default machine platform flavor when neither "
            f"bootstrapFlavor nor an explicit control plane type is configured"
        )

    def test_masters_use_default_machine_platform_flavor(self):
        """Assert that all master OpenStackMachines use the defaultMachinePlatform.type flavor.

        With no explicit controlPlane.platform.openstack.type, master machines
        must also resolve to defaultMachinePlatform.type.
        """
        for machine in self.master_machines:
            actual_flavor = machine["spec"]["flavor"]
            self.assertEqual(
                actual_flavor, self.default_flavor,
                f"Master machine {machine['metadata']['name']} flavor {actual_flavor!r} does not "
                f"match defaultMachinePlatform.type {self.default_flavor!r}"
            )

    def test_bootstrap_and_master_use_same_flavor(self):
        """Assert that bootstrap and master machines use the same inherited flavor.

        With defaultMachinePlatform.type as the single source of truth and no
        explicit overrides, both roles must resolve to the same flavor value.
        """
        bootstrap_flavor = self.bootstrap_machine["spec"]["flavor"]
        for machine in self.master_machines:
            master_flavor = machine["spec"]["flavor"]
            self.assertEqual(
                bootstrap_flavor, master_flavor,
                f"Bootstrap flavor {bootstrap_flavor!r} differs from master "
                f"{machine['metadata']['name']} flavor {master_flavor!r} — both "
                f"should inherit defaultMachinePlatform.type when no explicit "
                f"flavors are configured"
            )

    def test_inheritance_chain_resolves_to_default(self):
        """Assert the inheritance chain: defaultMachinePlatform -> bootstrap flavor.

        This test makes the inheritance chain explicit: the expected flavor is
        read directly from defaultMachinePlatform.type and compared against the
        bootstrap machine spec, verifying the full resolution path.
        """
        bootstrap_flavor = self.bootstrap_machine["spec"]["flavor"]
        self.assertEqual(
            bootstrap_flavor, self.default_flavor,
            f"Inheritance chain broken: platform.openstack.defaultMachinePlatform.type "
            f"is {self.default_flavor!r} but bootstrap OpenStackMachine spec.flavor "
            f"is {bootstrap_flavor!r}"
        )


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    with open(os.environ.get('JUNIT_FILE', '/dev/null'), 'wb') as output:
        unittest.main(
            testRunner=xmlrunner.XMLTestRunner(output=output),
            failfast=False,
            buffer=False,
            catchbreak=False,
            verbosity=2,
        )
