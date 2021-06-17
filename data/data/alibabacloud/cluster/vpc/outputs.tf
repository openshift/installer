output "vpc_id" {
  value = alicloud_vpc.vpc.id
}

output "vswitch_ids" {
  value = alicloud_vswitch.vswitchs.*.id
}

output "gw_id" {
  value = alicloud_nat_gateway.nat_gateway.id
}

output "eip_id" {
  value = alicloud_eip_address.eip.id
}

output "eip_ip" {
  value = alicloud_eip_address.eip.ip_address
}

output "slb_ids" {
  value = [alicloud_slb_load_balancer.slb_external.id, alicloud_slb_load_balancer.slb_internal.id]
}

output "slb_external_ip" {
  value = alicloud_slb_load_balancer.slb_external.address
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
