variable "external_rsg_id" {
  default = ""
  type    = "string"
}

variable "azure_location" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "cluster_id" {
  type = "string"
}

variable "boot_diagnostics" {
  type = "string"
}

# Storage ID
resource "random_id" "storage_id" {
  byte_length = 2
}

variable "extra_tags" {
  type = "map"
}
