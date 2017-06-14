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
  account_type        = "${var.storage_account_type}"

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

resource "azurerm_availability_set" "tectonic_masters" {
  name                = "${var.cluster_name}-masters"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"
}

resource "azurerm_network_interface" "tectonic_master" {
  count               = "${var.master_count}"
  name                = "${var.cluster_name}-master${count.index}"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"

  ip_configuration {
    private_ip_address_allocation           = "dynamic"
    name                                    = "${var.cluster_name}-MasterIPConfiguration"
    subnet_id                               = "${var.subnet}"
    load_balancer_backend_address_pools_ids = ["${azurerm_lb_backend_address_pool.api-lb.id}"]
  }
}

resource "azurerm_virtual_machine" "tectonic_master" {
  count                 = "${var.master_count}"
  name                  = "${var.cluster_name}-master${count.index}"
  location              = "${var.location}"
  resource_group_name   = "${var.resource_group_name}"
  network_interface_ids = ["${element(azurerm_network_interface.tectonic_master.*.id, count.index)}"]
  vm_size               = "${var.vm_size}"
  availability_set_id   = "${azurerm_availability_set.tectonic_masters.id}"

  storage_image_reference {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Stable"
    version   = "latest"
  }

  storage_os_disk {
    name          = "master-osdisk"
    caching       = "ReadWrite"
    create_option = "FromImage"
    os_type       = "linux"
    vhd_uri       = "${azurerm_storage_account.tectonic_master.primary_blob_endpoint}${azurerm_storage_container.tectonic_master.name}/${var.cluster_name}-master${count.index}.vhd"
  }

  os_profile {
    computer_name  = "${var.cluster_name}-master${count.index}"
    admin_username = "core"
    admin_password = ""

    custom_data = "${base64encode("${data.ignition_config.master.rendered}")}"
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
