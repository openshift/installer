variable "master_count" {
  type        = string
  description = "Number of masters"
  default     = 3
}

variable "ignition" {
  type        = string
  description = "The content of the master ignition file"
}

variable "hosts" {
  type        = list(map(string))
  description = "Hardware details for hosts"
}

variable "properties" {
  type        = list(map(string))
  description = "Properties for hosts"
}

variable "root_devices" {
  type        = list(map(string))
  description = "Root devices for hosts"
}

variable "driver_infos" {
  type        = list(map(string))
  description = "BMC information for hosts"
}

variable "instance_infos" {
  type        = list(map(string))
  description = "Instance information for hosts"
}
