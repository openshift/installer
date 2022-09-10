locals {
  description           = "Created By OpenShift Installer"
  vcenter_key           = keys(var.vsphere_vcenters)[0]
  failure_domains_count = length(var.vsphere_failure_domains)
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
  count = local.failure_domains_count
  name  = var.vsphere_failure_domains[count.index].topology.datacenter
}

data "vsphere_compute_cluster" "cluster" {
  count         = local.failure_domains_count
  name          = var.vsphere_failure_domains[count.index].topology.computeCluster
  datacenter_id = data.vsphere_datacenter.datacenter[count.index].id
}

data "vsphere_resource_pool" "resource_pool" {
  count = local.failure_domains_count
  name  = var.vsphere_failure_domains[count.index].topology.resourcePool
}

data "vsphere_datastore" "datastore" {
  count         = local.failure_domains_count
  name          = var.vsphere_failure_domains[count.index].topology.datastore
  datacenter_id = data.vsphere_datacenter.datacenter[count.index].id
}

data "vsphere_virtual_machine" "template" {
  count         = local.failure_domains_count
  name          = vsphereprivate_import_ova.import[count.index].name
  datacenter_id = data.vsphere_datacenter.datacenter[count.index].id
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
  count = local.failure_domains_count
  name  = format("%s-rhcos-%s-%s", var.cluster_id, var.vsphere_failure_domains[count.index].region, var.vsphere_failure_domains[count.index].zone)

  filename      = var.vsphere_ova_filepath
  cluster       = data.vsphere_compute_cluster.cluster[count.index].name
  resource_pool = data.vsphere_resource_pool.resource_pool[count.index].name
  datacenter    = data.vsphere_datacenter.datacenter[count.index].name
  datastore     = data.vsphere_datastore.datastore[count.index].name

  network = var.vsphere_networks[var.vsphere_failure_domains[count.index].name]


  folder    = var.vsphere_failure_domains[count.index].topology.folder
  tag       = vsphere_tag.tag.id
  disk_type = var.vsphere_disk_type
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
