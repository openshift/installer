resource "aws_security_group" "etcd_sec_group" {
  vpc_id = "${var.vpc_id}"
  count  = "${length(var.external_endpoints) == 0 ? 1 : 0}"

  tags {
    Name              = "${var.tectonic_cluster_name}_etcd_sg"
    KubernetesCluster = "${var.tectonic_cluster_name}"
  }

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }

  ingress {
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 22
    to_port     = 22
  }

  ingress {
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 2379
    to_port     = 2379
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
    cidr_blocks = ["0.0.0.0/0"]
  }
}
