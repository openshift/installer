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

// The fqdns of the etcd endpoints.
variable etcd_fqdns {
  type = "list"
}

// The amount of worker nodes to be created.
// Example: `3`
variable "master_count" {
  type = "string"
}

variable "worker_count" {
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

variable "tectonic_versions" {
  type = "map"
}

variable "tectonic_kube_dns_service_ip" {
  type = "string"
}
