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
module: os_delete_lb_resources
short_description: Delete load balancer resources
description:
    - Delete load balancer resources filtered by tag
options:
   tags:
     description:
        - The tags to filter the load balancer
     required: true
'''

from ansible.module_utils.basic import AnsibleModule
from ansible.module_utils.openstack import openstack_full_argument_spec, openstack_module_kwargs, openstack_cloud_from_module

_OCTAVIA_TAG_VERSION = 2, 5


def main():
    ''' Main module function '''
    argument_spec = openstack_full_argument_spec(
        tags=dict(type=list, required=True),
        region_name=dict(type=str, required=True))

    module_kwargs = openstack_module_kwargs()
    module = AnsibleModule(argument_spec, **module_kwargs)

    tags = module.params['tags']
    region_name = module.params['region_name']

    sdk, cloud = openstack_cloud_from_module(module)
    try:
        regions = cloud.load_balancer.get_all_version_data()
        endpoints = regions.get(region_name, list(regions.values())[0])
        services = list(endpoints.values())[0]
        versions = services.get('load-balancer', list(services.values())[0])

        tag_supported = any(
            tuple(map(int, v['version'].split('.'))) >= _OCTAVIA_TAG_VERSION
            for v in versions)

        resources = list(cloud.load_balancer.load_balancers(
            description=tags))
        if tag_supported:
            resources.extend(list(cloud.load_balancer.load_balancers(
                tags=tags)))
        for resource in resources:
            cloud.load_balancer.delete_load_balancer(resource, cascade=True)
    except sdk.exceptions.OpenStackCloudException as e:
        module.fail_json(msg=str(e))

    module.exit_json(
        changed=True)

if __name__ == '__main__':
    main()
