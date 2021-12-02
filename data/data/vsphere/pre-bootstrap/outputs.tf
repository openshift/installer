output "datacenter" {
  value = data.vsphere_datacenter.datacenter_zoning
}

output "datastore" {
  value = data.vsphere_datastore.datastore_zoning
}

output "cluster" {
  value = data.vsphere_compute_cluster.cluster_zoning
}

output "folder" {
  value = local.folders
}

output "ovaimport" {
  value = vsphereprivate_import_ova.import
}

output "template" {
  value = data.vsphere_virtual_machine.template
}

output "tags" {
  value = [vsphere_tag.tag.id]
}

output "cluster_domain" {
  value = var.cluster_domain
}

output "cluster_id" {
  value = var.cluster_id
}

output "vcenter_region_zone_map" {
  value = local.vcenter_region_zone_map
}
