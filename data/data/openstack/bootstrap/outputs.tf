output "bootstrap_ip" {
  value = var.openstack_external_network != "" ? openstack_networking_floatingip_v2.bootstrap_fip[0].address : openstack_compute_instance_v2.bootstrap.access_ip_v4
}

