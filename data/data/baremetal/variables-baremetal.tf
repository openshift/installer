variable "libvirt_uri" {
  type        = string
  description = "libvirt connection URI"
}

variable "bootstrap_os_image" {
  type        = string
  description = "The URL of the bootstrap OS disk image"
}

variable "bridges" {
  type        = list(map(string))
  description = "A list of network bridge maps, containing the interface name and optionally the MAC address"
}
