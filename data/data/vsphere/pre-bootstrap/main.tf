locals {
  folder      = var.vsphere_preexisting_folder ? var.vsphere_folder : vsphere_folder.folder[0].path
  description = "Created By OpenShift Installer"
  vcenter_region_zone_flatten = flatten([

    for v in var.vsphere_vcenters : [
      for r in v.regions : [
        for z in r.zones : {
          region         = r.name
          zone           = z.name
          datacenter     = r.datacenter
          cluster        = z.cluster
          datastore      = z.datastore
          network        = z.network
          vsphere_server = v.server
          user           = v.user
          password       = v.password
        }
      ]
    ]
  ])

  vcenter_region_zone_map = {
    for obj in local.vcenter_region_zone_flatten : "${obj.vsphere_server}-${obj.region}-${obj.zone}" => obj
  }
}

provider "vsphere" {
  user                 = var.vsphere_username
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_url
  allow_unverified_ssl = false
}

provider "vsphereprivate" {
  user                 = var.vsphere_username
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_url
  allow_unverified_ssl = false
}


//// Zoning

data "vsphere_datacenter" "datacenter_zoning" {
  for_each = local.vcenter_region_zone_map
  name     = each.value.datacenter
}

data "vsphere_compute_cluster" "cluster_zoning" {
  for_each      = local.vcenter_region_zone_map
  name          = each.value.cluster
  datacenter_id = data.vsphere_datacenter.datacenter_zoning[each.key].id
}

data "vsphere_datastore" "datastore_zoning" {

  for_each      = local.vcenter_region_zone_map
  name          = each.value.datastore
  datacenter_id = data.vsphere_datacenter.datacenter_zoning[each.key].id
}

data "vsphere_network" "network_zoning" {
  for_each      = local.vcenter_region_zone_map
  name          = each.value.network
  datacenter_id = data.vsphere_datacenter.datacenter_zoning[each.key].id
}

resource "vsphere_tag_category" "region_tag_category" {
  name        = "openshift-region"
  description = "Added by openshift-install do not remove"
  cardinality = "SINGLE"

  associable_types = [
    "Datacenter"
  ]
}

resource "vsphere_tag_category" "zone_tag_category" {
  name        = "openshift-zone"
  description = "Added by openshift-install do not remove"
  cardinality = "SINGLE"

  associable_types = [
    "Cluster"
  ]
}

resource "vsphere_tag" "region_tags" {
  for_each    = { for k, v in local.vcenter_region_zone_map : k => v.region }
  name        = each.value
  category_id = vsphere_tag_category.region_tag_category.id
  description = "Added by openshift-install do not remove"
}

resource "vsphere_tag" "zone_tags" {
  for_each    = { for k, v in local.vcenter_region_zone_map : k => v.zone }
  name        = each.value
  category_id = vsphere_tag_category.zone_tag_category.id
  description = "Added by openshift-install do not remove"
}

//// End Zoning

data "vsphere_datacenter" "datacenter" {
  name = var.vsphere_datacenter
}

data "vsphere_compute_cluster" "cluster" {
  name          = var.vsphere_cluster
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

data "vsphere_resource_pool" "resource_pool" {
  name = var.vsphere_resource_pool
}

data "vsphere_datastore" "datastore" {
  name          = var.vsphere_datastore
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

data "vsphere_network" "network" {
  name          = var.vsphere_network
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

data "vsphere_virtual_machine" "template" {
  name          = vsphereprivate_import_ova.import.name
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

resource "vsphereprivate_import_ova" "import" {
  name          = var.vsphere_template
  filename      = var.vsphere_ova_filepath
  cluster       = var.vsphere_cluster
  resource_pool = var.vsphere_resource_pool
  datacenter    = var.vsphere_datacenter
  datastore     = var.vsphere_datastore
  network       = var.vsphere_network
  folder        = local.folder
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

resource "vsphere_folder" "folder" {
  count = var.vsphere_preexisting_folder ? 0 : 1

  path          = var.vsphere_folder
  type          = "vm"
  datacenter_id = data.vsphere_datacenter.datacenter.id
  tags          = [vsphere_tag.tag.id]
}
