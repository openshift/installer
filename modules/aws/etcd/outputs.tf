# We have to do this join() & split() 'trick' because the ternary operator can't output lists.
output "endpoints" {
  value = ["${split(",", length(var.external_endpoints) == 0 ? join(",", aws_route53_record.etc_a_nodes.*.fqdn) : join(",", var.external_endpoints))}"]
}
