variable "cluster_name" {
  type = "string"
}

variable "compute_iam_role" {
  type        = "string"
  default     = ""
  description = "IAM role to use for the instance profiles of compute nodes."
}
