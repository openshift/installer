variable "external_rsg_name" {
  default = ""
  type    = "string"
}

variable "tectonic_azure_location" {
  type = "string"
}

variable "tectonic_cluster_name" {
  type = "string"
}

resource "azurerm_resource_group" "tectonic_cluster" {
  count    = "${var.external_rsg_name == "" ? 1 : 0}"
  location = "${var.tectonic_azure_location}"
  name     = "tectonic-cluster-${var.tectonic_cluster_name}"
}

output "name" {
  value = "${var.external_rsg_name == "" ? join("", azurerm_resource_group.tectonic_cluster.*.name) : var.external_rsg_name }"
}
