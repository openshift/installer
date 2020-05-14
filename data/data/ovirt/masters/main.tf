resource "ovirt_vm" "master" {
  count            = var.master_count
  name             = "${var.cluster_id}-master-${count.index}"
  cluster_id       = var.ovirt_cluster_id
  template_id      = var.ovirt_template_id
  instance_type_id = var.ovirt_master_instance_type_id
  type             = var.ovirt_master_vm_type
  cores            = var.ovirt_master_cores
  sockets          = var.ovirt_master_sockets
  // if instance type is declared then memory is redundant. Since terraform
  // doesn't allow to condionally omit it, it must be passed.
  // The number passed is multiplied by 4 and becomes the maximum memory the VM can have.
  memory = var.ovirt_master_instance_type_id != "" ? 16348 : var.ovirt_master_memory

  initialization {
    host_name     = "${var.cluster_id}-master-${count.index}"
    custom_script = var.ignition_master
  }

  block_device {
    interface = "virtio_scsi"
    size      = var.ovirt_master_os_disk_size_gb
  }
}

resource "ovirt_tag" "cluster_tag" {
  name   = var.cluster_id
  vm_ids = [for instance in ovirt_vm.master.* : instance.id]
}
