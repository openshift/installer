# Install Tectonic on Google Cloud Platform with Terraform

Use this guide to manually install a Tectonic cluster on a GCP account.

## Prerequsities

- **Terraform:** >= v0.10.7
- **DNS**: Ensure that the DNS zone is already created and available in [Cloud DNS](https://console.cloud.google.com/net-services/dns) for the account. For example if the `tectonic_base_domain` is set to `kube.example.com` a Cloud DNS zone must exist with the nameservers configured for that domain.
Tectonic Account: Register for a Tectonic Account, which is free for up to 10 nodes. You must provide the cluster license and pull secret during installation.
- **Tectonic Account** - Register for a [Tectonic Account](https://coreos.com/tectonic), which is free for up to 10 nodes. You must provide the cluster license and pull secret during installation.

## Getting Started

First, clone the Tectonic Installer repository:

```
$ git clone https://github.com/coreos/tectonic-installer.git
$ cd tectonic-installer
```
 
Initialise Terraform:

```
$ terraform init platforms/gcp
``` 

Configure your GCP credentials. See the [Terraform Google provider docs](https://www.terraform.io/docs/providers/google/index.html) for details.

```
$ export GOOGLE_APPLICATION_CREDENTIALS=/my-credentials.json
```

## Customize the deployment

Customizations to the base installation live in examples/terraform.tfvars.gcp. Export a variable that will be your cluster identifier:

```
$ export CLUSTER=my-cluster
$ export GOOGLE_PROJECT=my-project-id
```

Create a build directory to hold your customizations and copy the example file into it:

```
$ mkdir -p build/${CLUSTER}
$ cp examples/terraform.tfvars.gcp build/${CLUSTER}/terraform.tfvars
```

Edit the parameters with your GCP details: project id, domain name, license, etc. [View all of the GCP specific options](https://github.com/coreos/tectonic-installer/tree/master/Documentation/variables/gcp.md) and [the common Tectonic variables](https://github.com/coreos/tectonic-installer/tree/master/Documentation/variables/config.md)).

## Deploy the cluster

Test out the plan before deploying everything:

```
$ terraform plan -var-file=build/${CLUSTER}/terraform.tfvars platforms/gcp
```

Next, deploy the cluster:

```
$ terraform apply -var-file=build/${CLUSTER}/terraform.tfvars platforms/gcp
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
$ terraform destroy -var-file=build/${CLUSTER}/terraform.tfvars platforms/gcp
```
### Known issues and workarounds

This is a non stable version currently under heavy development. You can expect to find issues.
