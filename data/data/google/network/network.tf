resource "google_compute_network" "default" {
  name = "${var.cluster_name}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "${var.cluster_name}-default-${var.region}"
  region        = "${var.region}"
  network       = "${google_compute_network.default.self_link}"
  ip_cidr_range = "${var.cidr_block}"
}
