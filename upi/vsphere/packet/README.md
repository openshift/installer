# Installing OpenShift on VMware in Packet

This Terraform example is used to deploy RHCOS virtual machines in vSphere.
Our specific use-case is a Packet bare metal vSphere environment.

## Requirements

- AWS (Route53)
- Terraform 0.11 (0.12 is currently not supported)
- jq

## Terraform
You can create a Terraform variables file by copying the example provided:

```sh
cd upi/vsphere/packet
cp terraform.tfvars.example terraform.tfvars
```
Then modify `terraform.tfvars` file based your environment.

### Configure IPAM or Static IP Addressing

If you are using phpIPAM you can configure the following values:
* ipam
* ipam_token

Then phpIPAM will provide addresses that your virtual machines will use.

Otherwise configure static addressing via:

* bootstrap_ip_address
* control_plane_ip_addresses
* compute_ip_addresses

### Configuring Ignition

The Ignition files need to be appended with hostname and IP configuration
then added to each virtual machines' extra config - this is done by Terraform.

Either copy the Ignition files to `upi/vsphere/packet` or configure:
* bootstrap_ignition_path
* control_plane_ignition_path
* compute_ignition_path

### Configuring vSphere

The vSphere vCenter variables must be configured to your environment please modify:

* vsphere_server
* vsphere_user
* vsphere_password
* vsphere_cluster
* vsphere_datacenter
* vsphere_datastore
* vm_template

### Configuring OpenShift

The OpenShift-specific variables also must be configured please modify:

* cluster_id
* cluster_domain
* base_domain
* machine_cidr

**NOTE**: The `base_domain`, `cluster_id`, `machine_cidr` **must** match the values within `install-config.yaml`.

### Running Terraform commands

First let's ensure that you have you AWS profile set and a region specified.
In this example my AWS profile is named `openshift-dev` and uses the default region of `us-east-2`.
Provide your specific profile and region:

```
export AWS_PROFILE="openshift-dev"
export AWS_DEFAULT_REGION=us-east-2
```

Next we need to initialize a working directory containing Terraform configuration files:
```sh
terraform init
```
Once initialized continue with creating the infrastructure and waiting for bootstrap to complete:

```sh
terraform apply -auto-approve
openshift-install wait-for bootstrap-complete
```

Once bootstrap has completed remove the bootstrap node and release the IP address.
```sh
terraform destroy -target=module.rhcos_virtual_machines.module.bootstrap
terraform destroy -target=module.ipam_bootstrap.null_resource.ip_address
```

Wait for the cluster install to finish:
```sh
openshift-install wait-for install-complete
```
When the command completes it will provide authentication details and the URL to the OpenShift console.

### Delete OpenShift Cluster

```sh
terraform destroy -auto-approve
```
