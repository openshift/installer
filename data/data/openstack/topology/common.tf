locals {
  control_plane_port_ids = ["${coalescelist(openstack_networking_trunk_v2.control_plane.*.port_id,openstack_networking_port_v2.control_plane.*.id)}"]
}
