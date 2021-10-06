variable "base_volume_id" {
  type        = string
  description = "The ID of the base volume for the bootstrap node."
}

variable "network_id" {
  type        = string
  description = "The ID of a network resource containing the bootstrap node's addresses."
}

variable "pool" {
  type        = string
  description = "The name of the storage pool."
}
