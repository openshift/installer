# Openstack User Provided Infrastructure installation

This directory contains Ansible scripts that automate the command-line work in the [User-Provided-Infrastructure installation](../../docs/user/openstack/install_upi.md).

## Requirements

* Python
* Ansible
* Python modules required in the playbooks. Namely:
  * openstacksdk
  * netaddr
  * openstackclient


The included `requirements.txt` helps using `pip` for gathering the required dependencies in a Python virtual environment:

```shell
python3 -m venv venv
source venv/bin/activate
pip install -U pip
pip install -r requirements.txt
```

## How to use

Customize the cluster properties in the [Inventory](./inventory.yaml) file.

**NOTE:** To deploy with Kuryr SDN, update the `os_networking_type` field to `Kuryr`.

The playbooks in this directory are designed to reproduce an IPI installation, but are highly customizable. Please be aware of changes made to the install playbooks that may require changes to the teardown playbooks.

Every step can be run like this:

```shell
(venv)$ ansible-playbook -i inventory.yaml 01_network.yaml
```

For every script, a symmetrical teardown playbook is provided:

```shell
(venv)$ ansible-playbook -i inventory.yaml down-01_network.yaml
```

A full teardown can be achieved by running all the `down` scripts in reverse order.


Please refer to the [UPI documentation](../../docs/user/openstack/install_upi.md) for step-by-step instructions.
