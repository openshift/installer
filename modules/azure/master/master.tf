# TODO:
# Create global network tf file
# Add azurerm_route_table
# Add azurerm_network_security_group
# Add azurerm_availability_set

resource "azurerm_subnet" "tectonic_subnet" {
  name                 = "${var.cluster_name}_master_subnet"
  resource_group_name  = "${var.resource_group_name}"
  virtual_network_name = "${var.virtual_network}"
  address_prefix       = "10.0.2.0/24"
}


# create network interface
resource "azurerm_network_interface" "tectonic_nic" {
  name                = "tectonic_nic"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"

  ip_configuration {
    name                                    = "tectonic_configuration"
    subnet_id                               = "${azurerm_subnet.tectonic_subnet.id}"
    private_ip_address_allocation           = "dynamic"
    load_balancer_backend_address_pools_ids = ["${azurerm_lb_backend_address_pool.k8-lb.id}"]
  }
}

# Generate unique storage name
resource "random_id" "tectonic_storage_name" {
  byte_length = 4
}

resource "azurerm_storage_account" "tectonic_storage" {
  name                = "${random_id.tectonic_storage_name.hex}"
  resource_group_name = "${var.resource_group_name}"
  location            = "${var.location}"
  account_type        = "Premium_LRS"

  tags {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "tectonic_storage_container" {
  name                  = "vhd"
  resource_group_name   = "${var.resource_group_name}"
  storage_account_name  = "${azurerm_storage_account.tectonic_storage.name}"
  container_access_type = "private"
  depends_on            = ["azurerm_storage_account.tectonic_storage"]
}

resource "azurerm_virtual_machine" "tectonic_master_vm" {
  name                  = "tectonic_master_vm"
  location              = "${var.location}"
  resource_group_name   = "${var.resource_group_name}"
  network_interface_ids = ["${azurerm_network_interface.tectonic_nic.id}"]
  vm_size               = "${var.vm_size}"

  storage_image_reference {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Stable"
    version   = "latest"
  }

  storage_os_disk {
    name          = "master-osdisk"
    vhd_uri       = "${azurerm_storage_account.tectonic_storage.primary_blob_endpoint}${azurerm_storage_container.tectonic_storage_container.name}/myosdisk.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "tectonic-master-${count.index}"
    admin_username = "core"
    admin_password = ""
    custom_data    = "${base64encode("${ignition_config.master.rendered}")}"
  }

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/core/.ssh/authorized_keys"
      key_data = "${file(var.ssh_key)}"
    }
  }

  tags {
    environment = "staging"
  }
}
