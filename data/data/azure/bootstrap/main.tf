
resource "azurerm_storage_account" "ignition" {
  name = "ignition_storage_${var.cluster_id}"
  resource_group_name = "${azurerm_resource_group.ignition.name}"
  location                 = "${var.region}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_resource_group" "ignition" {
  name = "${var.resource_group_name}"
  location ="${var.region}"
}

resource "azurerm_storage_container" "ignition" {
  resource_group_name = "${azurerm_resource_group.ignition.name}"
  location                 = "${var.region}"
  name                  = "ignition"
  resource_group_name   = "${azurerm_resource_group.ignition.name}"
  storage_account_name  = "${azurerm_storage_account.ignition.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "ignition" {
  name = "bootstrap.ign"
  content = "${var.ignition}"
  resource_group_name    = "${azurerm_resource_group.ignition.name}"
  storage_account_name   = "${azurerm_storage_account.ignition.name}"
  storage_container_name = "${azurerm_storage_container.ignition.name}"

  type = "page"
  size = 5120
}

data "ignition_config" "redirect" {
  replace {
    source = "${azurerm_storage_blob.ignition.url}/bootstrap.ign"
  }
}

data "azurerm_subscription" "primary" {}

resource "azurerm_role_definition" "bootstrap" {
  name        = "${var.cluster_id}-bootstrap-role"
  scope       = "${data.azurerm_subscription.primary.id}"
  description = "bootstrap role for openshift installer"

  permissions {
    actions     = ["*"]
    not_actions = []
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}", # /subscriptions/00000000-0000-0000-0000-000000000000
  ]
}

resource "azurerm_network_interface" "ignition" {
  name                = "${var.cluster_id}-bootstrap-nic"
  location            = "${azurerm_resource_group.ignition.location}"
  resource_group_name = "${azurerm_resource_group.ignition.name}"

  ip_configuration {
    name                          = "bootstrap"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "bootstrap" {
  name                  = "${var.cluster_id}-bootstrap"
  location              = "${azurerm_resource_group.ignition.location}"
  resource_group_name   = "${azurerm_resource_group.ignition.name}"
  network_interface_ids = ["${azurerm_network_interface.ignition.id}"]
  vm_size               = "${var.vm_size}"

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

