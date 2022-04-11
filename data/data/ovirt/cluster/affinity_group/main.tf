resource "ovirt_affinity_group" "affinity_groups" {
  count       = length(var.ovirt_affinity_groups)
  name        = var.ovirt_affinity_groups[count.index]["name"]
  // TODO implement description
  description = var.ovirt_affinity_groups[count.index]["description"]
  cluster_id  = var.ovirt_cluster_id
  priority    = var.ovirt_affinity_groups[count.index]["priority"]
  vms_rule {
    enforcing = var.ovirt_affinity_groups[count.index]["enforcing"]
    affinity  = "negative"
  }
}
