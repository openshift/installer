output "bootstrap_floating_ip" {
  value = openstack_networking_floatingip_v2.bootstrap_fip.address
}
