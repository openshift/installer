# Libvirt Platform Customization

The following options are available when using libvirt:

- `platform.libvirt.network.if` - the network bridge attached to the libvirt network (`tt0` by default)

## Examples

An example `install-config.yaml` is shown below. This configuration has been modified to show the customization that is possible via the install config.

```yaml
apiVersion: v1
baseDomain: example.com
...
platform:
  libvirt:
    URI: qemu+tcp://192.168.122.1/system
    network:
      if: mybridge0
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

## Change cpu/memory

- Create manifests

```bash
$ openshift-install --log-level=debug --dir=. create manifest
```

- Edit openshift/99_openshift-cluster-api_master-machines-0.yaml to read (only relevant part shown)

```yaml
apiVersion: machine.openshift.io/v1beta1
kind: Machine
spec:
  providerSpec:
    value:
      apiVersion: libvirtproviderconfig.openshift.io/v1beta1
      domainMemory: 8192
      domainVcpu: 4
```

- Edit openshift/99_openshift-cluster-api_worker-machineset-0.yaml to read (only relevant part shown)

```yaml
apiVersion: machine.openshift.io/v1beta1
kind: MachineSet
spec:
  template:
    spec:
      providerSpec:
        value:
          apiVersion: libvirtproviderconfig.openshift.io/v1beta1
          domainMemory: 16384
          domainVcpu: 6
```

- Create cluster using the manifest

```bash
$ openshift-install --log-level=debug --dir=. create cluster
```
