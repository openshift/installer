# vSphere Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, vSphere-specific properties.

## Cluster-scoped properties

* `vCenter` (required string): The domain name or IP address of the vCenter.
* `username` (required string): The username to use to connect to the vCenter.
* `password` (required string): The password to use to connect to the vCenter.
* `datacenter` (required string): The name of the datacenter to use in the vCenter.
* `defaultDatastore` (required string): The default datastore to use for provisioning volumes.

## Machine pools

There are currently no configurable vSphere-specific machine-pool properties.

## Examples

Some example `install-config.yaml` are shown below.
For examples of platform-agnostic configuration fragments, see [here](../customization.md#examples).

### Minimal

An example minimal vSphere install config is:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
platform:
  vSphere:
    vCenter: your.vcenter.example.com
    username: username
    password: password
    datacenter: datacenter
    defaultDatastore: datastore
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```
