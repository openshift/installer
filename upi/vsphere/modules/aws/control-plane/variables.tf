variable "cluster_id" {
  type = "string"
}

variable "availability_zone" {
  type = "string"
}

variable "instance_count" {
  type = "string"
}

variable "tags" {
  type        = "map"
  default     = {}
  description = "AWS tags to be applied to created resources."
}

variable "target_group_arns" {
  type        = "list"
  default     = []
  description = "The list of target group ARNs for the load balancer."
}

variable "target_group_arns_length" {
  description = "The length of the 'target_group_arns' variable, to work around https://github.com/hashicorp/terraform/issues/12570."
}

variable "ip_addresses" {
  type = "list"
}
