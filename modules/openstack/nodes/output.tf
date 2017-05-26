output "user_data" {
  value = ["${data.ignition_config.node.*.rendered}"]
}

output "secgroup_master_name" {
  value = "${openstack_compute_secgroup_v2.master.name}"
}

output "secgroup_master_id" {
  value = "${openstack_compute_secgroup_v2.master.id}"
}

output "secgroup_node_name" {
  value = "${openstack_compute_secgroup_v2.node.name}"
}

output "secgroup_node_id" {
  value = "${openstack_compute_secgroup_v2.node.id}"
}

output "secgroup_self_hosted_etcd_name" {
  value = "${openstack_compute_secgroup_v2.self_hosted_etcd.name}"
}

output "secgroup_self_hosted_etcd_id" {
  value = "${openstack_compute_secgroup_v2.self_hosted_etcd.id}"
}
