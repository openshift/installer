locals {
  description         = "Created By OpenShift Installer"
  vcenter_key         = keys(var.vsphere_vcenters)[0]
  tag_category_key    = keys(var.vsphere_failure_zone)[0]
  region_tag_category = var.vsphere_failure_zone[local.tag_category_key].region.tagCategory
  zone_tag_category   = var.vsphere_failure_zone[local.tag_category_key].zone.tagCategory
}

provider "vsphere" {
  user                 = var.vsphere_vcenters[local.vcenter_key].user
  password             = var.vsphere_vcenters[local.vcenter_key].password
  vsphere_server       = var.vsphere_vcenters[local.vcenter_key].server
  allow_unverified_ssl = false
}

provider "vsphereprivate" {
  user                 = var.vsphere_vcenters[local.vcenter_key].user
  password             = var.vsphere_vcenters[local.vcenter_key].password
  vsphere_server       = var.vsphere_vcenters[local.vcenter_key].server
  allow_unverified_ssl = false
}

data "vsphere_datacenter" "datacenter" {
  for_each = var.vsphere_deployment_zone
  name     = var.vsphere_failure_zone[each.value.failureDomain].topology.datacenter
}

data "vsphere_compute_cluster" "cluster" {
  for_each      = var.vsphere_deployment_zone
  name          = var.vsphere_failure_zone[each.value.failureDomain].topology.computeCluster
  datacenter_id = data.vsphere_datacenter.datacenter[each.key].id
}

data "vsphere_resource_pool" "resource_pool" {
  for_each = var.vsphere_deployment_zone
  name     = each.value.placementConstraint.resourcePool
}

data "vsphere_datastore" "datastore" {
  for_each      = var.vsphere_deployment_zone
  name          = var.vsphere_failure_zone[each.value.failureDomain].topology.datastore
  datacenter_id = data.vsphere_datacenter.datacenter[each.key].id
}

data "vsphere_virtual_machine" "template" {
  for_each      = var.vsphere_deployment_zone
  name          = vsphereprivate_import_ova.import[each.key].name
  datacenter_id = data.vsphere_datacenter.datacenter[each.key].id
}

data "vsphere_datacenter" "folder_datacenter" {
  for_each = var.vsphere_folder_zone
  name     = each.value.vsphere_datacenter
}

resource "vsphere_folder" "folder" {
  for_each = var.vsphere_folder_zone

  path          = each.value.name
  type          = "vm"
  datacenter_id = data.vsphere_datacenter.folder_datacenter[each.key].id
  tags          = [vsphere_tag.tag.id]
}

resource "vsphereprivate_import_ova" "import" {
  for_each = var.vsphere_deployment_zone

  name          = format("%s-%s-%s", var.vsphere_template, var.vsphere_failure_zone[each.value.failureDomain].region.name, var.vsphere_failure_zone[each.value.failureDomain].zone.name)
  filename      = var.vsphere_ova_filepath
  cluster       = data.vsphere_compute_cluster.cluster[each.key].name
  resource_pool = data.vsphere_resource_pool.resource_pool[each.key].name
  datacenter    = data.vsphere_datacenter.datacenter[each.key].name
  datastore     = data.vsphere_datastore.datastore[each.key].name
  network       = var.vsphere_network_zone[each.key]
  folder        = each.value.placementConstraint.folder
  tag           = vsphere_tag.tag.id
  disk_type     = var.vsphere_disk_type
}

resource "vsphere_tag_category" "category" {
  name        = "openshift-${var.cluster_id}"
  description = "Added by openshift-install do not remove"
  cardinality = "SINGLE"

  associable_types = [
    "VirtualMachine",
    "ResourcePool",
    "Folder",
    "Datastore",
    "StoragePod"
  ]
}

resource "vsphere_tag" "tag" {
  name        = var.cluster_id
  category_id = vsphere_tag_category.category.id
  description = "Added by openshift-install do not remove"
}
