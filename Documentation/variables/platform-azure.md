# Terraform variables
This document gives an overview of the variables used in the different platforms of the Tectonic SDK.

## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| tectonic_azure_dns_resource_group |  | `tectonic-dns-group` | no |
| tectonic_azure_etcd_vm_size |  | `Standard_DS2` | no |
| tectonic_azure_external_vnet_id |  | `` | no |
| tectonic_azure_image_reference | The image ID as given in `azure image list`. Specifies the OS image of the VM. | `<map>` | no |
| tectonic_azure_location |  | - | yes |
| tectonic_azure_master_vm_size |  | `Standard_DS2` | no |
| tectonic_azure_ssh_key | Name of an Azure ssh key to use joe-sfo | - | yes |
| tectonic_azure_vm_size | The flavor ID as given in `azure flavor list`. Specifies the size (CPU/Memory/Drive) of the VM. | `Standard_DS2` | no |
| tectonic_azure_vnet_cidr_block |  | `10.0.0.0/16` | no |
| tectonic_azure_worker_vm_size |  | `Standard_DS2` | no |
| tectonic_ssh_key |  | `` | no |

