# Openstack User Provided Infrastructure installation

This directory contains the Ansible scripts that automate part of the work in the [UPI installation](../../docs/user/openstack/install_upi.md).

## Rationale

The tool for automated installation (IPI - Installer-Provided Infrastructure) cover the general case where all the OpenStack resources can be created ad-hoc.

The UPI case (User-Provided Infrastructure) instead lets complete freedom to the end-user of customizing the cluster and its underpinning OpenStack resources.

The installation process is detailed step by step in the relative [documentation](../../docs/user/openstack/install_upi.md).

These Ansible playbooks in this directory automate some of those steps. They are provided as a template: edit them to match your needs.

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

The playbooks are designed to reproduce an installation equivalent to IPI. Customize them as needed. It is advised to customize the teardown playbooks symmetrically.

Every step can be run like this:

```shell
(venv)$ ansible-playbook -i inventory.yaml 01_network.yaml
```

For every script, a symmetrical teardown playbook is provided:

```shell
(venv)$ ansible-playbook -i inventory.yaml down-01_network.yaml
```

A full teardown can be achieved by running all the `down` scripts in reverse order.
