#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest

import sys
import glob
import yaml

ASSETS_DIR = ""

EXPECTED_MACHINES_NUMBER = 10
EXPECTED_ZONE_NAMES = ["masterzone", "masterztwo", "masterzthree"]


class TestAZMachines(unittest.TestCase):
    def setUp(self):
        """Parse the Machines into a Python data structure."""
        self.machines = []
        for machine_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_master-machines-*.yaml'
        ):
            with open(machine_path) as f:
                self.machines.append(yaml.load(f, Loader=yaml.FullLoader))

    def test_zone_names(self):
        """Assert that all machines have one valid availability zone."""
        for machine in self.machines:
            zone = machine["spec"]["providerSpec"]["value"]["availabilityZone"]
            self.assertIn(zone, EXPECTED_ZONE_NAMES)

    def test_total_instance_number(self):
        """Assert that there are as many Machines as required ControlPlane replicas."""
        self.assertEqual(len(self.machines), EXPECTED_MACHINES_NUMBER)

    def test_replica_distribution(self):
        """Assert that machines are evenly distributed across zones."""
        zones = {}
        for machine in self.machines:
            zone = machine["spec"]["providerSpec"]["value"]["availabilityZone"]
            zones[zone] = zones.get(zone, 0) + 1

        setpoint = 0
        for replicas in zones.values():
            if setpoint == 0:
                setpoint = replicas
            else:
                self.assertTrue(-2 < replicas - setpoint < 2)


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    unittest.main(verbosity=2)
