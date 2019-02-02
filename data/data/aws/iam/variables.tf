variable "cluster_name" {
  type = "string"
}

variable "tags" {
  type        = "map"
  default     = {}
  description = "AWS tags to be applied to created resources."
}
