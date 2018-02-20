resource "azurerm_resource_group" "tectonic_cluster" {
  count    = "${var.external_rsg_id == "" ? 1 : 0}"
  location = "${var.azure_location}"
  name     = "tectonic-cluster-${var.cluster_name}"

  tags = "${merge(map(
    "Name", "tectonic-cluster-${var.cluster_name}",
    "tectonicClusterID", "${var.cluster_id}"),
    var.extra_tags)}"
}

resource "azurerm_storage_account" "tectonic_storage" {
  count                    = "${var.boot_diagnostics == "true" ? 1 : 0}"
  name                     = "tectonicstorage${random_id.storage_id.hex}"
  resource_group_name      = "${var.external_rsg_id == "" ?  join("", azurerm_resource_group.tectonic_cluster.*.name) : var.external_rsg_id}"
  location                 = "${var.azure_location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = "${merge(map(
    "Name", "${var.cluster_name}-storage",
    "tectonicClusterID", "${var.cluster_id}"),
    var.extra_tags)}"

  depends_on = ["azurerm_resource_group.tectonic_cluster"]
}
