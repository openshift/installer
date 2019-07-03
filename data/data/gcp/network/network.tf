resource "google_compute_network" "cluster_network" {
  name = "${var.cluster_id}-network"

  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "worker_subnet" {
  name          = "${var.cluster_id}-worker-subnet"
  network       = google_compute_network.cluster_network.self_link
  ip_cidr_range = var.worker_subnet_cidr
}

resource "google_compute_subnetwork" "master_subnet" {
  name          = "${var.cluster_id}-master-subnet"
  network       = google_compute_network.cluster_network.self_link
  ip_cidr_range = var.master_subnet_cidr
}
