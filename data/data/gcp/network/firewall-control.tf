resource "google_compute_firewall" "master_ingress_icmp" {
  name    = "${var.cluster_id}-master-ingress-icmp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "icmp"
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_ssh" {
  name    = "${var.cluster_id}-master-ingress-ssh"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_https" {
  name    = "${var.cluster_id}-master-ingress-https"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["6443"]
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_mcs" {
  name    = "${var.cluster_id}-master-ingress-mcs"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["22623"]
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_vxlan" {
  name    = "${var.cluster_id}-master-ingress-vxlan"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_vxlan_from_worker" {
  name    = "${var.cluster_id}-master-ingress-vxlan-from-worker"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_internal" {
  name    = "${var.cluster_id}-master-ingress-internal"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_internal_from_worker" {
  name    = "${var.cluster_id}-master-ingress-internal-from-worker"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_internal_udp" {
  name    = "${var.cluster_id}-master-ingress-internal-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_internal_from_worker_udp" {
  name    = "${var.cluster_id}-master-ingress-internal-from-worker-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kube_scheduler" {
  name    = "${var.cluster_id}-master-ingress-kube-scheduler"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10259"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kube_scheduler_from_worker" {
  name    = "${var.cluster_id}-master-ingress-kube-scheduler-from-worker"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10259"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kube_masterler_manager" {
  name    = "${var.cluster_id}-master-ingress-kube-masterler-manager"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10257"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kube_masterler_manager_from_worker" {
  name    = "${var.cluster_id}-master-ingress-kube-masterler-manager-from-worker"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10257"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kubelet_secure" {
  name    = "${var.cluster_id}-master-ingress-kubelet-secure"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kubelet_secure_from_worker" {
  name    = "${var.cluster_id}-master-ingress-kubelet-secure-from-worker"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_etcd" {
  name    = "${var.cluster_id}-master-ingress-etcd"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["2379-2380"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_services_tcp" {
  name    = "${var.cluster_id}-master-ingress-services-tcp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_services_udp" {
  name    = "${var.cluster_id}-master-ingress-services-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}
