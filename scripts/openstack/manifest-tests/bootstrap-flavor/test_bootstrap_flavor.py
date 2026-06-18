#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""Test that bootstrapFlavor generates correct OpenStackMachine manifests.

Verifies SC-001 (Bootstrap with Custom Flavor): when bootstrapFlavor is set
explicitly in the install-config, the bootstrap OpenStackMachine uses that
flavor while control plane machines use the controlPlane type flavor.
"""

import unittest
import xmlrunner

import os
import sys
import glob
import yaml

ASSETS_DIR = ""


class BootstrapFlavorMachine(unittest.TestCase):
    def setUp(self):
        """Parse the OpenStackMachine manifests and install-config into Python data structures."""
        # Read install-config from cluster-config.yaml to get expected flavors.
        with open(f'{ASSETS_DIR}/manifests/cluster-config.yaml') as f:
            cluster_config = yaml.load(f, Loader=yaml.FullLoader)
            self.install_config = yaml.load(
                cluster_config["data"]["install-config"], Loader=yaml.FullLoader
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

    def test_bootstrap_uses_bootstrap_flavor(self):
        """Assert that the bootstrap OpenStackMachine uses the bootstrapFlavor value."""
        expected_bootstrap_flavor = self.install_config["platform"]["openstack"]["bootstrapFlavor"]
        actual_flavor = self.bootstrap_machine["spec"]["flavor"]
        self.assertEqual(
            actual_flavor, expected_bootstrap_flavor,
            f"Bootstrap machine flavor {actual_flavor!r} does not match "
            f"bootstrapFlavor {expected_bootstrap_flavor!r}"
        )

    def test_masters_use_control_plane_flavor(self):
        """Assert that all master OpenStackMachines use the controlPlane type flavor."""
        expected_master_flavor = self.install_config["controlPlane"]["platform"]["openstack"]["type"]
        for machine in self.master_machines:
            actual_flavor = machine["spec"]["flavor"]
            self.assertEqual(
                actual_flavor, expected_master_flavor,
                f"Master machine {machine['metadata']['name']} flavor {actual_flavor!r} does not "
                f"match controlPlane type {expected_master_flavor!r}"
            )

    def test_bootstrap_and_master_flavors_differ(self):
        """Assert that bootstrap and control plane machines use different flavors."""
        bootstrap_flavor = self.bootstrap_machine["spec"]["flavor"]
        for machine in self.master_machines:
            master_flavor = machine["spec"]["flavor"]
            self.assertNotEqual(
                bootstrap_flavor, master_flavor,
                f"Bootstrap and master machine {machine['metadata']['name']} use the same "
                f"flavor {bootstrap_flavor!r}, but they should differ when bootstrapFlavor "
                f"is explicitly configured"
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
