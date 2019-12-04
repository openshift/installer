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
module: os_delete_network_resources
short_description: Delete network resources
description:
    - Delete network resources filtered by tag
options:
   resource:
     description:
        - The name of a network resource
     required: true
   tags:
     description:
        - The tags to filter the network resources
     required: true
'''

from ansible.module_utils.basic import AnsibleModule
from ansible.module_utils.openstack import openstack_full_argument_spec, openstack_module_kwargs, openstack_cloud_from_module


def main():
    ''' Main module function '''
    argument_spec = openstack_full_argument_spec(
        resource=dict(type=str, required=True),
        tags=dict(type=list, required=True))

    module_kwargs = openstack_module_kwargs()
    module = AnsibleModule(argument_spec, **module_kwargs)

    resource = module.params['resource']
    tags = module.params['tags']

    sdk, cloud = openstack_cloud_from_module(module)
    try:
        get_resources = getattr(cloud.network, '{0}s'.format(resource), None)
        if get_resources:
            resources = get_resources(tags=tags)
            network_attr = 'delete_{0}'.format(resource)
            delete_resource = getattr(cloud.network, network_attr, None)
            if delete_resource:
                for resource in resources:
                    delete_resource(resource)
    except sdk.exceptions.OpenStackCloudException as e:
        module.fail_json(msg=str(e))

    module.exit_json(
        changed=True)

if __name__ == '__main__':
    main()
