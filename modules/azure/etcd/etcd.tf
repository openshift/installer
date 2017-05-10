resource "azurerm_availability_set" "etcd" {
  name                = "${var.cluster_name}-etcd"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"
}

resource "azurerm_virtual_machine" "etcd_node" {
  count                 = "${var.etcd_count}"
  name                  = "${var.cluster_name}-etcd-${count.index}"
  resource_group_name   = "${var.resource_group_name}"
  network_interface_ids = ["${azurerm_network_interface.etcd_nic.*.id[count.index]}"]
  vm_size               = "${var.vm_size}"
  location              = "${var.location}"
  availability_set_id   = "${azurerm_availability_set.etcd.id}"

  storage_image_reference {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Stable"
    version   = "latest"
  }

  storage_os_disk {
    name          = "etcd-disk"
    vhd_uri       = "${azurerm_storage_account.etcd_storage.primary_blob_endpoint}${azurerm_storage_container.etcd_storage_container.name}/etcd-disk-${count.index}.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "etcd"
    admin_username = "core"
    admin_password = ""
    custom_data    = "${base64encode("${data.ignition_config.etcd.*.rendered[count.index]}")}"
  }

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/core/.ssh/authorized_keys"
      key_data = "${file(var.public_ssh_key)}"
    }
  }
}

resource "random_id" "storage" {
  byte_length = 4
}

resource "azurerm_storage_account" "etcd_storage" {
  name                = "${random_id.storage.hex}"
  resource_group_name = "${var.resource_group_name}"
  location            = "${var.location}"
  account_type        = "Premium_LRS"
}

resource "azurerm_storage_container" "etcd_storage_container" {
  name                  = "${var.cluster_name}-etcd-storage-container"
  resource_group_name   = "${var.resource_group_name}"
  storage_account_name  = "${azurerm_storage_account.etcd_storage.name}"
  container_access_type = "private"
  depends_on            = ["azurerm_storage_account.etcd_storage"]
}
