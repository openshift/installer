provider "ovirt" {
  url           = var.ovirt_url
  username      = var.ovirt_username
  password      = var.ovirt_password
  tls_ca_files  = var.ovirt_cafile == "" ? [] : [var.ovirt_cafile]
  tls_ca_bundle = var.ovirt_ca_bundle
  tls_insecure  = var.ovirt_insecure
}

resource "ovirt_vm" "bootstrap" {
  name        = "${var.cluster_id}-bootstrap"
  cluster_id  = var.ovirt_cluster_id
  template_id = var.release_image_template_id

  memory      = 8 * 1024 * 1024 * 1024
  cpu_cores   = 4
  cpu_threads = 1
  cpu_sockets = 1

  initialization_custom_script = var.ignition_bootstrap
}

resource "ovirt_tag" "cluster_bootstrap_tag" {
  name = "${var.cluster_id}-bootstrap"
}

resource "ovirt_vm_tag" "cluster_bootstrap_tag" {
  tag_id = ovirt_tag.cluster_bootstrap_tag.id
  vm_id  = ovirt_vm.bootstrap.id
}

resource "ovirt_vm_tag" "cluster_import_tag" {
  tag_id = ovirt_tag.cluster_bootstrap_tag.id
  vm_id  = var.tmp_import_vm_id
}

// ovirt_vm_start starts the master nodes.
resource "ovirt_vm_start" "bootstrap" {
  vm_id = ovirt_vm.bootstrap.id

  depends_on = [
    ovirt_vm_tag.cluster_bootstrap_tag,
    ovirt_vm_tag.cluster_import_tag,
  ]
}