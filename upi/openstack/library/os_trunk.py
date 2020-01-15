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
module: os_trunk
short_description: Create or delete a Network trunk
description:
    - Create or delete a Network Trunk for a parent port
options:
   name:
     description:
       - The name of the trunk that should be created/deleted.
     required: true
   port_id:
     description:
        - ID of the parent port
        - Required when I(state) is 'present'
   state:
     description:
        - Indicate desired state of the resource
     choices: ['present', 'absent']
     default: present
'''

from ansible.module_utils.basic import AnsibleModule
from ansible.module_utils.openstack import openstack_full_argument_spec, openstack_module_kwargs, openstack_cloud_from_module


def main():
    ''' Main module function '''
    argument_spec = openstack_full_argument_spec(
        name=dict(type='str', required=True),
        port_id=dict(type='str'),
        state=dict(default='present', choices=['absent', 'present']))

    module_kwargs = openstack_module_kwargs()
    module = AnsibleModule(argument_spec, **module_kwargs)

    name = module.params['name']
    port_id = module.params['port_id']
    state = module.params['state']

    sdk, cloud = openstack_cloud_from_module(module)
    try:
        trunk = cloud.network.find_trunk(name_or_id=name)
        if state == 'present':
            if not port_id:
                module.fail_json(msg='port_id required with present state')
            if not trunk:
                cloud.network.create_trunk(name=name, port_id=port_id)
        elif state == 'absent':
            if trunk:
                cloud.network.delete_trunk(trunk)
    except sdk.exceptions.OpenStackCloudException as e:
        module.fail_json(msg=str(e))

    module.exit_json(
        changed=True)

if __name__ == '__main__':
    main()
