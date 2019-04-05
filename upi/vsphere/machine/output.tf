output "ip_addresses" {
  value = ["${local.ip_addresses}"]
}

output "ips_exist" {
  value = "${local.ips_exist}"
}
