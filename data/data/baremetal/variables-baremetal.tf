variable "bootstrap_provisioning_ip" {
  type        = string
  description = "IP for the bootstrap VM provisioning nic"
}

variable "libvirt_uri" {
  type        = string
  description = "libvirt connection URI"
}

variable "bootstrap_os_image" {
  type        = string
  description = "The URL of the bootstrap OS disk image"
}

variable "external_bridge" {
  type        = string
  description = "The name of the external bridge"
}

variable "provisioning_bridge" {
  type        = string
  description = "The name of the provisioning bridge"
}

variable "ironic_username" {
  type        = string
  description = "Username for authentication to Ironic"
}

variable "ironic_password" {
  type        = string
  description = "Password for authentication to Ironic"
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
