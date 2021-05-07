output "ovirt_affinity_group_count" {
  value = length(ovirt_affinity_group.affinity_groups)
}
