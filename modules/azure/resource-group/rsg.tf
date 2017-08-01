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

variable "extra_tags" {
  type = "map"
}

resource "azurerm_resource_group" "tectonic_cluster" {
  count    = "${var.external_rsg_id == "" ? 1 : 0}"
  location = "${var.azure_location}"
  name     = "tectonic-cluster-${var.cluster_name}"

  tags = "${merge(map(
    "Name", "tectonic-cluster-${var.cluster_name}",
    "tectonicClusterID", "${var.cluster_id}"),
    var.extra_tags)}"
}

output "name" {
  value = "${var.external_rsg_id == "" ? join("", azurerm_resource_group.tectonic_cluster.*.name) : element(split("/", var.external_rsg_id), 4) }"
}
