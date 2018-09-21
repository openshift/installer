resource "aws_security_group" "api" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags = "${merge(map(
      "Name", "${var.cluster_name}_api_sg",
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "tectonicClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}

resource "aws_security_group_rule" "api_egress" {
  type              = "egress"
  security_group_id = "${aws_security_group.api.id}"

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "api_ingress_console" {
  type              = "ingress"
  security_group_id = "${aws_security_group.api.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 6443
  to_port     = 6443
}

resource "aws_security_group_rule" "mcs_ingress" {
  type              = "ingress"
  security_group_id = "${aws_security_group.api.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 49500
  to_port     = 49500
}

resource "aws_security_group" "console" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags = "${merge(map(
      "Name", "${var.cluster_name}_console_sg",
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "tectonicClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}

resource "aws_security_group_rule" "console_egress" {
  type              = "egress"
  security_group_id = "${aws_security_group.console.id}"

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "console_ingress_http" {
  type              = "ingress"
  security_group_id = "${aws_security_group.console.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 80
  to_port     = 80
}

resource "aws_security_group_rule" "console_ingress_https" {
  type              = "ingress"
  security_group_id = "${aws_security_group.console.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 443
  to_port     = 443
}
