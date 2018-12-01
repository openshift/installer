variable "tectonic_aws_config_version" {
  description = <<EOF
(internal) This declares the version of the AWS configuration variables.
It has no impact on generated assets but declares the version contract of the configuration.
EOF

  default = "1.0"
}

variable "tectonic_aws_master_ec2_type" {
  type        = "string"
  description = "Instance size for the master node(s). Example: `m4.large`."

  # FIXME: get this wired up to the machine default
  default = "m4.large"
}

variable "tectonic_aws_ec2_ami_override" {
  type        = "string"
  description = "(optional) AMI override for all nodes. Example: `ami-foobar123`."
  default     = ""
}

variable "tectonic_aws_master_extra_sg_ids" {
  description = <<EOF
(optional) List of additional security group IDs for master nodes.

Example: `["sg-51530134", "sg-b253d7cc"]`
EOF

  type    = "list"
  default = []
}

variable "tectonic_aws_vpc_cidr_block" {
  type = "string"

  description = <<EOF
Block of IP addresses used by the VPC.
This should not overlap with any other networks, such as a private datacenter connected via Direct Connect.
EOF
}

variable "tectonic_aws_external_vpc_id" {
  type = "string"

  description = <<EOF
(optional) ID of an existing VPC to launch nodes into.
If unset a new VPC is created.

Example: `vpc-123456`
EOF

  default = ""
}

variable "tectonic_aws_endpoints" {
  description = <<EOF
(optional) If set to "all", the default, then both public and private ingress resources (ELB, A-records) will be created.
If set to "private", then only create private-facing ingress resources (ELB, A-records). No public-facing ingress resources will be created.
If set to "public", then only create public-facing ingress resources (ELB, A-records). No private-facing ingress resources will be provisioned and all DNS records will be created in the public Route53 zone.
EOF
}

variable "tectonic_aws_external_private_zone" {
  default = ""

  description = <<EOF
(optional) If set, the given Route53 zone ID will be used as the internal (private) zone.
This zone will be used to create etcd DNS records as well as internal API and internal Ingress records.
If set, no additional private zone will be created.

Example: `"Z1ILINNUJGTAO1"`
EOF
}

variable "tectonic_aws_external_master_subnet_ids" {
  type = "list"

  description = <<EOF
(optional) List of subnet IDs within an existing VPC to deploy master nodes into.
Required to use an existing VPC, not applicable otherwise.

Example: `["subnet-111111", "subnet-222222", "subnet-333333"]`
EOF

  default = []
}

variable "tectonic_aws_external_worker_subnet_ids" {
  type = "list"

  description = <<EOF
(optional) List of subnet IDs within an existing VPC to deploy worker nodes into.
Required to use an existing VPC, not applicable otherwise.

Example: `["subnet-111111", "subnet-222222", "subnet-333333"]`
EOF

  default = []
}

variable "tectonic_aws_extra_tags" {
  type = "map"

  description = <<EOF
(optional) Extra AWS tags to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}

variable "tectonic_aws_master_root_volume_type" {
  type        = "string"
  default     = "gp2"
  description = "The type of volume for the root block device of master nodes."
}

variable "tectonic_aws_master_root_volume_size" {
  type        = "string"
  default     = "120"
  description = "The size of the volume in gigabytes for the root block device of master nodes."
}

variable "tectonic_aws_master_root_volume_iops" {
  type    = "string"
  default = "400"

  description = <<EOF
The amount of provisioned IOPS for the root block device of master nodes.
Ignored if the volume type is not io1.
EOF
}

variable "tectonic_aws_master_custom_subnets" {
  type    = "map"
  default = {}

  description = <<EOF
(optional) This configures master availability zones and their corresponding subnet CIDRs directly.

Example:
`{ eu-west-1a = "10.0.0.0/20", eu-west-1b = "10.0.16.0/20" }`
EOF
}

variable "tectonic_aws_worker_custom_subnets" {
  type    = "map"
  default = {}

  description = <<EOF
(optional) This configures worker availability zones and their corresponding subnet CIDRs directly.

Example: `{ eu-west-1a = "10.0.64.0/20", eu-west-1b = "10.0.80.0/20" }`
EOF
}

variable "tectonic_aws_region" {
  type        = "string"
  description = "The target AWS region for the cluster."
}

variable "tectonic_aws_installer_role" {
  type    = "string"
  default = ""

  description = <<EOF
(optional) Name of IAM role to use to access AWS in order to deploy the Tectonic Cluster.
The name is also the full role's ARN.

Example:
 * Role ARN  = arn:aws:iam::123456789012:role/tectonic-installer
EOF
}

variable "tectonic_aws_master_iam_role_name" {
  type    = "string"
  default = ""

  description = <<EOF
(optional) Name of IAM role to use for the instance profiles of master nodes.
The name is also the last part of a role's ARN.

Example:
 * Role ARN  = arn:aws:iam::123456789012:role/tectonic-installer
 * Role Name = tectonic-installer
EOF
}

variable "tectonic_aws_worker_iam_role_name" {
  type    = "string"
  default = ""

  description = <<EOF
(optional) Name of IAM role to use for the instance profiles of worker nodes.
The name is also the last part of a role's ARN.

Example:
 * Role ARN  = arn:aws:iam::123456789012:role/tectonic-installer
 * Role Name = tectonic-installer
EOF
}
