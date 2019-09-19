variable "base_image" {
  type = string
}

variable "cluster_id" {
  type        = string
  description = "The identifier for the cluster."
}

variable "flavor_name" {
  type = string
}

variable "instance_count" {
  type = string
}

variable "master_sg_ids" {
  type        = list(string)
  default     = ["default"]
  description = "The security group IDs to be applied to the master nodes."
}

variable "master_port_ids" {
  type        = list(string)
  description = "List of port ids for the master nodes"
}

variable "user_data_ign" {
  type = string
}

variable "root_volume_size" {
  type        = number
  description = "The size of the volume in gigabytes for the root block device."
}

variable "root_volume_type" {
  type        = string
  description = "The type of volume for the root block device."
}
