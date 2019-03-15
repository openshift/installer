
resource "azurerm_network_interface" "master" {
  name                = "${var.cluster_id}-master-nic"
  location            = "${var.region}"
  resource_group_name = "${var.resource_group_name}"

  ip_configuration {
    subnet_id                               = "${var.subnet_id}"
    name                                    = "master"
    private_ip_address_allocation           = "Dynamic"
  }
}

resource "azurerm_network_interface_backend_address_pool_association" "master" {
  network_interface_id = "${azurerm_network_interface.master.id}"
  backend_address_pool_id = "${var.elb_backend_pool_id}"
  ip_configuration_name = "master" #must be the same as nic's ip configuration name.
}

resource "azurerm_network_interface_nat_rule_association" "master" {
  network_interface_id  = "${azurerm_network_interface.master.id}"
  ip_configuration_name = "master"
  nat_rule_id           = "${var.elb_master_ssh_natrule_id}"
}

resource "azurerm_virtual_machine" "master" {
  name                  = "${var.cluster_id}-master"
  location              = "${var.region}"
  resource_group_name   = "${var.resource_group_name}"
  network_interface_ids = ["${azurerm_network_interface.master.id}"]
  vm_size               = "${var.vm_size}"

  delete_os_disk_on_termination = true
  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
    disk_size_gb      = 100
  }

  storage_image_reference {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Beta"
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
    custom_data = "${var.ignition_master}"
  }

  os_profile_linux_config {
    disable_password_authentication = false
    boot_diagnostics {
        enabled = true
        storage_uri = "${var.boot_diag_blob_endpoint}"
    }
  }

  tags = "${merge(map(
    "Name", "${var.cluster_id}-bootstrap",
  ), var.tags)}"

}