output "bootstrap_ip" {
  value = var.aws_public_ipv4_pool != "" ? aws_eip.bootstrap[0].public_ip : local.public_endpoints ? aws_instance.bootstrap.public_ip : aws_instance.bootstrap.private_ip
}
