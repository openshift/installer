variable "cluster_id" {
  type        = string
  description = "The name of the cluster."
}

variable "bootstrap_present" {
  type        = bool
  description = "If the bootstrap instance and instance_group should exist."
  default     = true
}

variable "ignition" {
  type        = string
  description = "The content of the bootstrap ignition file."
}

variable "image" {
  type        = string
  description = "The image for the bootstrap node."
}

variable "labels" {
  type        = map(string)
  description = "GCP labels to be applied to created resources."
  default     = {}
}

variable "machine_type" {
  type        = string
  description = "Machine type for the bootstrap node."
}

variable "network" {
  type        = string
  description = "The network the bootstrap node will be added to."
}

variable "subnet" {
  type        = string
  description = "The subnetwork the bootstrap node will be added to."
}

variable "root_volume_size" {
  type        = string
  description = "The volume size (in gibibytes) for the bootstrap node's root volume."
  default     = "128"
}

variable "root_volume_type" {
  type        = string
  description = "The volume type for the bootstrap node's root volume."
  default     = "pd-standard"
}

variable "zone" {
  type        = string
  description = "The zone for the bootstrap node."
}
