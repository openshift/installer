output "vpc_id" {
  value = "${data.aws_vpc.cluster_vpc.id}"
}

output "master_subnet_ids" {
  value = "${local.master_subnet_ids}"
}

output "etcd_sg_id" {
  value = "${aws_security_group.etcd.id}"
}

output "master_sg_id" {
  value = "${aws_security_group.master.id}"
}

output "worker_sg_id" {
  value = "${aws_security_group.worker.id}"
}

output "api_sg_id" {
  value = "${aws_security_group.api.id}"
}

output "console_sg_id" {
  value = "${aws_security_group.console.id}"
}

output "aws_lb_target_group_arns" {
  value = "${compact(concat(aws_lb_target_group.api_internal.*.arn, aws_lb_target_group.services.*.arn, aws_lb_target_group.api_external.*.arn))}"
}

output "aws_lb_target_group_arns_length" {
  // 2 for private endpoints and 1 for public endpoints
  value = "3"
}

output "aws_lb_api_external_dns_name" {
  value = "${aws_lb.api_external.dns_name}"
}

output "aws_lb_api_external_zone_id" {
  value = "${aws_lb.api_external.zone_id}"
}

output "aws_lb_api_internal_dns_name" {
  value = "${aws_lb.api_internal.dns_name}"
}

output "aws_lb_api_internal_zone_id" {
  value = "${aws_lb.api_internal.zone_id}"
}
