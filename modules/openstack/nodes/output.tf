output "master_user_data" {
  value = ["${ignition_config.master.*.rendered}"]
}

output "worker_user_data" {
  value = ["${ignition_config.worker.*.rendered}"]
}

output "master_secgroup_name" {
  value = "${openstack_compute_secgroup_v2.master.name}"
}

output "worker_secgroup_name" {
  value = "${openstack_compute_secgroup_v2.worker.name}"
}
