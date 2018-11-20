variable "cidr_block" {
  type = "string"
}

variable "cluster_id" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "external_master_subnet_ids" {
  type = "list"
}

variable "external_network" {
  description = "UUID of the external network providing Floating IP addresses."
  type        = "string"
  default     = ""
}

variable "masters_count" {
  type = "string"
}
