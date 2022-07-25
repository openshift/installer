locals {
  image_name = "${var.cluster_id}-rhcos"
}

provider "ovirt" {
  url           = var.ovirt_url
  username      = var.ovirt_username
  password      = var.ovirt_password
  tls_ca_files  = var.ovirt_cafile == "" ? [] : [var.ovirt_cafile]
  tls_ca_bundle = var.ovirt_ca_bundle
  tls_insecure  = var.ovirt_insecure
}

// We are creating a new disk from an image here. The process is a single step because a corrupted upload can cause the
// disk to be deleted and may need to be recreated.
resource "ovirt_disk_from_image" "releaseimage" {
  count = length(var.ovirt_base_image_name) == 0 ? 1 : 0

  // source_file provides the source file name to read from.
  source_file = var.ovirt_base_image_local_file_path

  alias             = local.image_name
  storage_domain_id = var.ovirt_storage_domain_id
  sparse            = true
  format            = "cow"
}

data "ovirt_blank_template" "blank" {}

resource "ovirt_vm" "tmp_import_vm" {
  // create the vm for import only when we don't have an existing template
  count = length(var.ovirt_base_image_name) == 0 ? 1 : 0

  name        = "tmpvm-for-${ovirt_disk_from_image.releaseimage.0.alias}"
  cluster_id  = var.ovirt_cluster_id
  template_id = data.ovirt_blank_template.blank.id
  os_type     = "rhcos_x64"
}

resource "ovirt_disk_attachment" "tmp_import_vm" {
  count          = length(var.ovirt_base_image_name) == 0 ? 1 : 0
  vm_id          = ovirt_vm.tmp_import_vm.0.id
  disk_id        = ovirt_disk_from_image.releaseimage.0.id
  disk_interface = "virtio_scsi"
  bootable       = true
  active         = true
}

resource "ovirt_nic" "tmp_import_vm" {
  count           = length(var.ovirt_base_image_name) == 0 ? 1 : 0
  vm_id           = ovirt_vm.tmp_import_vm.0.id
  vnic_profile_id = var.ovirt_vnic_profile_id
  name            = "tmpnic-for-${ovirt_disk_from_image.releaseimage.0.alias}"
}