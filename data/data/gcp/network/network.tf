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

resource "google_compute_router" "router" {
  name    = "${var.cluster_id}-router"
  network = google_compute_network.cluster_network.self_link
}

resource "google_compute_router_nat" "master_nat" {
  name                               = "${var.cluster_id}-nat-master"
  router                             = google_compute_router.router.name
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = google_compute_subnetwork.master_subnet.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}

resource "google_compute_router_nat" "worker_nat" {
  name                               = "${var.cluster_id}-nat-worker"
  router                             = google_compute_router.router.name
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = google_compute_subnetwork.worker_subnet.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}
