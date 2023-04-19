output "control_plane_ips" {
  value = concat(
    openstack_compute_instance_v2.master_conf_0.*.access_ip_v4,
    openstack_compute_instance_v2.master_conf_1.*.access_ip_v4,
    openstack_compute_instance_v2.master_conf_2.*.access_ip_v4,
  )
}

output "master_sg_ids" {
  value = concat(
    var.openstack_master_extra_sg_ids,
    [openstack_networking_secgroup_v2.master.id],
  )
}

output "master_port_ids" {
  value = local.master_port_ids
}

output "private_network_id" {
  value = local.nodes_default_port.network_id
}

output "nodes_default_port" {
  value = local.nodes_default_port
}
