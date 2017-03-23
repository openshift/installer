output "ingress_external_fqdn" {
  value = "${aws_route53_record.ingress-public.name}"
}

output "ingress_internal_fqdn" {
  value = "${aws_route53_record.ingress-private.name}"
}

output "api_external_fqdn" {
  value = "${aws_route53_record.api-external.name}"
}

output "api_internal_fqdn" {
  value = "${aws_route53_record.api-internal.name}"
}
