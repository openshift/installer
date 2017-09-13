output "etcd_a_nodes" {
  value = "${aws_route53_record.etcd_a_nodes.*.fqdn}"
}

output "worker_nodes" {
  value = "${aws_route53_record.worker_nodes.*.fqdn}"
}

output "worker_nodes_public" {
  value = "${aws_route53_record.worker_nodes_public.*.fqdn}"
}

output "master_nodes" {
  value = "${aws_route53_record.master_nodes.*.fqdn}"
}
