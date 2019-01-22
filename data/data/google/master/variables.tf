variable "cluster_name" {
  type = "string"
  description = "The name of the cluster."
}

variable "extra_labels" {
  type        = "map"
  default     = {}
  description = "Extra GCP labels to be applied to created resources."
}

variable "ignition" {
  type = "string"
  description = "The content of the master ignition file."
}

variable "image_name" {
  type    = "string"
  default = ""
  description = "The image for the master nodes."
}

variable "instance_count" {
  type = "string"
  description = "The number of masters to launch."
}

variable "instance_type" {
  type = "string"
  description = "The instance type for the master nodes."
}

variable "network" {
  type = "string"
  description = "The network the masters will be added to."
}

variable "subnetwork" {
  type = "string"
  description = "The subnetwork the masters will be added to."
}

variable "root_volume_size" {
  type        = "string"
  description = "The size of the volume in gigabytes for the root block device."
}

variable "root_volume_type" {
  type        = "string"
  description = "The type of volume for the root block device."
}

variable "zones" {
  type = "list"
}
