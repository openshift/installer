variable "aws_config_version" {
  description = <<EOF
(internal) This declares the version of the AWS configuration variables.
It has no impact on generated assets but declares the version contract of the configuration.
EOF

  default = "1.0"
}

variable "aws_bootstrap_instance_type" {
  type = string
  description = "Instance type for the bootstrap node. Example: `m4.large`."
}

variable "aws_master_instance_type" {
  type = string
  description = "Instance type for the master node(s). Example: `m4.large`."
}

variable "aws_ami" {
  type = string
  description = "AMI for all nodes.  Example: `ami-foobar123`."
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
  type        = string
  description = "The type of volume for the root block device of master nodes."
}

variable "aws_master_root_volume_size" {
  type        = string
  description = "The size of the volume in gigabytes for the root block device of master nodes."
}

variable "aws_master_root_volume_iops" {
  type = string

  description = <<EOF
The amount of provisioned IOPS for the root block device of master nodes.
Ignored if the volume type is not io1.
EOF

}

variable "aws_region" {
  type = string
  description = "The target AWS region for the cluster."
}

variable "aws_master_availability_zones" {
  type = list(string)
  description = "The availability zones in which to create the masters. The length of this list must match master_count."
}

variable "aws_worker_availability_zones" {
  type = list(string)
  description = "The availability zones to provision for workers.  Worker instances are created by the machine-API operator, but this variable controls their supporting infrastructure (subnets, routing, etc.)."
}

