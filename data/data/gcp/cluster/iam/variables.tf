variable "project_id" {
  type        = string
  description = "The target GCP project for the cluster."
}

variable "cluster_id" {
  type = string
}

variable "service_account" {
  type        = string
  description = "The service account used by the instances."
}