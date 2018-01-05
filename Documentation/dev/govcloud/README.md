# Install Tectonic on AWS GovCloud Platform with Terraform

Use this guide to manually install a Tectonic cluster on a AWS GovCloud account.

## Prerequsities

- **Terraform:** >= v0.10.7
- **Tectonic Account:** Register for a [Tectonic Account](https://coreos.com/tectonic), which is free for up to 10 nodes. You must provide the cluster license and pull secret during installation.
- **AWS GovCloud:** Obtain credentials for [GovCloud](http://docs.aws.amazon.com/govcloud-us/latest/UserGuide/govcloud-differences.html)    
- **DNS:** The Tectonic Installer assumes that a PowerDNS server instance is running and reachable from the VPC where the cluster is running.

See [contrib/govcloud](../../../contrib/govcloud) for an example of a prebuilt VPC with restricted VPN access and a PowerDNS server.

## Getting Started

First, clone the Tectonic Installer repository:

```
$ git clone https://github.com/coreos/tectonic-installer.git
$ cd tectonic-installer
```
 
Initialise Terraform:

```
$ terraform init platforms/govcloud
``` 

Configure your AWS GovCloud credentials.

```
$ export AWS_ACCESS_KEY_ID=my-id
$ export AWS_SECRET_ACCESS_KEY=secret-key
```

## Customize the deployment

Customizations to the base installation live in examples/terraform.tfvars.govcloud. Export a variable that will be your cluster identifier:

```
$ export CLUSTER=my-cluster
```

Create a build directory to hold your customizations and copy the example file into it:

```
$ mkdir -p build/${CLUSTER}
$ cp examples/terraform.tfvars.govcloud build/${CLUSTER}/terraform.tfvars
```

Edit the parameters with your VPC details:
```
tectonic_govcloud_external_vpc_id
tectonic_govcloud_external_master_subnet_ids
tectonic_govcloud_external_worker_subnet_ids
tectonic_govcloud_dns_server_ip

```

## Deploy the cluster

If you are following the [contrib/govcloud](../../../contrib/govcloud) example and deploying from an external machine, connect to the VPN now. 
Add the `tectonic_govcloud_dns_server_ip` to your local DNS resolver.
 
Test out the plan before deploying everything:

```
$ terraform plan -var-file=build/${CLUSTER}/terraform.tfvars platforms/govcloud
```

Next, deploy the cluster:

```
$ terraform apply -var-file=build/${CLUSTER}/terraform.tfvars platforms/govcloud
```

This should run for a little bit, and when complete, your Tectonic cluster should be ready.

### Access the cluster

The Tectonic Console should be up and running after the containers have downloaded. You can access it at the DNS name configured in your variables file prefixed by the cluster name, i.e ```https://cluster_name.tectonic_base_domain```.

Inside of the /generated folder you should find any credentials, including the CA if generated, and a kubeconfig. You can use this to control the cluster with kubectl:

```
$ export KUBECONFIG=generated/auth/kubeconfig
$ kubectl cluster-info
```
### Delete the cluster

```
$ terraform destroy -var-file=build/${CLUSTER}/terraform.tfvars platforms/govcloud
```

## Known issues and workarounds

At the moment because of the [AWS user data limit](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html#instancedata-add-user-data) and [Ignition not supporting the S3 protocol for replacing the content](https://github.com/coreos/bugs/issues/2216), the Ignition config for the nodes is stored in a public bucket and it has to be removed manually.
