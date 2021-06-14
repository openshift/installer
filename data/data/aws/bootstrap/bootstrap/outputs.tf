output "bootstrap_ip" {
  value = local.public_endpoints ? aws_instance.bootstrap.public_ip : aws_instance.bootstrap.private_ip
}
