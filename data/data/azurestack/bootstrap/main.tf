locals {
  bootstrap_nic_ip_v4_configuration_name = "bootstrap-nic-ip-v4"
  description                            = "Created By OpenShift Installer"
}
/*
// TODO: Implement this data type
data "azurestack_storage_account_sas" "ignition" {
  connection_string = var.storage_account.primary_connection_string
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

  start  = timestamp()
  expiry = timeadd(timestamp(), "24h")

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
*/
resource "azurestack_storage_container" "ignition" {
  name                 = "ignition"
  resource_group_name  = var.resource_group_name
  storage_account_name = var.storage_account.name
  // TODO: make this "private" once we have the storage account sas token
  container_access_type = "blob"
}

resource "local_file" "ignition_bootstrap" {
  content  = var.ignition
  filename = "${path.module}/ignition_bootstrap.ign"
}

resource "azurestack_storage_blob" "ignition" {
  name                   = "bootstrap.ign"
  source                 = local_file.ignition_bootstrap.filename
  resource_group_name    = var.resource_group_name
  storage_account_name   = var.storage_account.name
  storage_container_name = azurestack_storage_container.ignition.name
  type                   = "block"
}

data "ignition_config" "redirect" {
  replace {
    // TODO: use sas url once we have the storage account sas token
    //source = "${azurestack_storage_blob.ignition.url}${data.azurestack_storage_account_sas.ignition.sas}"
    source = azurestack_storage_blob.ignition.url
  }
}

resource "azurestack_public_ip" "bootstrap_public_ip_v4" {
  count = var.private ? 0 : 1

  location                     = var.region
  name                         = "${var.cluster_id}-bootstrap-pip-v4"
  resource_group_name          = var.resource_group_name
  public_ip_address_allocation = "Static"
}

data "azurestack_public_ip" "bootstrap_public_ip_v4" {
  count = var.private ? 0 : 1

  name                = azurestack_public_ip.bootstrap_public_ip_v4[0].name
  resource_group_name = var.resource_group_name
}

resource "azurestack_network_interface" "bootstrap" {
  name                = "${var.cluster_id}-bootstrap-nic"
  location            = var.region
  resource_group_name = var.resource_group_name

  ip_configuration {
    primary                       = true
    name                          = local.bootstrap_nic_ip_v4_configuration_name
    subnet_id                     = var.subnet_id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = var.private ? null : azurestack_public_ip.bootstrap_public_ip_v4[0].id
    load_balancer_backend_address_pools_ids = concat(
      [var.ilb_backend_pool_v4_id],
      ! var.private ? [var.elb_backend_pool_v4_id] : null
    )
  }
}

resource "azurestack_virtual_machine" "bootstrap" {
  name                  = "${var.cluster_id}-bootstrap"
  location              = var.region
  resource_group_name   = var.resource_group_name
  network_interface_ids = [azurestack_network_interface.bootstrap.id]
  vm_size               = var.vm_size
  availability_set_id   = var.availability_set_id

  os_profile {
    computer_name  = "${var.cluster_id}-bootstrap-vm"
    admin_username = "core"
    # The password is normally applied by WALA (the Azure agent), but this
    # isn't installed in RHCOS. As a result, this password is never set. It is
    # included here because it is required by the Azure ARM API.
    admin_password = "NotActuallyApplied!"
    custom_data    = base64encode(data.ignition_config.redirect.rendered)
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  storage_image_reference {
    id = var.vm_image_uri
  }

  storage_os_disk {
    name          = "${var.cluster_id}-bootstrap_OSDisk" # os disk name needs to match cluster-api convention
    create_option = "FromImage"
    disk_size_gb  = 100
    managed_disk_type = "Standard_LRS"
  }

  boot_diagnostics {
    enabled     = true
    storage_uri = var.storage_account.primary_blob_endpoint
  }
}

resource "azurestack_network_security_rule" "bootstrap_ssh_in" {
  name                        = "bootstrap_ssh_in"
  priority                    = 103
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = var.resource_group_name
  network_security_group_name = var.nsg_name
  description                 = local.description
}
