variable "addresses" {
  type        = "list"
  default     = []
  description = "IP addresses to assign to the boostrap node."
}

variable "base_volume_id" {
  type        = "string"
  description = "The ID of the base volume for the bootstrap node."
}

variable "cluster_name" {
  type        = "string"
  description = "The name of the cluster."
}

variable "ignition" {
  type        = "string"
  description = "The content of the bootstrap ignition file."
}

variable "network_id" {
  type        = "string"
  description = "The ID of a network resource containing the bootstrap node's addresses."
}

variable "memory" {
  type        = "string"
  description = "The memory configured for the bootstrap instance."
  default     = 2048
}

variable "cpus" {
  type        = "string"
  description = "The number of CPUs configured for the bootstrap instance."
  default     = 2
}
