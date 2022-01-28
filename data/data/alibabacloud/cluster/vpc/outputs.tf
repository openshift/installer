output "vpc_id" {
  value = local.vpc_id
}

output "vswitch_ids" {
  value = local.vswitch_ids
}

output "az_to_vswitch_id" {
  value = zipmap(data.alicloud_vswitches.vswitches.vswitches.*.zone_id, data.alicloud_vswitches.vswitches.vswitches.*.id)
}

output "slb_ids" {
  value = concat(alicloud_slb_load_balancer.slb_external[*].id, [alicloud_slb_load_balancer.slb_internal.id])
}

output "slb_group_length" {
  value = length(concat(alicloud_slb_load_balancer.slb_external[*].id, [alicloud_slb_load_balancer.slb_internal.id]))
}

output "slb_external_ip" {
  value = local.is_external ? alicloud_slb_load_balancer.slb_external[0].address : null
}

output "slb_internal_ip" {
  value = alicloud_slb_load_balancer.slb_internal.address
}

output "sg_master_id" {
  value = alicloud_security_group.sg_master.id
}

output "sg_worker_id" {
  value = alicloud_security_group.sg_worker.id
}
