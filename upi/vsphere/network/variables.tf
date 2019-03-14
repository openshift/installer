variable "machine_cidr" {
  type        = "string"
  description = "This is the public network netmask."
}

variable "master_count" {
  type        = "string"
  description = "The number of master IP addresses to obtain from the machine_cidr."
}

variable "worker_count" {
  type        = "string"
  description = "The number of worker IP addresses to obtain from the machine_cidr."
}

variable "base_domain" {
  type        = "string"
  description = "The base domain to check for DNS records against."
}

variable "cluster_domain" {
  type        = "string"
  description = "The cluster domain to check for DNS records against."
}
