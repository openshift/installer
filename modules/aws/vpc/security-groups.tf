resource "aws_security_group" "cluster_default" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = "${merge(map(
    "Name","${var.cluster_name}-sg-cluster_default",
    "KubernetesCluster", "${var.cluster_name}"
  ), var.extra_tags)}"
}
