variable "machine_cidr" {
  type        = "string"
  description = "This is the public network netmask."
}

variable "cluster_domain" {
  type        = "string"
  description = "This is the cluster domain where the API record is created"
}

variable "master_count" {
  type        = "string"
  description = "The number of master IP addresses to obtain from the machine_cidr."
}

variable "worker_count" {
  type        = "string"
  description = "The number of worker IP addresses to obtain from the machine_cidr."
}
