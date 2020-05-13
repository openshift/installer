// phpIPAM variables

variable "ipam" {
  type        = string
  description = "The IPAM server to use for IP management."
  default     = ""
}

variable "ipam_token" {
  type        = string
  description = "The IPAM token to use for requests."
  default     = ""
}

//////
// vSphere variables
//////

variable "vsphere_server" {
  type        = string
  description = "This is the vSphere server for the environment."
}

variable "vsphere_user" {
  type        = string
  description = "vSphere server user for the environment."
}

variable "vsphere_password" {
  type        = string
  description = "vSphere server password"
}

variable "vsphere_cluster" {
  type        = string
  description = "This is the name of the vSphere cluster."
}

variable "vsphere_datacenter" {
  type        = string
  description = "This is the name of the vSphere data center."
}

variable "vsphere_datastore" {
  type        = string
  description = "This is the name of the vSphere data store."
}

variable "vm_template" {
  type        = string
  description = "This is the name of the VM template to clone."
}

variable "vm_network" {
  type        = string
  description = "This is the name of the publicly accessible network for cluster ingress and access."
  default     = "VM Network"
}

variable "vm_dns_addresses" {
  type    = list(string)
  default = ["1.1.1.1", "9.9.9.9"]
}

/////////
// OpenShift cluster variables
/////////

variable "cluster_id" {
  type        = string
  description = "This cluster id must be of max length 27 and must have only alphanumeric or hyphen characters."
}

variable "base_domain" {
  type        = string
  description = "The base DNS zone to add the sub zone to."
}

variable "cluster_domain" {
  type        = string
  description = "The base DNS zone to add the sub zone to."
}

variable "machine_cidr" {
  type = string
}

/////////
// Bootstrap machine variables
/////////

variable "bootstrap_ignition_path" {
  type    = string
  default = "./bootstrap.ign"
}

variable "bootstrap_complete" {
  type    = string
  default = "false"
}

variable "bootstrap_ip_address" {
  type    = string
  default = ""
}

variable "lb_ip_address" {
  type    = string
  default = ""
}

///////////
// control-plane machine variables
///////////

variable "control_plane_ignition_path" {
  type    = string
  default = "./master.ign"
}

variable "control_plane_count" {
  type    = string
  default = "3"
}

variable "control_plane_ip_addresses" {
  type    = list(string)
  default = []
}
variable "control_plane_memory" {
  type    = string
  default = "16384"
}

variable "control_plane_num_cpus" {
  type    = string
  default = "4"
}

//////////
// compute machine variables
//////////

variable "compute_ignition_path" {
  type    = string
  default = "./worker.ign"
}

variable "compute_count" {
  type    = string
  default = "3"
}

variable "compute_ip_addresses" {
  type    = list(string)
  default = []
}

variable "compute_memory" {
  type    = string
  default = "8192"
}

variable "compute_num_cpus" {
  type    = string
  default = "4"
}

variable "ssh_public_key_path" {
  type    = string
  default = "~/.ssh/id_rsa.pub"
}
