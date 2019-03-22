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

variable "extra_user_names" {
  type    = "list"
  default = []
}

variable "extra_user_password_hashes" {
  type    = "list"
  default = []
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
  type        = "string"
  description = "This is the public network netmask."
}

/////////
// Bootstrap machine variables
/////////

variable "bootstrap_ignition_url" {
  type = "string"
}

variable "bootstrap_complete" {
  type    = "string"
  default = "false"
}

variable "bootstrap_ip" {
  type        = "string"
  description = "The IP address in the machine_cidr to apply to the bootstrap."
  default     = ""
}

///////////
// Control Plane machine variables
///////////

variable "control_plane_instance_count" {
  type        = "string"
  description = "The number of control plane instances to deploy."
  default     = 3
}

variable "control_plane_ignition" {
  type = "string"
}

variable "control_plane_ips" {
  type        = "list"
  description = "The IP addresses in the machine_cidr to apply to the control plane machines."
  default     = []
}

//////////
// Compute machine variables
//////////

variable "compute_instance_count" {
  type        = "string"
  description = "The number of compute instances to deploy."
  default     = 3
}

variable "compute_ignition" {
  type = "string"
}
