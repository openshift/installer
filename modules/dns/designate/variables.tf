variable "etcd_tls_enabled" {
  description = "Indicates whether TLS is used for etcd"
  type        = "string"
  default     = "1"
}

variable "cluster_name" {
  description = "The name of the cluster"
  type        = "string"
}

variable "base_domain" {
  description = "The base domain used in records"
  type        = "string"
}

variable "admin_email" {
  description = "The admin email for the DNS zone"
  type        = "string"
}

variable "master_count" {
  description = "The number of masters"
  type        = "string"
}

variable "worker_count" {
  description = "The number of workers"
  type        = "string"
}

variable "etcd_count" {
  description = "The number of etcd nodes"
  type        = "string"
}

variable "etcd_ips" {
  description = "List of string IPs for etcd nodes"
  type        = "list"
}

variable "master_ips" {
  description = "List of string IPs for masters"
  type        = "list"
}

variable "worker_ips" {
  description = "List of string IPs for workers"
  type        = "list"
}
