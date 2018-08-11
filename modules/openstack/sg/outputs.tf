output "api_sg_id" {
  value = "${openstack_networking_secgroup_v2.api.id}"
}

output "console_sg_id" {
  value = "${openstack_networking_secgroup_v2.console.id}"
}

output "etcd_sg_id" {
  value = "${openstack_networking_secgroup_v2.etcd.id}"
}

output "master_sg_id" {
  value = "${openstack_networking_secgroup_v2.master.id}"
}

output "worker_sg_id" {
  value = "${openstack_networking_secgroup_v2.worker.id}"
}
