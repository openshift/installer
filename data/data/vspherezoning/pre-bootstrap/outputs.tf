output "resource_pool" {
  value = data.vsphere_resource_pool.resource_pool
}

output "datastore" {
  value = data.vsphere_datastore.datastore
}

output "datacenter" {
  value = data.vsphere_datacenter.datacenter
}

output "template" {
  value = data.vsphere_virtual_machine.template
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
