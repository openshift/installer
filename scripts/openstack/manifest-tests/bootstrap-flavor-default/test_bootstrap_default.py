#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""Test that bootstrap defaults to the control plane flavor when bootstrapFlavor is unset.

Verifies SC-002 (Bootstrap Defaults to Control Plane Flavor): when
bootstrapFlavor is not specified in the install-config, the bootstrap
OpenStackMachine uses the same flavor as the control plane machines.
"""

import unittest
import xmlrunner

import os
import sys
import glob
import yaml

ASSETS_DIR = ""


class BootstrapFlavorDefault(unittest.TestCase):
    def setUp(self):
        """Parse the OpenStackMachine manifests and install-config into Python data structures."""
        # Read install-config from cluster-config.yaml to get expected flavors.
        with open(f'{ASSETS_DIR}/manifests/cluster-config.yaml') as f:
            cluster_config = yaml.load(f, Loader=yaml.FullLoader)
            self.install_config = yaml.load(
                cluster_config["data"]["install-config"], Loader=yaml.FullLoader
            )

        # Verify that bootstrapFlavor is absent from the install-config (the
        # whole point of this test case).
        bootstrap_flavor = (
            self.install_config
            .get("platform", {})
            .get("openstack", {})
            .get("bootstrapFlavor")
        )
        self.assertIsNone(
            bootstrap_flavor,
            f"Expected bootstrapFlavor to be absent from install-config, "
            f"but got {bootstrap_flavor!r}"
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

    def test_bootstrap_uses_control_plane_flavor(self):
        """Assert that the bootstrap OpenStackMachine uses the control plane flavor.

        When bootstrapFlavor is not set, the bootstrap machine must fall back
        to controlPlane.platform.openstack.type.
        """
        expected_flavor = self.install_config["controlPlane"]["platform"]["openstack"]["type"]
        actual_flavor = self.bootstrap_machine["spec"]["flavor"]
        self.assertEqual(
            actual_flavor, expected_flavor,
            f"Bootstrap machine flavor {actual_flavor!r} does not match "
            f"controlPlane type {expected_flavor!r} — bootstrap should default "
            f"to the control plane flavor when bootstrapFlavor is not specified"
        )

    def test_masters_use_control_plane_flavor(self):
        """Assert that all master OpenStackMachines use the controlPlane type flavor."""
        expected_flavor = self.install_config["controlPlane"]["platform"]["openstack"]["type"]
        for machine in self.master_machines:
            actual_flavor = machine["spec"]["flavor"]
            self.assertEqual(
                actual_flavor, expected_flavor,
                f"Master machine {machine['metadata']['name']} flavor {actual_flavor!r} does not "
                f"match controlPlane type {expected_flavor!r}"
            )

    def test_bootstrap_and_master_use_same_flavor(self):
        """Assert that bootstrap and control plane machines use the same flavor.

        With no bootstrapFlavor set, both should resolve to
        controlPlane.platform.openstack.type.
        """
        expected_flavor = self.install_config["controlPlane"]["platform"]["openstack"]["type"]
        bootstrap_flavor = self.bootstrap_machine["spec"]["flavor"]

        self.assertEqual(
            bootstrap_flavor, expected_flavor,
            f"Bootstrap machine flavor {bootstrap_flavor!r} does not match "
            f"expected control plane flavor {expected_flavor!r}"
        )

        for machine in self.master_machines:
            master_flavor = machine["spec"]["flavor"]
            self.assertEqual(
                bootstrap_flavor, master_flavor,
                f"Bootstrap flavor {bootstrap_flavor!r} differs from master machine "
                f"{machine['metadata']['name']} flavor {master_flavor!r} — both should "
                f"use the control plane flavor when bootstrapFlavor is not specified"
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
