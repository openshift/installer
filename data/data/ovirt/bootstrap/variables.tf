variable "cluster_id" {
  description = "The ID of Openshift cluster"
}

variable "ovirt_cluster_id" {
  type        = string
  description = "The ID of oVirt's cluster"
}

variable "ovirt_template_id" {
  type        = string
  description = "The ID of oVirt's VM template"
}

variable "ignition_bootstrap" {
  type        = string
  description = "bootstrap ignition config"
}
