locals {
  bootstrap_nic_ip_v4_configuration_name = "bootstrap-nic-ip-v4"
  bootstrap_nic_ip_v6_configuration_name = "bootstrap-nic-ip-v6"
  description                            = "Created By OpenShift Installer"

  # At this time min_tls_version is only supported in the Public Cloud and US Government Cloud.
  environments_with_min_tls_version = ["public", "usgovernment"]
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

resource "azurerm_storage_account" "cluster" {
  name                            = "cluster${var.random_storage_account_suffix}"
  resource_group_name             = var.resource_group_name
  location                        = var.azure_region
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  min_tls_version                 = contains(local.environments_with_min_tls_version, var.azure_environment) ? "TLS1_2" : null
  allow_nested_items_to_be_public = false
  tags                            = var.azure_extra_tags
}

# copy over the vhd to cluster resource group and create an image using that
resource "azurerm_storage_container" "vhd" {
  name                 = "vhd"
  storage_account_name = azurerm_storage_account.cluster.name
}

resource "azurerm_storage_blob" "rhcos_image" {
  name                   = "rhcos${var.random_storage_account_suffix}.vhd"
  storage_account_name   = azurerm_storage_account.cluster.name
  storage_container_name = azurerm_storage_container.vhd.name
  type                   = "Page"
  source_uri             = var.azure_image_url
  metadata               = tomap({ source_uri = var.azure_image_url })
}

# Creates Shared Image Gallery
# https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/shared_image_gallery
resource "azurerm_shared_image_gallery" "sig" {
  name                = "gallery_${replace(var.cluster_id, "-", "_")}"
  resource_group_name = var.resource_group_name
  location            = var.azure_region
  tags                = var.azure_extra_tags
}

# Creates image definition
# https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/shared_image
resource "azurerm_shared_image" "cluster" {
  name                = var.cluster_id
  gallery_name        = azurerm_shared_image_gallery.sig.name
  resource_group_name = var.resource_group_name
  location            = var.azure_region
  os_type             = "Linux"
  architecture        = var.azure_vm_architecture

  identifier {
    publisher = "RedHat"
    offer     = "rhcos"
    sku       = "basic"
  }

  tags = var.azure_extra_tags
}

resource "azurerm_shared_image" "clustergen2" {
  name                = "${var.cluster_id}-gen2"
  gallery_name        = azurerm_shared_image_gallery.sig.name
  resource_group_name = var.resource_group_name
  location            = var.azure_region
  os_type             = "Linux"
  hyper_v_generation  = "V2"
  architecture        = var.azure_vm_architecture

  identifier {
    publisher = "RedHat-gen2"
    offer     = "rhcos-gen2"
    sku       = "gen2"
  }

  tags = var.azure_extra_tags
}

resource "azurerm_shared_image_version" "cluster_image_version" {
  name                = var.azure_image_release
  gallery_name        = azurerm_shared_image.cluster.gallery_name
  image_name          = azurerm_shared_image.cluster.name
  resource_group_name = azurerm_shared_image.cluster.resource_group_name
  location            = azurerm_shared_image.cluster.location

  blob_uri           = azurerm_storage_blob.rhcos_image.url
  storage_account_id = azurerm_storage_account.cluster.id

  target_region {
    name                   = azurerm_shared_image.cluster.location
    regional_replica_count = 1
  }
}

resource "azurerm_shared_image_version" "clustergen2_image_version" {
  name                = var.azure_image_release
  gallery_name        = azurerm_shared_image.clustergen2.gallery_name
  image_name          = azurerm_shared_image.clustergen2.name
  resource_group_name = azurerm_shared_image.clustergen2.resource_group_name
  location            = azurerm_shared_image.clustergen2.location

  blob_uri           = azurerm_storage_blob.rhcos_image.url
  storage_account_id = azurerm_storage_account.cluster.id

  target_region {
    name                   = azurerm_shared_image.clustergen2.location
    regional_replica_count = 1
  }
}

data "azurerm_storage_account_sas" "ignition" {
  connection_string = azurerm_storage_account.cluster.primary_connection_string
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
  storage_account_name  = azurerm_storage_account.cluster.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "ignition" {
  name                   = "bootstrap.ign"
  source                 = var.ignition_bootstrap_file
  storage_account_name   = azurerm_storage_account.cluster.name
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
  tags                = var.azure_extra_tags
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
  tags                = var.azure_extra_tags
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

  tags = var.azure_extra_tags
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

  source_image_id = var.azure_hypervgeneration_version == "V2" ? azurerm_shared_image_version.clustergen2_image_version.id : azurerm_shared_image_version.cluster_image_version.id

  computer_name = "${var.cluster_id}-bootstrap-vm"
  custom_data   = base64encode(data.ignition_config.redirect.rendered)

  boot_diagnostics {
    storage_account_uri = null # null enables managed storage account for boot diagnostics
  }

  depends_on = [
    azurerm_network_interface_backend_address_pool_association.public_lb_bootstrap_v4,
    azurerm_network_interface_backend_address_pool_association.public_lb_bootstrap_v6,
    azurerm_network_interface_backend_address_pool_association.internal_lb_bootstrap_v4,
    azurerm_network_interface_backend_address_pool_association.internal_lb_bootstrap_v6
  ]

  tags = var.azure_extra_tags
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
