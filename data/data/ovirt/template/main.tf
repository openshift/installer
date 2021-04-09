// this module is responsible to create the unique template
// for the openshift cluster and has this properties
// 1. the name of the template will be always set after the name
// of the openshift cluster (var.cluster_id) i.e 'clustername-4t9hs2'
// which the CLUSTER.INFRA_ID
// 2. the disk.alias (the disk name) will be set to the releaseImage name
// as set by the installer, and in terraform is var.openstack_base_image_name.

// make this search return at least 1 result to signal we need to create
// the template.
data "ovirt_templates" "osImage" {
  search = {
    criteria       = "name=${var.openstack_base_image_name} or name=Blank"
    case_sensitive = true
  }
}

data "ovirt_clusters" "clusters" {
  search = {
    criteria       = ""
    case_sensitive = false
  }
}

// work around the missing regexall in terraform < 0.12.9
// if length(regexall("^Blank.*$", t.name)
locals {
  existing_id = [for t in data.ovirt_templates.osImage.templates : t.id if substr(t.name, 0, 5) != "Blank"]
  // if we don't find the cluster this should fail
  cluster    = [for c in data.ovirt_clusters.clusters.clusters : c if c.id == var.ovirt_cluster_id][0]
  network_id = [for n in local.cluster.networks : n.id if n.name == var.ovirt_network_name][0]
}

// upload the disk if we don't have an existing template
resource "ovirt_image_transfer" "releaseimage" {
  count             = length(local.existing_id) == 0 && var.bootstrapping ? 1 : 0
  alias             = var.openstack_base_image_name
  source_url        = var.openstack_base_image_local_file_path
  storage_domain_id = var.ovirt_storage_domain_id
  sparse            = true
  timeouts {
    create = "20m"
  }
}

resource "ovirt_vm" "tmp_import_vm" {
  // create the vm for import only when we don't have an existing template
  count      = length(local.existing_id) == 0 && var.bootstrapping ? 1 : 0
  name       = "tmpvm-for-${ovirt_image_transfer.releaseimage.0.alias}"
  cluster_id = var.ovirt_cluster_id
  auto_start = false
  block_device {
    disk_id   = ovirt_image_transfer.releaseimage.0.disk_id
    interface = "virtio_scsi"
  }
  os {
    type = "rhcos_x64"
  }
  nics {
    name            = "nic1"
    vnic_profile_id = var.ovirt_vnic_profile_id
  }
  timeouts {
    create = "20m"
  }
  depends_on = [ovirt_image_transfer.releaseimage]
}

data "ovirt_vms" "tmp_import_vm_data" {
  count = length(local.existing_id) == 0 ? 1 : 0
  search = {
    criteria       = "name=tmpvm-for-${ovirt_image_transfer.releaseimage.0.alias}"
    case_sensitive = true
  }
  depends_on = [ovirt_vm.tmp_import_vm]
}

resource "ovirt_template" "releaseimage_template" {
  // create the template only when we don't have an existing template
  count = length(local.existing_id) == 0 ? 1 : 0
  // name the template after the openshift cluster id
  name       = var.openstack_base_image_name
  cluster_id = data.ovirt_vms.tmp_import_vm_data.0.vms.0.cluster_id
  // create from vm
  vm_id = data.ovirt_vms.tmp_import_vm_data.0.vms.0.id
  timeouts {
    create = "20m"
  }
}

// finally get the template by name(should be unique), fail if it doesn't exist
data "ovirt_templates" "finalTemplate" {
  search = {
    criteria       = "name=${var.openstack_base_image_name}"
    case_sensitive = true
  }
  depends_on = [ovirt_template.releaseimage_template]
}
