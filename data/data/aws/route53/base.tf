locals {

  // Because of the issue https://github.com/hashicorp/terraform/issues/12570, the consumers cannot count 0/1
  // based on if api_external_lb_dns_name for example, which will be null when there is no external lb for API.
  // So publish_strategy serves an coordinated proxy for that decision.
  public_endpoints = var.publish_strategy == "External" ? true : false
}

data "aws_route53_zone" "public" {
  count = local.public_endpoints ? 1 : 0

  name = var.base_domain
}

resource "aws_route53_zone" "int" {
  name          = var.cluster_domain
  force_destroy = true

  vpc {
    vpc_id = var.vpc_id
  }

  tags = merge(
    {
      "Name" = "${var.cluster_id}-int"
    },
    var.tags,
  )

  depends_on = [aws_route53_record.api_external]
}

resource "aws_route53_record" "api_external" {
  count = local.public_endpoints ? 1 : 0

  zone_id = data.aws_route53_zone.public[0].zone_id
  name    = "api.${var.cluster_domain}"
  type    = "A"

  alias {
    name                   = var.api_external_lb_dns_name
    zone_id                = var.api_external_lb_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api_internal" {
  zone_id = aws_route53_zone.int.zone_id
  name    = "api-int.${var.cluster_domain}"
  type    = "A"

  alias {
    name                   = var.api_internal_lb_dns_name
    zone_id                = var.api_internal_lb_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api_external_internal_zone" {
  zone_id = aws_route53_zone.int.zone_id
  name    = "api.${var.cluster_domain}"
  type    = "A"

  alias {
    name                   = var.api_internal_lb_dns_name
    zone_id                = var.api_internal_lb_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "etcd_a_nodes" {
  count   = var.use_ipv6 == false ? var.etcd_count : 0
  type    = "A"
  ttl     = "60"
  zone_id = aws_route53_zone.int.zone_id
  name    = "etcd-${count.index}.${var.cluster_domain}"
  records = [var.etcd_ip_addresses[count.index]]
}

resource "aws_route53_record" "etcd_aaaa_nodes" {
  count   = var.use_ipv6 == true ? var.etcd_count : 0
  type    = "AAAA"
  ttl     = "60"
  zone_id = aws_route53_zone.int.zone_id
  name    = "etcd-${count.index}.${var.cluster_domain}"
  records = [var.etcd_ipv6_addresses[count.index]]
}

resource "aws_route53_record" "etcd_cluster" {
  type    = "SRV"
  ttl     = "60"
  zone_id = aws_route53_zone.int.zone_id
  name    = "_etcd-server-ssl._tcp"
  records = formatlist("0 10 2380 %s", var.use_ipv6 == false ? aws_route53_record.etcd_a_nodes.*.fqdn : aws_route53_record.etcd_aaaa_nodes.*.fqdn)
}

