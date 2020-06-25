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

variable "bridges" {
  type        = list(string)
  description = "A list of network interfaces to attach to the bootstrap virtual machine"
}
