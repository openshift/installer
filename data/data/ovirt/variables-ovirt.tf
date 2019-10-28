variable "bootstrap_dns" {
  type        = string
  default     = true
  description = "Whether to include DNS entries for the bootstrap node or not."
}

variable "ovirt_url" {
  type        = string
  description = "The oVirt engine URL"
}

variable "ovirt_username" {
  type        = string
  description = "The name of user to access oVirt engine API"
}

variable "ovirt_password" {
  type        = string
  description = "The plain password of user to access oVirt engine API"
}

variable "ovirt_cluster_id" {
  type        = string
  description = "The ID of oVirt's cluster"
}

variable "ovirt_template_id" {
  type        = string
  description = "The ID of oVirt's template"
}
