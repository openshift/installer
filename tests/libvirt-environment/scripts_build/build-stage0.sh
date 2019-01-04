#!/usr/bin/env bash
set -xe

dnf update -y
dnf install -y \
    w3m \
    unzip \
    deltarpm pki-ca \
    kubernetes-client \
    git gcc gcc-c++ libtool golang \
    qemu-img libvirt libvirt-python libvirt-client libvirt-devel @virtualization \
    dnsmasq kubernetes-client

dnf clean all
rm -rf /var/cache/dnf/*

curl -OL https://github.com/openshift/origin/releases/download/v3.10.0/openshift-origin-client-tools-v3.10.0-dd10d17-linux-64bit.tar.gz
tar -zxf openshift-origin-client-tools-v3.10.0-dd10d17-linux-64bit.tar.gz
mv -f ./openshift-origin-client-tools-v3.10.0-dd10d17-linux-64bit/oc /usr/local/bin

curl -OL https://releases.hashicorp.com/terraform/0.11.8/terraform_0.11.8_linux_amd64.zip
unzip terraform_0.11.8_linux_amd64.zip
mv -f ./terraform /usr/local/bin

rm -rf ./openshift-origin-client-tools-v3.10.0-dd10d17-linux-64bit.tar.gz \
  ./openshift-origin-client-tools-v3.10.0-dd10d17-linux-64bit \
  ./terraform_0.11.8_linux_amd64.zip
