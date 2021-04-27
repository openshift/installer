data "aws_security_group" "master" {
  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-master-sg"
    },
    var.tags,
  )
}
