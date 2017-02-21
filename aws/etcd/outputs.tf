output "endpoints" {
  value = "${join(",",formatlist("http://%s:2379",aws_route53_record.etc_a_nodes.*.fqdn))}"
}
