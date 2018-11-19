locals {
  arn = "aws"
}

resource "aws_instance" "master" {
  count = "${var.instance_count}"
  ami   = "${var.ec2_ami}"

  instance_type = "${var.ec2_type}"
  subnet_id     = "${element(var.subnet_ids, count.index)}"
  user_data     = "${var.user_data_ign}"

  vpc_security_group_ids      = ["${var.master_sg_ids}"]
  associate_public_ip_address = "${var.public_endpoints}"

  lifecycle {
    # Ignore changes in the AMI which force recreation of the resource. This
    # avoids accidental deletion of nodes whenever a new CoreOS Release comes
    # out.
    ignore_changes = ["ami"]
  }

  tags = "${merge(map(
      "Name", "${var.cluster_name}-master-${count.index}",
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "tectonicClusterID", "${var.cluster_id}",
      "clusterid", "${var.cluster_name}"
    ), var.extra_tags)}"

  root_block_device {
    volume_type = "${var.root_volume_type}"
    volume_size = "${var.root_volume_size}"
    iops        = "${var.root_volume_type == "io1" ? var.root_volume_iops : 0}"
  }

  volume_tags = "${merge(map(
    "Name", "${var.cluster_name}-master-${count.index}-vol",
    "kubernetes.io/cluster/${var.cluster_name}", "owned",
    "tectonicClusterID", "${var.cluster_id}"
  ), var.extra_tags)}"
}

resource "aws_lb_target_group_attachment" "master" {
  count = "${var.instance_count * var.target_group_arns_length}"

  target_group_arn = "${var.target_group_arns[count.index % var.target_group_arns_length]}"
  target_id        = "${aws_instance.master.*.private_ip[count.index / var.target_group_arns_length]}"
}
