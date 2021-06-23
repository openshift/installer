output "control_plane_ips" {
  value = module.masters.ip_addresses
}

output "lb_target_group_arns" {
  value = module.vpc.aws_lb_target_group_arns
}

output "lb_target_group_arns_length" {
  value = module.vpc.aws_lb_target_group_arns_length
}

output "vpc_id" {
  value = module.vpc.vpc_id
}

output "public_subnet_ids" {
  value = values(module.vpc.az_to_public_subnet_id)
}

output "private_subnet_ids" {
  value = values(module.vpc.az_to_private_subnet_id)
}

output "master_sg_id" {
  value = module.vpc.master_sg_id
}

output "ami_id" {
  value = var.aws_region == var.aws_ami_region ? var.aws_ami : aws_ami_copy.imported[0].id
}

output "bootstrap_instance_profile_name" {
  value = module.iam.bootstrap_instance_profile_name
}