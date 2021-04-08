#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest

import sys
import glob
import yaml

ASSETS_DIR = ""

EXPECTED_MACHINES_NUMBER = 1000
EXPECTED_COMPUTE_ZONE_NAMES = ["ComputeAZ1", "ComputeAZ2", "ComputeAZ3"]
EXPECTED_VOLUME_ZONE_NAMES = ["VolumeAZ1", "VolumeAZ2", "VolumeAZ3"]


class TestVolumeAZMachinesets(unittest.TestCase):
    def setUp(self):
        """Parse the MachineSets into a Python data structure."""
        self.machinesets = []
        for machineset_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_worker-machineset-*.yaml'
        ):
            with open(machineset_path) as f:
                self.machinesets.append(yaml.load(f, Loader=yaml.FullLoader))

    def test_machineset_zone_name(self):
        """Assert that there is exactly one MachineSet per volume availability zone."""
        found = []
        for machineset in self.machinesets:
            master_zone = machineset["spec"]["template"]["spec"]["providerSpec"]["value"]["availabilityZone"]
            volume_zone = machineset["spec"]["template"]["spec"]["providerSpec"]["value"]["rootVolume"]["availabilityZone"]
            self.assertIn(volume_zone, EXPECTED_VOLUME_ZONE_NAMES)
            self.assertNotIn(volume_zone, found)
            self.assertEqual(master_zone[-3:], volume_zone[-3:])
            found.append(volume_zone)
        self.assertEqual(len(self.machinesets), len(EXPECTED_VOLUME_ZONE_NAMES))

    def test_total_replica_number(self):
        """Assert that replicas spread across the MachineSets add up to the expected number."""
        total_found = 0
        for machineset in self.machinesets:
            total_found += machineset["spec"]["replicas"]
        self.assertEqual(total_found, EXPECTED_MACHINES_NUMBER)

    def test_replica_distribution(self):
        """Assert that replicas are evenly distributed across machinesets."""
        setpoint = 0
        for machineset in self.machinesets:
            replicas = machineset["spec"]["replicas"]
            if setpoint == 0:
                setpoint = replicas
            else:
                self.assertTrue(-2 < replicas - setpoint < 2)


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    unittest.main(verbosity=2)
