variable "aws_config_version" {
  description = <<EOF
(internal) This declares the version of the AWS configuration variables.
It has no impact on generated assets but declares the version contract of the configuration.
EOF

  default = "1.0"
}

variable "aws_extra_tags" {
  type = "map"

  description = <<EOF
(optional) Extra AWS tags to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}

variable "aws_region" {
  type        = "string"
  description = "The target AWS region for the cluster."
}

variable "aws_control_plane_availability_zones" {
  type        = "list"
  description = "The availability zones in which to create the control-plane. The length of this list must match control_plane_count."
}

variable "aws_compute_availability_zones" {
  type        = "list"
  description = "The availability zones to provision for computes.  Compute instances are created by the machine-API operator, but this variable controls their supporting infrastructure (subnets, routing, etc.)."
}

variable "bootstrap_ip_address" {
  type = "string"
}

variable "bootstrap_count" {
  type = "string"
}

variable "control_plane_ip_addresses" {
  type = "list"
}

variable "control_plane_count" {
  type = "string"
}

variable "compute_ip_addresses" {
  type = "list"
}

variable "compute_count" {
  type = "string"
}

variable "cluster_id" {
  type        = "string"
  description = "This cluster id must be of max length 27 and must have only alphanumeric or hyphen characters."
}

variable "base_domain" {
  type        = "string"
  description = "The base DNS zone to add the sub zone to."
}

variable "cluster_domain" {
  type        = "string"
  description = "The base DNS zone to add the sub zone to."
}

variable "vpc_id" {
  type = "string"
}

variable "machine_cidr" {
  type = "string"
}

variable "aws_public_subnet_id" {
  type = "list"
}

variable "aws_private_subnet_id" {
  type = "list"
}

variable "aws_availability_zone" {
  type = "string"
}
