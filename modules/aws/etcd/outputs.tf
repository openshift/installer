output "endpoints" {
  value = "${length(var.external_endpoints) == 0 ? join(",", aws_route53_record.etc_a_nodes.*.fqdn) : join(",", var.external_endpoints)}"
}
