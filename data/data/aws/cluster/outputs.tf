output "control_plane_ips" {
  value = module.masters.ip_addresses
}

output "lb_target_group_arns" {
  value = var.aws_lb_target_group_arns
}

output "lb_target_group_arns_length" {
  value = var.aws_lb_target_group_arns_length
}
