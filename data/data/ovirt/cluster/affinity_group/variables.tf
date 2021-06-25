variable "ovirt_cluster_id" {
  type        = string
  description = "The ID of Cluster"
}

variable "ovirt_affinity_groups" {
  type        = list(object({ name = string, priority = number, description = string, enforcing = string }))
  description = "Control plane affinity groups that will be created."
  default     = []
}

