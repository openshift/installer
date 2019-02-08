output "service_port_id" {
  value = "${openstack_networking_port_v2.service_port.id}"
}

output "bootstrap_port_id" {
  value = "${openstack_networking_port_v2.bootstrap_port.id}"
}

output "control_plane_ips" {
  value = "${flatten(openstack_networking_port_v2.control_plane.*.all_fixed_ips)}"
}

output "control_plane_port_names" {
  value = "${openstack_networking_port_v2.control_plane.*.name}"
}

output "service_vm_fixed_ip" {
  value = "${openstack_networking_port_v2.service_port.all_fixed_ips[0]}"
}

output "control_plane_sg_id" {
  value = "${openstack_networking_secgroup_v2.control_plane.id}"
}

output "control_plane_port_ids" {
  value = "${local.control_plane_port_ids}"
}
