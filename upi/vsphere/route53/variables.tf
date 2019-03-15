variable "cluster_domain" {
  description = "The domain for the cluster that all DNS records must belong"
  type        = "string"
}

variable "etcd_count" {
  description = "The number of etcd members."
  type        = "string"
}

variable "etcd_ip_addresses" {
  description = "List of string IPs for machines running etcd members."
  type        = "list"
  default     = []
}

variable "bootstrap_ip" {
  description = "A String IP for bootstrap machine."
  type        = "string"
}

variable "worker_ips" {
  description = "List of string IPs for worker machines."
  type        = "list"
  default     = []
}

variable "base_domain" {
  description = "The base domain used for public records."
  type        = "string"
}

variable "cluster_id" {
  type        = "string"
  description = "The identifier for the cluster."
}

variable "tags" {
  type        = "map"
  default     = {}
  description = "AWS tags to be applied to created resources."
}
