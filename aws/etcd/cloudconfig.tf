data "template_file" "userdata" {
  count    = "${var.node_count}"
  template = "${file("${path.module}/userdata.yaml")}"

  vars {
    node_name   = "node-${count.index}.${var.etcd_domain}"
    etcd_domain = "${var.etcd_domain}"
  }
}
