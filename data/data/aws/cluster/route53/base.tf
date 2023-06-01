locals {

  // Because of the issue https://github.com/hashicorp/terraform/issues/12570, the consumers cannot count 0/1
  // based on if api_external_lb_dns_name for example, which will be null when there is no external lb for API.
  // So publish_strategy serves an coordinated proxy for that decision.
  public_endpoints = var.publish_strategy == "External" ? true : false

  use_cname = contains(["us-gov-west-1", "us-gov-east-1", "us-iso-east-1"], var.region)
  use_alias = ! local.use_cname
}

provider "aws" {
  alias = "private_hosted_zone"

  assume_role {
    role_arn = var.internal_zone_role
  }

  region = var.region

  skip_region_validation = true

  endpoints {
    ec2     = lookup(var.custom_endpoints, "ec2", null)
    elb     = lookup(var.custom_endpoints, "elasticloadbalancing", null)
    iam     = lookup(var.custom_endpoints, "iam", null)
    route53 = lookup(var.custom_endpoints, "route53", null)
    s3      = lookup(var.custom_endpoints, "s3", null)
    sts     = lookup(var.custom_endpoints, "sts", null)
  }
}

data "aws_route53_zone" "public" {
  count = local.public_endpoints ? 1 : 0

  name = var.base_domain

  depends_on = [aws_route53_record.api_external_internal_zone_alias, aws_route53_record.api_external_internal_zone_cname]
}

data "aws_route53_zone" "int" {
  provider = aws.private_hosted_zone

  zone_id = var.internal_zone == null ? aws_route53_zone.new_int[0].id : var.internal_zone
}

resource "aws_route53_zone" "new_int" {
  count = var.internal_zone == null ? 1 : 0

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
}

resource "aws_route53_record" "api_external_alias" {
  count = local.use_alias && local.public_endpoints ? 1 : 0

  zone_id = data.aws_route53_zone.public[0].zone_id
  name    = "api.${var.cluster_domain}"
  type    = "A"

  alias {
    name                   = var.api_external_lb_dns_name
    zone_id                = var.api_external_lb_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api_internal_alias" {
  provider = aws.private_hosted_zone
  count    = local.use_alias ? 1 : 0

  zone_id = data.aws_route53_zone.int.zone_id
  name    = "api-int.${var.cluster_domain}"
  type    = "A"

  alias {
    name                   = var.api_internal_lb_dns_name
    zone_id                = var.api_internal_lb_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api_external_internal_zone_alias" {
  provider = aws.private_hosted_zone
  count    = local.use_alias ? 1 : 0

  zone_id = data.aws_route53_zone.int.zone_id
  name    = "api.${var.cluster_domain}"
  type    = "A"

  alias {
    name                   = var.api_internal_lb_dns_name
    zone_id                = var.api_internal_lb_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api_external_cname" {
  count = local.use_cname && local.public_endpoints ? 1 : 0

  zone_id = data.aws_route53_zone.public[0].zone_id
  name    = "api.${var.cluster_domain}"
  type    = "CNAME"
  ttl     = 10

  records = [var.api_external_lb_dns_name]
}

resource "aws_route53_record" "api_internal_cname" {
  provider = aws.private_hosted_zone
  count    = local.use_cname ? 1 : 0

  zone_id = data.aws_route53_zone.int.zone_id
  name    = "api-int.${var.cluster_domain}"
  type    = "CNAME"
  ttl     = 10

  records = [var.api_internal_lb_dns_name]
}

resource "aws_route53_record" "api_external_internal_zone_cname" {
  provider = aws.private_hosted_zone
  count    = local.use_cname ? 1 : 0

  zone_id = data.aws_route53_zone.int.zone_id
  name    = "api.${var.cluster_domain}"
  type    = "CNAME"
  ttl     = 10

  records = [var.api_internal_lb_dns_name]
}
