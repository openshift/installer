locals {
  bootstrap_nic_ip_configuration_name = "bootstrap-nic-ip"
  ssh_nat_rule_id                     = "${var.ssh_nat_rule_id}"
}

resource "random_string" "storage_suffix" {
  length  = 5
  upper   = false
  special = false

  keepers = {
    # Generate a new ID only when a new resource group is defined
    resource_group = "${var.resource_group_name}"
  }
}

resource "azurerm_storage_account" "ignition" {
  name                     = "ignitiondata${random_string.storage_suffix.result}"
  resource_group_name      = "${var.resource_group_name}"
  location                 = "${var.region}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

data "azurerm_storage_account_sas" "ignition" {
  connection_string = "${azurerm_storage_account.ignition.primary_connection_string}"
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "${timestamp()}"
  expiry = "${timeadd(timestamp(), "24h")}"

  permissions {
    read    = true
    list    = true
    create  = false
    add     = false
    delete  = false
    process = false
    write   = false
    update  = false
  }
}

resource "azurerm_storage_container" "ignition" {
  resource_group_name   = "${var.resource_group_name}"
  name                  = "ignition"
  storage_account_name  = "${azurerm_storage_account.ignition.name}"
  container_access_type = "private"
}

resource "local_file" "ignition_bootstrap" {
  content  = "${var.ignition}"
  filename = "${path.module}/ignition_bootstrap.ign"
}

resource "azurerm_storage_blob" "ignition" {
  name                   = "bootstrap.ign"
  source                 = "${local_file.ignition_bootstrap.filename}"
  resource_group_name    = "${var.resource_group_name}"
  storage_account_name   = "${azurerm_storage_account.ignition.name}"
  storage_container_name = "${azurerm_storage_container.ignition.name}"
  type                   = "block"
}

data "ignition_config" "redirect" {
  replace {
    source = "${azurerm_storage_blob.ignition.url}${data.azurerm_storage_account_sas.ignition.sas}"
  }
}

resource "azurerm_network_interface" "bootstrap" {
  name                = "${var.cluster_id}-bootstrap-nic"
  location            = "${var.region}"
  resource_group_name = "${var.resource_group_name}"

  ip_configuration {
    subnet_id                     = "${var.subnet_id}"
    name                          = "${local.bootstrap_nic_ip_configuration_name}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_nat_rule_association" "bootstrap_ssh" {
  network_interface_id  = "${azurerm_network_interface.bootstrap.id}"
  ip_configuration_name = "${local.bootstrap_nic_ip_configuration_name}"
  nat_rule_id           = "${local.ssh_nat_rule_id}"
}

resource "azurerm_network_interface_backend_address_pool_association" "public_lb_bootstrap" {
  network_interface_id    = "${azurerm_network_interface.bootstrap.id}"
  backend_address_pool_id = "${var.elb_backend_pool_id}"
  ip_configuration_name   = "${local.bootstrap_nic_ip_configuration_name}"
}

resource "azurerm_network_interface_backend_address_pool_association" "internal_lb_bootstrap" {
  network_interface_id    = "${azurerm_network_interface.bootstrap.id}"
  backend_address_pool_id = "${var.ilb_backend_pool_id}"
  ip_configuration_name   = "${local.bootstrap_nic_ip_configuration_name}"
}

data "azurerm_subscription" "current" {}

resource "azurerm_virtual_machine" "bootstrap" {
  name                  = "${var.cluster_id}-bootstrap"
  location              = "${var.region}"
  resource_group_name   = "${var.resource_group_name}"
  network_interface_ids = ["${azurerm_network_interface.bootstrap.id}"]
  vm_size               = "${var.vm_size}"

  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  identity {
    type         = "UserAssigned"
    identity_ids = ["${var.identity}"]
  }

  storage_os_disk {
    name              = "${var.cluster_id}-bootstrap_OSDisk" # os disk name needs to match cluster-api convention
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
    disk_size_gb      = 100
  }

  storage_image_reference {
    id = "${data.azurerm_subscription.current.id}${var.vm_image}"
  }

  os_profile {
    computer_name  = "${var.cluster_id}-bootstrap-vm"
    admin_username = "core"
    admin_password = "P@ssword1234!"
    custom_data    = "${data.ignition_config.redirect.rendered}"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  boot_diagnostics {
    enabled     = true
    storage_uri = "${var.boot_diag_blob_endpoint}"
  }
}
