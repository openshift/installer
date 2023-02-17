#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest

import os
import sys
import glob
import yaml

ASSETS_DIR = ""
INSTALLCONFIG_PATH = ""

class FailureDomain:
    def __init__(self, failure_domain):
        self.failure_domain = failure_domain

    def __hash__(self):
        return hash(str(sorted(self.failure_domain.items())))

    def __eq__(self, other):
        return type(self) == type(other) and self.__hash__() == other.__hash__()

EMPTY_FAILURE_DOMAIN = FailureDomain({})

class TestMachinesets(unittest.TestCase):
    def setUp(self):
        """Parse the expected values from install-config and collect MachineSet resources."""
        self.machinesets = []
        for machineset_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_worker-machineset-*.yaml'
        ):
            with open(machineset_path) as f:
                self.machinesets.append(yaml.load(f, Loader=yaml.FullLoader))

        self.expected_failure_domains = set()
        with open(INSTALLCONFIG_PATH) as f:
            installconfig = yaml.load(f, Loader=yaml.FullLoader)
            installconfig_failure_domains = installconfig["compute"][0]["platform"]["openstack"]["failureDomains"]
            for failure_domain in installconfig_failure_domains:
                self.expected_failure_domains.add(FailureDomain(failure_domain))
            self.expected_machine_replicas = installconfig["compute"][0]["replicas"]

    def test_valid_failure_domain(self):
        """Assert that there is exactly one MachineSet per failure domain."""
        found = set()
        for machineset in self.machinesets:
            failure_domain = FailureDomain(machineset["spec"]["template"]["spec"]["providerSpec"][
                "value"]["failureDomain"])
            self.assertIn(failure_domain, self.expected_failure_domains)
            self.assertNotIn(failure_domain, found)
            self.assertNotEqual(failure_domain, EMPTY_FAILURE_DOMAIN)
            found.add(failure_domain)

        self.assertEqual(len(self.machinesets), len(self.expected_failure_domains))

    def test_total_replica_number(self):
        """Assert that replicas spread across the MachineSets add up to the expected number."""
        total_found = 0
        for machineset in self.machinesets:
            total_found += machineset["spec"]["replicas"]
        self.assertEqual(total_found, self.expected_machine_replicas)

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
    INSTALLCONFIG_PATH = os.path.join(os.path.dirname(os.path.realpath(__file__)), 'install-config.yaml')
    unittest.main(verbosity=2)
