data "template_disk_attachments" "master" {
  template_id = var.ovirt_template_id
}

resource "ovirt_vm" "master" {
  count               = var.master_count
  name                = "${var.cluster_id}-master-${count.index}"
  cluster_id          = var.ovirt_cluster_id
  template_id         = var.ovirt_template_id
  // TODO implement instance_type_id
  instance_type_id    = var.ovirt_master_instance_type_id
  // TODO implement vm_type
  type                = var.ovirt_master_vm_type
  cpu_cores           = var.ovirt_master_cores
  cpu_sockets         = var.ovirt_master_sockets
  cpu_threads         = 1
  // if instance type is declared then memory is redundant. Since terraform
  // doesn't allow to condionally omit it, it must be passed.
  // The number passed is multiplied by 4 and becomes the maximum memory the VM can have.
  memory              = var.ovirt_master_instance_type_id != "" ? 16348 : var.ovirt_master_memory
  auto_pinning_policy = var.ovirt_master_auto_pinning_policy != "" ? var.ovirt_master_auto_pinning_policy : null
  hugepages           = var.ovirt_master_hugepages > 0 ? var.ovirt_master_hugepages : null

  # Here we check if the ovirt_master_clone is set and use that as a bool if yes, default to the VM type otherwise.
  #
  # Clone set explicitly -> clone = var.ovirt_master_clone
  # VM type desktop -> clone = false
  # VM type server or high performance -> clone = true
  clone = var.ovirt_master_clone != null ? tobool(var.ovirt_master_clone) : (var.ovirt_master_vm_type == "desktop" ? false : true)

  // TODO implement initialization
  initialization {
    host_name     = "${var.cluster_id}-master-${count.index}"
    custom_script = var.ignition_master
  }

  // TODO implement template_disk_attachment_override
  dynamic "template_disk_attachment_override" {
    for_each       = template_disk_attachments.disk
    content {
      disk_id        = template_block_device.value["disk_id"]
      interface      = "virtio_scsi"
      size           = var.ovirt_master_os_disk_size_gb
      format         = var.ovirt_master_format != "" ? var.ovirt_master_format : null
      sparse         = tobool(var.ovirt_master_sparse)
    }
  }
  depends_on = [var.ovirt_affinity_group_count]
}

// TODO implement ovirt_affinity_group
data "ovirt_affinity_group" "master" {
  count = length(var.ovirt_master_affinity_groups)
  name  = var.ovirt_master_affinity_groups[count.index]
}

// TODO implement ovirt_vm_affinity_group
resource "ovirt_vm_affinity_group" "master" {
  count             = length(var.ovirt_master_affinity_groups)
  vm_id             = ovirt_vm.master.id
  affinity_group_id = data.ovirt_affinity_group.master[count.index].id
}

resource "ovirt_vm_tag" "master" {
  count = var.master_count

  vm_id  = ovirt_vm.master[count.index].id
  tag_id = var.tag_id
}

resource "ovirt_vm_start" "master" {
  count = var.master_count
  vm_id = ovirt_vm.master[0].id

  depends_on = [ovirt_vm_tag.master, ovirt_vm_affinity_group.master]
}
