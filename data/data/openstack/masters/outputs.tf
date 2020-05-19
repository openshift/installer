output "master_ips" {
  value = openstack_compute_instance_v2.master_conf.*.access_ip_v4
}
