# oVirt/RHV User Provided Infrastructure

This folder contains the Ansible scripts to help automate as possible the oVirt/RHV UPI
in a step by step installation process documented [here](../../docs/users/ovirt/install_upi.md).

## Getting started

Inspect and customize the [inventory.yml](./inventory.yml) variables.

Execute every step of the installation like below:

```sh
$ ansible-playbook -i inventory.yml bootstrap.yml
```

Please refer to the [oVirt/RHV documentation](../../docs/users/ovirt/install_upi.md)
for the full step-by-step process.
