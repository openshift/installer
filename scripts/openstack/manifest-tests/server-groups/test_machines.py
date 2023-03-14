#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest
import xmlrunner

import os
import sys
import glob
import yaml

ASSETS_DIR = ""


class ServerGroupMachines(unittest.TestCase):
    def setUp(self):
        """Parse the Machines into a Python data structure."""
        self.machines = []
        for machine_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_master-machines-*.yaml'
        ):
            with open(machine_path) as f:
                self.machines.append(yaml.load(f, Loader=yaml.FullLoader))

    def test_consistent_group_name(self):
        """Assert that all machines bear the same server group name."""
        group_name = None
        for machine in self.machines:
            name = machine["spec"]["providerSpec"]["value"]["serverGroupName"]
            if group_name is None:
                group_name = name

            self.assertEqual(name, group_name)


class ServerGroupMachinesets(unittest.TestCase):
    def setUp(self):
        """Parse the MachineSets into a Python data structure."""
        self.machinesets = []
        for machineset_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-cluster-api_worker-machineset-*.yaml'
        ):
            with open(machineset_path) as f:
                self.machinesets.append(yaml.load(f, Loader=yaml.FullLoader))

    def test_consistent_group_names(self):
        """Assert that server group names are unique across machinesets."""
        found = []
        for machineset in self.machinesets:
            name = machineset["spec"]["template"]["spec"]["providerSpec"][
                "value"]["serverGroupName"]
            self.assertNotIn(name, found)
            found.append(name)


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    with open(os.environ.get('JUNIT_FILE', '/dev/null'), 'wb') as output:
        unittest.main(testRunner=xmlrunner.XMLTestRunner(output=output), failfast=False, buffer=False, catchbreak=False, verbosity=2)
