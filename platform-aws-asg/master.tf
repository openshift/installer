resource "aws_autoscaling_group" "masters" {
  name                 = "${var.tectonic_cluster_name}-masters"
  desired_capacity     = "${var.tectonic_master_count}"
  max_size             = "${var.tectonic_master_count * 3}"
  min_size             = "1"
  launch_configuration = "${aws_launch_configuration.master_conf.id}"
  vpc_zone_identifier  = ["${aws_subnet.master_subnet.*.id}"]

  load_balancers = ["${aws_elb.api-internal.id}", "${aws_elb.api-external.id}"]

  tag {
    key                 = "Name"
    value               = "${var.tectonic_cluster_name}-master"
    propagate_at_launch = true
  }

  tag {
    key                 = "KubernetesCluster"
    value               = "${var.tectonic_cluster_name}"
    propagate_at_launch = true
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_launch_configuration" "master_conf" {
  instance_type        = "${var.tectonic_aws_master_ec2_type}"
  image_id             = "${data.aws_ami.coreos_ami.image_id}"
  name_prefix          = "${var.tectonic_cluster_name}-master-"
  key_name             = "${aws_key_pair.ssh-key.key_name}"
  security_groups      = ["${aws_security_group.master_sec_group.id}", "${aws_security_group.cluster_default.id}"]
  iam_instance_profile = "${aws_iam_instance_profile.master_profile.arn}"
  user_data            = "${ignition_config.master.rendered}"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "master_sec_group" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags {
    Name = "${var.tectonic_cluster_name}_master_sg"
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
    from_port   = 443
    to_port     = 443
  }

  ingress {
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 10255
    to_port     = 10255
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_iam_instance_profile" "master_profile" {
  name  = "${var.tectonic_cluster_name}-master-profile"
  roles = ["${aws_iam_role.master_role.name}"]
}

resource "aws_iam_role" "master_role" {
  name = "${var.tectonic_cluster_name}-master-role"
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

resource "aws_iam_role_policy" "master_policy" {
  name = "${var.tectonic_cluster_name}_master_policy"
  role = "${aws_iam_role.master_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "ec2:*",
      "Resource": "*",
      "Effect": "Allow"
    },
    {
      "Action": "elasticloadbalancing:*",
      "Resource": "*",
      "Effect": "Allow"
    },
    {
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:GetRepositoryPolicy",
        "ecr:DescribeRepositories",
        "ecr:ListImages",
        "ecr:BatchGetImage"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}
