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
        self.assertIn("controlPlaneLoadBalancer", cluster_infra["spec"]["platformSpec"]["openstack"])

        controlPlaneLoadBalancer = cluster_infra["spec"]["platformSpec"]["openstack"]["controlPlaneLoadBalancer"]

        self.assertIn("type", controlPlaneLoadBalancer)
        self.assertEqual("BGP", controlPlaneLoadBalancer["type"])
        self.assertIn("bgp", controlPlaneLoadBalancer)
        self.assertIn("speakers", controlPlaneLoadBalancer["bgp"])
        self.assertEqual(1, len(controlPlaneLoadBalancer["bgp"]["speakers"]))
        self.assertIn("subnetCIDR", controlPlaneLoadBalancer["bgp"]["speakers"][0])
        self.assertEqual("192.168.0.0/10", controlPlaneLoadBalancer["bgp"]["speakers"][0]["subnetCIDR"])
        self.assertIn("peers", controlPlaneLoadBalancer["bgp"]["speakers"][0])
        self.assertEqual(1, len(controlPlaneLoadBalancer["bgp"]["speakers"][0]["peers"]))
        self.assertIn("ip", controlPlaneLoadBalancer["bgp"]["speakers"][0]["peers"][0])
        self.assertEqual("192.168.0.1", controlPlaneLoadBalancer["bgp"]["speakers"][0]["peers"][0]["ip"])
        self.assertIn("password", controlPlaneLoadBalancer["bgp"]["speakers"][0]["peers"][0])
        self.assertEqual("changeme", controlPlaneLoadBalancer["bgp"]["speakers"][0]["peers"][0]["password"])


if __name__ == '__main__':
    ASSETS_DIR = sys.argv.pop()
    unittest.main(verbosity=2)
