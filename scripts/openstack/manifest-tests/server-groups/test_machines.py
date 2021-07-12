#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest

import sys
import glob
import yaml

ASSETS_DIR = ""


class TestMachinesServerGroup(unittest.TestCase):
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


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    unittest.main(verbosity=2)
