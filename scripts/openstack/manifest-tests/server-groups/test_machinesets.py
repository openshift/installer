#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest

import sys
import glob
import yaml

ASSETS_DIR = ""


class TestMachinesetsServerGroup(unittest.TestCase):
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
    unittest.main(verbosity=2)
