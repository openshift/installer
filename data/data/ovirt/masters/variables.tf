variable "cluster_id" {
  description = "The ID of Openshift cluster"
}

variable "cluster_domain" {
  description = "The domain name of Openshift cluster"
}

variable "master_count" {
  type        = string
  description = "Number of masters"
  default     = 3
}

variable "ovirt_cluster_id" {
  type        = string
  description = "The ID of oVirt's cluster"
}

variable "ovirt_template_id" {
  type        = string
  description = "The ID of oVirt's VM template"
}

variable "ignition_master" {
  type        = string
  description = "master ignition config"
}
