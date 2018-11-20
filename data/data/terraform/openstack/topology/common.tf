locals {
  master_subnet_ids = ["${coalescelist(openstack_networking_port_v2.masters.*.id,var.external_master_subnet_ids)}"]
}
