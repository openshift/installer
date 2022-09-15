locals {
  bootstrap_nic_ip_v4_configuration_name = "bootstrap-nic-ip-v4"
  bootstrap_nic_ip_v6_configuration_name = "bootstrap-nic-ip-v6"
  description                            = "Created By OpenShift Installer"
  tags = merge(
    {
      "kubernetes.io_cluster.${var.cluster_id}" = "owned"
    },
    var.azure_extra_tags,
  )
}

provider "azurerm" {
  features {}
  subscription_id             = var.azure_subscription_id
  client_id                   = var.azure_client_id
  client_secret               = var.azure_client_secret
  client_certificate_password = var.azure_certificate_password
  client_certificate_path     = var.azure_certificate_path
  tenant_id                   = var.azure_tenant_id
  environment                 = var.azure_environment
}

data "azurerm_storage_account" "storage_account" {
  name                = var.storage_account_name
  resource_group_name = var.resource_group_name
}

data "azurerm_storage_account_sas" "ignition" {
  connection_string = data.azurerm_storage_account.storage_account.primary_connection_string
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
    filter  = false
    tag     = false
  }
}

resource "azurerm_storage_container" "ignition" {
  name                  = "ignition"
  storage_account_name  = var.storage_account_name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "ignition" {
  name                   = "bootstrap.ign"
  source                 = var.ignition_bootstrap_file
  storage_account_name   = var.storage_account_name
  storage_container_name = azurerm_storage_container.ignition.name
  type                   = "Block"
}

data "ignition_config" "redirect" {
  replace {
    source = "${azurerm_storage_blob.ignition.url}${data.azurerm_storage_account_sas.ignition.sas}"
  }
}

resource "azurerm_public_ip" "bootstrap_public_ip_v4" {
  count = var.azure_private || ! var.use_ipv4 ? 0 : 1

  sku                 = "Standard"
  location            = var.azure_region
  name                = "${var.cluster_id}-bootstrap-pip-v4"
  resource_group_name = var.resource_group_name
  allocation_method   = "Static"
}

data "azurerm_public_ip" "bootstrap_public_ip_v4" {
  count = var.azure_private ? 0 : 1

  name                = azurerm_public_ip.bootstrap_public_ip_v4[0].name
  resource_group_name = var.resource_group_name
}

resource "azurerm_public_ip" "bootstrap_public_ip_v6" {
  count = var.azure_private || ! var.use_ipv6 ? 0 : 1

  sku                 = "Standard"
  location            = var.azure_region
  name                = "${var.cluster_id}-bootstrap-pip-v6"
  resource_group_name = var.resource_group_name
  allocation_method   = "Static"
  ip_version          = "IPv6"
}

data "azurerm_public_ip" "bootstrap_public_ip_v6" {
  count = var.azure_private || ! var.use_ipv6 ? 0 : 1

  name                = azurerm_public_ip.bootstrap_public_ip_v6[0].name
  resource_group_name = var.resource_group_name
}

resource "azurerm_network_interface" "bootstrap" {
  name                = "${var.cluster_id}-bootstrap-nic"
  location            = var.azure_region
  resource_group_name = var.resource_group_name

  dynamic "ip_configuration" {
    for_each = [for ip in [
      {
        // LIMITATION: azure does not allow an ipv6 address to be primary today
        primary : var.use_ipv4,
        name : local.bootstrap_nic_ip_v4_configuration_name,
        ip_address_version : "IPv4",
        public_ip_id : var.azure_private ? null : azurerm_public_ip.bootstrap_public_ip_v4[0].id,
        include : var.use_ipv4 || var.use_ipv6,
      },
      {
        primary : ! var.use_ipv4,
        name : local.bootstrap_nic_ip_v6_configuration_name,
        ip_address_version : "IPv6",
        public_ip_id : var.azure_private || ! var.use_ipv6 ? null : azurerm_public_ip.bootstrap_public_ip_v6[0].id,
        include : var.use_ipv6,
      },
      ] : {
      primary : ip.primary
      name : ip.name
      ip_address_version : ip.ip_address_version
      public_ip_id : ip.public_ip_id
      include : ip.include
      } if ip.include
    ]
    content {
      primary                       = ip_configuration.value.primary
      name                          = ip_configuration.value.name
      subnet_id                     = var.master_subnet_id
      private_ip_address_version    = ip_configuration.value.ip_address_version
      private_ip_address_allocation = "Dynamic"
      public_ip_address_id          = ip_configuration.value.public_ip_id
    }
  }
}

resource "azurerm_network_interface_backend_address_pool_association" "public_lb_bootstrap_v4" {
  // This is required because terraform cannot calculate counts during plan phase completely and therefore the `vnet/public-lb.tf`
  // conditional need to be recreated. See https://github.com/hashicorp/terraform/issues/12570
  count = (! var.azure_private || ! var.azure_outbound_user_defined_routing) ? 1 : 0

  network_interface_id    = azurerm_network_interface.bootstrap.id
  backend_address_pool_id = var.elb_backend_pool_v4_id
  ip_configuration_name   = local.bootstrap_nic_ip_v4_configuration_name
}

resource "azurerm_network_interface_backend_address_pool_association" "public_lb_bootstrap_v6" {
  // This is required because terraform cannot calculate counts during plan phase completely and therefore the `vnet/public-lb.tf`
  // conditional need to be recreated. See https://github.com/hashicorp/terraform/issues/12570
  count = var.use_ipv6 && (! var.azure_private || ! var.azure_outbound_user_defined_routing) ? 1 : 0

  network_interface_id    = azurerm_network_interface.bootstrap.id
  backend_address_pool_id = var.elb_backend_pool_v6_id
  ip_configuration_name   = local.bootstrap_nic_ip_v6_configuration_name
}

resource "azurerm_network_interface_backend_address_pool_association" "internal_lb_bootstrap_v4" {
  count = var.use_ipv4 ? 1 : 0

  network_interface_id    = azurerm_network_interface.bootstrap.id
  backend_address_pool_id = var.ilb_backend_pool_v4_id
  ip_configuration_name   = local.bootstrap_nic_ip_v4_configuration_name
}

resource "azurerm_network_interface_backend_address_pool_association" "internal_lb_bootstrap_v6" {
  count = var.use_ipv6 ? 1 : 0

  network_interface_id    = azurerm_network_interface.bootstrap.id
  backend_address_pool_id = var.ilb_backend_pool_v6_id
  ip_configuration_name   = local.bootstrap_nic_ip_v6_configuration_name
}

resource "azurerm_linux_virtual_machine" "bootstrap" {
  name                  = "${var.cluster_id}-bootstrap"
  location              = var.azure_region
  resource_group_name   = var.resource_group_name
  network_interface_ids = [azurerm_network_interface.bootstrap.id]
  size                  = var.azure_master_vm_type
  admin_username        = "core"
  # The password is normally applied by WALA (the Azure agent), but this
  # isn't installed in RHCOS. As a result, this password is never set. It is
  # included here because it is required by the Azure ARM API.
  admin_password                  = "NotActuallyApplied!"
  disable_password_authentication = false
  encryption_at_host_enabled      = var.azure_master_encryption_at_host_enabled

  identity {
    type         = "UserAssigned"
    identity_ids = [var.identity]
  }

  os_disk {
    name                   = "${var.cluster_id}-bootstrap_OSDisk" # os disk name needs to match cluster-api convention
    caching                = "ReadWrite"
    storage_account_type   = var.azure_master_root_volume_type
    disk_size_gb           = 100
    disk_encryption_set_id = var.azure_master_disk_encryption_set_id
  }

  source_image_id = var.vm_image

  computer_name = "${var.cluster_id}-bootstrap-vm"
  custom_data   = base64encode(data.ignition_config.redirect.rendered)

  boot_diagnostics {
    storage_account_uri = data.azurerm_storage_account.storage_account.primary_blob_endpoint
  }

  depends_on = [
    azurerm_network_interface_backend_address_pool_association.public_lb_bootstrap_v4,
    azurerm_network_interface_backend_address_pool_association.public_lb_bootstrap_v6,
    azurerm_network_interface_backend_address_pool_association.internal_lb_bootstrap_v4,
    azurerm_network_interface_backend_address_pool_association.internal_lb_bootstrap_v6
  ]
}

resource "azurerm_network_security_rule" "bootstrap_ssh_in" {
  count = var.azure_private ? 0 : 1

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
