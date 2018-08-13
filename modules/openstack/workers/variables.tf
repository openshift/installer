variable "base_domain" {
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

variable "image_name" {
  type = "string"
}

variable "instance_count" {
  type = "string"
}

variable "key_pair" {
  type = "string"
}

variable "worker_sg_ids" {
  type = "list"
}

variable "user_data_ign" {
  type = "string"
}
