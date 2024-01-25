variable "name" {
  type = string
}

variable "image_filepath" {
  type = string
}

variable "cluster_id" {
  type = string
}

variable "resource_group_id" {
  type = string
}

variable "region" {
  type = string
}

variable "tags" {
  type = list(string)
}

variable "cos_resource_instance_crn" {
  type = string
}

variable "endpoint_visibility" {
  type = string
}
