variable "machine_cidr" {
  type        = "string"
  description = "This is the public network netmask."
}

variable "cluster_domain" {
  type        = "string"
  description = "This is the cluster domain where the API record is created"
}

variable "control_plane_count" {
  type        = "string"
  description = "The number of master IP addresses to obtain from the machine_cidr."
}

variable "compute_count" {
  type        = "string"
  description = "The number of worker IP addresses to obtain from the machine_cidr."
}

variable "ipam" {
  type        = "string"
  description = "The IPAM server to use for IP management."
}

variable "ipam_token" {
  type        = "string"
  description = "The IPAM token to use for requests."
}
