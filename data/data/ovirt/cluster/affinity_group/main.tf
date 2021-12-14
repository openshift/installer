resource "ovirt_affinity_group" "affinity_groups" {
  count        = length(var.ovirt_affinity_groups)
  name         = var.ovirt_affinity_groups[count.index]["name"]
  description  = var.ovirt_affinity_groups[count.index]["description"]
  cluster_id   = var.ovirt_cluster_id
  priority     = var.ovirt_affinity_groups[count.index]["priority"]
  vm_positive  = false
  vm_enforcing = var.ovirt_affinity_groups[count.index]["enforcing"]
  lifecycle {
    ignore_changes = [vm_list]
  }
}
