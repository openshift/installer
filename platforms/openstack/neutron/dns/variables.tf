variable "tectonic_console_records" {
  type = "list"
}

variable "tectonic_api_records" {
  type = "list"
}

// The name of the cluster.
variable "cluster_name" {
  type = "string"
}

// The base DNS domain of the cluster.
// Example: `openstack.dev.coreos.systems`
variable "base_domain" {
  type = "string"
}
