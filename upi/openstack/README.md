# Ansible Playbooks for Openstack UPI

This directory contains the Ansible scripts expected to automate most of the command-line work in the [User-Provided-Infrastructure installation](../../docs/user/openstack/install_upi.md).

## How to use

Customize the cluster properties in the [Inventory](./inventory.yaml) file.

**NOTE:** To deploy with Kuryr SDN, update the `os_networking_type` field to `Kuryr`.

The playbooks in this directory are designed to reproduce an IPI installation, but are highly customizable. Please be aware of changes made to the install playbooks that may require changes to the teardown playbooks.

Every step can be run like this:

```shell
(venv)$ ansible-playbook -i inventory.yaml network.yaml
```

For every script, a symmetrical teardown playbook is provided:

```shell
(venv)$ ansible-playbook -i inventory.yaml down-network.yaml
```

A full teardown can be achieved by running all the `down` scripts in reverse order.

Please refer to the [UPI documentation](../../docs/user/openstack/install_upi.md) for step-by-step instructions.
