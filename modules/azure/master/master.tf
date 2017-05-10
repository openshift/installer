# TODO:
# Create global network tf file
# Add azurerm_route_table
# Add azurerm_network_security_group
# Add azurerm_availability_set

# Generate unique storage name
resource "random_id" "tectonic_master_storage_name" {
  byte_length = 4
}

resource "azurerm_storage_account" "tectonic_master" {
  name                = "${random_id.tectonic_master_storage_name.hex}"
  resource_group_name = "${var.resource_group_name}"
  location            = "${var.location}"
  account_type        = "Premium_LRS"

  tags {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "tectonic_master" {
  name                  = "${var.cluster_name}-vhd-master"
  resource_group_name   = "${var.resource_group_name}"
  storage_account_name  = "${azurerm_storage_account.tectonic_master.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine_scale_set" "tectonic_masters" {
  name                = "${var.cluster_name}-masters"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "${var.vm_size}"
    tier     = "${element(split("_", var.vm_size),0)}"
    capacity = "${var.master_count}"
  }

  network_profile {
    name    = "${var.cluster_name}-MasterNetworkProfile"
    primary = true

    ip_configuration {
      name                                   = "${var.cluster_name}-MasterIPConfiguration"
      subnet_id                              = "${var.subnet}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.api-lb.id}"]
    }
  }

  storage_profile_image_reference {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Stable"
    version   = "latest"
  }

  storage_profile_os_disk {
    name           = "master-osdisk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    os_type        = "linux"
    vhd_containers = ["${azurerm_storage_account.tectonic_master.primary_blob_endpoint}${azurerm_storage_container.tectonic_master.name}"]
  }

  os_profile {
    computer_name_prefix = "tectonic-master-"
    admin_username       = "core"
    admin_password       = ""

    custom_data = "${base64encode("${data.ignition_config.master.*.rendered[count.index]}")}"
  }

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/core/.ssh/authorized_keys"
      key_data = "${file(var.public_ssh_key)}"
    }
  }

  tags {
    environment = "staging"
  }
}
