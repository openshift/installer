output "cluster_domain" {
  value = var.cluster_domain
}
output "cluster_id" {
  value = var.cluster_id
}
output "tags" {
  value = [vsphere_tag.tag.id]
}
