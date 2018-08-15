# Libvirt HOWTO

Launching clusters via libvirt is especially useful for operator development.

## One-time setup
It's expected that you will create and destroy clusters often in the course of development. These steps only need to be run once.

Before you begin, install the [build dependencies](dependencies.md).

### Enable KVM

Make sure you have KVM enabled by checking for the device:

```console
$ ls -l /dev/kvm 
crw-rw-rw-+ 1 root kvm 10, 232 Oct 31 09:22 /dev/kvm
```

If it is missing, try some of the ideas [here][kvm-install].

### Install and Enable Libvirt
On Fedora, CentOS/RHEL:

```sh
sudo yum install libvirt libvirt-devel
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

### Get a pull secret
Go to https://account.coreos.com/ and obtain a *pull secret*.

### Make sure you have permissions for `qemu:///system`
You may want to grant yourself permissions to use libvirt as a non-root user. You could allow all users in the wheel group by doing the following:
```sh
cat <<EOF >> /etc/polkit-1/rules.d/80-libvirt.rules
polkit.addRule(function(action, subject) {
  if (action.id == "org.libvirt.unix.manage" && subject.local && subject.active && subject.isInGroup("wheel")) {
      return polkit.Result.YES;
  }
});
EOF
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

In order for this to work, you'll need to enable TCP connections for libvirt.

#### Configure libvirtd.conf
To do this, first modify your `/etc/libvirt/libvirtd.conf` and set the
following:
```
listen_tls = 0
listen_tcp = 1
auth_tcp="none"
tcp_port = "16509"
```

Note that authentication is not currently supported, but should be soon.

#### Configure the service runner to pass `--listen` to libvirtd
In addition to the config, you'll have to pass an additional command-line
argument to libvirtd. On Fedora, modify `/etc/sysconfig/libvirtd` and set:

```
LIBVIRTD_ARGS="--listen"
```

On Debian based distros, modify `/etc/default/libvirtd` and set:

```
libvirtd_opts="--listen"
```

Next, restart libvirt: `systemctl restart libvirtd`

#### Firewall
Finally, if you have a firewall, you may have to allow connections to the
libvirt daemon from the IP range used by your cluster nodes.

The following examples use the default cluster IP range of `192.168.126.0/24` (which is currently not configurable) and a libvirt `default` subnet of `192.168.124.0/24`, which might be different in your configuration.
If you're uncertain about the libvirt *default* subnet you should be able to see its address using the command `ip -4 a show dev virbr0` or by inspecting `virsh --connect qemu:///system net-dumpxml default`.
Ensure the cluster IP range does not overlap your `virbr0` IP address.

#### iptables

```sh
iptables -I INPUT -p tcp -s 192.168.126.0/24 -d 192.168.124.1 --dport 16509 -j ACCEPT -m comment --comment "Allow insecure libvirt clients"
```

#### Firewalld

If using `firewalld`, simply obtain the name of the existing active zone which
can be used to integrate the appropriate source and ports to allow connections from
the IP range used by your cluster nodes. An example is shown below.

```console
$ sudo firewall-cmd --get-active-zones
FedoraWorkstation
  interfaces: enp0s25 tun0
```
With the name of the active zone, include the source and port to allow connections
from the IP range used by your cluster nodes.

```sh
sudo firewall-cmd --zone=FedoraWorkstation --add-source=192.168.126.0/24
sudo firewall-cmd --zone=FedoraWorkstation --add-port=16509/tcp
```

Verification of the source and port can be done listing the zone

```sh
sudo firewall-cmd --zone=FedoraWorkstation --list-ports
sudo firewall-cmd --zone=FedoraWorkstation --list-sources
```

NOTE: When the firewall rules are no longer needed, `sudo firewall-cmd --reload`
will remove the changes made as they were not permanently added. For persistence,
add `--permanent` to the `firewall-cmd` commands and run them a second time.

### Configure default libvirt storage pool

Check to see if a default storage pool has been defined in Libvirt by running
`virsh --connect qemu:///system pool-list`.  If it does not exist, create it:

```sh
sudo virsh pool-define /dev/stdin <<EOF
<pool type='dir'>
  <name>default</name>
  <target>
    <path>/var/lib/libvirt/images</path>
  </target>
</pool>
EOF

sudo virsh pool-start default
sudo virsh pool-autostart default
```

### Set up NetworkManager DNS overlay

This step is optional, but useful for being able to resolve cluster-internal hostnames from your host.
1. Edit `/etc/NetworkManager/NetworkManager.conf` and set `dns=dnsmasq` in section `[main]`
2. Tell dnsmasq to use your cluster. The syntax is `server=/<baseDomain>/<firstIP>`.

    For this example:

    ```sh
    echo server=/tt.testing/192.168.126.1 | sudo tee /etc/NetworkManager/dnsmasq.d/tectonic.conf
    ```
3. `systemctl restart NetworkManager`

### Install the Terraform provider

1. Make sure you have the `virsh` binary installed: `sudo dnf install libvirt-client libvirt-devel`
2. Install the libvirt terraform provider:
```sh
GOBIN=~/.terraform.d/plugins go get -u github.com/dmacvicar/terraform-provider-libvirt
```

### Cache Terrafrom plugins (optional, but makes subsequent runs a bit faster)

```sh
cat <<EOF > $HOME/.terraformrc
plugin_cache_dir = "$HOME/.terraform.d/plugin-cache"
EOF
```

## Build and run the installer

With [libvirt configured](#install-and-enable-libvirt), you can proceed with [the usual quick-start](../../README.md#quick-start).
Set `TAGS` when building if you need `destroy cluster` support for libvirt; this is not enabled by default because it requires [cgo][]:

```sh
TAGS=libvirt_destroy hack/build.sh
```

To avoid being prompted repeatedly, you can set [environment variables](../user/environment-variables.md) to reflect your libvirt choices.  For example, selecting libvirt, setting [our earlier name choices](#pick-names), [our pull secret](#get-a-pull-secret), and telling both the installer and the machine-API operator to contact `libvirtd` at [the usual libvirt IP](#firewall), you can use:

```sh
export OPENSHIFT_INSTALL_PLATFORM=libvirt
export OPENSHIFT_INSTALL_BASE_DOMAIN=tt.testing
export OPENSHIFT_INSTALL_CLUSTER_NAME=test1
export OPENSHIFT_INSTALL_PULL_SECRET_PATH=path/to/your/pull-secret.json
export OPENSHIFT_INSTALL_LIBVIRT_URI=qemu+tcp://192.168.122.1/system
```

## Cleanup

If you compiled with `libvirt_destroy`, you can use:

```sh
openshift-install destroy cluster
```

If you did not compile with `libvirt_destroy`, you can use [`virsh-cleanup.sh`](../../scripts/maintenance/virsh-cleanup.sh), but note it will currently destroy *all* libvirt resources.

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
ssh core@$OPENSHIFT_INSTALL_CLUSTER_NAME-bootstrap.$OPENSHIFT_INSTALL_BASE_DOMAIN
sudo journalctl -f -u bootkube -u tectonic
```

You'll have to wait for etcd to reach quorum before this makes any progress.

Using the domain names above will only work if you [set up the DNS overlay](#set-up-networkmanager-dns-overlay) or have otherwise configured your system to resolve cluster domain names.
Alternatively, if you didn't set up DNS on the host, you can use:

```sh
virsh -c "${OPENSHIFT_INSTALL_LIBVIRT_URI}" domifaddr "${OPENSHIFT_INSTALL_CLUSTER_NAME}-master-0"  # to get the master IP
ssh core@$MASTER_IP
```

Here `OPENSHIFT_INSTALL_LIBVIRT_URI` is the libvirt connection URI which you [passed to the installer](#build-and-run-the-installer).

### Inspect the cluster with kubectl

You'll need a `kubectl` binary on your path and [the kubeconfig from your `cluster` call](../../README.md#kubeconfig).

```sh
export KUBECONFIG="${DIR}/auth/kubeconfig-admin"
kubectl get --all-namespaces pods
```

Alternatively, you can run `kubectl` from the bootstrap or master nodes.
Use `scp` or similar to transfer your local `${DIR}/auth/kubeconfig-admin`, then [SSH in](#ssh-access) and run:

```sh
export KUBECONFIG=/where/you/put/your/kubeconfig
kubectl get --all-namespaces pods
```

## FAQ

### Libvirt vs. AWS
1. There isn't a load balancer on libvirt.

## Troubleshooting
If following the above steps hasn't quite worked, please review this section for well known issues.

### Install throws an `Unable to resolve address 'localhost'` error

If you're seeing an error similar to

```
Error: Error refreshing state: 1 error(s) occurred:

* provider.libvirt: virError(Code=38, Domain=7, Message='Unable to resolve address 'localhost' service '-1': Servname not supported for ai_socktype')


FATA[0019] failed to run Terraform: exit status 1
```

it is likely that your install configuration contains three backslashes after the protocol (e.g. `qemu+tcp:///...`), when it should only be two.

### SELinux might prevent access to image files
Configuring the storage pool to store images in a path incompatible with the SELinux policies (e.g. your home directory) might lead to the following errors:

```
Error: Error applying plan:

1 error(s) occurred:

* libvirt_domain.etcd: 1 error(s) occurred:

* libvirt_domain.etcd: Error creating libvirt domain: virError(Code=1, Domain=10, Message='internal error: process exited while connecting to monitor: 2018-07-30T22:52:54.865806Z qemu-kvm: -fw_cfg name=opt/com.coreos/config,file=/home/user/VirtualMachines/etcd.ign: can't load /home/user/VirtualMachines/etcd.ign')
```

[As described here][libvirt_selinux_issues] you can workaround by disabling SELinux, or store the images in a place well-known to work, e.g. by using the default pool.

### Random domain creation errors due to libvirt race conditon
Depending on your libvirt version you might encounter [a race condition][bugzilla_libvirt_race] leading to an error similar to:

```
* libvirt_domain.master.0: Error creating libvirt domain: virError(Code=43, Domain=19, Message='Network not found: no network with matching name 'tectonic'')
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
[cgo]: https://golang.org/cmd/cgo/
[issues_libvirt]: https://github.com/openshift/installer/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+libvirt
[libvirt_selinux_issues]: https://github.com/dmacvicar/terraform-provider-libvirt/issues/142#issuecomment-409040151
[tfprovider_libvirt_race]: https://github.com/dmacvicar/terraform-provider-libvirt/issues/402#issuecomment-419500064
[kvm-install]: http://www.linux-kvm.org/page/FAQ#Preparing_to_use_KVM
