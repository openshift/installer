data "aws_route53_zone" "base" {
  name = var.base_domain
}

resource "aws_route53_zone" "cluster" {
  name          = var.cluster_domain
  force_destroy = true

  tags = {
    "Name"     = var.cluster_domain
    "Platform" = "vSphere"
  }
}

resource "aws_route53_record" "name_server" {
  name    = var.cluster_domain
  type    = "NS"
  ttl     = "300"
  zone_id = data.aws_route53_zone.base.zone_id
  records = aws_route53_zone.cluster.name_servers
}

