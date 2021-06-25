provider "ovirt" {
  url       = var.ovirt_url
  username  = var.ovirt_username
  password  = var.ovirt_password
  cafile    = var.ovirt_cafile
  ca_bundle = var.ovirt_ca_bundle
  insecure  = var.ovirt_insecure
}

resource "ovirt_vm" "bootstrap" {
  name        = "${var.cluster_id}-bootstrap"
  memory      = "8192"
  cores       = "4"
  cluster_id  = var.ovirt_cluster_id
  template_id = var.release_image_template_id

  initialization {
    custom_script = var.ignition_bootstrap
  }
}

resource "ovirt_tag" "cluster_bootstrap_tag" {
  name   = "${var.cluster_id}-bootstrap"
  vm_ids = concat([ovirt_vm.bootstrap.id], [var.tmp_import_vm_id])
}
