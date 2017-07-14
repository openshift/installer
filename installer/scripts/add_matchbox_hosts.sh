#!/usr/bin/env bash
set -ex

{
  echo "172.18.0.10 node1 node1.example.com"
  echo "172.18.0.11 node2 node2.example.com"
  echo "172.18.0.12 node3 node3.example.com"
  echo "127.0.0.1 matchbox.example.com"
} >> /etc/hosts
