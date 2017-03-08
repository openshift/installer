# Terraform variables
This document gives an overview of the variables used in the different platforms of the Tectonic SDK.

## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| tectonic_base_domain | The base DNS domain of the cluster | `openstack.dev.coreos.systems` | no |
| tectonic_cluster_name | The name of the cluster. This will be prepended to "tectonic_base_domain" resulting in the URL to the Tectonic console. | `demo` | no |
| tectonic_etcd_count | The amount of etcd nodes to be created. | `1` | no |
| tectonic_kube_version | The hyperkube "quay.io/coreos/hyperkube" image version. | - | yes |
| tectonic_master_count | The amount of master nodes to be created | `1` | no |
| tectonic_openstack_external_gateway_id | The ID of the network to be used as the external internet gateway as given in `openstack network list`. | `6d6357ac-0f70-4afa-8bd7-c274cc4ea235` | no |
| tectonic_openstack_flavor_id | The flavor ID as given in `openstack flavor list`. Specifies the size (CPU/Memory/Drive) of the VM. | `5cf64088-893b-46b5-9bb1-ee020277635d` | no |
| tectonic_openstack_image_id | The image ID as given in `openstack image list`. Specifies the OS image of the VM. | `acdcd535-5408-40f3-8e88-ad8ebb6507e6` | no |
| tectonic_worker_count | The amount of worker nodes to be created | `3` | no |

