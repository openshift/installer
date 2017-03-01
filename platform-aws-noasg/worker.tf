resource "aws_instance" "worker-node" {
  count                  = "${var.worker_count}"
  instance_type          = "${var.worker_ec2_type}"
  ami                    = "${data.aws_ami.coreos_ami.image_id}"
  key_name               = "${aws_key_pair.ssh-key.key_name}"
  vpc_security_group_ids = ["${aws_security_group.worker_sec_group.id}", "${aws_security_group.cluster_default.id}"]
  source_dest_check      = false
  iam_instance_profile   = "${aws_iam_instance_profile.worker_profile.id}"
  user_data              = "${ignition_config.worker.rendered}"
  subnet_id              = "${aws_subnet.worker_subnet.*.id[count.index % var.az_count]}"

  tags {
    Name              = "${var.cluster_name}-worker-${count.index}"
    KubernetesCluster = "${var.cluster_name}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "worker_sec_group" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags {
    Name = "${var.cluster_name}_worker_sg"
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
    from_port   = 443
    to_port     = 443
  }

  ingress {
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 10255
    to_port     = 10255
  }

  ingress {
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 10250
    to_port     = 10250
  }

  ingress {
    protocol    = "tcp"
    self        = true
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 30000
    to_port     = 32767
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_iam_instance_profile" "worker_profile" {
  name  = "${var.cluster_name}-worker-profile"
  roles = ["${aws_iam_role.worker_role.name}"]
}

resource "aws_iam_role" "worker_role" {
  name = "${var.cluster_name}-worker-role"
  path = "/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}

resource "aws_iam_role_policy" "worker_policy" {
  name = "${var.cluster_name}_worker_policy"
  role = "${aws_iam_role.worker_role.id}"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:Describe*"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "ec2:AttachVolume",
                "ec2:DetachVolume"
            ],
            "Resource": [
                "arn:aws:ec2:*:*:instance/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:*"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
EOF
}

resource "aws_elb_attachment" "console" {
  count    = "${var.worker_count}"
  elb      = "${aws_elb.console.id}"
  instance = "${aws_instance.worker-node.*.id[count.index]}"
}
