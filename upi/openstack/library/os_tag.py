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
module: os_tag
short_description: Set a tag to a resource
description:
    - Set a tag to a Network resource
options:
   resource:
     description:
       - The type of the Network resource.
     required: true
   name:
     description:
        - The name of the resource to be tagged
     required: true
   tags:
     description:
        - Tags to be set on the resources
     required: true
'''

from ansible.module_utils.basic import AnsibleModule
from ansible.module_utils.openstack import openstack_full_argument_spec, openstack_module_kwargs, openstack_cloud_from_module


def main():
    ''' Main module function '''
    argument_spec = openstack_full_argument_spec(
        resource=dict(type=str, required=True),
        name=dict(type=str, required=True),
        tags=dict(type='list', required=True))

    module_kwargs = openstack_module_kwargs()
    module = AnsibleModule(argument_spec, **module_kwargs)

    resource = module.params['resource']
    name = module.params['name']
    tags = module.params['tags']

    sdk, cloud = openstack_cloud_from_module(module)
    try:
        network_attr = 'find_{0}'.format(resource)
        get_resource = getattr(cloud.network, network_attr, None)
        if get_resource:
            resource = get_resource(name)
            cloud.network.set_tags(resource, tags)
    except sdk.exceptions.OpenStackCloudException as e:
        module.fail_json(msg=str(e))

    module.exit_json(
        changed=True)

if __name__ == '__main__':
    main()
