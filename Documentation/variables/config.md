# Common Tectonic Terraform variables
All the common Tectonic SDK variables used for *all* platforms.

## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| tectonic_base_domain | The base DNS domain of the cluster. Example: `openstack.dev.coreos.systems` | - | yes |
| tectonic_cluster_name | The name of the cluster. This will be prepended to `tectonic_base_domain` resulting in the URL to the Tectonic console. Example: `demo` | - | yes |
| tectonic_etcd_count | The amount of etcd nodes to be created. Example: `1` | `1` | no |
| tectonic_kube_version | The hyperkube "quay.io/coreos/hyperkube" image version. | - | yes |
| tectonic_master_count | The amount of master nodes to be created. Example: `1` | - | yes |
| tectonic_worker_count | The amount of worker nodes to be created. Example: `3` | - | yes |

