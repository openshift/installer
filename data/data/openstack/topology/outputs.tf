output "master_sg_id" {
  value = openstack_networking_secgroup_v2.master.id
}

output "master_port_ids" {
  value = local.master_port_ids
}

output "private_network_id" {
  value = openstack_networking_network_v2.openshift-private.id
}

output "nodes_subnet_id" {
  value = openstack_networking_subnet_v2.nodes.id
}

