data "external" "ping" {
  program = ["bash", "${path.root}/network/cidr_to_ip.sh"]

  query = {
    cidr = "${var.machine_cidr}"
    master_count = "${var.master_count}"
    worker_count = "${var.worker_count}"
  }
}
