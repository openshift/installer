variable "depends" {
  type    = any
  default = null
}

variable "ssh_private_key_path" {}
variable "cluster_name" {}
variable "cluster_basedomain" {}
variable "cf_zone_id" {}
variable "ocp_version" {}
variable "ocp_version_zstream" {}
variable "nodes" {
  description = "Generic list of OpenShift node types"
  type        = list(string)
  default     = ["bootstrap", "master", "worker"]
}

