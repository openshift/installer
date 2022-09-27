#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest

import sys
import glob
import yaml

ASSETS_DIR = ""

class TestClusterInfraObject(unittest.TestCase):
    def setUp(self):
        """Parse the Cluster Infrastructure object into a Python data structure."""
        self.machines = []
        cluster_infra = f'{ASSETS_DIR}/manifests/cluster-infrastructure-02-config.yaml'
        with open(cluster_infra) as f:
            self.cluster_infra = yaml.load(f, Loader=yaml.FullLoader)


    def failure_domains(self):
        """Assert that the Cluster infrastructure object contains failure domain information."""
        failureDomains = machine["spec"]["platformSpec"]["openstack"]["availabilityZone"]
        failureDomainNames = []

        for machine in self.machines:
            computeZone = machine["spec"]["platformSpec"]["openstack"]["availabilityZone"]
            computeZones[computeZone] = computeZones.get(computeZone, 0) + 1
            if computeZones[computeZone] > highest_number_of_replicas:
                highest_number_of_replicas = computeZones[computeZone]

        self.assertIn("nova-1", computeZones)
        self.assertIn("nova-2", computeZones)

        for name in computeZones:
            replicas = computeZones[name]
            self.assertTrue(highest_number_of_replicas - replicas < 2, msg=f'Zone {name} has too few machines: {replicas} (highest number found: {highest_number_of_replicas})')


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    unittest.main(verbosity=2)
