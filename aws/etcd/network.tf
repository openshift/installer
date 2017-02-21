data "aws_vpc" "etcd_vpc" {
  id = "${var.vpc_id}"
}

data "aws_subnet" "az_subnet" {
  count  = "${var.node_count}"
  vpc_id = "${data.aws_vpc.etcd_vpc.id}"

  filter = {
    name   = "availabilityZone"
    values = ["${data.aws_availability_zones.zones.names[count.index]}"]
  }
}

resource "aws_default_security_group" "default_sec_group" {
  vpc_id = "${data.aws_vpc.etcd_vpc.id}"

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

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
    cidr_blocks = ["0.0.0.0/0"]
  }
}
