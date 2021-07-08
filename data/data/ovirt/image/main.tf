locals {
  image_name = "${var.cluster_id}-rhcos"
}

provider "ovirt" {
  url       = var.ovirt_url
  username  = var.ovirt_username
  password  = var.ovirt_password
  cafile    = var.ovirt_cafile
  ca_bundle = var.ovirt_ca_bundle
  insecure  = var.ovirt_insecure
}

// upload the disk if we don't have an existing template
resource "ovirt_image_transfer" "releaseimage" {
  count = length(var.ovirt_base_image_name) == 0 ? 1 : 0

  alias             = local.image_name
  source_url        = var.ovirt_base_image_local_file_path
  storage_domain_id = var.ovirt_storage_domain_id
  sparse            = true
  timeouts {
    create = "20m"
  }
}

resource "ovirt_vm" "tmp_import_vm" {
  // create the vm for import only when we don't have an existing template
  count = length(var.ovirt_base_image_name) == 0 ? 1 : 0

  name       = "tmpvm-for-${ovirt_image_transfer.releaseimage.0.alias}"
  cluster_id = var.ovirt_cluster_id
  auto_start = false
  block_device {
    disk_id   = ovirt_image_transfer.releaseimage.0.disk_id
    interface = "virtio_scsi"
  }
  os {
    type = "rhcos_x64"
  }
  nics {
    name            = "nic1"
    vnic_profile_id = var.ovirt_vnic_profile_id
  }
  timeouts {
    create = "20m"
  }
  depends_on = [ovirt_image_transfer.releaseimage]
}
