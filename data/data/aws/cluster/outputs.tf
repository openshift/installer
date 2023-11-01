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

output "edge_public_subnet_ids" {
  value = values(module.vpc.az_to_edge_public_subnet_id)
}

output "edge_private_subnet_ids" {
  value = values(module.vpc.az_to_edge_private_subnet_id)
}

output "master_sg_id" {
  value = module.vpc.master_sg_id
}

output "ami_id" {
  value = var.aws_region == var.aws_ami_region ? var.aws_ami : aws_ami_copy.imported[0].id
}

output "aws_external_api_lb_dns_name" {
  value = module.vpc.aws_lb_api_external_dns_name
}

output "aws_internal_api_lb_dns_name" {
  value = module.vpc.aws_lb_api_internal_dns_name
}
