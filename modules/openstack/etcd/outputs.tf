output "ip_addresses" {
  value = "${openstack_compute_instance_v2.etcd_node.*.access_ip_v4}"
}
