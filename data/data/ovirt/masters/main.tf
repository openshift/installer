resource "ovirt_vm" "master" {
  count       = var.master_count
  name        = "${var.cluster_id}-master-${count.index}"
  cluster_id  = var.ovirt_cluster_id
  template_id = var.ovirt_template_id
  cores       = var.ovirt_master_cpu
  memory      = var.ovirt_master_mem

  initialization {
    host_name     = "${var.cluster_id}-master-${count.index}"
    custom_script = var.ignition_master
  }
}

resource "ovirt_tag" "cluster_tag" {
  name   = var.cluster_id
  vm_ids = [for instance in ovirt_vm.master.* : instance.id]
}
