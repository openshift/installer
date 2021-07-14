locals {
  master_port_ids = coalescelist(
    openstack_networking_trunk_v2.masters.*.port_id,
    openstack_networking_port_v2.masters.*.id,
  )
  master_sg_ids = concat(
    var.openstack_master_extra_sg_ids,
    [openstack_networking_secgroup_v2.master.id],
  )
}

