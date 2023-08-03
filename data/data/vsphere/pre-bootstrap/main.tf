locals {
  description = "Created By OpenShift Installer"


  // todo: multiple vcenter change
  vcenter_key = keys(var.vsphere_vcenters)[0]
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

// todo: multiple vcenter change, maybe?
data "vsphere_datacenter" "datacenter" {
  for_each = var.vsphere_failure_domain_map
  name     = var.vsphere_failure_domain_map[each.key].topology.datacenter
}

data "vsphere_compute_cluster" "cluster" {
  for_each      = var.vsphere_failure_domain_map
  name          = var.vsphere_failure_domain_map[each.key].topology.computeCluster
  datacenter_id = data.vsphere_datacenter.datacenter[each.key].id
}

data "vsphere_resource_pool" "resource_pool" {
  for_each = var.vsphere_failure_domain_map
  name     = var.vsphere_failure_domain_map[each.key].topology.resourcePool
}

data "vsphere_datastore" "datastore" {
  for_each      = var.vsphere_failure_domain_map
  name          = var.vsphere_failure_domain_map[each.key].topology.datastore
  datacenter_id = data.vsphere_datacenter.datacenter[each.key].id
}

// Why is there two datacenters?
// The vm folder object is defined at the datacenter
// level. Each failure domain has a datacenter folder pair/
// We need to get only the unique datacenter-folder pair
// and create those folders. See vsphere.go
// createDatacenterFolderMap

data "vsphere_datacenter" "folder_datacenter" {
  for_each = var.vsphere_folders
  name     = each.value.vsphere_datacenter
}

resource "vsphere_folder" "folder" {
  for_each = var.vsphere_folders

  path          = each.value.name
  type          = "vm"
  datacenter_id = data.vsphere_datacenter.folder_datacenter[each.key].id
  tags          = [vsphere_tag.tag.id]
}

resource "vsphereprivate_import_ova" "import" {
  for_each = var.vsphere_import_ova_failure_domain_map

  name = format("%s-rhcos-%s-%s", var.cluster_id, var.vsphere_failure_domain_map[each.key].region, var.vsphere_failure_domain_map[each.key].zone)

  filename      = var.vsphere_ova_filepath
  cluster       = data.vsphere_compute_cluster.cluster[each.key].name
  resource_pool = data.vsphere_resource_pool.resource_pool[each.key].name
  datacenter    = data.vsphere_datacenter.datacenter[each.key].name
  datastore     = data.vsphere_datastore.datastore[each.key].name

  network   = var.vsphere_networks[each.key]
  folder    = var.vsphere_failure_domain_map[each.key].topology.folder
  tag       = vsphere_tag.tag.id
  disk_type = var.vsphere_disk_type

  // Since the folder resource might not be ran because there could be
  // user defined folder per failure domain if a folder is created
  // the import resource is not waiting. Adding
  // this depends_on so the import happens after creating folder(s).
  depends_on = [vsphere_folder.folder]
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
