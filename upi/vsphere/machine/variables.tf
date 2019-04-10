variable "name" {
  type = "string"
}

variable "instance_count" {
  type = "string"
}

variable "ignition" {
  type    = "string"
  default = ""
}

variable "ignition_url" {
  type    = "string"
  default = ""
}

variable "resource_pool_id" {
  type = "string"
}

variable "folder" {
  type = "string"
}

variable "datastore" {
  type = "string"
}

variable "network" {
  type = "string"
}

variable "cluster_domain" {
  type = "string"
}

variable "datacenter_id" {
  type = "string"
}

variable "template" {
  type = "string"
}

variable "machine_cidr" {
  type = "string"
}

variable "ipam" {
  type = "string"
}

variable "ipam_token" {
  type = "string"
}

variable "ip_addresses" {
  type = "list"
}
