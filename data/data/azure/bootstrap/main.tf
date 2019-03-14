
resource "random_string" "storage_suffix" {
  length = 5
  upper = false
  special = false
}

resource "azurerm_storage_account" "ignition" {
  name = "ignitiondata${random_string.storage_suffix.result}"
  resource_group_name      = "${azurerm_resource_group.ignition.name}"
  location                 = "${var.region}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_resource_group" "ignition" {
  name      = "${var.resource_group_name}"
  location  = "${var.region}"
}

resource "azurerm_storage_container" "ignition" {
  resource_group_name   = "${azurerm_resource_group.ignition.name}"
  name                  = "ignition"
  storage_account_name  = "${azurerm_storage_account.ignition.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "ignition" {
  name                    = "bootstrap.ign"
  source                  = "${var.ignition}"
  resource_group_name     = "${azurerm_resource_group.ignition.name}"
  storage_account_name    = "${azurerm_storage_account.ignition.name}"
  storage_container_name  = "${azurerm_storage_container.ignition.name}"
  type = "block"
}

data "ignition_config" "redirect" {
  replace {
    source = "${azurerm_storage_blob.ignition.url}/bootstrap.ign"
  }
}

data "azurerm_subscription" "primary" {}

resource "azurerm_network_interface" "ignition" {
  name                = "${var.cluster_id}-bootstrap-nic"
  location            = "${azurerm_resource_group.ignition.location}"
  resource_group_name = "${azurerm_resource_group.ignition.name}"

  ip_configuration {
    subnet_id                               = "${var.subnet_id}"
    name                                    = "bootstrap"
    private_ip_address_allocation           = "Dynamic"
  }
}

resource "azurerm_network_interface_backend_address_pool_association" "ignition" {
  network_interface_id = "${azurerm_network_interface.ignition.id}"
  backend_address_pool_id = "${var.elb_backend_pool_id}"
  ip_configuration_name = "bootstrap" #must be the same as nic's ip configuration name.
}

resource "azurerm_virtual_machine" "bootstrap" {
  name                  = "${var.cluster_id}-bootstrap"
  location              = "${azurerm_resource_group.ignition.location}"
  resource_group_name   = "${azurerm_resource_group.ignition.name}"
  network_interface_ids = ["${azurerm_network_interface.ignition.id}"]
  vm_size               = "${var.vm_size}"

  delete_os_disk_on_termination = true
  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_image_reference {
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7.5"
    version   = "latest"
  }

  #  lifecycle {
  #   # Ignore changes in the AMI which force recreation of the resource. This
  #   # avoids accidental deletion of nodes whenever a new OS release comes out.
  #   ignore_changes = ["ami"]
  # }

  os_profile {
    computer_name  = "${var.cluster_id}-bootstrap-vm"
    admin_username = "king"
    admin_password = "P@ssword1234!"
    custom_data = "${data.ignition_config.redirect.rendered}"
  }

  os_profile_linux_config {
    disable_password_authentication = false

  }

  tags = "${merge(map(
    "Name", "${var.cluster_id}-bootstrap",
  ), var.tags)}"

}

// TODO :
// LB?
// Public IP
// NSG for SSH (tcp 22) + Journald gateway (tcp 19531)
// TESTS :)

