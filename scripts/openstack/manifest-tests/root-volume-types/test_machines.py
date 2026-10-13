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
EXPECTED_MASTER_ROOT_VOLUME_TYPE_NAMES = ["type-1", "type-2", "type-3"]
EXPECTED_MASTER_ROOT_VOLUME_ZONE_NAMES = ["VolumeAZ1", "VolumeAZ2", "VolumeAZ3"]
EXPECTED_MASTER_ZONE_NAMES = ["MasterAZ1", "MasterAZ2", "MasterAZ3"]

EXPECTED_WORKER_REPLICAS = 1000
EXPECTED_WORKER_ROOT_VOLUME_TYPES = ["type-A", "type-B", "type-C"]



class RootVolumeTypeMachines(unittest.TestCase):
    def setUp(self):
        """Parse the Machines into a Python data structure."""
        self.machines = []
        for machine_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_master-machines-*.yaml'
        ):
            with open(machine_path) as f:
                self.machines.append(yaml.load(f, Loader=yaml.FullLoader))

    def test_zone_names(self):
        """Assert that all machines have one valid compute az that matches volume az."""
        for machine in self.machines:
            master_zone = machine["spec"]["providerSpec"]["value"]["availabilityZone"]
            volume_zone = machine["spec"]["providerSpec"]["value"]["rootVolume"]["availabilityZone"]
            self.assertIn(master_zone, EXPECTED_MASTER_ZONE_NAMES)
            self.assertIn(volume_zone, EXPECTED_MASTER_ROOT_VOLUME_ZONE_NAMES)
            self.assertEqual(master_zone[-3:], volume_zone[-3:])

    def test_type_names(self):
        """Assert that all machines have one valid volume type."""
        for machine in self.machines:
            volume_type = machine["spec"]["providerSpec"]["value"]["rootVolume"]["volumeType"]
            self.assertIn(volume_type, EXPECTED_MASTER_ROOT_VOLUME_TYPE_NAMES)

    def test_total_instance_number(self):
        """Assert that there are as many Machines as required ControlPlane replicas."""
        self.assertEqual(len(self.machines), EXPECTED_MASTER_REPLICAS)

    def test_replica_distribution(self):
        """Assert that machines are evenly distributed across failure domains."""
        volume_zones = {}
        volume_types = {}
        for machine in self.machines:
            volume_zone = machine["spec"]["providerSpec"]["value"]["rootVolume"]["availabilityZone"]
            volume_zones[volume_zone] = volume_zones.get(volume_zone, 0) + 1
            volume_type = machine["spec"]["providerSpec"]["value"]["rootVolume"]["volumeType"]
            volume_types[volume_type] = volume_types.get(volume_zone, 0) + 1

        setpoint = 0
        for replicas in volume_zones.values():
            if setpoint == 0:
                setpoint = replicas
            else:
                self.assertTrue(-2 < replicas - setpoint < 2, msg="volume zone mismatch")

        setpoint = 0
        for replicas in volume_types.values():
            if setpoint == 0:
                setpoint = replicas
            else:
                self.assertTrue(-2 < replicas - setpoint < 2, msg="volume type mismatch")


class RootVolumeTypesMachinesets(unittest.TestCase):
    def setUp(self):
        """Parse the MachineSets into a Python data structure."""
        self.machinesets = []
        for machineset_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_worker-machineset-*.yaml'
        ):
            with open(machineset_path) as f:
                self.machinesets.append(yaml.load(f, Loader=yaml.FullLoader))

    def test_machineset_zone_name(self):
        """Assert that there are as many MachineSets as failure domains."""
        self.assertEqual(len(self.machinesets), len(EXPECTED_WORKER_ROOT_VOLUME_TYPES))

    def test_machineset_volume_type_name(self):
        """Assert that each MachineSet has a valid volume type."""
        for machineset in self.machinesets:
            volume_type = machineset["spec"]["template"]["spec"]["providerSpec"]["value"]["rootVolume"]["volumeType"]
            self.assertIn(volume_type, EXPECTED_WORKER_ROOT_VOLUME_TYPES)

    def test_total_replica_number(self):
        """Assert that replicas spread across the MachineSets add up to the expected number."""
        total_found = 0
        for machineset in self.machinesets:
            total_found += machineset["spec"]["replicas"]
        self.assertEqual(total_found, EXPECTED_WORKER_REPLICAS)

    def test_replica_distribution(self):
        """Assert that MachineSets are evenly distributed across failure domains."""
        setpoint = 0
        for machineset in self.machinesets:
            replicas = machineset["spec"]["replicas"]
            if setpoint == 0:
                setpoint = replicas
            else:
                self.assertTrue(-2 < replicas - setpoint < 2)


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    with open(os.environ.get('JUNIT_FILE', '/dev/null'), 'wb') as output:
        unittest.main(testRunner=xmlrunner.XMLTestRunner(output=output), failfast=False, buffer=False, catchbreak=False, verbosity=2)
