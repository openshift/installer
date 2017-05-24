// The hyperkube image tag.
variable kube_image_tag {
  type = "string"
}

// The hyperkube image url.
variable kube_image_url {
  type = "string"
}

// The content of the kubeconfig file.
variable kubeconfig_content {
  type = "string"
}

// The content of the /etc/resolv.conf file.
variable resolv_conf_content {
  type = "string"
}

// The amount of nodes to be created.
// Example: `3`
variable "instance_count" {
  type = "string"
}

variable "core_public_keys" {
  type = "list"
}

// The name of the cluster.
// The master hostnames will be prefixed with this.
variable "cluster_name" {
  type = "string"
}

variable "tectonic_kube_dns_service_ip" {
  type = "string"
}

variable "node_labels" {
  type = "string"
}

variable "node_taints" {
  type = "string"
}

variable "hostname_infix" {
  type = "string"
}

variable "bootkube_service" {
  type = "string"
}

variable "tectonic_service" {
  type = "string"
}

variable "tectonic_experimental" {
  default = false
}
