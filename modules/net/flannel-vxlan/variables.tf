variable "flannel_image" {
  description = "Container image for flanneld"
  type        = "string"
}

variable "flannel_cni_image" {
  description = "Container image for flannel cni"
  type        = "string"
}

variable "cluster_cidr" {
  description = "A CIDR notation IP range from which to assign pod IPs"
  type        = "string"
}

variable "enabled" {
  description = "If set to true, flannel networking will be deployed"
  default     = true
}
