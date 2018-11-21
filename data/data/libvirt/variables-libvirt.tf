variable "bootstrap_dns" {
  default     = true
  description = "Whether to include DNS entries for the bootstrap node or not."
}

variable "libvirt_uri" {
  type        = "string"
  description = "libvirt connection URI"
}

variable "libvirt_network_if" {
  type        = "string"
  description = "The name of the bridge to use"
}

variable "libvirt_ip_range" {
  type        = "string"
  description = "IP range for the libvirt machines"
}

variable "os_image" {
  type        = "string"
  description = "The URL of the OS disk image"
}

variable "libvirt_bootstrap_ip" {
  type        = "string"
  description = "the desired bootstrap ip"
}

variable "libvirt_master_ips" {
  type        = "list"
  description = "the list of desired master ips. Must match master_count"
}
