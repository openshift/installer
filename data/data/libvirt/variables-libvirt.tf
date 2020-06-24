variable "bootstrap_dns" {
  default     = true
  description = "Whether to include DNS entries for the bootstrap node or not."
}

variable "libvirt_uri" {
  type        = string
  description = "libvirt connection URI"
}

variable "libvirt_network_if" {
  type        = string
  description = "The name of the bridge to use"
}

variable "os_image" {
  type        = string
  description = "The URL of the OS disk image"
}

variable "libvirt_bootstrap_ip" {
  type        = string
  description = "the desired bootstrap ip"
}

variable "libvirt_master_ips" {
  type        = list(string)
  description = "the list of desired master ips. Must match master_count"
}

# It's definitely recommended to bump this if you can.
variable "libvirt_master_memory" {
  type        = string
  description = "RAM in MiB allocated to masters"
  default     = "6144"
}

# At some point this one is likely to default to the number
# of physical cores you have.  See also
# https://pagure.io/standard-test-roles/pull-request/223
variable "libvirt_master_vcpu" {
  type        = string
  description = "CPUs allocated to masters"
  default     = "4"
}

variable "libvirt_bootstrap_memory" {
  type        = number
  description = "RAM in MiB allocated to the bootstrap node"
  default     = 2048
}

# Currently RHCOS maintain its default 16G size if that
# changes we need to change it here also
# https://github.com/coreos/coreos-assembler/pull/924
variable "libvirt_master_size" {
  type        = string
  description = "Size of the volume in bytes"
  default     = "17179869184"
}
