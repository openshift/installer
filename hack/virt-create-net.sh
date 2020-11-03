#!/bin/bash

sudo virsh net-create ./hack/net.xml

echo server=/api.test-cluster.redhat.com/192.168.126.1 | sudo tee /etc/NetworkManager/dnsmasq.d/aio.conf
echo -e "[main]\ndns=dnsmasq" | sudo tee /etc/NetworkManager/conf.d/aio.conf
systemctl reload NetworkManager.service
