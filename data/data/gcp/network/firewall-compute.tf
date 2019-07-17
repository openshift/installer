resource "google_compute_firewall" "worker_ingress_icmp" {
  name    = "${var.cluster_id}-worker-in-icmp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "icmp"
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_ssh" {
  name    = "${var.cluster_id}-worker-in-ssh"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_vxlan" {
  name    = "${var.cluster_id}-worker-in-vxlan"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_vxlan_from_master" {
  name    = "${var.cluster_id}-worker-in-vxlan-from-master"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_internal" {
  name    = "${var.cluster_id}-worker-in-internal"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_internal_from_master" {
  name    = "${var.cluster_id}-worker-in-internal-from-master"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_internal_udp" {
  name    = "${var.cluster_id}-worker-in-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_internal_from_master_udp" {
  name    = "${var.cluster_id}-worker-in-from-master-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_kubelet_insecure" {
  name    = "${var.cluster_id}-worker-in-kubelet-insecure"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_kubelet_insecure_from_master" {
  name    = "${var.cluster_id}-worker-in-kubelet-insecure-fr-mast"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_services_tcp" {
  name    = "${var.cluster_id}-worker-in-services-tcp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_services_udp" {
  name    = "${var.cluster_id}-worker-in-services-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}
