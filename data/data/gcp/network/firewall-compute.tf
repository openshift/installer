resource "google_compute_firewall" "worker_ingress_icmp" {
  name    = "${var.cluster_id}-worker-in-icmp"
  network = local.cluster_network

  allow {
    protocol = "icmp"
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_ssh" {
  name    = "${var.cluster_id}-worker-in-ssh"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_overlay" {
  name    = "${var.cluster_id}-worker-in-overlay"
  network = local.cluster_network

  # allow VXLAN and GENEVE
  allow {
    protocol = "udp"
    ports    = ["4789", "6081"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_overlay_from_master" {
  name    = "${var.cluster_id}-worker-in-overlay-from-master"
  network = local.cluster_network

  # allow VXLAN and GENEVE
  allow {
    protocol = "udp"
    ports    = ["4789", "6081"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_internal" {
  name    = "${var.cluster_id}-worker-in-internal"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_internal_from_master" {
  name    = "${var.cluster_id}-worker-in-internal-from-master"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_internal_udp" {
  name    = "${var.cluster_id}-worker-in-udp"
  network = local.cluster_network

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_internal_from_master_udp" {
  name    = "${var.cluster_id}-worker-in-from-master-udp"
  network = local.cluster_network

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_kubelet_insecure" {
  name    = "${var.cluster_id}-worker-in-kubelet-insecure"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_kubelet_insecure_from_master" {
  name    = "${var.cluster_id}-worker-in-kubelet-insecure-fr-mast"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_services_tcp" {
  name    = "${var.cluster_id}-worker-in-services-tcp"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}

resource "google_compute_firewall" "worker_ingress_services_udp" {
  name    = "${var.cluster_id}-worker-in-services-udp"
  network = local.cluster_network

  allow {
    protocol = "udp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-worker"]
}
