variable "cluster_id" {
  type        = string
  description = "The name of the cluster."
}

variable "ignition" {
  type        = string
  description = "The content of the masters ignition file."
}

variable "image" {
  type        = string
  description = "The image for the master instances."
}

variable "instance_count" {
  type        = string
  description = "The number of master instances to launch."
}

variable "labels" {
  type        = map(string)
  description = "GCP labels to be applied to created resources."
  default     = {}
}

variable "machine_type" {
  type        = string
  description = "The machine type for the master instances."
}

variable "subnet" {
  type        = string
  description = "The subnetwork the master instances will be added to."
}

variable "root_volume_size" {
  type        = string
  description = "The size of the volume in gigabytes for the root block device."
}

variable "root_volume_type" {
  type        = string
  description = "The type of volume for the root block device."
}

variable "zones" {
  type = list
}
