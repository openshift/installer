variable "cluster_id" {
  type = string
}

variable "tags" {
  type        = map(string)
  description = "Tags to be applied to created resources."
}
