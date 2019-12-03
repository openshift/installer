resource "ovirt_image_transfer" "releaseimage_transfer" {
  alias             = var.openstack_base_image_name
  source_url        = var.openstack_base_image_local_file_path
  storage_domain_id = var.ovirt_storage_domain_id
  sparse            = true
}

resource "ovirt_vm" "tmp_import_vm" {
  name = "tmpvm-for-${ovirt_image_transfer.releaseimage_transfer.alias}"
  cluster_id = var.ovirt_cluster_id
  block_device {
    disk_id   = ovirt_image_transfer.releaseimage_transfer.disk_id
    interface = "virtio_scsi"
  }
}
resource "ovirt_vnic" "vm_nic1" {
  vm_id           = ovirt_vm.tmp_import_vm.id
  name            = "nic1"
  // default profile id
  vnic_profile_id = "0000000a-000a-000a-000a-000000000398"
}

resource "ovirt_template" "releaseimage_template" {
  name = "template-for-${ovirt_image_transfer.releaseimage_transfer.alias}"
  cluster_id = ovirt_vm.tmp_import_vm.cluster_id
  memory      = "16384"
  cores       = "4"
  // create from vm
  vm_id = ovirt_vm.tmp_import_vm.id
}

resource "ovirt_vm" "bootstrap" {
  name        = "${var.cluster_id}-bootstrap"
  memory      = "8192"
  cores       = "4"
  cluster_id  = var.ovirt_cluster_id
  template_id = var.ovirt_template_id

  initialization {
    host_name     = "bootstrap-${var.cluster_id}"
    custom_script = var.ignition_bootstrap
  }
}
