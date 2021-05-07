variable "cluster_id" {
  type        = string
  description = "The identifier for the cluster."
}

variable "image" {
  description = "The URL of the OS disk image"
  type        = string
}

variable "pool" {
  type        = string
  description = "The name of the storage pool."
}
