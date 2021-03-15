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

variable "metal_facility" {
  type = string
}

variable "metal_project_id" {
  type = string
}

variable "metal_plan" {
  type = string
}

variable "metal_hardware_reservation_id" {
  type = string
}
