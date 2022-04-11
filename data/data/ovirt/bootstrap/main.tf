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
  cpu_cores   = "4"
  cluster_id  = var.ovirt_cluster_id
  template_id = var.release_image_template_id

  // TODO implement initialization
  initialization {
    custom_script = var.ignition_bootstrap
  }
}

resource "ovirt_tag" "cluster_bootstrap_tag" {
  name = "${var.cluster_id}-bootstrap"
}

resource "ovirt_vm_tag" "cluster_bootstrap_tag" {
  vm_id  = ovirt_vm.bootstrap.id
  tag_id = ovirt_tag.cluster_bootstrap_tag.id
}
