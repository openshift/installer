variable "pxe_kernel" {
  type = string
}

variable "pxe_initrd" {
  type = string
}

variable "pxe_kernel_args" {
  type = list(string)
}

variable "matchbox_http_endpoint" {
  type = string
}

variable "cluster_id" {
  type = string
}

variable "igntion_config_content" {
  type = string
}

variable "packet_facility" {
  type = string
}

variable "packet_project_id" {
  type = string
}
