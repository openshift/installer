resource "aws_route53_record" "a_record" {
  for_each = var.records

  type    = "A"
  ttl     = "60"
  zone_id = var.zone_id
  name    = each.key
  records = [each.value]
}
