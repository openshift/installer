variable "vpc_id" {
  type        = string
  description = "The VPC id of the bootstrap ECS."
}

variable "vswitch_ids" {
  type        = list(string)
  description = "The VSwitch id of the bootstrap ECS."
}

variable "slb_ids" {
  type        = list(string)
  description = "The load balancer IDs of the bootstrap ECS."
}

variable "sg_master_id" {
  type        = string
  description = "The security group ID of the master ECS."
}
