output "int_zone_id" {
  value = "${aws_route53_zone.tectonic-int.id}"
}

output "ingress_external_fqdn" {
  value = "${aws_route53_record.ingress-public.name}"
}

output "ingress_internal_fqdn" {
  value = "${aws_route53_record.ingress-public.name}"
}

output "api_external_fqdn" {
  value = "${aws_route53_record.api-external.name}"
}

output "api_internal_fqdn" {
  value = "${aws_route53_record.api-external.name}"
}
