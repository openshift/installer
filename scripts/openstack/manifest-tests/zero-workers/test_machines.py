#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest
import xmlrunner

import os
import sys
import glob
import yaml

ASSETS_DIR = ""

EXPECTED_MASTER_REPLICAS = 3
EXPECTED_WORKER_REPLICAS = 0


class Machines(unittest.TestCase):
    def setUp(self):
        """Parse the Machines into a Python data structure."""
        self.machines = []
        for machine_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_master-machines-*.yaml'
        ):
            with open(machine_path) as f:
                self.machines.append(yaml.load(f, Loader=yaml.FullLoader))

    def test_total_instance_number(self):
        """Assert that there are as many Machines as required ControlPlane replicas."""
        self.assertEqual(len(self.machines), EXPECTED_MASTER_REPLICAS)


class Machinesets(unittest.TestCase):
    def setUp(self):
        """Parse the MachineSets into a Python data structure."""
        self.machinesets = []
        for machineset_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_worker-machineset-*.yaml'
        ):
            with open(machineset_path) as f:
                self.machinesets.append(yaml.load(f, Loader=yaml.FullLoader))

    def test_total_replica_number(self):
        """Assert that there is at least one MachineSet."""
        self.assertGreater(len(self.machinesets), 0)

    def test_total_replica_number(self):
        """Assert that replicas spread across the MachineSets add up to the expected number."""
        total_found = 0
        for machineset in self.machinesets:
            total_found += machineset["spec"]["replicas"]
        self.assertEqual(total_found, EXPECTED_WORKER_REPLICAS)


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    with open(os.environ.get('JUNIT_FILE', '/dev/null'), 'wb') as output:
        unittest.main(testRunner=xmlrunner.XMLTestRunner(output=output), failfast=False, buffer=False, catchbreak=False, verbosity=2)
