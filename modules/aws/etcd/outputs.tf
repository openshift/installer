output "endpoints" {
  value = ["${aws_route53_record.etc_a_nodes.*.fqdn}"]
}
