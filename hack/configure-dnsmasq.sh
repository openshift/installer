#!/bin/sh

# Configure the wildcard dns record in NetworkManager's dnsmasq configuration file.

domain=$1

echo -e "[main]\ndns=dnsmasq" | sudo tee /etc/NetworkManager/conf.d/openshift.conf
echo -e "server=/$domain/192.168.126.1\naddress=/.apps.$domain/192.168.126.51" | sudo tee /etc/NetworkManager/dnsmasq.d/openshift.conf

sudo systemctl restart NetworkManager
