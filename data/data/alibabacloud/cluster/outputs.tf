output "vpc_id" {
  value = module.vpc.vpc_id
}

output "vswitch_ids" {
  value = module.vpc.vswitch_ids
}

output "slb_ids" {
  value = module.vpc.slb_ids
}

output "sg_master_id" {
  value = module.vpc.sg_master_id
}
