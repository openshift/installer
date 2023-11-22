#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest
import xmlrunner

import os
import sys
import yaml

ASSETS_DIR = ""

class ControlPlaneMachineSet(unittest.TestCase):
    def setUp(self):
        """Parse the CPMS into a Python data structure."""
        with open(f'{ASSETS_DIR}/openshift/99_openshift-machine-api_master-control-plane-machine-set.yaml') as f:
            self.cpms = yaml.load(f, Loader=yaml.FullLoader)

    def test_compute_zones(self):
        """Assert that the OpenStack CPMS failureDomains value is empty."""
        self.assertIsNone(self.cpms["spec"]["template"]["machines_v1beta1_machine_openshift_io"].get("failureDomains"))


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    with open(os.environ.get('JUNIT_FILE', '/dev/null'), 'wb') as output:
        unittest.main(testRunner=xmlrunner.XMLTestRunner(output=output), failfast=False, buffer=False, catchbreak=False, verbosity=2)
