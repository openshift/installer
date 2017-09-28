variable "container_images" {
  description = "Container images to use"
  type        = "map"
}

variable "image_re" {
  description = <<EOF
(internal) Regular expression used to extract repo and tag components from image strings
EOF

  type = "string"
}

variable "kubelet_cni_bin_dir" {
  type = "string"
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
  description = "(optional) Taints that Kubelet will apply on the node"
  default     = ""
}

variable "kubeconfig_fetch_cmd" {
  type        = "string"
  description = "(optional) The command that fetches `/etc/kubernetes/kubeconfig`."
  default     = ""
}

variable "cloud_provider" {
  type        = "string"
  description = "(optional) The cloud provider to be used for the kubelet."
  default     = ""
}

variable "cloud_provider_config" {
  type        = "string"
  description = "(optional) The cloud provider config to be used for the kubelet."
  default     = ""
}

variable "bootstrap_upgrade_cl" {
  type        = "string"
  description = "(optional) Whether to trigger a ContainerLinux OS upgrade during the bootstrap process."
  default     = "true"
}

variable "torcx_store_url" {
  type        = "string"
  description = "(optional) URL template for torcx store. Leave empty to use the default CoreOS endpoint."
  default     = ""
}

variable "tectonic_vanilla_k8s" {
  description = <<EOF
If set to true, a vanilla Kubernetes cluster will be deployed, omitting any Tectonic assets.
EOF
}

variable "assets_location" {
  type        = "string"
  description = "(optional) The storing location to retrieve the assets."
  default     = ""
}

variable "etcd_count" {
  type    = "string"
  default = 0
}

variable "etcd_tls_enabled" {
  type    = "string"
  default = true
}

variable "etcd_advertise_name_list" {
  type    = "list"
  default = []
}

variable "etcd_initial_cluster_list" {
  type    = "list"
  default = []
}

variable "base_domain" {
  type    = "string"
  default = ""
}

variable "cluster_name" {
  type    = "string"
  default = ""
}

variable "metadata_provider" {
  type    = "string"
  default = ""
}

variable "use_metadata" {
  default = true
}
