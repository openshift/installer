output "endpoints" {
  value = ["${length(var.external_endpoints) == 0 ? aws_route53_record.etc_a_nodes.*.fqdn : var.external_endpoints}"]
}
