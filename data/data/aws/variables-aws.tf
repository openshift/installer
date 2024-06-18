variable "aws_config_version" {
  description = <<EOF
(internal) This declares the version of the AWS configuration variables.
It has no impact on generated assets but declares the version contract of the configuration.
EOF

  default = "1.0"
}

variable "custom_endpoints" {
  type = map(string)

  description = <<EOF
(optional) Custom AWS endpoints to override existing services.
Check - https://www.terraform.io/docs/providers/aws/guides/custom-service-endpoints.html

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}

variable "aws_bootstrap_instance_type" {
  type        = string
  description = "Instance type for the bootstrap node. Example: `m4.large`."
}

variable "aws_master_instance_type" {
  type        = string
  description = "Instance type for the master node(s). Example: `m4.large`."
}

variable "aws_ami" {
  type        = string
  description = "AMI for all nodes.  An encrypted copy of this AMI will be used.  Example: `ami-foobar123`."
}

variable "aws_ami_region" {
  type        = string
  description = "Region for the AMI for all nodes.  An encrypted copy of this AMI will be used.  Example: `ami-foobar123`."
}

variable "aws_extra_tags" {
  type = map(string)

  description = <<EOF
(optional) Extra AWS tags to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}

variable "aws_master_root_volume_type" {
  type = string
  description = "The type of volume for the root block device of master nodes."
}

variable "aws_master_root_volume_size" {
  type = string
  description = "The size of the volume in gigabytes for the root block device of master nodes."
}

variable "aws_master_root_volume_iops" {
  type = string

  description = <<EOF
The amount of provisioned IOPS for the root block device of master nodes.
Ignored if the volume type is not io1.
EOF

}

variable "aws_master_root_volume_encrypted" {
  type = bool

  description = <<EOF
Indicates whether the root EBS volume for master is encrypted. Encrypted Amazon EBS volumes
may only be attached to machines that support Amazon EBS encryption.
EOF

}

variable "aws_master_root_volume_kms_key_id" {
  type = string

  description = <<EOF
(optional) Indicates the KMS key that should be used to encrypt the Amazon EBS volume.
If not set and root volume has to be encrypted, the default KMS key for the account will be used.
EOF

  default = ""
}

variable "aws_master_instance_metadata_authentication" {
  type        = string
  default     = "optional"
  description = "The session tokens requirement, also referred to as Instance Metadata Service Version 2 (IMDSv2). Values are optional or required. Defaults to optional."
}

variable "aws_bootstrap_instance_metadata_authentication" {
  type        = string
  default     = "optional"
  description = "The session tokens requirement, also referred to as Instance Metadata Service Version 2 (IMDSv2). Values are optional or required. Defaults to optional."
}

variable "aws_region" {
  type        = string
  description = "The target AWS region for the cluster."
}

variable "aws_master_availability_zones" {
  type        = list(string)
  description = "The availability zones in which to create the masters. The length of this list must match master_count."
}

variable "aws_worker_availability_zones" {
  type        = list(string)
  description = "The availability zones to provision for workers.  Worker instances are created by the machine-API operator, but this variable controls their supporting infrastructure (subnets, routing, etc.)."
}

variable "aws_edge_local_zones" {
  type    = list(string)
  default = []

  description = "The zones to provision subnets for the edge pool. Edge instances are created by the machine-API operator, but this variable controls their supporting infrastructure (subnets, routing, etc.)."
}

variable "aws_edge_parent_zones_index" {
  type    = map(string)
  default = {}

  description = <<EOF
(optional) A map of the Edge's Parent Zone indexes to discover the private's route table ID used by the private subnet from the parent's Zone in the Region.
Each Local or Wavelength Zone is connected to a parent zone in the Region. If the parent zone has private Route tables, the installer
uses the index to associate with the private edge subnets.

Example: `{ "us-east-1-nyc-1a"=5, "us-east-1-wl1-nyc-wlz-1"=2 }`
EOF
}

variable "aws_vpc" {
  type = string
  default = null
  description = "(optional) An existing network (VPC ID) into which the cluster should be installed."
}

variable "aws_public_subnets" {
  type = list(string)
  default = null
  description = "(optional) Existing public subnets into which the cluster should be installed."
}

variable "aws_private_subnets" {
  type = list(string)
  default = null
  description = "(optional) Existing private subnets into which the cluster should be installed."
}

variable "aws_internal_zone" {
  type = string
  default = null
  description = "(optional) An existing hosted zone (zone ID) to use for the internal API."
}

variable "aws_internal_zone_role" {
  type = string
  default = null
  description = "(optional) A role to assume when using an existing hosted zone from another account."
}


variable "aws_publish_strategy" {
  type = string
  description = "The cluster publishing strategy, either Internal or External"
}

variable "aws_ignition_bucket" {
  type = string
  description = "The S3 bucket where the ignition configuration is stored"
}

variable "aws_bootstrap_stub_ignition" {
  type = string
  description = <<EOF
The stub Ignition config that should be used to boot the bootstrap instance. This already points to the presigned URL for the s3 bucket
specified in aws_ignition_bucket.
EOF
}

variable "aws_master_iam_role_name" {
  type        = string
  description = "The name of the IAM role that will be attached to master instances."
  default     = ""
}

variable "aws_worker_iam_role_name" {
  type        = string
  description = "The name of the IAM role that will be attached to worker instances."
  default     = ""
}

variable "aws_preserve_bootstrap_ignition" {
  type        = bool
  description = "The variable that needs to be set to avoid destuction of S3 objects during bootstrap destroy."
}

variable "aws_master_security_groups" {
  type        = list(string)
  description = "(optional) List of additional security group IDs to attach to the master nodes"
  default     = []
}

variable "aws_master_use_spot_instance" {
  type        = bool
  description = "(optional) Whether to use instances from the Spot market for masters"
  default     = false
}
