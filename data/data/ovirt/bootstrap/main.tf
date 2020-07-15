resource "ovirt_vm" "bootstrap" {
  name        = "${var.cluster_id}-bootstrap"
  memory      = "8192"
  cores       = "4"
  cluster_id  = var.ovirt_cluster_id
  template_id = var.ovirt_template_id

  initialization {
    custom_script = var.ignition_bootstrap
  }
}

resource "ovirt_tag" "cluster_bootstrap_tag" {
  name   = "${var.cluster_id}-bootstrap"
  vm_ids = concat([ovirt_vm.bootstrap.id], [var.ovirt_tmp_template_vm_id])
}
