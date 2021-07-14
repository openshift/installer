output "control_plane_ips" {
  value = <<EOT
    var.master_count > 2 ? [
      openstack_compute_instance_v2.master_conf_0.access_ip_v4,
      openstack_compute_instance_v2.master_conf_1.access_ip_v4,
      openstack_compute_instance_v2.master_conf_2.access_ip_v4 ] :
        var.master_count > 1 ? [
          openstack_compute_instance_v2.master_conf_0.access_ip_v4,
          openstack_compute_instance_v2.master_conf_1.access_ip_v4 ] :
            var.master_count > 0 ? [
              openstack_compute_instance_v2.master_conf_0.access_ip_v4 ] : 0
EOT
}
