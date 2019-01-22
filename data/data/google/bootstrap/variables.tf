variable "cluster_name" {
  type        = "string"
  description = "The name of the cluster."
}

variable "extra_labels" {
  type        = "map"
  default     = {}
  description = "Extra GCP labels to be applied to created resources."
}

variable "ignition" {
  type        = "string"
  description = "The content of the bootstrap ignition file."
}

variable "image_name" {
  type        = "string"
  description = "The image for the bootstrap node."
}

variable "instance_type" {
  type        = "string"
  default     = "n1-standard-2"
  description = "The instance type for the bootstrap node."
}

variable "network" {
  type = "string"
  description = "The network the bootstrap node will be added to."
}

variable "subnetwork" {
  type = "string"
  description = "The subnetwork the bootstrap node will be added to."
}

variable "root_volume_size" {
  type        = "string"
  default     = "30"
  description = "The volume size (in gibibytes) for the bootstrap node's root volume."
}

variable "root_volume_type" {
  type        = "string"
  default     = "pd-standard"
  description = "The volume type for the bootstrap node's root volume."
}

variable "zone" {
  type        = "string"
  description = "The zone for the bootstrap node."
}
