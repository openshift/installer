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
module: os_create_subnet_pool
short_description: Create a subnet pool
description:
    - Create a subnet pool
options:
   name:
     description:
       - name of the new subnet pool
     required: true
   prefix_len:
     description:
        - subnet pool default prefix length
   prefixes:
     description:
        - subnet pool prefixes
'''

from ansible.module_utils.basic import AnsibleModule
from ansible.module_utils.openstack import openstack_full_argument_spec, openstack_module_kwargs, openstack_cloud_from_module


def main():
    ''' Main module function '''
    argument_spec = openstack_full_argument_spec(
        name=dict(type='str', required=True),
        prefix_lenght=dict(type='str'),
        prefixes=dict(type=list))

    module_kwargs = openstack_module_kwargs()
    module = AnsibleModule(argument_spec, **module_kwargs)

    name = module.params['name']
    prefix_lenght = module.params['prefix_lenght']
    prefixes = module.params['prefixes']

    sdk, cloud = openstack_cloud_from_module(module)
    try:
        subnet_pool = cloud.network.find_subnet_pool(name_or_id=name)
        if not subnet_pool:
            cloud.network.create_subnet_pool(
                name=name, default_prefixlen=prefix_lenght, prefixes=prefixes)
            module.exit_json(changed=True)
    except sdk.exceptions.OpenStackCloudException as e:
        module.fail_json(msg=str(e))

    module.exit_json(
        changed=False)

if __name__ == '__main__':
    main()
