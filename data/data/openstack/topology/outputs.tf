output "master_sg_id" {
  value = var.openstack_disable_sg ? null : openstack_networking_secgroup_v2.master[0].id
}

output "master_port_ids" {
  value = local.master_port_ids
}

output "private_network_id" {
  value = local.nodes_network_id
}

output "nodes_subnet_id" {
  value = local.nodes_subnet_id
}

