#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest
import xmlrunner

import os
import sys
import yaml

ASSETS_DIR = ""

EXPECTED_MASTER_ZONE_NAMES = ["masterzone", "masterztwo", "masterzthree"]

class ControlPlaneMachineSet(unittest.TestCase):
    def setUp(self):
        """Parse the CPMS into a Python data structure."""
        with open(f'{ASSETS_DIR}/openshift/99_openshift-machine-api_master-control-plane-machine-set.yaml') as f:
            self.cpms = yaml.load(f, Loader=yaml.FullLoader)

    def test_providerspec_failuredomain_fields(self):
        """Assert that the failure-domain-managed fields in the CPMS providerSpec are omitted."""
        provider_spec = self.cpms["spec"]["template"]["machines_v1beta1_machine_openshift_io"]["spec"]["providerSpec"]["value"]
        self.assertNotIn("availabilityZone", provider_spec)

    def test_compute_zones(self):
        """Assert that the CPMS failure domain zones match the expected machine-pool zones."""
        self.assertIn("failureDomains", self.cpms["spec"]["template"]["machines_v1beta1_machine_openshift_io"])
        failure_domains = self.cpms["spec"]["template"]["machines_v1beta1_machine_openshift_io"]["failureDomains"]["openstack"]

        compute_zones = []
        for failure_domain in failure_domains:
            zone = failure_domain["availabilityZone"]
            compute_zones.append(zone)
            self.assertIn(zone, EXPECTED_MASTER_ZONE_NAMES)

        for expected_zone in EXPECTED_MASTER_ZONE_NAMES:
            self.assertIn(expected_zone, compute_zones)


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    with open(os.environ.get('JUNIT_FILE', '/dev/null'), 'wb') as output:
        unittest.main(testRunner=xmlrunner.XMLTestRunner(output=output), failfast=False, buffer=False, catchbreak=False, verbosity=2)
