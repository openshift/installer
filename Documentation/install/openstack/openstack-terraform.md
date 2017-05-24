# Install Tectonic on OpenStack with Terraform

Following this guide will deploy a Tectonic cluster within your OpenStack account.

Generally, the OpenStack platform templates adhere to the standards defined by the project [conventions][conventions] and [generic platform requirements][generic]. This document aims to document the implementation details specific to the OpenStack platform.

<p style="background:#d9edf7; padding: 10px;" class="text-info"><strong>Alpha:</strong> These modules and instructions are currently considered alpha. See the <a href="../../platform-lifecycle.md">platform life cycle</a> for more details.</p>

## Prerequsities

 - **CoreOS Container Linux** - The latest Container Linux Beta (1353.2.0 or later) [uploaded into Glance](https://coreos.com/os/docs/latest/booting-on-openstack.html) and its OpenStack image ID.
 - **Tectonic Account** - Register for a [Tectonic Account][register], which is free for up to 10 nodes. You will need to provide the cluster license and pull secret below.

## Getting Started
OpenStack is a highly customizable environment where different components can be enabled/disabled. This installation includes the following two flavors:

- `nova`: Only Nova computing nodes are being created for etcd, master and worker nodes, assuming the nodes get public IPs assigned.
- `neutron`: A private Neutron network is being created with master/worker nodes exposed via floating IPs connected to an etcd instance via an internal network.

Replace `<flavor>` with either option in the following commands. Now we're ready to specify our cluster configuration.

### Download and extract Tectonic Installer

Open a new terminal, and run the following commands to download and extract Tectonic Installer.

```bash
$ curl -O https://releases.tectonic.com/tectonic-1.6.2-tectonic.1.tar.gz # download
$ tar xzvf tectonic-1.6.2-tectonic.1.tar.gz # extract the tarball
$ cd tectonic
```

### Initialize and configure Terraform

Start by setting the `INSTALLER_PATH` to the location of your platform's Tectonic installer. The platform should be `darwin` or `linux`.

```bash
$ export INSTALLER_PATH=$(pwd)/tectonic-installer/linux/installer # Edit the platform name.
$ export PATH=$PATH:$(pwd)/tectonic-installer/linux # Put the `terraform` binary in our PATH
```

Make a copy of the Terraform configuration file for our system. Do not share this configuration file as it is specific to your machine.

```bash
$ sed "s|<PATH_TO_INSTALLER>|$INSTALLER_PATH|g" terraformrc.example > .terraformrc
$ export TERRAFORM_CONFIG=$(pwd)/.terraformrc
```

Next, get the modules that Terraform will use to create the cluster resources:

```
$ terraform get platforms/openstack/<flavor>
```

Configure your AWS credentials for setting up Route 53 DNS record entries. See the [AWS docs][env] for details.

```
$ export AWS_ACCESS_KEY_ID=
$ export AWS_SECRET_ACCESS_KEY=
```

Set your desired region:

```
$ export AWS_REGION=
```

Configure your OpenStack credentials.

```
$ export OS_TENANT_NAME=
$ export OS_USERNAME=
$ export OS_PASSWORD=
$ export OS_AUTH_URL=
$ export OS_REGION_NAME=
```

## Customize the deployment

Customizations to the base installation live in `examples/terraform.tfvars.<flavor>`. Export a variable that will be your cluster identifier:

```
$ export CLUSTER=my-cluster
```

Create a build directory to hold your customizations and copy the example file into it:

```
$ mkdir -p build/${CLUSTER}
# for Neutron:
$ cp examples/terraform.tfvars.openstack-neutron build/${CLUSTER}/terraform.tfvars
# for Nova:
$ cp examples/terraform.tfvars.openstack-nova build/${CLUSTER}/terraform.tfvars
```

Edit the parameters with your OpenStack details. View all of the [OpenStack Nova][openstack-nova-vars] and [OpenStack Neutron][openstack-neutron-vars] specific options and [the common Tectonic variables][vars].

## Deploy the cluster

Test out the plan before deploying everything:

```
$ terraform plan -var-file=build/${CLUSTER}/terraform.tfvars platforms/openstack/<flavor>
```

Next, deploy the cluster:

```
$ terraform apply -var-file=build/${CLUSTER}/terraform.tfvars platforms/openstack/<flavor>
```

This should run for a little bit, and when complete, your Tectonic cluster should be ready.

If you encounter any issues, check the known issues and workarounds below.

### Access the cluster

The Tectonic Console should be up and running after the containers have downloaded. You can access it at the DNS name configured in your variables file.

Inside of the `/generated` folder you should find any credentials, including the CA if generated, and a kubeconfig. You can use this to control the cluster with `kubectl`:

```
$ KUBECONFIG=generated/auth/kubeconfig
$ kubectl cluster-info
```

### Delete the cluster

Deleting your cluster will remove only the infrastructure elements created by Terraform. If you selected an existing VPC and subnets, these items are not touched. To delete, run:

```
$ terraform destroy -var-file=build/${CLUSTER}/terraform.tfvars platforms/openstack/<flavor>
```

### Known issues and workarounds

If you experience pod-to-pod networking issues, try lowering the MTU setting of the CNI bridge.
Change the contents of `modules/bootkube/resources/manifests/kube-flannel.yaml` and configure the following settings:

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: kube-flannel-cfg
  namespace: kube-system
  labels:
    tier: node
    k8s-app: flannel
data:
  cni-conf.json: |
    {
      "name": "cbr0",
      "type": "flannel",
      "delegate": {
        "mtu": 1400,
        "isDefaultGateway": true
      }
    }
  net-conf.json: |
    {
      "Network": "${cluster_cidr}",
      "Backend": {
        "Type": "vxlan",
        "Port": 4789
      }
    }
```

Setting the IANA standard port `4789` can help debugging when using `tcpdump -vv -i eth0` on the worker/master nodes as encapsulated VXLAN packets will be shown.

See the [troubleshooting][troubleshooting] document for workarounds for bugs that are being tracked.

[conventions]: ../../conventions.md
[generic]: ../../generic-platform.md
[env]: http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-environment
[register]: https://account.coreos.com/signup/summary/tectonic-2016-12
[account]: https://account.coreos.com
[vars]: ../../variables/config.md
[troubleshooting]: ../../troubleshooting/faq.md
[openstack-nova-vars]: ../../variables/openstack-nova.md
[openstack-neutron-vars]: ../../variables/openstack-neutron.md
