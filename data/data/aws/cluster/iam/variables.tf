variable "cluster_id" {
  type = string
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "AWS tags to be applied to created resources."
}

variable "worker_iam_role_name" {
  type        = string
  description = "The name of the existing role to use for the instance profile for workers"
}
