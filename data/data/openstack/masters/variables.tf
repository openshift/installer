variable "base_image_id" {
  type        = string
  description = "The identifier of the Glance image for master nodes."
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

variable "server_group_name" {
  type        = string
  description = "Name of the server group for the master nodes."
}

variable "additional_network_ids" {
  type        = list(string)
  description = "IDs of additional networks for master nodes."
}
