# Terraform variables
This document gives an overview of the variables used in the different platforms of the Tectonic SDK.

## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| tectonic_aws_az_count | Number of availability zones the cluster should span. Example: `3` | - | yes |
| tectonic_aws_external_vpc_id | ID of existing VPC to build the cluster into. Example: `vpc-5c73a334` | `` | no |
| tectonic_aws_master_ec2_type | EC2 instance type to use for master nodes. Example: `m4.large` | - | yes |
| tectonic_aws_vpc_cidr_block | IP address range to use when creating the cluster VPC. Example: `10.0.0.0/16` | `10.0.0.0/16` | no |
| tectonic_aws_worker_ec2_type | EC2 instance type to use for worker nodes. Example: `m4.large` | - | yes |
| tectonic_base_domain | The base DNS domain of the cluster. Example: `openstack.dev.coreos.systems` | - | yes |
| tectonic_cluster_name | The name of the cluster. This will be prepended to `tectonic_base_domain` resulting in the URL to the Tectonic console. Example: `demo` | - | yes |
| tectonic_etcd_count | The amount of etcd nodes to be created. Example: `1` | - | yes |
| tectonic_kube_version | The hyperkube "quay.io/coreos/hyperkube" image version. | - | yes |
| tectonic_master_count | The amount of master nodes to be created. Example: `1` | - | yes |
| tectonic_openstack_external_gateway_id | The ID of the network to be used as the external internet gateway as given in `openstack network list`. | `6d6357ac-0f70-4afa-8bd7-c274cc4ea235` | no |
| tectonic_openstack_flavor_id | The flavor ID as given in `openstack flavor list`. Specifies the size (CPU/Memory/Drive) of the VM. | `5cf64088-893b-46b5-9bb1-ee020277635d` | no |
| tectonic_openstack_image_id | The image ID as given in `openstack image list`. Specifies the OS image of the VM. | `edd9e119-a2db-4ccd-a205-5290682254e9` | no |
| tectonic_worker_count | The amount of worker nodes to be created. Example: `3` | - | yes |

