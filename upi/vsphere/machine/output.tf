output "ip_addresses" {
  value = "${data.external.ip_address.*.result.ip_address}"
}
