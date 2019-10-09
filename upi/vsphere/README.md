# Installing OpenShift on VMware

The Terraform examples provided require a AWS account.  AWS is used
for its API interface to network load balancing and DNS.
The `./modules/rhcos_virtual_machines` can be used independently if so desired.

We will begin with using `openshift-install` as each Terraform example
starts with creating the Ignition files.

## Create Ignition Configs

You can create a `install-config.yaml` file by copying the example provided:

```sh
cd upi/vsphere/
cp install-config.yaml.example install-config.yaml
```

All of the following variables must be configured:
* baseDomain
* metadata.name
* networking.machineCIDR
* platform.vsphere.vCenter
* platform.vsphere.username
* platform.vsphere.password
* platform.vsphere.datacenter
* platform.vsphere.defaultDatastore
* pullSecret
* sshKey

**NOTE**: Keep track of the `baseDomain`, `metadata.name`, `networking.machineCIDR` values.
They **must** match the values within `terraform.tfvars`.

Once configured please run:

```sh
openshift-install create ignition-configs
```

This will generate the Ignition files: `bootstrap.ign`, `master.ign` and
`compute.ign` that you will use in the Terraform examples below.

## VMware Cloud on AWS Example
Please review the [README](./vmc/README.md)

## Packet w/Route53 Example
Please review the [README](./packet/README.md)
