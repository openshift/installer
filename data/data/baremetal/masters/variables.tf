variable "master_count" {
  type        = string
  description = "Number of masters"
  default     = 3
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

variable "master_ignition_url" {
  type        = string
  description = "The URL of the full ignition"
}

variable "master_ignition_url_ca_cert" {
  type        = string
  description = "Root CA cert of the full ignition URL"
}

variable "master_ignition_url_headers" {
  type        = map(string)
  description = "Headers to use when retrieving master_ignition_url"
}
