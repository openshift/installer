output "resource_pool" {
  value = data.vsphere_resource_pool.resource_pool.id
}

output "datastore" {
  value = data.vsphere_datastore.datastore.id
}

output "folder" {
  value = local.folder
}

output "datacenter" {
  value = data.vsphere_datacenter.datacenter.id
}

output "template" {
  value = data.vsphere_virtual_machine.template.id
}

output "guest_id" {
  value = data.vsphere_virtual_machine.template.guest_id
}

output "thin_disk" {
  value = data.vsphere_virtual_machine.template.disks.0.thin_provisioned
}

output "scrub_disk" {
  value = data.vsphere_virtual_machine.template.disks.0.eagerly_scrub
}

output "cluster_domain" {
  value = var.cluster_domain
}

output "cluster_id" {
  value = var.cluster_id
}

output "tags" {
  value = [vsphere_tag.tag.id]
}
