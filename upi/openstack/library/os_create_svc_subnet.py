#!/usr/bin/python
# -*- coding: utf-8 -*-

# Copyright 2019 Red Hat, Inc. and/or its affiliates
# and other contributors as indicated by the @author tags.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


DOCUMENTATION = '''
---
module: os_create_svc_subnet
short_description: Create a subnet
description:
    - Create a subnet for the services
options:
   subnet_name:
     description:
       - The name of the subnet to be created
     required: true
   network_name:
     description:
        - Network this subnet belongs to
     required: true
   cidr:
     description:
       - Subnet range in CIDR notation
     required: true
   openstack_svc_network:
     description:
        -  CIDR of network to allocate IPs for OpenStack
           Octavia's Amphora VMs. Note that the value of this
           option must be at least two times the size of the
           cidr option, as octavia uses two IPs from this network
           for each loadbalancer, one given by OpenShift and other
           for VRRP connections.
     required: true
'''

from ansible.module_utils.basic import AnsibleModule
from ansible.module_utils.openstack import openstack_full_argument_spec, openstack_module_kwargs, openstack_cloud_from_module

import ipaddress

_IP_VERSION = 4
_ENABLE_DHCP = False


def main():
    ''' Main module function '''
    argument_spec = openstack_full_argument_spec(
        subnet_name=dict(type='str', required=True),
        network_name=dict(type='str', required=True),
        cidr=dict(type='str', required=True),
        openstack_svc_network=dict(type='str', required=True))

    module_kwargs = openstack_module_kwargs()
    module = AnsibleModule(argument_spec, **module_kwargs)

    subnet_name = module.params['subnet_name']
    network_name = module.params['network_name']
    cidr = module.params['cidr']
    openstack_cidr = module.params['openstack_svc_network']

    svc_network = ipaddress.ip_network(cidr, strict=False)
    openstack_svc_cidr = ipaddress.ip_network(openstack_cidr, strict=False)

    # Parts of openstack_svc_cidr not overlapping with svc_network
    allocation_ranges = []
    # If the first IPs are different the last part
    # of the openstack network is already allocated
    if svc_network[0] != openstack_svc_cidr[0]:
        allocation_ranges.extend([{
            "start": str(openstack_svc_cidr[0]+1),
            "end": str(svc_network[0]-1)}])

    # If the last IPs are different the first
    # part of the openstack network is already allocated
    if svc_network[-1] != openstack_svc_cidr[-1]:
        allocation_ranges.extend([{
            "start": str(svc_network[-1]+1),
            "end": str(openstack_svc_cidr[-2])}])

    if allocation_ranges:
        gateway_ip = str(allocation_ranges[-1]["end"])
        allocation_ranges[-1]["end"] = str(ipaddress.IPv4Network(gateway_ip)[0]-1)

    sdk, cloud = openstack_cloud_from_module(module)
    try:
        subnet = cloud.network.find_subnet(name_or_id=subnet_name)
        network = cloud.network.find_network(name_or_id=network_name)
        if not network:
            module.fail_json(msg='Network not Found')
        if not allocation_ranges:
            module.fail_json(msg='Unable to retrieve not overlapping ranges')
        if not subnet:
            cloud.network.create_subnet(
                network_id=network.id, cidr=openstack_cidr,
                name=subnet_name, gateway_ip=gateway_ip,
                ip_version=_IP_VERSION, enable_dhcp=_ENABLE_DHCP,
                allocation_pools=allocation_ranges)
            module.exit_json(changed=True)
    except sdk.exceptions.OpenStackCloudException as e:
        module.fail_json(msg=str(e))

    module.exit_json(
        changed=False)

if __name__ == '__main__':
    main()
