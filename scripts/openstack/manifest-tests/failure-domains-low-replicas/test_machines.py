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

    def hasPorts(self):
        return "ports" in self.failure_domain and len(self.failure_domain["ports"]) > 0


class TestMachines(unittest.TestCase):
    def setUp(self):
        """Parse the expected values from install-config and collect Machine resources."""
        self.machines = []
        for machine_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_master-machines-*.yaml'
        ):
            with open(machine_path) as f:
                self.machines.append(yaml.load(f, Loader=yaml.FullLoader))

        self.expected_failure_domains = set()
        self.expected_machine_replicas = 0
        with open(INSTALLCONFIG_PATH) as f:
            installconfig = yaml.load(f, Loader=yaml.FullLoader)
            installconfig_failure_domains = installconfig["controlPlane"]["platform"]["openstack"]["failureDomains"]
            for failure_domain in installconfig_failure_domains:
                self.expected_failure_domains.add(FailureDomain(failure_domain))
            self.expected_machine_replicas = installconfig["controlPlane"]["replicas"]

    def test_valid_failure_domain(self):
        """Assert that all machines have one valid failure domain."""
        for machine in self.machines:
            failure_domain = FailureDomain(machine["spec"]["providerSpec"]["value"]["failureDomain"])
            self.assertIn(failure_domain, self.expected_failure_domains)

    def test_total_instance_number(self):
        """Assert that there are as many Machines as required ControlPlane replicas."""
        self.assertEqual(len(self.machines), self.expected_machine_replicas)

    def test_failure_domain_number(self):
        """Assert that there are as many failure domains as required."""
        failure_domains = set()
        for machine in self.machines:
            failure_domains.add(FailureDomain(machine["spec"]["providerSpec"]["value"]["failureDomain"]))

        # Since we extract failure domains from machines, there can't be more
        # failure domains than machines in the cluster
        if self.expected_machine_replicas >= len(self.expected_failure_domains):
            self.assertEqual(len(failure_domains), len(self.expected_failure_domains))
        else:
            self.assertEqual(len(failure_domains), self.expected_machine_replicas)

    def test_replica_distribution(self):
        """Assert that machines are evenly distributed across failure domains."""
        failure_domains = {}
        for machine in self.machines:
            failure_domain = FailureDomain(machine["spec"]["providerSpec"]["value"]["failureDomain"])
            failure_domains[failure_domain] = failure_domains.get(failure_domain, 0) + 1

        setpoint = 0
        for replicas in failure_domains.values():
            if setpoint == 0:
                setpoint = replicas
            else:
                self.assertTrue(-2 < replicas - setpoint < 2)

    def test_machines_primary_network(self):
        """Assert that machines with ports in their failure domain don't have a primaryNetwork set."""
        for machine in self.machines:
            if FailureDomain(machine["spec"]["providerSpec"]["value"]["failureDomain"]).hasPorts():
                self.assertNotIn("primaryNetwork", machine["spec"]["providerSpec"]["value"])


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    INSTALLCONFIG_PATH = os.path.join(os.path.dirname(os.path.realpath(__file__)), 'install-config.yaml')
    unittest.main(verbosity=2)
