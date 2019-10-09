resource "aws_lb_target_group_attachment" "bootstrap" {
  count = "${var.instance_count * var.target_group_arns_length}"

  target_group_arn = "${var.target_group_arns[count.index]}"
  target_id        = "${var.ip_address}"

  availability_zone = "${var.availability_zone}"
}
