data "aws_security_group" "worker" {
  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-worker-sg"
    },
    var.tags,
  )
}
