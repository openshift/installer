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
