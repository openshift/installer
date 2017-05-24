# Uninstall Tectonic

## Deleting a cluster created by Tectonic Installer

Clusters created with the graphical Tectonic Installer can be deleted with a few commands using cluster configuration data found either in the location where Tectonic Installer was extracted and executed, or in the cluster's `assets.zip` file downloaded at the end of installation.

### Destroying a cluster with Installer state

Recent users of Tectonic Installer will have downloaded a release `.tar.gz` file, extracted a `tectonic` directory from it, and executed the installation from there. Locate the directory where Tectonic Installer was extracted and executed.

#### Export AWS credentials

First, open a terminal and set AWS credentials for the destroy operation. The AWS credentials used to perform the installation are recommended. Replace `<ACCESSKEYID>` and `<SECRETACCESSKEY>` in the commands below with appropriate values for the AWS account.

```bash
$ export AWS_ACCESS_KEY_ID=<ACCESSKEYID>
$ export AWS_SECRET_ACCESS_KEY=<SECRETACCESSKEY>
```

#### Destroy the cluster

Next, navigate to the cluster state directory written to the extracted `tectonic` directory during the install process. Cluster state directories are stored beneath the same parent directory as the `installer` and `terraform` binaries, in a child directory called `clusters`. Each state directory beneath `clusters` is named by a cluster's name suffixed with the date and time of the install.

```bash
# Replace <os> with darwin or linux
# Replace <CLUSTERNAME> with a string like mytectonic_2017-05-03_11-41-02
$ cd tectonic/tectonic-installer/<os>/clusters/<CLUSTERNAME>
$ export PATH=$(pwd)/../..:$PATH	# Add Installer's terraform binary to PATH
$ TERRAFORM_CONFIG=$(pwd)/.terraformrc terraform destroy --force
```

`terraform destroy` will itemize the destruction of the cluster's resources, producing many lines of output like the following:

```bash
tls_private_key.ingress: Refreshing state... (ID: 38aaae42623d255797e70602cf81b27574496fdf)
[...]
module.vpc.aws_security_group.master: Destroying... (ID: sg-33693754)
module.vpc.aws_security_group.worker: Destroying... (ID: sg-0b66386c)
module.vpc.aws_security_group.worker: Destruction complete
module.vpc.aws_security_group.master: Destruction complete

Destroy complete! Resources: 48 destroyed.
```

### Destroying a cluster with state from assets.zip

If using a state directory extracted from `assets.zip`, add a recent Tectonic Installer's `terraform` binary to the shell's $PATH before changing to the directory extracted from `assets.zip` and invoking `terraform destroy` as above.

## Deleting other cluster resources

This process will not delete resources not provisioned by Tectonic Installer or subsequent configuration with Terraform.

### Delete Kubernetes LoadBalancer services

Kubernetes [`LoadBalancer` services][k8s-lb] expose cluster facilities by manipulating the AWS Elastic Load Balancer (ELB) API. These ELBs are created dynamically by Kubernetes, after Tectonic Installer creates the cluster nodes and network environment. Such ELBs that are not part of the Tectonic Installer and Terraform state will not be deleted by the process described in this document.

Before destroying the cluster, delete all services of `type=LoadBalancer` with `kubectl delete`, or in Tectonic Console, to avoid this issue.

For example, to delete *every* LoadBalancer service in a cluster with `kubectl`:

```sh
$ kubectl delete services -l type=LoadBalancer
```

### Delete other AWS resources

Kubernetes LoadBalancer service ELBs are the most common resource added after Tectonic installation, but other resources not created by Tectonic Installer will not be deleted by the process described in this document.

These resources may include ELBs manually created with endpoints in the cluster subnet, Elastic Block Store (EBS) volumes that were attached to nodes or pods after installing Tectonic, and any other resources associated with Tectonic's Virtual Private Cloud (VPC) but not created by the Tectonic installer. Delete these resources by using their AWS service's UI or API.

## Reinstall

Use the latest [Tectonic AWS][install-aws] documentation to create a new cluster and deploy your applications and services.


[assets]: ../../admin/assets-zip.md
[install-aws]: index.md
[k8s-lb]: https://kubernetes.io/docs/user-guide/load-balancer/
