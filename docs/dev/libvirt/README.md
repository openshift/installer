# Libvirt HOWTO

Launching clusters via libvirt is especially useful for operator development.

**NOTE:** Some aspects of the installation can be customized through the
`install-config.yaml` file. See
[how to create an install-config.yaml file](../../user/overview.md#multiple-invocations) and [the libvirt platform customization](customization.md) documents.

## One-time setup

It's expected that you will create and destroy clusters often in the course of development. These steps only need to be run once.

Before you begin, install the [build dependencies](../dependencies.md).

### Enable KVM

Make sure you have KVM enabled by checking for the device:

```console
$ ls -l /dev/kvm
crw-rw-rw-+ 1 root kvm 10, 232 Oct 31 09:22 /dev/kvm
```

If it is missing, try some of the ideas [here][kvm-install].

### Install and Enable Libvirt

On CentOS 7, first enable the
[kvm-common](http://mirror.centos.org/centos/7/virt/x86_64/kvm-common/)
repository to ensure you get a new enough version of qemu-kvm.

On Fedora, CentOS/RHEL:

```sh
sudo yum install libvirt-devel libvirt-daemon-kvm libvirt-client
```

Then start libvirtd:

```sh
sudo systemctl enable --now libvirtd
```

### Pick names

In this example, we'll set the base domain to `tt.testing` and the cluster name to `test1`.

### Clone the project

```sh
git clone https://github.com/openshift/installer.git
cd installer
```

### Enable IP Forwarding

Libvirt creates a bridged connection to the host machine, but in order for the
network bridge to work IP forwarding needs to be enabled. The following command
will tell you if forwarding is enabled:

```sh
sysctl net.ipv4.ip_forward
```

If the command output is:

```sh
net.ipv4.ip_forward = 0
```

then forwarding is disabled and proceed with the rest of this section. If IP
forwarding is enabled then skip the rest of this section.

To enable IP forwarding for the current boot:

```sh
sysctl net.ipv4.ip_forward=1
```

or to persist the setting across reboots (recommended):

```sh
echo "net.ipv4.ip_forward = 1" | sudo tee /etc/sysctl.d/99-ipforward.conf
sudo sysctl -p /etc/sysctl.d/99-ipforward.conf
```

### Configure libvirt to accept TCP connections

The Kubernetes [cluster-api](https://github.com/kubernetes-sigs/cluster-api)
components drive deployment of worker machines.  The libvirt cluster-api
provider will run inside the local cluster, and will need to connect back to
the libvirt instance on the host machine to deploy workers.

In order for this to work, you'll need to enable unauthenticated TCP
connections for libvirt.

**NOTE:** The following configuration disables all encryption and authentication
options in libvirtd and causes it to listen on all network interfaces and IP
addresses. **A connection to this privileged libvirtd gives the client
privileges equivalent to those of a root shell.** This configuration has a
security impact on a par with running a telnet server with no root password set.
It is critical to follow the steps below to **configure the firewall to prevent
access to libvirt from other hosts on the LAN/WAN**.

#### For systemd activated libvirt

This applies only if the libvirt daemon is configured to use socket activation.
This is currently the case on Fedora 31 (and later) and Arch Linux.

First, you need to start the libvirtd TCP socket, which is managed by systemd:

```sh
sudo systemctl start libvirtd-tcp.socket
```

To make this change persistent accross reboots you can optionally enable it:

```sh
sudo systemctl enable libvirtd-tcp.socket
```

Then to enable TCP access to libvirtd, modify `/etc/libvirt/libvirtd.conf` and
set the following:

```
auth_tcp = "none"
```

Then restart libvirt:

```sh
sudo systemctl restart libvirtd
```

#### For permanently running libvirt daemon

This applies only if the libvirt daemon is started only through
`libvirtd.service` and without making use of systemd socket activation (through
`libvirtd.socket` and similar systemd units).


For RHEL/CentOS, make sure that the following is set in
`/etc/sysconfig/libvirtd`:

```
LIBVIRTD_ARGS="--listen"
```

For Debian based distros, make sure that the following is set in
`/etc/default/libvirtd`:

```
libvirtd_opts="--listen"
```

Then to enable TCP access to libvirtd, modify `/etc/libvirt/libvirtd.conf` and
set the following:

```
listen_tls = 0
listen_tcp = 1
auth_tcp = "none"
tcp_port = "16509"
```

Then restart libvirt:

```sh
sudo systemctl restart libvirtd
```

#### Configure qemu.conf

On Debian/Ubuntu it might be needed to configure security driver for qemu.
Installer uses terraform libvirt, and it has a known issue, that might cause
unexpected `Could not open '/var/lib/libvirt/images/<FILE_NAME>': Permission denied`
errors. Double check that `security_driver = "none"` line is present in
`/etc/libvirt/qemu.conf` and not commented out.

#### Firewall

Finally, if you have a firewall, you may have to allow connections to the
libvirt daemon from the IP range used by your cluster nodes.

The following examples use the default cluster IP range of `192.168.126.0/24` (which is currently not configurable) and a libvirt `default` subnet of `192.168.122.0/24`, which might be different in your configuration.
If you're uncertain about the libvirt *default* subnet you should be able to see its address using the command `ip -4 a show dev virbr0` or by inspecting `virsh --connect qemu:///system net-dumpxml default`.
Ensure the cluster IP range does not overlap your `virbr0` IP address.

#### iptables

```sh
iptables -I INPUT -p tcp -s 192.168.126.0/24 -d 192.168.122.1 --dport 16509 -j ACCEPT -m comment --comment "Allow insecure libvirt clients"
```

#### Firewalld

If using `firewalld`, the specifics will depend on how your distribution has set
up the various zones. The following instructions should work as is for Fedora,
CentOS, RHEL and Arch Linux.

First, as we don't want to expose the libvirt port externally, we will need to
actively block it:

```sh
sudo firewall-cmd --add-rich-rule "rule service name="libvirt" reject"
```

For systems with libvirt version 5.1.0 and later, libvirt will set new bridged
network interfaces in the `libvirt` zone. We thus need to allow `libvirt`
traffic from the VMs to reach the host:

```sh
sudo firewall-cmd --zone=libvirt --add-service=libvirt
```

For system with an older libvirt, we will move the new bridge interface to a
dedicated network zone and enable incoming libvirt, DNS & DHCP traffic:

```sh
sudo firewall-cmd --zone=dmz --change-interface=tt0
sudo firewall-cmd --zone=dmz --add-service=libvirt
sudo firewall-cmd --zone=dmz --add-service=dns
sudo firewall-cmd --zone=dmz --add-service=dhcp
```

NOTE: When the firewall rules are no longer needed, `sudo firewall-cmd --reload`
will remove the changes made as they were not permanently added. For persistence,
add `--permanent` to the `firewall-cmd` commands and run them a second time.

### Set up NetworkManager DNS overlay

This step allows installer and users to resolve cluster-internal hostnames from your host.

1. Tell NetworkManager to use `dnsmasq`:

    ```sh
    echo -e "[main]\ndns=dnsmasq" | sudo tee /etc/NetworkManager/conf.d/openshift.conf
    ```

2. Tell dnsmasq to use your cluster. The syntax is `server=/<baseDomain>/<firstIP>`.

    For this example:

    ```sh
    echo server=/tt.testing/192.168.126.1 | sudo tee /etc/NetworkManager/dnsmasq.d/openshift.conf
    ```

3. Reload NetworkManager to pick up the `dns` configuration change: `sudo systemctl reload NetworkManager`

## Build the installer

Set `TAGS=libvirt` to add support for libvirt; this is not enabled by default because libvirt is [development only](../../../README.md#supported-platforms).

```sh
TAGS=libvirt hack/build.sh
```

## Run the installer

With [libvirt configured](#install-and-enable-libvirt), you can proceed with [the usual quick-start](../../../README.md#quick-start).

## Cleanup

To remove resources associated with your cluster, run:

```sh
openshift-install destroy cluster
```

You can also use [`virsh-cleanup.sh`](../../../scripts/maintenance/virsh-cleanup.sh), but note that it will currently destroy *all* libvirt resources.

### Firewall

With the cluster removed, you no longer need to allow libvirt nodes to reach your `libvirtd`. Restart
`firewalld` to remove your temporary changes as follows:

```sh
sudo firewall-cmd --reload
```

## Exploring your cluster

Some things you can do:

### SSH access

The bootstrap node, e.g. `test1-bootstrap.tt.testing`, runs the bootstrap process. You can watch it:

```sh
ssh "core@${CLUSTER_NAME}-bootstrap.${BASE_DOMAIN}"
sudo journalctl -f -u bootkube -u openshift
```

You'll have to wait for etcd to reach quorum before this makes any progress.

Using the domain names above will only work if you [set up the DNS overlay](#set-up-networkmanager-dns-overlay) or have otherwise configured your system to resolve cluster domain names.
Alternatively, if you didn't set up DNS on the host, you can use:

```sh
virsh -c "${LIBVIRT_URI}" domifaddr "${CLUSTER_NAME}-master-0"  # to get the master IP
ssh core@$MASTER_IP
```

Here `LIBVIRT_URI` is the libvirt connection URI which you [passed to the installer](../../../README.md#quick-start).

### Inspect the cluster with kubectl

You'll need a `kubectl` binary on your path and [the kubeconfig from your `cluster` call](../../../README.md#connect-to-the-cluster).

```sh
export KUBECONFIG="${DIR}/auth/kubeconfig"
kubectl get --all-namespaces pods
```

Alternatively, you can run `kubectl` from the bootstrap or master nodes.
Use `scp` or similar to transfer your local `${DIR}/auth/kubeconfig`, then [SSH in](#ssh-access) and run:

```sh
export KUBECONFIG=/where/you/put/your/kubeconfig
kubectl get --all-namespaces pods
```

## FAQ

### Libvirt vs. AWS

1. There isn't a load balancer on libvirt.

## Troubleshooting

If following the above steps hasn't quite worked, please review this section for well known issues.

### Console doesn't come up

In case of libvirt there is no wildcard DNS resolution and console depends on the route which is created by auth operator ([Issue #1007](https://github.com/openshift/installer/issues/1007)).
To make it work we need to first create the manifests and edit the `domain` for ingress config, before directly creating the cluster.

- Add another domain entry in the openshift.conf which used by dnsmasq.
Here `tt.testing` is the domain which I choose when running the installer.
Here the IP in the address belong to one of the worker node.

```console
$ cat /etc/NetworkManager/dnsmasq.d/openshift.conf
server=/tt.testing/192.168.126.1
address=/.apps.tt.testing/192.168.126.51
```

- Make sure you restart the NetworkManager after change in `openshift.conf`:

```console
$ sudo systemctl restart NetworkManager
```

- Create the manifests:

```console
$ openshift-install --dir $INSTALL_DIR create manifests
```

- Domain entry in cluster-ingress-02-config.yml file should not contain cluster name:

```console
# Assuming `test1` as cluster name
$ sed -i 's/test1.//' $INSTALL_DIR/manifests/cluster-ingress-02-config.yml
```

- Start the installer to create the cluster:

```console
$ openshift-install --dir $INSTALL_DIR create cluster
```

### Install throws an `Unable to resolve address 'localhost'` error

If you're seeing an error similar to

```
Error: Error refreshing state: 1 error(s) occurred:

* provider.libvirt: virError(Code=38, Domain=7, Message='Unable to resolve address 'localhost' service '-1': Servname not supported for ai_socktype')


FATA[0019] failed to run Terraform: exit status 1
```

it is likely that your install configuration contains three backslashes after the protocol (e.g. `qemu+tcp:///...`), when it should only be two.

### Random domain creation errors due to libvirt race conditon

Depending on your libvirt version you might encounter [a race condition][bugzilla_libvirt_race] leading to an error similar to:

```
* libvirt_domain.master.0: Error creating libvirt domain: virError(Code=43, Domain=19, Message='Network not found: no network with matching name 'test1'')
```
This is also being [tracked on the libvirt-terraform-provider][tfprovider_libvirt_race] but is likely not fixable on the client side, which is why you should upgrade libvirt to >=4.5 or a patched version, depending on your environment.

### MacOS support currently broken

* Support for libvirt on Mac OS [is currently broken and being worked on][brokenmacosissue201].

### Error with firewall initialization on Arch Linux

If you're on Arch Linux and get an error similar to

```
libvirt: “Failed to initialize a valid firewall backend”
```

or

```
error: Failed to start network default
error: internal error: Failed to initialize a valid firewall backend
```

please check out [this thread on superuser][arch_firewall_superuser].

### Github Issue Tracker

You might find other reports of your problem in the [Issues tab for this repository][issues_libvirt] where we ask you to provide any additional information.
If your issue is not reported, please do.

[arch_firewall_superuser]: https://superuser.com/questions/1063240/libvirt-failed-to-initialize-a-valid-firewall-backend
[brokenmacosissue201]: https://github.com/openshift/installer/issues/201
[bugzilla_libvirt_race]: https://bugzilla.redhat.com/show_bug.cgi?id=1576464
[issues_libvirt]: https://github.com/openshift/installer/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+libvirt
[libvirt_selinux_issues]: https://github.com/dmacvicar/terraform-provider-libvirt/issues/142#issuecomment-409040151
[tfprovider_libvirt_race]: https://github.com/dmacvicar/terraform-provider-libvirt/issues/402#issuecomment-419500064
[kvm-install]: http://www.linux-kvm.org/page/FAQ#Preparing_to_use_KVM
