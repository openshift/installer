data "external" "ping" {
  program = ["bash", "${path.root}/network/cidr_to_ip.sh"]

  query = {
    cidr                = "${var.machine_cidr}"
    control_plane_count = "${var.control_plane_count}"
    compute_count       = "${var.compute_count}"
    cluster_domain      = "${var.cluster_domain}"
    ipam                = "${var.ipam}"
    ipam_token          = "${var.ipam_token}"
  }
}
