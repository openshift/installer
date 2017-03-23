# We have to do this join() & split() 'trick' because the ternary operator can't output lists.
output "endpoints" {
  value = ["${split(",", join(",", var.external_endpoints) == "" ? join(",", aws_route53_record.etc_a_nodes.*.fqdn) :  join(",", var.external_endpoints))}"]
}