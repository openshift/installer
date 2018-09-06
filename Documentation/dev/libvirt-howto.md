# Libvirt howto

Tectonic has limited support for installing a Libvirt cluster. This is useful especially
for operator development.

## HOW TO:
### 1. One-time setup
It's expected that you will create and destroy clusters often in the course of development. These steps only need to be run once (or once per Container Linux or RHCOS update).

#### 1.1 Pick a name and ip range
In this example, we'll set the baseDomain to `tt.testing`, the name to `test1` and the ipRange to `192.168.124.0/24`

#### 1.2 Clone the repo
```sh
git clone https://github.com/openshift/installer.git
cd installer
```

#### 1.3 Download and prepare the operating system image

#### 1.3a RHCOS

Download the latest RHCOS image (you will need access to the Red Hat internal build systems):

```sh
wget http://aos-ostree.rhev-ci-vms.eng.rdu2.redhat.com/rhcos/images/cloud/latest/rhcos-qemu.qcow2.gz
gunzip rhcos-qemu.qcow2.gz
```

#### 1.3b Container Linux

Download the latest stable Container Linux image:
```sh
wget https://stable.release.core-os.net/amd64-usr/current/coreos_production_qemu_image.img.bz2
bunzip2 coreos_production_qemu_image.img.bz2
```

Because of the greater disk requirements of OpenShift, you'll need to expand the root drive with the following:
```sh
qemu-img resize coreos_production_qemu_image.img +8G
```

#### 1.4 Get a pull secret
Go to https://account.coreos.com/ and obtain a Tectonic *pull secret*.

#### 1.5 Make sure you have permisions for `qemu:///system`
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

#### 1.6 Prepare the configuration file
1. `cp examples/tectonic.libvirt.yaml ./`
1. Edit the configuration file:
    1. Set an email and password in the `admin` section
    1. Set a `baseDomain` (to `tt.testing`)
    1. Set the `sshKey` in the `admin` section to the **contents** of an ssh key (e.g. `ssh-rsa AAAA...`)
    1. Set the `imagePath` to the **absolute** path of the operating system image you downloaded
    1. Set the `name` (e.g. test1)
    1. Look at the `podCIDR` and `serviceCIDR` fields in the `networking` section. Make sure they don't conflict with anything important.
    1. Set the `pullSecret` to your JSON pull secret.

#### 1.7 Set up NetworkManager DNS overlay
This step is optional, but useful for being able to resolve cluster-internal hostnames from your host.
1. Edit `/etc/NetworkManager/NetworkManager.conf` and set `dns=dnsmasq` in section `[main]`
2. Tell dnsmasq to use your cluster. The syntax is `server=/<baseDomain>/<firstIP>`. For this example:
```sh
echo server=/tt.testing/192.168.124.1 | sudo tee /etc/NetworkManager/dnsmasq.d/tectonic.conf
```
3. `systemctl restart NetworkManager`

#### 1.8 Install the terraform provider
1. Make sure you have the `virsh` binary installed: `sudo dnf install libvirt-client libvirt-devel`
2. Install the libvirt terraform provider:
```sh
GOBIN=~/.terraform.d/plugins go get github.com/crawford/terraform-provider-libvirt
```

#### 1.9 Cache terrafrom plugins (optional, but makes subsequent runs a bit faster)
```sh
cat <<EOF > $HOME/.terraformrc
plugin_cache_dir = "$HOME/.terraform.d/plugin-cache"
EOF
```

### 2. Build the installer
Following the instructions in the root README:

```sh
bazel build tarball
```

### 3. Create a cluster
```sh
tar -zxf bazel-bin/tectonic-dev.tar.gz
cd tectonic-dev
export PATH="${PWD}/installer:${PATH}"
```

Initialize (the environment variables are a convenience):
```sh
tectonic init --config=../tectonic.libvirt.yaml
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
Be sure to destroy, or else you will need to manually use virsh to clean up the leaked resources. The [`virsh-cleanup`](../../scripts/maintenance/virsh-cleanup) script may help with this, but note it will currently destroy *all* libvirt resources.

# Exploring your cluster
Some things you can do:

## Watch the bootstrap process
The first master node, e.g. test1-master-0.tt.testing, runs the tectonic bootstrap process. You can watch it:

```sh
ssh core@$CLUSTER_NAME-master-0.$BASE_DOMAIN
sudo journalctl -f -u bootkube -u tectonic
```
You'll have to wait for etcd to reach quorum before this makes any progress.

## Inspect the cluster with kubectl
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

## Connect to the cluster console
This will take ~30 minutes to be available. Simply go to `https://${CLUSTER_NAME}-api.${BASE_DOMAIN}:6443/console/` (e.g. `test1.tt.testing`) and log in using the credentials above.


# Libvirt vs. AWS
1. There isn't a load balancer. This means:
    1. We need to manually remap ports that the loadbalancer would
    2. Only the first server (e.g. master) is actually used. If you want to reach another, you have to manually update the domain name.
