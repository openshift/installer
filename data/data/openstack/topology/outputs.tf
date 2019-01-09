output "lb_port_id" {
  value = "${openstack_networking_port_v2.lb_port.id}"
}

output "bootstrap_port_id" {
  value = "${openstack_networking_port_v2.bootstrap_port.id}"
}

output "controlplane_ips" {
  value = "${flatten(openstack_networking_port_v2.controlplane.*.all_fixed_ips)}"
}

output "controlplane_port_names" {
  value = "${openstack_networking_port_v2.controlplane.*.name}"
}

output "service_vm_fixed_ip" {
  value = "${openstack_networking_port_v2.lb_port.all_fixed_ips[0]}"
}

output "controlplane_sg_id" {
  value = "${openstack_networking_secgroup_v2.controlplane.id}"
}

output "controlplane_port_ids" {
  value = "${local.controlplane_port_ids}"
}
