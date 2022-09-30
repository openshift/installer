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
        """Assert that the Cluster infrastructure object contains the BGP configuration."""
        self.assertIn("apiLoadBalancer", cluster_infra["spec"]["platformSpec"]["openstack"])

        apiLoadBalancer = cluster_infra["spec"]["platformSpec"]["openstack"]["apiLoadBalancer"]

        self.assertIn("type", apiLoadBalancer)
        self.assertEqual("BGP", apiLoadBalancer["type"])
        self.assertIn("bgp", apiLoadBalancer)
        self.assertIn("speakers", apiLoadBalancer["bgp"])
        self.assertEqual(1, len(apiLoadBalancer["bgp"]["speakers"]))
        self.assertIn("failureDomain", apiLoadBalancer["bgp"]["speakers"][0])
        self.assertEqual("default", apiLoadBalancer["bgp"]["speakers"][0]["failureDomain"])
        self.assertIn("peers", apiLoadBalancer["bgp"]["speakers"][0])
        self.assertEqual(1, len(apiLoadBalancer["bgp"]["speakers"][0]["peers"]))
        self.assertIn("ip", apiLoadBalancer["bgp"]["speakers"][0]["peers"][0])
        self.assertEqual("192.168.0.1", apiLoadBalancer["bgp"]["speakers"][0]["peers"][0]["ip"])
        self.assertIn("password", apiLoadBalancer["bgp"]["speakers"][0]["peers"][0])
        self.assertEqual("changeme", apiLoadBalancer["bgp"]["speakers"][0]["peers"][0]["password"])


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    unittest.main(verbosity=2)
