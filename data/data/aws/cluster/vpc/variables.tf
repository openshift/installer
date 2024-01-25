variable "availability_zones" {
  type        = list(string)
  description = "The availability zones in which to provision subnets."
}

variable "edge_zones" {
  type        = list(string)
  default     = []
  description = "The local zones to provision subnets."
}

variable "edge_parent_gw_map" {
  type        = map(string)
  default     = {}
  description = "The parent zone index used to lookup the NAT gateway for private subnets in Local Zone."
}

variable "edge_zones_type" {
  type        = map(string)
  default     = {}
  description = "A map with types of Edge (Local or Wavelength) zones."
}

variable "cidr_blocks" {
  type        = list(string)
  description = "A list of IPv4 CIDRs with 0 index being the main CIDR."
}

variable "cluster_id" {
  type = string
}

variable "publish_strategy" {
  type        = string
  description = "The publishing strategy for endpoints like load balancers"
}

variable "region" {
  type        = string
  description = "The target AWS region for the cluster."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "AWS tags to be applied to created resources."
}

variable "vpc" {
  type        = string
  description = "An existing network (VPC ID) into which the cluster should be installed."
}

variable "public_subnets" {
  type        = list(string)
  description = "Existing public subnets into which the cluster should be installed."
}

variable "private_subnets" {
  type        = list(string)
  description = "Existing private subnets into which the cluster should be installed."
}