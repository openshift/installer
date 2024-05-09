locals {
  // The name of the masters' ipconfiguration is hardcoded to "pipconfig". It needs to match cluster-api
  // https://github.com/openshift/cluster-api-provider-azure/blob/master/pkg/cloud/azure/services/networkinterfaces/networkinterfaces.go#L131
  ip_v4_configuration_name = "pipConfig"
  // TODO: Azure machine provider probably needs to look for pipConfig-v6 as well (or a different name like pipConfig-secondary)
  ip_v6_configuration_name = "pipConfig-v6"
}

resource "azurerm_network_interface" "master" {
  count = var.instance_count

  name                          = "${var.cluster_id}-master-${count.index}-nic"
  location                      = var.region
  resource_group_name           = var.resource_group_name
  enable_accelerated_networking = var.vm_networking_type

  dynamic "ip_configuration" {
    for_each = [for ip in [
      {
        // LIMITATION: azure does not allow an ipv6 address to be primary today
        primary : var.use_ipv4,
        name : local.ip_v4_configuration_name,
        ip_address_version : "IPv4",
        include : var.use_ipv4 || var.use_ipv6
      },
      {
        primary : ! var.use_ipv4,
        name : local.ip_v6_configuration_name,
        ip_address_version : "IPv6",
        include : var.use_ipv6
      },
      ] : {
      primary : ip.primary
      name : ip.name
      ip_address_version : ip.ip_address_version
      include : ip.include
      } if ip.include
    ]
    content {
      primary                       = ip_configuration.value.primary
      name                          = ip_configuration.value.name
      subnet_id                     = var.subnet_id
      private_ip_address_version    = ip_configuration.value.ip_address_version
      private_ip_address_allocation = "Dynamic"
    }
  }

  tags = var.azure_extra_tags
}

resource "azurerm_network_interface_backend_address_pool_association" "master_v4" {
  // This is required because terraform cannot calculate counts during plan phase completely and therefore the `vnet/public-lb.tf`
  // conditional need to be recreated. See https://github.com/hashicorp/terraform/issues/12570
  count = (! var.private || var.outbound_type != "UserDefinedRouting") ? var.instance_count : 0

  network_interface_id    = element(azurerm_network_interface.master.*.id, count.index)
  backend_address_pool_id = var.elb_backend_pool_v4_id
  ip_configuration_name   = local.ip_v4_configuration_name
}

resource "azurerm_network_interface_backend_address_pool_association" "master_v6" {
  // This is required because terraform cannot calculate counts during plan phase completely and therefore the `vnet/public-lb.tf`
  // conditional need to be recreated. See https://github.com/hashicorp/terraform/issues/12570
  count = var.use_ipv6 && (! var.private || var.outbound_type != "UserDefinedRouting") ? var.instance_count : 0

  network_interface_id    = element(azurerm_network_interface.master.*.id, count.index)
  backend_address_pool_id = var.elb_backend_pool_v6_id
  ip_configuration_name   = local.ip_v6_configuration_name
}

resource "azurerm_network_interface_backend_address_pool_association" "master_internal_v4" {
  count = var.use_ipv4 ? var.instance_count : 0

  network_interface_id    = element(azurerm_network_interface.master.*.id, count.index)
  backend_address_pool_id = var.ilb_backend_pool_v4_id
  ip_configuration_name   = local.ip_v4_configuration_name
}

resource "azurerm_network_interface_backend_address_pool_association" "master_internal_v6" {
  count = var.use_ipv6 ? var.instance_count : 0

  network_interface_id    = element(azurerm_network_interface.master.*.id, count.index)
  backend_address_pool_id = var.ilb_backend_pool_v6_id
  ip_configuration_name   = local.ip_v6_configuration_name
}

resource "azurerm_linux_virtual_machine" "master" {
  count = var.instance_count

  name                  = "${var.cluster_id}-master-${count.index}"
  location              = var.region
  zone                  = var.availability_zones[count.index] != "" ? var.availability_zones[count.index] : null
  resource_group_name   = var.resource_group_name
  network_interface_ids = [element(azurerm_network_interface.master.*.id, count.index)]
  size                  = var.vm_size
  admin_username        = "core"
  # The password is normally applied by WALA (the Azure agent), but this
  # isn't installed in RHCOS. As a result, this password is never set. It is
  # included here because it is required by the Azure ARM API.
  admin_password                  = "NotActuallyApplied!"
  disable_password_authentication = false
  encryption_at_host_enabled      = var.encryption_at_host_enabled
  secure_boot_enabled             = var.secure_boot == "Enabled"
  vtpm_enabled                    = var.virtualized_trusted_platform_module == "Enabled"

  additional_capabilities {
    ultra_ssd_enabled = var.ultra_ssd_enabled
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [var.identity]
  }

  os_disk {
    name                   = "${var.cluster_id}-master-${count.index}_OSDisk" # os disk name needs to match cluster-api convention
    caching                = "ReadWrite"
    storage_account_type   = var.os_volume_type
    disk_size_gb           = var.os_volume_size
    disk_encryption_set_id = var.disk_encryption_set_id

    security_encryption_type         = var.security_encryption_type
    secure_vm_disk_encryption_set_id = var.secure_vm_disk_encryption_set_id
  }

  # Either source_image_id or source_image_reference must be defined
  source_image_id = ! var.use_marketplace_image ? var.vm_image : null

  dynamic "source_image_reference" {
    for_each = var.use_marketplace_image ? [1] : []

    content {
      publisher = var.vm_image_publisher
      offer     = var.vm_image_offer
      sku       = var.vm_image_sku
      version   = var.vm_image_version
    }
  }

  dynamic "plan" {
    for_each = var.use_marketplace_image && var.vm_image_has_plan ? [1] : []

    content {
      publisher = var.vm_image_publisher
      product   = var.vm_image_offer
      name      = var.vm_image_sku
    }
  }

  //we don't provide a ssh key, because it is set with ignition.
  //it is required to provide at least 1 auth method to deploy a linux vm
  computer_name = "${var.cluster_id}-master-${count.index}"
  custom_data   = base64encode(var.ignition)

  boot_diagnostics {
    storage_account_uri = null # null enables managed storage account for boot diagnostics
  }

  tags = var.azure_extra_tags
}
