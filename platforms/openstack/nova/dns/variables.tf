variable "tectonic_console_records" {
  type = "list"
}

variable "tectonic_api_records" {
  type = "list"
}

variable "etcd_records" {
  type = "list"
}

variable "master_records" {
  type = "list"
}

variable "master_count" {
  type = "string"
}

variable "worker_records" {
  type = "list"
}

variable "worker_count" {
  type = "string"
}

// The name of the cluster.
// The etcd hostnames will be prefixed with this.
variable "cluster_name" {
  type = "string"
}

// The base DNS domain of the cluster.
// Example: `openstack.dev.coreos.systems`
variable "base_domain" {
  type = "string"
}

output "etcd_fqdn" {
  value = "${aws_route53_record.etcd.fqdn}"
}
