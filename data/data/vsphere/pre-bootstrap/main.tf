locals {
  folder      = var.vsphere_preexisting_folder ? var.vsphere_folder : vsphere_folder.folder[0].path
  description = "Created By OpenShift Installer"
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

data "vsphere_datacenter" "datacenter" {
  name = var.vsphere_datacenter
}

data "vsphere_compute_cluster" "cluster" {
  name          = var.vsphere_cluster
  datacenter_id = data.vsphere_datacenter.datacenter.id
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
  name       = var.vsphere_template
  filename   = var.vsphere_ova_filepath
  cluster    = var.vsphere_cluster
  datacenter = var.vsphere_datacenter
  datastore  = var.vsphere_datastore
  network    = var.vsphere_network
  folder     = local.folder
  tag        = vsphere_tag.tag.id
  disk_type  = var.vsphere_disk_type
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

