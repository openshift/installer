output "user_data" {
  value = ["${data.ignition_config.etcd.*.rendered}"]
}

output "secgroup_name" {
  value = "${openstack_compute_secgroup_v2.etcd.name}"
}
