# VMware UPI install with Terraform EXAMPLE

## Pre-Requisites

* terraform
* jq

## Build a Cluster

1. Create an install-config.yaml.
The machine CIDR for the dev cluster is 139.178.89.192/26.

```yaml
apiVersion: v1
baseDomain: devcluster.openshift.com
metadata:
  name: mstaeble
networking:
  machineNetwork:
  - cidr: "139.178.89.192/26"
platform:
  vsphere:
    vCenter: vcsa.vmware.devcluster.openshift.com
    username: YOUR_VSPHERE_USER
    password: YOUR_VSPHERE_PASSWORD
    datacenter: dc1
    defaultDatastore: nvme-ds1
pullSecret: YOUR_PULL_SECRET
sshKey: YOUR_SSH_KEY
```

1. Run `openshift-install create ignition-configs`.

1. Fill out a terraform.tfvars file with the ignition configs generated.

There is an example terraform.tfvars file in this directory named terraform.tfvars.example. The example file is set up for use with the dev cluster running at vcsa.vmware.devcluster.openshift.com. At a minimum, you need to set values for the following variables.

* cluster_id
* cluster_domain
* vsphere_user
* vsphere_password
* ipam_token OR bootstrap_ip_address, lb_ip_address, control_plane_ip_addresses, and compute_ip_addresses

The bootstrap ignition config must be placed in a location that will be accessible by the bootstrap machine. For example, you could store the bootstrap ignition config in a gist.

Even if declaring static IPs a DHCP server is still required early in the boot process to download the ignition files.

1. Fetch the latest terraform ignition plugin from [here](https://github.com/community-terraform-providers/terraform-provider-ignition/releases), unpack it and place it to `/bin/terraform-provider-ignition`s

1. Run `terraform init`.

1. Ensure that you have you AWS profile set and a region specified. The installation will use create AWS route53 resources for routing to the OpenShift cluster.

1. Run `terraform apply -auto-approve`.
This will reserve IP addresses for the VMs.

1. Run `openshift-install wait-for bootstrap-complete`. Wait for the bootstrapping to complete.

1. Run `terraform apply -auto-approve -var 'bootstrap_complete=true'`.
This will destroy the bootstrap VM.

1. Run `openshift-install wait-for install-complete`. Wait for the cluster install to finish.

1. Enjoy your new OpenShift cluster.

1. Run `terraform destroy -auto-approve`.
