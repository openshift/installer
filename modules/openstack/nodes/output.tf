output "user_data" {
  value = ["${data.ignition_config.node.*.rendered}"]
}

output "secgroup_name" {
  value = "${openstack_compute_secgroup_v2.node.name}"
}

output "secgroup_id" {
  value = "${openstack_compute_secgroup_v2.node.id}"
}
