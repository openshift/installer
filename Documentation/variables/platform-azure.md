# Terraform variables: platform-azure
The Tectonic SDK variables used for: platform-azure.

## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| tectonic_azure_dns_resource_group |  | `tectonic-dns-group` | no |
| tectonic_azure_image_reference | The image ID as given in `azure image list`. Specifies the OS image of the VM. | `<map>` | no |
| tectonic_azure_location |  | `East US` | no |
| tectonic_azure_vm_size | The flavor ID as given in `azure flavor list`. Specifies the size (CPU/Memory/Drive) of the VM. | `Standard_DS2` | no |

