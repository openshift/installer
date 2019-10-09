output "aws_lb_target_group_arns" {
  value = "${compact(concat(aws_lb_target_group.ingress_external_http.*.arn,aws_lb_target_group.ingress_external_https.*.arn, aws_lb_target_group.api_internal.*.arn, aws_lb_target_group.services.*.arn, aws_lb_target_group.api_external.*.arn))}"
}

output "aws_lb_target_group_arns_length" {
  // 2 for private endpoints and 1 for public endpoints
  // 2 ingress endpoints
  value = "5"
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
