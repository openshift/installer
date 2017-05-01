variable "container_images" {
  description = "Container images to use"
  type        = "map"
}

variable "assets_s3_location" {
  type        = "string"
  description = "Location on S3 of the Bootkube/Tectonic assets to use (bucket/key)"
}

variable "kubeconfig_s3_location" {
  type        = "string"
  description = "Location on S3 of the kubeconfig file to use (bucket/key)"
}

variable "kube_dns_service_ip" {
  type        = "string"
  description = "Service IP used to reach kube-dns"
}

variable "kubelet_node_label" {
  type        = "string"
  description = "Label that Kubelet will apply on the node"
}

variable "kubelet_node_taints" {
  type        = "string"
  description = "Taints that Kubelet will apply on the node"
}

variable "etcd_endpoints" {
  type        = "list"
  description = "List of etcd endpoints"
}

variable "etcd_gateway_enabled" {
  description = "Specifies whether the etcd gateway should be enabled or not."
  default     = true
}

variable "bootkube_service" {
  type        = "string"
  description = "The content of the bootkube systemd service unit"
}

variable "tectonic_service" {
  type        = "string"
  description = "The content of the tectonic installer systemd service unit"
}

variable "tectonic_service_disabled" {
  description = "Specifies whether the tectonic installer systemd unit will be disabled. If true, no tectonic assets will be deployed"
  default     = false
}

variable "locksmithd_disabled" {
  description = "Specifies whether locksmith will be disabled or not"
  default     = false
}
