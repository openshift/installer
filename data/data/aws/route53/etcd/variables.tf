variable "cluster_domain" {
  description = "The domain for the cluster that all DNS records must belong"
  type        = "string"
}

variable "zone_id" {
  description = "The zone_id of the internal route53 zone"
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
