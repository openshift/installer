locals {
  controlplane_port_ids = ["${coalescelist(openstack_networking_trunk_v2.controlplane.*.port_id,openstack_networking_port_v2.controlplane.*.id)}"]
}
