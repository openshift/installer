variable "bootkube_id" {
  type = "string"
}

variable "calico_image" {
  description = "Container image for calico node"
  type        = "string"
}

variable "calico_cni_image" {
  description = "Container image for calico cni"
  type        = "string"
}

variable "cluster_cidr" {
  description = "A CIDR notation IP range from which to assign pod IPs"
  type        = "string"
}

variable "enabled" {
  description = "If set true, calico network policy will be deployed"
}

variable "cni_version" {
  default = "0.3.0"
}

variable "log_level" {
  default = "WARNING"
}

variable kube_apiserver_url {
  type = "string"
}
