resource "google_compute_network" "cluster_network" {
  name = "${var.cluster_id}-network"

  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "compute_subnet" {
  name          = "${var.cluster_id}-compute-subnet"
  network       = google_compute_network.cluster_network.self_link
  ip_cidr_range = var.compute_subnet_cidr
}

resource "google_compute_subnetwork" "control_subnet" {
  name          = "${var.cluster_id}-control-subnet"
  network       = google_compute_network.cluster_network.self_link
  ip_cidr_range = var.control_subnet_cidr
}
