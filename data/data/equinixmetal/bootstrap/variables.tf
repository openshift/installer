variable "ignition" {
  type        = string
  description = "The content of the bootstrap ignition file."
  default = "TODO"
}






variable "depends" {
  type    = any
  default = null
}

variable "plan" {}
variable "facility" {}
variable "operating_system" {}
variable "project_id" {}
variable "billing_cycle" {}
variable "ssh_private_key_path" { default = "TODO" }
variable "cluster_name" {default = "TODO" }
variable "cluster_basedomain" {default = "TODO" }
// variable "cf_zone_id" {}
variable "ocp_version" {default = "TODO" }
variable "ocp_version_zstream" {default = "TODO" }
variable "nodes" {
  description = "Generic list of OpenShift node types"
  type        = list(string)
  default     = ["bootstrap", "master", "worker"]
}

