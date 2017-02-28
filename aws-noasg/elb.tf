resource "aws_elb" "api-external" {
  name            = "${var.cluster_name}-api-external"
  subnets         = ["${aws_subnet.az_subnet_pub.*.id}"]
  internal        = false
  security_groups = ["${aws_security_group.master_sec_group.id}"]

  listener {
    instance_port     = 22
    instance_protocol = "tcp"
    lb_port           = 22
    lb_protocol       = "tcp"
  }

  listener {
    instance_port     = 443
    instance_protocol = "tcp"
    lb_port           = 443
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "TCP:22"
    interval            = 5
  }
}

resource "aws_elb" "console" {
  name            = "${var.cluster_name}-console"
  subnets         = ["${aws_subnet.az_subnet_pub.*.id}"]
  internal        = false
  security_groups = ["${aws_security_group.master_sec_group.id}"]

  listener {
    instance_port     = 32001
    instance_protocol = "tcp"
    lb_port           = 80
    lb_protocol       = "tcp"
  }

  listener {
    instance_port     = 32000
    instance_protocol = "tcp"
    lb_port           = 443
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:32002/healthz"
    interval            = 5
  }
}

resource "aws_route53_record" "api-internal" {
  zone_id = "${aws_route53_zone.tectonic-int.zone_id}"
  name    = "${var.cluster_name}-k8s.${var.tectonic_domain}"
  type    = "A"
  ttl     = "60"

  records = [
    "${aws_instance.master-node.*.private_ip}",
  ]
}

resource "aws_route53_record" "api-external" {
  zone_id = "${data.aws_route53_zone.tectonic-ext.zone_id}"
  name    = "${var.cluster_name}-k8s.${var.tectonic_domain}"
  type    = "A"

  alias {
    name                   = "${aws_elb.api-external.dns_name}"
    zone_id                = "${aws_elb.api-external.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "ingress-public" {
  zone_id = "${data.aws_route53_zone.tectonic-ext.zone_id}"
  name    = "${var.cluster_name}.${var.tectonic_domain}"
  type    = "A"

  alias {
    name                   = "${aws_elb.console.dns_name}"
    zone_id                = "${aws_elb.console.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "ingress-private" {
  zone_id = "${aws_route53_zone.tectonic-int.zone_id}"
  name    = "${var.cluster_name}.${var.tectonic_domain}"
  type    = "A"

  # ttl     = "60"
  # records = ["${aws_instance.worker-node.*.private_ip}"]
  alias {
    name                   = "${aws_elb.console.dns_name}"
    zone_id                = "${aws_elb.console.zone_id}"
    evaluate_target_health = true
  }
}
