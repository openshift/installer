locals {
  folders = { for k, v in vsphere_folder.folder : k => v.path }
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

/*
  region_map = {
    for obj in local.vcenter_region_zone_flatten : "${obj.vsphere_server}-${obj.region}-${obj.zone}" => obj.
  }
*/
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


// TODO: can we keep this consistant?
// with local.vcenter_region_zone_map

data "vsphere_datacenter" "datacenter_zoning" {
  for_each = var.regions_map
  name     = each.value.datacenter
}

data "vsphere_compute_cluster" "cluster_zoning" {
  for_each      = local.vcenter_region_zone_map
  name          = each.value.cluster
  datacenter_id = data.vsphere_datacenter.datacenter_zoning[each.value.region].id
}

data "vsphere_datastore" "datastore_zoning" {
  for_each      = local.vcenter_region_zone_map
  name          = each.value.datastore
  datacenter_id = data.vsphere_datacenter.datacenter_zoning[each.value.region].id
}

data "vsphere_virtual_machine" "template" {
  for_each      = local.vcenter_region_zone_map
  name          = vsphereprivate_import_ova.import[each.key].name
  datacenter_id = data.vsphere_datacenter.datacenter_zoning[each.value.region].id
}

resource "vsphereprivate_import_ova" "import" {
  for_each = local.vcenter_region_zone_map

  name       = "${var.vsphere_template}-${each.value.zone}"
  filename   = var.vsphere_ova_filepath
  cluster    = each.value.cluster
  datacenter = each.value.datacenter
  datastore  = each.value.datastore
  network    = each.value.network
  folder     = local.folders[each.value.region]
  tag        = vsphere_tag.tag.id
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
    "ClusterComputeResource"
  ]
}

resource "vsphere_tag" "regions_tags" {
  for_each    = var.regions_map
  name        = each.key
  category_id = vsphere_tag_category.region_tag_category.id
  description = "Added by openshift-install do not remove"
}

resource "vsphere_tag" "zones_tags" {
  for_each    = var.zones_map
  name        = each.key
  category_id = vsphere_tag_category.zone_tag_category.id
  description = "Added by openshift-install do not remove"
}

resource "vsphereprivate_tag_attach" "regions_datacenters" {
  for_each = var.regions_map

  objectid   = data.vsphere_datacenter.datacenter_zoning[each.key].id
  objecttype = "Datacenter"

  tagid = vsphere_tag.regions_tags[each.key].id
}

resource "vsphereprivate_tag_attach" "zones_clusters" {
  for_each = local.vcenter_region_zone_map

  objectid   = data.vsphere_compute_cluster.cluster_zoning[each.key].id
  objecttype = "ClusterComputeResource"

  tagid = vsphere_tag.zones_tags[each.value.zone].id
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
  for_each = var.regions_map

  path          = "${var.vsphere_folder}-${each.value.name}"
  type          = "vm"
  datacenter_id = data.vsphere_datacenter.datacenter_zoning[each.key].id
  tags          = [vsphere_tag.tag.id]
}
