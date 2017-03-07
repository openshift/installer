resource "aws_instance" "etcd_node" {
  count                  = "${var.node_count}"
  ami                    = "${var.coreos_ami}"
  instance_type          = "t2.medium"
  subnet_id              = "${var.etcd_subnets[count.index]}"
  key_name               = "${var.ssh_key}"
  user_data              = "${ignition_config.etcd.*.rendered[count.index]}"
  vpc_security_group_ids = ["${aws_security_group.etcd_sec_group.id}"]

  tags {
    Name = "${var.cluster_name}-etcd-${count.index}"
    KubernetesCluster = "${var.cluster_name}"
  }
}
