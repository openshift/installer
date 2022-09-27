#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest

import sys
import glob
import yaml

ASSETS_DIR = ""


class TestMachinesFailureDomain(unittest.TestCase):
    def setUp(self):
        """Parse the Machines into a Python data structure."""
        self.machines = []
        for machine_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_master-machines-*.yaml'
        ):
            with open(machine_path) as f:
                self.machines.append(yaml.load(f, Loader=yaml.FullLoader))


    def test_masters_spread(self):
        """Assert that Control plane machines are spread across provided compute zones."""
        computeZones = {}
        highest_number_of_replicas = 0
        for machine in self.machines:
            computeZone = machine["spec"]["providerSpec"]["value"]["availabilityZone"]
            computeZones[computeZone] = computeZones.get(computeZone, 0) + 1
            if computeZones[computeZone] > highest_number_of_replicas:
                highest_number_of_replicas = computeZones[computeZone]

        self.assertIn("nova-1", computeZones)
        self.assertIn("nova-2", computeZones)

        for name in computeZones:
            replicas = computeZones[name]
            self.assertTrue(highest_number_of_replicas - replicas < 2, msg=f'Zone {name} has too few machines: {replicas} (highest number found: {highest_number_of_replicas})')

    def test_subnet_id(self):
        """Assert that the machines in rack-2 have the proper subnetID set."""
        machines_in_rack_2 = []
        for machine in self.machines:
            computeZone = machine["spec"]["providerSpec"]["value"]["availabilityZone"]
            if computeZone == "rack-2":
                machines_in_rack_2.append(machine)

        for machine in machines_in_rack_2:
            network = machine["spec"]["providerSpec"]["value"].get("network")
            self.assertIsNotNone(network)
            subnets = network.get("subnets")
            self.assertTrue(len(subnets)>0)
            self.assertEqual(subnets[0].get("uuid"), "d7ffd3d8-87e1-481f-a818-ff2f17787d40")


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    unittest.main(verbosity=2)
