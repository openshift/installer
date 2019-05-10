data "aws_route53_zone" "public" {
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
  zone_id = data.aws_route53_zone.public.zone_id
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
  count   = var.etcd_count
  type    = "A"
  ttl     = "60"
  zone_id = aws_route53_zone.int.zone_id
  name    = "etcd-${count.index}.${var.cluster_domain}"
  # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
  # force an interpolation expression to be interpreted as a list by wrapping it
  # in an extra set of list brackets. That form was supported for compatibilty in
  # v0.11, but is no longer supported in Terraform v0.12.
  #
  # If the expression in the following list itself returns a list, remove the
  # brackets to avoid interpretation as a list of lists. If the expression
  # returns a single list item then leave it as-is and remove this TODO comment.
  records = [var.etcd_ip_addresses[count.index]]
}

resource "aws_route53_record" "etcd_cluster" {
  type    = "SRV"
  ttl     = "60"
  zone_id = aws_route53_zone.int.zone_id
  name    = "_etcd-server-ssl._tcp"
  records = formatlist("0 10 2380 %s", aws_route53_record.etcd_a_nodes.*.fqdn)
}

