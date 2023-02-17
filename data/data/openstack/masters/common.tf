locals {
  master_port_ids = coalescelist(
    openstack_networking_trunk_v2.masters.*.port_id,
    openstack_networking_port_v2.masters.*.id,
  )
  master_sg_ids = concat(
    var.openstack_master_extra_sg_ids,
    [openstack_networking_secgroup_v2.master.id],
  )
  master_failuredomain_0 = jsondecode(var.openstack_master_ports_json[0])
  master_failuredomain_1 = jsondecode(var.openstack_master_ports_json[1])
  master_failuredomain_2 = jsondecode(var.openstack_master_ports_json[2])

  master_failuredomain_0_ports = local.master_failuredomain_0.ports
  master_failuredomain_1_ports = local.master_failuredomain_1.ports
  master_failuredomain_2_ports = local.master_failuredomain_2.ports

  master_failuredomain_0_trunks = local.master_failuredomain_0.trunks
  master_failuredomain_1_trunks = local.master_failuredomain_1.trunks
  master_failuredomain_2_trunks = local.master_failuredomain_2.trunks
}
