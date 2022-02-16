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
  memory              = var.ovirt_master_instance_type_id != "" ? 16348 : var.ovirt_master_memory
  affinity_groups     = var.ovirt_master_affinity_groups
  auto_pinning_policy = var.ovirt_master_auto_pinning_policy != "" ? var.ovirt_master_auto_pinning_policy : null
  hugepages           = var.ovirt_master_hugepages > 0 ? var.ovirt_master_hugepages : null

  # Here we check if the ovirt_master_clone is set and use that as a bool if yes, default to the VM type otherwise.
  #
  # Clone set explicitly -> clone = var.ovirt_master_clone
  # VM type desktop -> clone = false
  # VM type server or high performance -> clone = true
  clone = var.ovirt_master_clone != null ? tobool(var.ovirt_master_clone) : (var.ovirt_master_vm_type == "desktop" ? false : true)

  initialization {
    host_name     = "${var.cluster_id}-master-${count.index}"
    custom_script = var.ignition_master
  }

  block_device {
    interface      = "virtio_scsi"
    size           = var.ovirt_master_os_disk_size_gb
    format         = var.ovirt_master_format != "" ? var.ovirt_master_format : null
    sparse         = tobool(var.ovirt_master_sparse)
    storage_domain = var.ovirt_storage_domain_id
  }
  depends_on = [var.ovirt_affinity_group_count]
}


resource "ovirt_tag" "cluster_tag" {
  name   = var.cluster_id
  vm_ids = [for instance in ovirt_vm.master.* : instance.id]
}
