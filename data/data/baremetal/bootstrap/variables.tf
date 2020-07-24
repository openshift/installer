variable "cluster_id" {
  type        = string
  description = "The identifier for the cluster."
}

variable "image" {
  description = "The URL of the OS disk image"
  type        = string
}

variable "ignition" {
  type        = string
  description = "The content of the bootstrap ignition file."
}

variable "external_bridge" {
  type        = string
  description = "The name of the bridge providing external access"
}

variable "provisioning_bridge" {
  type        = string
  description = "The name of the bridge used for provisioning"
}
