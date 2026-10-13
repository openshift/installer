output "fqdns" {
  value = values(aws_route53_record.a_record)[*].name
}
