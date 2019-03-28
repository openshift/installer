variable "name" {
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

variable "folder_id" {
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

variable "extra_user_names" {
  type = "list"
}

variable "extra_user_password_hashes" {
  type = "list"
}

variable "datacenter_id" {
  type = "string"
}

variable "template" {
  type = "string"
}

variable "pull_secret" {
  type    = "string"
  default = ""
}

variable "machine_cidr" {
  type = "string"
}

variable "ips" {
  type = "list"
}
