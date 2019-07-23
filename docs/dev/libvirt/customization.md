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
