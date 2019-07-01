resource "google_compute_firewall" "control_ingress_icmp" {
  name    = "${var.cluster_id}-control-ingress-icmp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "icmp"
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_ssh" {
  name    = "${var.cluster_id}-control-ingress-ssh"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_https" {
  name    = "${var.cluster_id}-control-ingress-https"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["6443"]
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_mcs" {
  name    = "${var.cluster_id}-control-ingress-mcs"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["22623"]
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_vxlan" {
  name    = "${var.cluster_id}-control-ingress-vxlan"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_vxlan_from_compute" {
  name    = "${var.cluster_id}-control-ingress-vxlan-from-compute"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_internal" {
  name    = "${var.cluster_id}-control-ingress-internal"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_internal_from_compute" {
  name    = "${var.cluster_id}-control-ingress-internal-from-compute"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_internal_udp" {
  name    = "${var.cluster_id}-control-ingress-internal-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_internal_from_compute_udp" {
  name    = "${var.cluster_id}-control-ingress-internal-from-compute-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_kube_scheduler" {
  name    = "${var.cluster_id}-control-ingress-kube-scheduler"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10259"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_kube_scheduler_from_compute" {
  name    = "${var.cluster_id}-control-ingress-kube-scheduler-from-compute"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10259"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_kube_controller_manager" {
  name    = "${var.cluster_id}-control-ingress-kube-controller-manager"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10257"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_kube_controller_manager_from_compute" {
  name    = "${var.cluster_id}-control-ingress-kube-controller-manager-from-compute"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10257"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_kubelet_secure" {
  name    = "${var.cluster_id}-control-ingress-kubelet-secure"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_kubelet_secure_from_compute" {
  name    = "${var.cluster_id}-control-ingress-kubelet-secure-from-compute"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_etcd" {
  name    = "${var.cluster_id}-control-ingress-etcd"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["2379-2380"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_services_tcp" {
  name    = "${var.cluster_id}-control-ingress-services-tcp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "tcp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "control_ingress_services_udp" {
  name    = "${var.cluster_id}-control-ingress-services-udp"
  network = google_compute_network.cluster_network.self_link

  allow {
    protocol = "udp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}
