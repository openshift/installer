variable "base_image" {
  type = "string"
}

variable "cluster_id" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "flavor_name" {
  type = "string"
}

variable "instance_count" {
  type = "string"
}

variable "master_sg_ids" {
  type        = "list"
  default     = ["default"]
  description = "The security group IDs to be applied to the master nodes."
}

variable "subnet_ids" {
  type = "list"
}

variable "user_data_ign" {
  type = "string"
}

variable "service_vm_fixed_ip" {
  type = "string"
}
