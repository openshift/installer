// phpIPAM variables

variable "ipam" {
  type        = "string"
  description = "The IPAM server to use for IP management."
  default     = ""
}

variable "ipam_token" {
  type        = "string"
  description = "The IPAM token to use for requests."
  default     = ""
}

// AWS-Specific variables

variable "aws_region" {
  type        = "string"
  description = "The target AWS region for the cluster."
}

variable "aws_control_plane_availability_zones" {
  type        = "list"
  description = "The availability zones in which to create the control plane. The length of this list must match control_plane_count."
}

variable "aws_compute_availability_zones" {
  type        = "list"
  description = "The availability zones to provision for computes.  Compute instances are created by the machine-API operator, but this variable controls their supporting infrastructure (subnets, routing, etc.)."
}

variable "vpc_id" {
  type = "string"
}

variable "aws_public_subnet_id" {
  type = "list"
}

variable "aws_private_subnet_id" {
  type = "list"
}

// TODO: this might need to be a list...but may not work right with vsphere on aws
variable "aws_availability_zone" {
  type = "string"
}

//////
// vSphere variables
//////

variable "vsphere_server" {
  type        = "string"
  description = "This is the vSphere server for the environment."
}

variable "vsphere_user" {
  type        = "string"
  description = "vSphere server user for the environment."
}

variable "vsphere_password" {
  type        = "string"
  description = "vSphere server password"
}

variable "vsphere_cluster" {
  type        = "string"
  description = "This is the name of the vSphere cluster."
}

variable "vsphere_datacenter" {
  type        = "string"
  description = "This is the name of the vSphere data center."
}

variable "vsphere_datastore" {
  type        = "string"
  description = "This is the name of the vSphere data store."
}

variable "vm_template" {
  type        = "string"
  description = "This is the name of the VM template to clone."
}

variable "vm_network" {
  type        = "string"
  description = "This is the name of the publicly accessible network for cluster ingress and access."
  default     = "VM Network"
}

/////////
// OpenShift cluster variables
/////////

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

variable "machine_cidr" {
  type = "string"
}

/////////
// Bootstrap machine variables
/////////

variable "bootstrap_ignition_path" {
  type = "string"
}

variable "bootstrap_complete" {
  type    = "string"
  default = "false"
}

variable "bootstrap_ip_address" {
  type    = "string"
  default = ""
}

///////////
// control plane machine variables
///////////

variable "control_plane_ignition_path" {
  type = "string"
}

variable "control_plane_count" {
  type    = "string"
  default = "3"
}

variable "control_plane_ip_addresses" {
  type    = "list"
  default = []
}

//////////
// compute machine variables
//////////

variable "compute_ignition_path" {
  type = "string"
}

variable "compute_count" {
  type    = "string"
  default = "3"
}

variable "compute_ip_addresses" {
  type    = "list"
  default = []
}

variable "vm_dns_addresses" {
  type    = "list"
  default = ["1.1.1.1", "9.9.9.9"]
}
