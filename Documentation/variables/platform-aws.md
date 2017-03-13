# Terraform variables: platform-aws
The Tectonic SDK variables used for: platform-aws.

## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| tectonic_aws_az_count | Number of availability zones the cluster should span. Example: `3` | - | yes |
| tectonic_aws_external_vpc_id | ID of existing VPC to build the cluster into. Example: `vpc-5c73a334` | `` | no |
| tectonic_aws_master_ec2_type | EC2 instance type to use for master nodes. Example: `m4.large` | - | yes |
| tectonic_aws_vpc_cidr_block | IP address range to use when creating the cluster VPC. Example: `10.0.0.0/16` | `10.0.0.0/16` | no |
| tectonic_aws_worker_ec2_type | EC2 instance type to use for worker nodes. Example: `m4.large` | - | yes |

## Outputs

| Name | Description |
|------|-------------|
| endpoints |  |
| foo |  |

