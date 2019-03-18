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

variable "datastore_id" {
  type = "string"
}

variable "network_id" {
  type = "string"
}

variable "vm_template_id" {
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
