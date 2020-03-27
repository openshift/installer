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


