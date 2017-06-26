# Installing Tectonic workers on Red Hat Enterprise Linux

The flexibility and ease of adding drivers to Red Hat Enterprise Linux allows Kubernetes worker components to be deployed with specialized workloads and hardware. This guide describes using Tectonic Installer to launch these components on a RHEL environment.

## Architecture

### Deployment Ideology

Deployment of Tectonic workers atop Red Hat Enterprise Linux is modeled after the traditional methods of installing software on RHEL. It is expected that users will be familiar with the Red Hat package management system, RPM, as well as its common transport mechanism, YUM/DNF.

Like other Red Hat infrastructure deployments, technical steps may be performed according to user preference. This deployment may be manually executed (as described in this guide), orchestrated through Kickstart/Anaconda, or using other orchestration/configuration management systems like Ansible.

### Execution Ideology

While the *installation* of the Tectonic worker components is designed to fit within a traditional Red Hat focused environment, the *execution* of the binaries are intended to mirror that of CoreOS Container Linux. As such, a utility called `kubelet-wrapper` will spin up a copy of `hyperkube` inside `rkt`. This containerized Kubernetes binary reads its configuration from a combination of configuration files managed by both the administrator and by CoreOS. CoreOS managed files are deployed either in RPM files or via Tectonic operators. When files are deployed via RPM, local overrides are possible (but discouraged). For files deployed via the Tectonic operators, the entire lifecycle is expected to be managed by the Tectonic platform.

## Technical Steps

Deploy a Tectonic worker atop Red Hat Enterprise Linux using the process outlined below.

### Deploy Red Hat Enterprise Linux

Deploy RHEL. Any standard deployment technique may be used, including an optical disk installation, a netbooted installation, or an image based deployment (standard for VMWare and OpenStack). For more information, see the [Red Hat Enterprise Linux Install Documentation][rhel-install].

### Enable "extras" repo

Once basic installation of a host is complete, ensure that the additional Red Hat Enterprise Linux repository `extras` is included. For users of `subscription-manager` this can completed with the command:

```
$ subscription-manager repos --enable=rhel-7-server-extras-rpms
```

For users not leveraging subscription-manager, ensure that the correct URL for the mirror of extras that is to be used is placed in the corresponding file in `/etc/yum.repos.d` and set to `enabled`.

### Install the `tectonic-release` RPM

The `tectonic-release` RPM includes the repo definition for the Tectonic software as well as relevant signing keys. The GPG signing key fingerprint for CoreOS shipped RPMs is:

`3681 363D B1AA 55E0 33A2  7699 CF86 6CFE 1643 1E6A`

Download the RPM from the CoreOS `yum` repository:

```
$ curl -LJO http://yum.prod.coreos.systems/repo/tectonic-rhel/el7/x86_64/Packages/tectonic-release-1.6.2-4.el7.noarch.rpm
```

Verify the signature:

```
$ rpm -qip tectonic-release-1.6.2-4.el7.noarch.rpm
Name        : tectonic-release
Version     : 1.6.2
Release     : 4.el7
Architecture: noarch
Install Date: (not installed)
Group       : System Environment/Base
Size        : 22899
License     : ASL 2.0
Signature   : RSA/SHA256, Fri 02 Jun 2017 01:01:53 PM PDT, Key ID cf866cfe16431e6a
Source RPM  : tectonic-release-1.6.2-4.el7.src.rpm
Build Date  : Fri 02 Jun 2017 01:01:04 PM PDT
Build Host  : buildhost.tectonic.coreos.systems
Relocations : (not relocatable)
URL         : https://coreos.com/tectonic
Summary     : Tectonic release files and repository configuration
Description : Tectonic release files including the /etc/tectonic-version file, signing keys and RPM repository files.
```

Confirm that the signature on the RPM matches the last 16 characters of the fingerprint ID above.

After verifying the signature, install the `tectonic-release` RPM:

```
$ yum localinstall tectonic-release-1.6.2-4.el7.noarch.rpm
```

### Install the `tectonic-worker` RPM

After the `tectonic-release` RPM is installed, complete the installation of the `tectonic-worker` RPM:

```
$ yum install tectonic-worker
```

This will download the relevant dependencies and then prompt to validate the
GPG key installed by the `tectonic-release` RPM.

### Copy the `kubeconfig` file from the Tectonic Installer to the host

The [Tectonic installer][tectonic-installer] generates a `kubeconfig` file which is used by all Tectonic workers to authenticate to the API server. Because this file is identical on all hosts, it can be retrieved from an existing worker, a node in the control plane, or from the assets bundle created by the installer.

To use the `kubeconfig` from the assets bundle, extract the bundle to disk and then change to the root directory of the extracted bundle. The file will be located at the path `generated/auth/kubeconfig`. Copy the file to the worker and place it in the path `/etc/kubernetes/kubeconfig`.

### Configure the DNS service address

A cluster-wide DNS service will be deployed as part of the Tectonic system. To allow the kubelet to discover the location of other pods and services, inform the system of the DNS service address.

The DNS service address can be manually extracted from the file `terraform.tfvars` located in the installer assets directory. It is located under the key `tectonic_kube_dns_service_ip`.

Because the file `terraform.tfvars` is intended for machine consumption, it is often easier to retrieve this value using the utility [jq][jq-util]. If available, this can be done with the command:

```
$ jq .tectonic_kube_dns_service_ip terraform.tfvars
```

Once this value has been retrieved it will be placed in the user managed file `/etc/sysconfig/tectonic-worker` on the host in the field `KUBERNETES_DNS_SERVICE_IP=`.

### Configure Firewalld

The default CNI installation for Tectonic uses VXLAN for its communications with [flannel][flannel-repo], which requires communications between hosts on UDP port 4789. The Kubernetes API also communicates with hosts on TCP port 10250. To simplify the configuration of these options, either allow all communications between cluster members, place the relevant ethernet interfaces into the "trusted" zone using FirewallD, or at a minimum allow `4789/udp` and `10250/tcp`. These last steps can be completed with the commands:

```
$ firewall-cmd --add-port 10250/tcp
$ firewall-cmd --add-port 10250/tcp --permanent
$ firewall-cmd --add-port 4789/udp
$ firewall-cmd --add-port 4789/udp --permanent
```

Note: These settings may not be all inclusive and will not represent relative node ports or other communications which may need to be performed. For more information consult the [Kubernetes Networking][k8s-networking] documentation.

### Enable and start the service

This process is the same as with all systemd hosts. The service as installed by the `tectonic-worker` RPM is called `kubelet`. It can be started with the command:

```
$ systemctl start kubelet.service
```

It will take a number of minutes for the worker to retrieve the relevant assets from Quay.io, bootstrap, and join the cluster. Use `journalctl` to monitor progress:

```
$ journalctl -u kubelet.service
```

NOTE: PolicyKit requires the user to be in a relevant group with access to the journal. By default, Red Hat provides the groups `adm` and `systemd-journal` for this purpose. The command may also be run as the root user.

To ensure the service starts on each boot run the command:

```
$ systemctl enable kubelet.service
```

### SELinux

A policy allowing the Tectonic Worker has not yet been completed, therefore users must run SELinux in Permissive mode. The ability to run in Enforcing mode may be completed in the future.


## Troubleshooting

[rhel-install]: https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/7/html/Installation_Guide/index.html
[tectonic-installer]: https://github.com/coreos/tectonic-installer
[jq-util]: https://stedolan.github.io/jq/
[flannel-repo]: https://github.com/coreos/flannel
[k8s-networking]: https://coreos.com/kubernetes/docs/latest/kubernetes-networking.html
<!-- vim: ts=2 sw=2 tw=80 expandtab:
-->
