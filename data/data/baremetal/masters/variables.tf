variable "master_count" {
  type        = string
  description = "Number of masters"
  default     = 3
}

variable "ignition" {
  type        = string
  description = "The content of the master ignition file"
}

variable "masters" {
  type        = list(map(string))
  description = "Hardware details for masters"
}

variable "properties" {
  type        = list(map(string))
  description = "Properties for masters"
}

variable "root_devices" {
  type        = list(map(string))
  description = "Root devices for masters"
}

variable "driver_infos" {
  type        = list(map(string))
  description = "BMC information for masters"
}

variable "instance_infos" {
  type        = list(map(string))
  description = "Instance information for masters"
}
