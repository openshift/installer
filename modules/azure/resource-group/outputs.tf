output "name" {
  value = "${var.external_rsg_id == "" ? element(concat(azurerm_resource_group.tectonic_cluster.*.name, list("")), 0) : element(split("/", var.external_rsg_id), 4)}"
}

output "storage_id" {
  value = "${random_id.storage_id.hex}"
}

output "storage_blob_endpoint" {
  value = "${var.boot_diagnostics == "true" ? element(concat(azurerm_storage_account.tectonic_storage.*.primary_blob_endpoint, list("")), 0) : ""}"
}

output "storage_blob_apikey" {
  value = "${var.boot_diagnostics == "true" ? element(concat(azurerm_storage_account.tectonic_storage.*.primary_access_key, list("")), 0) : ""}"
}

output "storage_name" {
  value = "${var.boot_diagnostics == "true" ? element(concat(azurerm_storage_account.tectonic_storage.*.name, list("")), 0) : ""}"
}
