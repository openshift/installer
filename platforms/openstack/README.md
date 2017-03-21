## OpenStack

Prerequsities:

1. The latest Container Linux Alpha (1339.0.0 or later) [uploaded into Glance](https://coreos.com/os/docs/latest/booting-on-openstack.html) and its OpenStack image ID.
1. Since OpenStack nova doesn't provide any DNS registration service, AWS Route53 is being used in this example.
Ensure you have a configured `aws` CLI installation.
1. Ensure you have OpenStack credentials set up, i.e. the environment variables `OS_TENANT_NAME`, `OS_USERNAME`, `OS_PASSWORD`, `OS_AUTH_URL`, `OS_REGION_NAME` are set.
1. Create a folder with the cluster's name under `./build` (e.g. `./build/<cluster-name>`)

Either
- Copy the `assets-<cluster-name>.zip` to `./boot/<cluster-name>`
or
- Create a `config.tfvars` file in the root of the repo having the following properties set.
All other options can be left empty.

```
tectonic_worker_count = <worker count>

tectonic_master_count = <master count>

tectonic_base_domain = "my.domain"

tectonic_cluster_name = "test"

tectonic_admin_email = "me@here.com"

tectonic_admin_password_hash = "<bcrypt-tool encoded password>"

tectonic_license = "..."

tectonic_pull_secret = "..."
```

Note: The `bcrypt-tool` is available at https://github.com/coreos/bcrypt-tool.

### Flavors

The following example flavors are included:

- `nova`: Only Nova computing nodes are being created for etcd, master and worker nodes, assuming the nodes get public IPs assigned.
- `neutron`: A private Neutron network is being created with master/worker nodes exposed via floating IPs connected to an etcd instance via an internal network.

Replace `<flavor>` with either option in the following commands.

Please refer to the variable [Documentation](../../Documentation/generic-platform.md) documentation for flavor specific properties which can be overriden in the `config.tfvars` file.

*Note:* If you are using Nova and experience networking issues between pods on different hosts, you might need to disable spoofing protection.

### Installation

Ensure all *prerequsities* are met.

To create the necessary configuration in the `build/<name>` folder, execute
```
$ make PLATFORM=openstack/<flavor> CLUSTER=<cluster-name> localconfig
```

To generate the Terraform plan, execute
```
$ make PLATFORM=openstack/<flavor> CLUSTER=<cluster-name> plan
```

To apply and create the cluster using Terraform execute
```
$ make PLATFORM=openstack/<flavor> CLUSTER=<cluster-name> apply
```

To destroy the cluster invoke
```
$ make PLATFORM=openstack/<flavor> CLUSTER=<cluster-name> destroy
```
