# Terraform variables
This document gives an overview of the variables used in the different platforms of the Tectonic SDK.

## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| tectonic_openstack_external_gateway_id | The ID of the network to be used as the external internet gateway as given in `openstack network list`. | `6d6357ac-0f70-4afa-8bd7-c274cc4ea235` | no |
| tectonic_openstack_flavor_id | The flavor ID as given in `openstack flavor list`. Specifies the size (CPU/Memory/Drive) of the VM. | `5cf64088-893b-46b5-9bb1-ee020277635d` | no |
| tectonic_openstack_floatingip_pool | The name name of the floating IP pool as given in `openstack floating ip list`. This pool will be used to assign floating IPs to worker and master nodes. | `public` | no |
| tectonic_openstack_image_id | The image ID as given in `openstack image list`. Specifies the OS image of the VM. | `acdcd535-5408-40f3-8e88-ad8ebb6507e6` | no |

