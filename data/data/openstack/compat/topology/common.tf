locals {
  master_port_ids = coalescelist(
    openstack_networking_trunk_v2.masters.*.port_id,
    openstack_networking_port_v2.masters.*.id,
  )
  description = "Created By OpenShift Installer"
}

