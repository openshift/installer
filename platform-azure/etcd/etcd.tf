resource "azurerm_virtual_machine" "etcd_node" {
  count                 = "${var.tectonic_etcd_count}"
  name                  = "${var.tectonic_cluster_name}_etcd_node_${count.index}"
  resource_group_name   = "${var.tectonic_azure_resource_group_name}"
  network_interface_ids = ["${azurerm_network_interface.etcd_nic.id}"]
  vm_size               = "${var.tectonic_azure_vm_size}"
  location              = "${var.tectonic_azure_location}"

  storage_image_reference {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Stable"
    version   = "latest"
  }

  storage_os_disk {
    name          = "etcd-disk"
    vhd_uri       = "${azurerm_storage_account.etcd_storage.primary_blob_endpoint}${azurerm_storage_container.etcd_storage_container.name}/etcd-disk.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "etcd"
    admin_username = "core"
    admin_password = "Microsoft123!"
    custom_data    = "${base64encode("${ignition_config.etcd.rendered}")}"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "random_id" "storage" {
  byte_length = 4
}

resource "azurerm_storage_account" "etcd_storage" {
  name                = "${random_id.storage.hex}"
  resource_group_name = "${var.tectonic_azure_resource_group_name}"
  location            = "${var.tectonic_azure_location}"
  account_type        = "Standard_LRS"
}

resource "azurerm_storage_container" "etcd_storage_container" {
  name                  = "${var.tectonic_cluster_name}-etcd-storage-container"
  resource_group_name   = "${var.tectonic_azure_resource_group_name}"
  storage_account_name  = "${azurerm_storage_account.etcd_storage.name}"
  container_access_type = "private"
  depends_on            = ["azurerm_storage_account.etcd_storage"]
}
