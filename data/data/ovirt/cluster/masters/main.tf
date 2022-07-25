data "ovirt_template_disk_attachments" "master" {
  template_id = var.ovirt_template_id
}

data "ovirt_cluster_hosts" "master" {
  cluster_id = var.ovirt_cluster_id
}

data "ovirt_affinity_group" "master" {
  count      = var.ovirt_master_affinity_groups == null ? 0 : length(var.ovirt_master_affinity_groups)
  cluster_id = var.ovirt_cluster_id
  name       = var.ovirt_master_affinity_groups[count.index]
  depends_on = [var.ovirt_affinity_group_count]
}

locals {
  vm_affinity_groups = [
    for pair in setproduct(data.ovirt_affinity_group.master.*.id, ovirt_vm.master.*.id) : {
      affinity_group_id = pair[0]
      vm_id             = pair[1]
    }
  ]
}

resource "ovirt_vm_affinity_group" "master" {
  count             = length(local.vm_affinity_groups)
  vm_id             = local.vm_affinity_groups[count.index].vm_id
  cluster_id        = var.ovirt_cluster_id
  affinity_group_id = local.vm_affinity_groups[count.index].affinity_group_id
}

// ovirt_vm creates the master nodes
resource "ovirt_vm" "master" {
  count            = var.master_count
  name             = "${var.cluster_id}-master-${count.index}"
  cluster_id       = var.ovirt_cluster_id
  template_id      = var.ovirt_template_id
  instance_type_id = var.ovirt_master_instance_type_id != "" ? var.ovirt_master_instance_type_id : null
  vm_type          = var.ovirt_master_vm_type
  cpu_cores        = var.ovirt_master_cores
  cpu_sockets      = var.ovirt_master_sockets
  cpu_threads      = var.ovirt_master_threads

  // if instance type is declared then memory is redundant. Since terraform
  // doesn't allow to conditionally omit it, it must be passed.
  // The number passed is multiplied by 4 and becomes the maximum memory the VM can have.
  memory = var.ovirt_master_instance_type_id != "" || var.ovirt_master_memory == "" ? 16348 * 1024 * 1024 : tonumber(var.ovirt_master_memory) * 1024 * 1024

  huge_pages        = var.ovirt_master_hugepages > 0 ? var.ovirt_master_hugepages : null
  serial_console    = var.ovirt_master_vm_type == "high_performance" ? true : null
  soundcard_enabled = var.ovirt_master_vm_type == "high_performance" ? false : null
  memory_ballooning = var.ovirt_master_vm_type == "high_performance" ? false : null
  cpu_mode          = var.ovirt_master_vm_type == "high_performance" ? "host_passthrough" : null

  # Here we check if the ovirt_master_clone is set and use that as a bool if yes, default to the VM type otherwise.
  #
  # Clone set explicitly -> clone = var.ovirt_master_clone
  # VM type desktop -> clone = false
  # VM type server or high performance -> clone = true
  clone = var.ovirt_master_clone != null ? tobool(var.ovirt_master_clone) : (var.ovirt_master_vm_type == "desktop" ? false : true)

  # Initialization sets the host name and script run when the machine first starts.
  initialization_hostname      = "${var.cluster_id}-master-${count.index}"
  initialization_custom_script = var.ignition_master

  # Placement policy dictates which hosts this master can run on.
  #
  # TODO there may be a bug here since we are pinning the masters to the existing detected hosts and this is never
  #      updated.
  placement_policy_affinity = var.ovirt_master_auto_pinning_policy != "" && var.ovirt_master_auto_pinning_policy != "none" ? "migratable" : null
  placement_policy_host_ids = var.ovirt_master_auto_pinning_policy != "" && var.ovirt_master_auto_pinning_policy != "none" ? data.ovirt_cluster_hosts.master.hosts.*.id : null

  # This section overrides the format and sparse option for the disks from the template.
  dynamic "template_disk_attachment_override" {
    for_each = data.ovirt_template_disk_attachments.master.disk_attachments
    content {
      disk_id      = template_disk_attachment_override.value.disk_id
      format       = var.ovirt_master_format != "" ? var.ovirt_master_format : null
      provisioning = var.ovirt_master_sparse == null ? null : (tobool(var.ovirt_master_sparse) ? "sparse" : "non-sparse")
    }
  }
  depends_on = [var.ovirt_affinity_group_count]
}

data "ovirt_disk_attachments" "master" {
  count = var.master_count
  vm_id = ovirt_vm.master.*.id[count.index]
}

// ovirt_vm_disks_resize resizes the master disks to the specified size.
resource "ovirt_vm_disks_resize" "master" {
  count = var.master_count
  vm_id = ovirt_vm.master.*.id[count.index]
  size  = var.ovirt_master_os_disk_size_gb * 1024 * 1024 * 1024
}

// ovirt_vm_graphic_consoles removes the graphic consoles from non-desktop machines.
resource "ovirt_vm_graphics_consoles" "master" {
  count = var.ovirt_master_vm_type == "high_performance" ? var.master_count : 0
  vm_id = ovirt_vm.master.*.id[count.index]
}

// ovirt_vm_optimize_cpu_settings auto-optimizes CPU and NUMA alignment on server and HP types
resource "ovirt_vm_optimize_cpu_settings" "master" {
  count = var.ovirt_master_auto_pinning_policy != "" && var.ovirt_master_auto_pinning_policy != "none" ? var.master_count : 0
  vm_id = ovirt_vm.master.*.id[count.index]
}

// ovirt_vm_start starts the master nodes.
resource "ovirt_vm_start" "master" {
  count = var.master_count
  vm_id = ovirt_vm.master.*.id[count.index]

  depends_on = [
    ovirt_vm_graphics_consoles.master,
    ovirt_vm_optimize_cpu_settings.master,
    ovirt_vm_disks_resize.master,
    ovirt_vm_tag.master,
    ovirt_vm_affinity_group.master,
  ]
}

resource "ovirt_tag" "cluster_tag" {
  name = var.cluster_id
}

resource "ovirt_vm_tag" "master" {
  count  = length(ovirt_vm.master)
  tag_id = ovirt_tag.cluster_tag.id
  vm_id  = ovirt_vm.master.*.id[count.index]
}
