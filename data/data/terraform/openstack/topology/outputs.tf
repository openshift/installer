output "bootstrap_port_id" {
  value = "${openstack_networking_port_v2.bootstrap_port.id}"
}

output "master_sg_id" {
  value = "${openstack_networking_secgroup_v2.master.id}"
}

output "master_subnet_ids" {
  value = "${local.master_subnet_ids}"
}
