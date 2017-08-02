# TODO
# Add to availabilityset

# Generate unique storage name
resource "random_id" "tectonic_storage_name" {
  byte_length = 4
}

resource "azurerm_storage_account" "tectonic_worker" {
  name                = "worker-${var.cluster_name}-${random_id.tectonic_storage_name.hex}"
  resource_group_name = "${var.resource_group_name}"
  location            = "${var.location}"
  account_type        = "${var.storage_type}"

  tags {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "tectonic_worker" {
  name                  = "vhd"
  resource_group_name   = "${var.resource_group_name}"
  storage_account_name  = "${azurerm_storage_account.tectonic_worker.name}"
  container_access_type = "private"
}

# resource "azurerm_lb_backend_address_pool" "workers" {
#   name                = "workers-lb-pool"
#   resource_group_name = "${var.resource_group_name}"
#   loadbalancer_id     = "${azurerm_lb.tectonic_lb.id}"
# }

resource "azurerm_virtual_machine_scale_set" "tectonic_workers" {
  name                = "${var.cluster_name}-workers"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "${var.vm_size}"
    tier     = "${element(split("_", var.vm_size),0)}"
    capacity = "${var.worker_count}"
  }

  network_profile {
    name    = "${var.cluster_name}-WorkerNetworkProfile"
    primary = true

    ip_configuration {
      name      = "${var.cluster_name}-WorkerIPConfiguration"
      subnet_id = "${var.subnet}"

      # load_balancer_backend_address_pool_ids = ["azurerm_lb_backend_address_pool.workers.id"]
    }
  }

  storage_profile_image_reference {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Stable"
    version   = "${var.versions["container_linux"]}"
  }

  storage_profile_os_disk {
    name           = "worker-osdisk"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    os_type        = "linux"
    vhd_containers = ["${azurerm_storage_account.tectonic_worker.primary_blob_endpoint}${azurerm_storage_container.tectonic_worker.name}"]
  }

  os_profile {
    computer_name_prefix = "tectonic-worker-"
    admin_username       = "core"
    admin_password       = ""
    custom_data          = "${base64encode("${data.ignition_config.worker.rendered}")}"
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
