resource "google_compute_firewall" "compute_ingress_icmp" {
  name    = "${var.cluster_id}-compute-ingress-icmp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "icmp"
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_ssh" {
  name    = "${var.cluster_id}-compute-ingress-ssh"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_vxlan" {
  name    = "${var.cluster_id}-compute-ingress-vxlan"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_vxlan_from_master" {
  name    = "${var.cluster_id}-compute-ingress-vxlan-from-master"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_internal" {
  name    = "${var.cluster_id}-compute-ingress-internal"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_internal_from_master" {
  name    = "${var.cluster_id}-compute-ingress-internal-from-master"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_internal_udp" {
  name    = "${var.cluster_id}-compute-ingress-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_internal_from_master_udp" {
  name    = "${var.cluster_id}-compute-ingress-udp-from-master"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_kubelet_insecure" {
  name    = "${var.cluster_id}-compute-ingress-kubelet-insecure"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_kubelet_insecure_from_master" {
  name    = "${var.cluster_id}-compute-ingress-kubelet-insecure-from-master"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_services_tcp" {
  name    = "${var.cluster_id}-compute-ingress-services-tcp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "compute_ingress_services_udp" {
  name    = "${var.cluster_id}-compute-ingress-services-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}
