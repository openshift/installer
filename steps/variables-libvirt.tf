variable "tectonic_libvirt_uri" {
  type        = "string"
  description = "libvirt connection URI"
}

variable "tectonic_libvirt_tls_ca_path" {
  type        = "string"
  description = "path to the libvirt CA certificate"
}

variable "tectonic_libvirt_tls_cert_path" {
  type        = "string"
  description = "path to the libvirt client certificate"
}

variable "tectonic_libvirt_tls_key_path" {
  type        = "string"
  description = "path to the libvirt client private key"
}

variable "tectonic_libvirt_network_name" {
  type        = "string"
  description = "Name of the libvirt network to create"
}

variable "tectonic_libvirt_network_if" {
  type        = "string"
  description = "The name of the bridge to use"
}

variable "tectonic_libvirt_ip_range" {
  type        = "string"
  description = "IP range for the libvirt machines"
}

variable "tectonic_os_image" {
  type        = "string"
  description = "The URL of the OS disk image"
}

variable "tectonic_libvirt_bootstrap_ip" {
  type        = "string"
  description = "the desired bootstrap ip"
}

variable "tectonic_libvirt_master_ips" {
  type        = "list"
  description = "the list of desired master ips. Must match tectonic_master_count"
}

variable "tectonic_libvirt_worker_ips" {
  type        = "list"
  description = "the list of desired worker ips. Must match tectonic_worker_count"
}
