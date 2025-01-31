#!/usr/bin/env python
# -*- coding: utf-8 -*-

import unittest
import xmlrunner

import os
import sys
import glob
import yaml

ASSETS_DIR = ""

class GenerateMachineConfig(unittest.TestCase):
    def setUp(self):
        self.machine_configs = []
        for machine_config_path in glob.glob(
                f'{ASSETS_DIR}/openshift/99_openshift-machineconfig_99-dual-stack-*.yaml'
        ):
            with open(machine_config_path) as f:
                self.machine_configs.append(yaml.load(f, Loader=yaml.FullLoader))

    def test_kernel_args(self):
        """Assert there are machine configs configuring the kernel args for masters and workers"""
        for machine_config in self.machine_configs:
            kernel_args = machine_config["spec"]["kernelArguments"]
            self.assertIn("ip=enp3s0:dhcp,dhcp6", kernel_args)


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    with open(os.environ.get('JUNIT_FILE', '/dev/null'), 'wb') as output:
        unittest.main(testRunner=xmlrunner.XMLTestRunner(output=output), failfast=False, buffer=False, catchbreak=False, verbosity=2)
