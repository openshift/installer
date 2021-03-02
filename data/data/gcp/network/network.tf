resource "google_compute_network" "cluster_network" {
  count       = var.preexisting_network ? 0 : 1
  description = local.description

  name = var.cluster_network

  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "worker_subnet" {
  count       = var.preexisting_network ? 0 : 1
  description = local.description

  name          = var.worker_subnet
  network       = google_compute_network.cluster_network[0].self_link
  ip_cidr_range = var.worker_subnet_cidr
}

resource "google_compute_subnetwork" "master_subnet" {
  count       = var.preexisting_network ? 0 : 1
  description = local.description

  name          = var.master_subnet
  network       = google_compute_network.cluster_network[0].self_link
  ip_cidr_range = var.master_subnet_cidr
}

resource "google_compute_router" "router" {
  count       = var.preexisting_network ? 0 : 1
  description = local.description

  name    = "${var.cluster_id}-router"
  network = google_compute_network.cluster_network[0].self_link
}

resource "google_compute_router_nat" "master_nat" {
  count = var.preexisting_network ? 0 : 1

  name                               = "${var.cluster_id}-nat-master"
  router                             = google_compute_router.router[0].name
  nat_ip_allocate_option             = "AUTO_ONLY"
  min_ports_per_vm                   = 7168
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = google_compute_subnetwork.master_subnet[0].self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}

resource "google_compute_router_nat" "worker_nat" {
  count = var.preexisting_network ? 0 : 1

  name                               = "${var.cluster_id}-nat-worker"
  router                             = google_compute_router.router[0].name
  nat_ip_allocate_option             = "AUTO_ONLY"
  min_ports_per_vm                   = 4096
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = google_compute_subnetwork.worker_subnet[0].self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}
