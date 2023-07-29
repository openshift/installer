#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest
import xmlrunner

import os
import sys
import glob
import yaml

ASSETS_DIR = ""

class ConvertMachine(unittest.TestCase):
    def setUp(self):
        """Parse the Machines into a Python data structure."""
        self.masters = []
        for machine_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_master-machines-*.yaml'
        ):
            with open(machine_path) as f:
                self.masters.append(yaml.load(f, Loader=yaml.FullLoader))

        with open(f'{ASSETS_DIR}/manifests/cluster-config.yaml') as f:
            cluster_config = yaml.load(f, Loader=yaml.FullLoader)
            self.install_config = yaml.load(cluster_config["data"]["install-config"], Loader=yaml.FullLoader)

    def test_flavor(self):
        """Assert that all machines take flavor from computeFlavor."""
        for machine in self.masters:
            master_flavor = machine["spec"]["providerSpec"]["value"]["flavor"]
            expected_master_flavor = self.install_config["platform"]["openstack"]["computeFlavor"]
            self.assertEqual(master_flavor, expected_master_flavor)

class ConvertMachineSet(unittest.TestCase):
    def setUp(self):
        """Parse the MachineSets into a Python data structure."""
        self.machinesets = []
        for machineset_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_worker-machineset-*.yaml'
        ):
            with open(machineset_path) as f:
                self.machinesets.append(yaml.load(f, Loader=yaml.FullLoader))

        with open(f'{ASSETS_DIR}/manifests/cluster-config.yaml') as f:
            cluster_config = yaml.load(f, Loader=yaml.FullLoader)
            self.install_config = yaml.load(cluster_config["data"]["install-config"], Loader=yaml.FullLoader)

    def test_flavor(self):
        """Assert that worker machinesets take flavor from machinepool."""
        for machineset in self.machinesets:
            flavor = machineset["spec"]["template"]["spec"]["providerSpec"]["value"]["flavor"]
            compute_flavor = self.install_config["platform"]["openstack"]["computeFlavor"]
            expected_flavor = self.install_config["compute"][0]["platform"]["openstack"]["type"]
            self.assertEqual(flavor, expected_flavor)
            self.assertNotEqual(flavor, compute_flavor)


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    with open(os.environ.get('JUNIT_FILE', '/dev/null'), 'wb') as output:
        unittest.main(testRunner=xmlrunner.XMLTestRunner(output=output), failfast=False, buffer=False, catchbreak=False, verbosity=2)
