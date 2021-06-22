locals {
  // The name of the masters' ipconfiguration is hardcoded to "pipconfig". It needs to match cluster-api
  // https://github.com/openshift/cluster-api-provider-azure/blob/master/pkg/cloud/azure/services/networkinterfaces/networkinterfaces.go#L131
  ip_v4_configuration_name = "pipConfig"
}

resource "azurestack_network_interface" "master" {
  count = var.instance_count

  name                = "${var.cluster_id}-master${count.index}-nic"
  location            = var.region
  resource_group_name = var.resource_group_name

  ip_configuration {
    primary                       = true
    name                          = local.ip_v4_configuration_name
    subnet_id                     = var.subnet_id
    private_ip_address_allocation = "Dynamic"
    load_balancer_backend_address_pools_ids = concat(
      [var.ilb_backend_pool_v4_id],
      ! var.private ? [var.elb_backend_pool_v4_id] : null
    )
  }
}

resource "azurestack_virtual_machine" "master" {
  count = var.instance_count

  name                  = "${var.cluster_id}-master-${count.index}"
  location              = var.region
  resource_group_name   = var.resource_group_name
  network_interface_ids = [element(azurestack_network_interface.master.*.id, count.index)]
  vm_size               = var.vm_size
  availability_set_id   = var.availability_set_id

  os_profile {
    computer_name  = "${var.cluster_id}-master-${count.index}"
    admin_username = "core"
    # The password is normally applied by WALA (the Azure agent), but this
    # isn't installed in RHCOS. As a result, this password is never set. It is
    # included here because it is required by the Azure ARM API.
    admin_password = "NotActuallyApplied!"
    custom_data    = base64encode(var.ignition)
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  storage_image_reference {
    id = var.vm_image_uri
  }

  storage_os_disk {
    name          = "${var.cluster_id}-master-${count.index}_OSDisk" # os disk name needs to match cluster-api convention
    create_option = "FromImage"
    // Only Disk CachingType 'None' is supported for disk with size greater than 1023 GB.
    // caching       = "ReadOnly"
    // TODO: Remove this 1023 GB limit once we are creating own own image
    disk_size_gb         = min(var.os_volume_size, 1023)
    managed_disk_type = "Standard_LRS"
  }

  boot_diagnostics {
    enabled     = true
    storage_uri = var.storage_account.primary_blob_endpoint
  }
}
