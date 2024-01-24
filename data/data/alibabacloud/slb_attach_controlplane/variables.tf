variable "slb_ids" {
  type        = list(string)
  description = "The IDs of the load balancers to which to attach the bootstrap and control plane VMs."
}

variable "master_ecs_ids" {
  type        = list(string)
  description = "The list of control plane instance ids."
}

variable "bootstrap_ecs_id" {
  type        = string
  description = "The ID for the bootstrap instance."
}
