# Terraform variables
This document gives an overview of the variables used in the different platforms of the Tectonic SDK.

## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| tectonic_aws_az_count |  | - | yes |
| tectonic_aws_etcd_ec2_type |  | - | yes |
| tectonic_aws_external_master_subnet_ids |  | `<list>` | no |
| tectonic_aws_external_vpc_id |  | `` | no |
| tectonic_aws_external_worker_subnet_ids |  | `<list>` | no |
| tectonic_aws_master_ec2_type |  | - | yes |
| tectonic_aws_ssh_key |  | - | yes |
| tectonic_aws_vpc_cidr_block |  | `10.0.0.0/16` | no |
| tectonic_aws_worker_ec2_type |  | - | yes |

