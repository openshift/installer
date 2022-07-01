resource "ovirt_affinity_group" "affinity_groups" {
  count       = length(var.ovirt_affinity_groups)
  cluster_id  = var.ovirt_cluster_id
  name        = var.ovirt_affinity_groups[count.index]["name"]
  description = var.ovirt_affinity_groups[count.index]["description"]
  priority    = var.ovirt_affinity_groups[count.index]["priority"]
  enforcing   = var.ovirt_affinity_groups[count.index]["enforcing"]
  vms_rule {
    affinity  = "negative"
    enforcing = var.ovirt_affinity_groups[count.index]["enforcing"]
  }
}
