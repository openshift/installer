// this module is responsible to create the unique template
// for the openshift cluster and has this properties
// 1. the name of the template will be always set after the name
// of the openshift cluster (var.cluster_id) i.e 'clustername-4t9hs2'
// which the CLUSTER.INFRA_ID
// 2. the disk.alias (the disk name) will be set to the releaseImage name
// as set by the installer, and in terraform is var.openstack_base_image_name.

locals {
  image_name = "${var.cluster_id}-rhcos"
}

// template created using the uploaded image
resource "ovirt_template" "releaseimage_template" {
  count = var.tmp_import_vm_id != "" ? 1 : 0

  // name the template after the openshift cluster id
  name        = local.image_name
  description = "Template in use by OpenShift. Do not delete!"
  // create from vm
  vm_id = var.tmp_import_vm_id
}

// existing template provided by the user
data "ovirt_templates" "finalTemplate" {
  count = var.tmp_import_vm_id == "" ? 1 : 0

  fail_on_empty = true
  name          = var.openstack_base_image_name
}
