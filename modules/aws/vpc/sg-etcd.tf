resource "aws_security_group" "etcd" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags = "${merge(map(
      "Name", "${var.cluster_name}_etcd_sg",
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "tectonicClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    to_port     = 0
  }

  ingress {
    protocol  = "tcp"
    from_port = 0
    to_port   = 65535
    self      = true
  }

  ingress {
    protocol  = "tcp"
    from_port = 0
    to_port   = 65535
    security_groups = ["${aws_security_group.master.id}"]
  }
}
