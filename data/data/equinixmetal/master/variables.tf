variable "ignition" {
  type        = string
  description = "The content of the bootstrap ignition file."
}


variable "plan" {}
variable "node_count" {
  type = number
}
variable "facility" { default = "" }
variable "metro" { default = "" }
variable "cluster_domain" {}
variable "base_domain" {}
variable "ssh_private_key_path" {}
variable "project_id" {}
variable "bootstrap_ip" {}
variable "operating_system" {}
//variable "ocp_version" {default = "TODO" }
//variable "ocp_version_zstream" {default = "TODO" }
variable "depends" {
  type    = any
  default = null
}

variable "ip_addresses" {
  description = "IP Reservation IPv4 addresses to assign to masters"
  type        = list(string)
}
