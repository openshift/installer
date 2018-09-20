# Libvirt HOWTO

Tectonic has limited support for installing a Libvirt cluster. This is useful especially
for operator development.

## 1. One-time setup
It's expected that you will create and destroy clusters often in the course of development. These steps only need to be run once (or once per RHCOS update).

Before you begin, install the [build dependencies](dependencies.md).

### 1.1 Install and Enable Libvirt
On Fedora:

```sh
sudo dnf install libvirt-daemon
```

or on CentOS / RHEL:

```sh
sudo yum install libvirt-daemon
```

Then start libvirtd:

```sh
sudo systemctl start libvirtd
sudo systemctl enable libvirtd
```

### 1.2 Pick a name and ip range
In this example, we'll set the baseDomain to `tt.testing`, the name to `test1` and the ipRange to `192.168.124.0/24`

### 1.3 Clone the repo
```sh
git clone https://github.com/openshift/installer.git
cd installer
```

### 1.4 (Optional) Download and prepare the operating system image

*By default, the installer will download the latest RHCOS image every time it is invoked. This may be problematic for users who create a large number of clusters or who have limited network bandwidth. The installer allows a local image to be used instead.*

Download the latest RHCOS image (you will need access to the Red Hat internal build systems):

```sh
curl http://aos-ostree.rhev-ci-vms.eng.rdu2.redhat.com/rhcos/images/cloud/latest/rhcos-qemu.qcow2.gz | gunzip > rhcos-qemu.qcow2
```

### 1.5 Get a pull secret
Go to https://account.coreos.com/ and obtain a Tectonic *pull secret*.

### 1.6 Make sure you have permissions for `qemu:///system`
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

### 1.7 Configure libvirt to accept TLS connections

The Kubernetes [cluster-api](https://github.com/kubernetes-sigs/cluster-api)
components drive deployment of worker machines.  The libvirt cluster-api
provider will run inside the local cluster, and will need to connect back to
the libvirt instance on the host machine to deploy workers.

In order for this to work, you'll need to enable TLS connections for libvirt.
To do this, first generate the TLS assets:

```
$ go run ./hack/libvirt-ca/main.go --network="192.168.124.0/24" --out $HOME
```

You can omit the `--network` flag if you're using the default
`192.168.124.0/24` network, and of course store the resulting
certificates and keys wherever you like.

Next, modify your `/etc/libvirt/libvirtd.conf` and set the following:

```
listen_tls = 1
tls_port = "16514"
key_file = "/path/to/serverkey.pem"
cert_file = "/path/to/servercert.pem"
ca_file = "/path/to/cacert.pem"
```

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

### 1.8 Configure default libvirt storage pool

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


### 1.9 Prepare the installer configuration file
1. `cp examples/libvirt.yaml ./`
2. Edit the configuration file:
    1. Set an email and password in the `admin` section
    2. Set a `baseDomain` (to `tt.testing`)
    3. Set the `sshKey` in the `admin` section to the **contents** of an ssh key (e.g. `ssh-rsa AAAA...`)
    4. Set the `name` (e.g. test1)
    5. Look at the `podCIDR` and `serviceCIDR` fields in the `networking` section. Make sure they don't conflict with anything important.
    6. Set the `pullSecret` to your JSON pull secret.
    7. Ensure the `libvirt.uri` IP address matches your virbr0 interface IP address which belongs to the libvirt *default* network.
       If you're uncertain about the libvirt *default* subnet you should be able to see its address using the command `ip -4 a show dev virbr0` or by inspecting `virsh --connect qemu:///system net-dumpxml default`.
    8. Ensure the `libvirt.network.ipRange` does not overlap your virbr0 IP address
    9. (Optional) Change the `image` to the file URL of the operating system image you downloaded (e.g. `file:///home/user/Downloads/rhcos-qemu.qcow2`). This will allow the installer to re-use that image instead of having to download it every time.

### 1.10 Set up NetworkManager DNS overlay
This step is optional, but useful for being able to resolve cluster-internal hostnames from your host.
1. Edit `/etc/NetworkManager/NetworkManager.conf` and set `dns=dnsmasq` in section `[main]`
2. Tell dnsmasq to use your cluster. The syntax is `server=/<baseDomain>/<firstIP>`.

    For this example:

    ```sh
    echo server=/tt.testing/192.168.124.1 | sudo tee /etc/NetworkManager/dnsmasq.d/tectonic.conf
    ```
3. `systemctl restart NetworkManager`

### 1.11 Install the terraform provider
1. Make sure you have the `virsh` binary installed: `sudo dnf install libvirt-client libvirt-devel`
2. Install the libvirt terraform provider:
```sh
GOBIN=~/.terraform.d/plugins go get -u github.com/dmacvicar/terraform-provider-libvirt
```

### 1.12 Cache terrafrom plugins (optional, but makes subsequent runs a bit faster)
```sh
cat <<EOF > $HOME/.terraformrc
plugin_cache_dir = "$HOME/.terraform.d/plugin-cache"
EOF
```

## 2. Build the installer
Following the instructions in the root README:

```sh
bazel build tarball
```

## 3. Create a cluster
```sh
tar -zxf bazel-bin/tectonic-dev.tar.gz
alias tectonic="${PWD}/tectonic-dev/installer/tectonic"
```

Initialize (the environment variables are a convenience):
```sh
tectonic init --config=libvirt.yaml
export CLUSTER_NAME=<the cluster name>
export BASE_DOMAIN=<the base domain>
```

Install ($CLUSTER_NAME is `test1`):
```sh
tectonic install --dir=$CLUSTER_NAME
```

When you're done, destroy:
```sh
tectonic destroy --dir=$CLUSTER_NAME
```
Be sure to destroy, or else you will need to manually use virsh to clean up the leaked resources. The [`virsh-cleanup`](../../scripts/maintenance/virsh-cleanup.sh) script may help with this, but note it will currently destroy *all* libvirt resources.

With the cluster removed, you no longer need to allow libvirt nodes to reach your `libvirtd`. Restart
`firewalld` to remove your temporary changes as follows:

```sh
sudo firewall-cmd --reload
```

## 4. Exploring your cluster
Some things you can do:

### Watch the bootstrap process
The bootstrap node, e.g. test1-bootstrap.tt.testing, runs the tectonic bootstrap process. You can watch it:

```sh
ssh core@$CLUSTER_NAME-bootstrap.$BASE_DOMAIN
sudo journalctl -f -u bootkube -u tectonic
```
You'll have to wait for etcd to reach quorum before this makes any progress.

### Inspect the cluster with kubectl
You'll need a kubectl binary on your path.
```sh
export KUBECONFIG="${PWD}/${CLUSTER_NAME}/generated/auth/kubeconfig"
kubectl get -n tectonic-system pods
```

Alternatively, if you didn't set up DNS on the host (or you want to
do things from the node for other reasons), on the master you can
```sh
host# virsh domifaddr master0  # to get the master IP
host# ssh core@<ip of master>
master0# export KUBECONFIG=/var/opt/tectonic/auth/kubeconfig
master0# kubectl get -n tectonic-system pods
```

### Connect to the cluster console
This will take ~30 minutes to be available. Simply go to `https://${CLUSTER_NAME}-api.${BASE_DOMAIN}:6443/console/` (e.g. `test1.tt.testing`) and log in using the credentials above.



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

it is likely that your install configuration contains three backslashes after the protocol (i.e. `qemu+tcp:///...`), when it should only be two.

### Init throws an `unsupported protocol scheme` error
If you're seeing an error similar to

```
$ tectonic init --config ~/tectonic.libvirt.yaml
FATA[0000] Get : unsupported protocol scheme ""
```

then you're probably missing the `file:///` in the value for `image:` in the install configuration.

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
[issues_libvirt]: https://github.com/openshift/installer/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+libvirt
[libvirt_selinux_issues]: https://github.com/dmacvicar/terraform-provider-libvirt/issues/142#issuecomment-409040151
[tfprovider_libvirt_race]: https://github.com/dmacvicar/terraform-provider-libvirt/issues/402#issuecomment-419500064
