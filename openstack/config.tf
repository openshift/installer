variable "flavor_id" {
  type    = "string"
  default = "bbcb7eb5-5c8d-498f-9d7e-307c575d3566"
}

variable "image_id" {
  type    = "string"
  default = "3a0c0bac-fa91-4c96-bfcb-ee215ba1cd4d"
}

variable "external_gateway_id" {
  type = "string"
  default = "6d6357ac-0f70-4afa-8bd7-c274cc4ea235"
}

variable "controller_count" {
  type    = "string"
  default = "1"
}

variable "worker_count" {
  type    = "string"
  default = "3"
}
