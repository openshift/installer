# Pre-create server groups for the Compute MachineSets, with the given policy.
resource "openstack_compute_servergroup_v2" "server_groups" {
  for_each = var.server_group_names
  name     = each.key
  policies = [var.server_group_policy]
}
