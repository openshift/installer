output "ip_addresses" {
  value = [aws_instance.bootstrap.private_ip]
}

