#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest

import os
import sys
import glob
import yaml

ASSETS_DIR = ""
INSTALLCONFIG_PATH = ""


class TestMachines(unittest.TestCase):
    def setUp(self):
        """Parse the expected values from install-config and collect Machine resources."""
        self.machines = []
        for machine_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_master-machines-*.yaml'
        ):
            with open(machine_path) as f:
                self.machines.append(yaml.load(f, Loader=yaml.FullLoader))

        self.expected_failure_domains = []
        with open(INSTALLCONFIG_PATH) as f:
            installconfig = yaml.load(f, Loader=yaml.FullLoader)
            installconfig_failure_domains = installconfig["controlPlane"]["platform"]["openstack"]["failureDomains"]
            for failure_domain in installconfig_failure_domains:
                self.expected_failure_domains.append(failure_domain)
            self.expected_machine_replicas = installconfig["controlPlane"]["replicas"]


    def test_total_instance_number(self):
        """Assert that there are as many Machines as required ControlPlane replicas."""
        self.assertEqual(len(self.machines), self.expected_machine_replicas)


    def test_control_plane_ports(self):
        """Assert that "control-plane" ports are listed as the first network."""

        # Build a dict of expected networks in the form {network_uuid: [subnet_uuids]}
        expected_control_plane_networks = {}
        for failure_domain in self.expected_failure_domains:
            for portTarget in failure_domain["portTargets"]:
                if portTarget["id"] == "control-plane":
                    expected_control_plane_network_uuid = portTarget["network"]["id"]
                    expected_control_plane_network_subnet_uuids = [fixed_ip["subnet"]["id"] for fixed_ip in portTarget["fixedIPs"]]
                    expected_control_plane_networks[expected_control_plane_network_uuid] = expected_control_plane_network_subnet_uuids

        # Build a dict of actual networks in the form {network_uuid: [subnet_uuids]}
        actual_first_networks = {}
        for machine in self.machines:
            actual_first_network = machine["spec"]["providerSpec"]["value"]["networks"][0]
            actual_first_network_uuid = actual_first_network.get("uuid")
            actual_first_network_subnet_uuids = [fixed_ip.get("filter", {}).get("id") for fixed_ip in actual_first_network.get("subnets", [{}])]
            actual_first_networks[actual_first_network_uuid] = actual_first_network_subnet_uuids

        # Assert that all expected networks and subnets are represented in machines
        for expected_control_plane_network, expected_subnet_uuids in expected_control_plane_networks.items():
            self.assertIn(expected_control_plane_network, actual_first_networks)
            actual_subnet_uuids = actual_first_networks[expected_control_plane_network]
            for expected_subnet_uuid in expected_subnet_uuids:
                self.assertIn(expected_subnet_uuid, actual_subnet_uuids)

        # Assert that all actual networks and subnets were among the expected values
        for actual_first_network, actual_subnet_uuids in actual_first_networks.items():
            self.assertIn(actual_first_network, expected_control_plane_networks)
            expected_subnet_uuids = expected_control_plane_networks[actual_first_network]
            for actual_subnet_uuid in actual_subnet_uuids:
                self.assertIn(actual_subnet_uuid, expected_subnet_uuids)


    def test_additional_ports(self):
        """Assert that additional ports are listed in the machines networks."""
        expected_additional_network_uuid_sets = []
        for failure_domain in self.expected_failure_domains:
            expected_additional_network_uuid_sets.append({portTarget["network"]["id"] for portTarget in failure_domain["portTargets"] if portTarget["id"] != "control-plane"})

        actual_additional_network_uuid_sets = []
        for machine in self.machines:
            actual_additional_network_uuid_sets.append({network.get("uuid") for network in machine["spec"]["providerSpec"]["value"]["networks"][1:]})

        for actual_additional_network_uuids in actual_additional_network_uuid_sets:
            self.assertIn(actual_additional_network_uuids, expected_additional_network_uuid_sets)

        for expected_additional_network_uuids in expected_additional_network_uuid_sets:
            self.assertIn(expected_additional_network_uuids, actual_additional_network_uuid_sets)


    def test_computeAvailabilityZone_names(self):
        """Assert that all machines have one valid compute availability zone."""
        valid_compute_zones = {failureDomain["computeAvailabilityZone"] for failureDomain in self.expected_failure_domains}
        for machine in self.machines:
            zone = machine["spec"]["providerSpec"]["value"]["availabilityZone"]
            self.assertIn(zone, valid_compute_zones)


    def test_replica_distribution(self):
        """Assert that machines are evenly distributed across compute availability zones."""
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


    def test_storage_zone_names(self):
        """Assert that all machines have one valid storage availability zone that matches the given compute az."""
        valid_storage_zones = {failureDomain["storageAvailabilityZone"] for failureDomain in self.expected_failure_domains}
        for machine in self.machines:
            compute_zone = machine["spec"]["providerSpec"]["value"]["availabilityZone"]
            storage_zone = machine["spec"]["providerSpec"]["value"]["rootVolume"]["availabilityZone"]
            self.assertIn(storage_zone, valid_storage_zones)
            # availability zones have matching names in the test-case install-config
            self.assertEqual(compute_zone[-1], storage_zone[-1])


    def test_no_primarySubnet_if_control_plane_portTarget(self):
        """Since install-config sets a control-plane portTarget, assert that primarySubnet is empty."""
        for machine in self.machines:
            self.assertIsNone(machine["spec"]["providerSpec"]["value"].get("primarySubnet"))


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    INSTALLCONFIG_PATH = os.path.join(os.path.dirname(os.path.realpath(__file__)), 'install-config.yaml')
    unittest.main(verbosity=2)
