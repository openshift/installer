output "etcd_names" {
  value = ["${concat(aws_route53_record.etcd_a_nodes.*.fqdn)}"]
}
