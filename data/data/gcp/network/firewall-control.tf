resource "google_compute_firewall" "master_ingress_icmp" {
  name    = "${var.cluster_id}-master-in-icmp"
  network = local.cluster_network

  allow {
    protocol = "icmp"
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_ssh" {
  name    = "${var.cluster_id}-master-in-ssh"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.network_cidr]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_https" {
  name    = "${var.cluster_id}-master-in-https"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["6443"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_from_health_checks" {
  name    = "${var.cluster_id}-master-in-from-health-checks"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["6080", "22624"]
  }

  source_ranges = ["35.191.0.0/16", "130.211.0.0/22"]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_mcs" {
  name    = "${var.cluster_id}-master-in-mcs"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["22623"]
  }

  source_ranges = var.preexisting_network ? ["0.0.0.0/0"] : [var.network_cidr, google_compute_address.master_nat_ip[0].address, google_compute_address.worker_nat_ip[0].address]
  target_tags   = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_overlay" {
  name    = "${var.cluster_id}-master-in-overlay"
  network = local.cluster_network

  # allow VXLAN and GENEVE
  allow {
    protocol = "udp"
    ports    = ["4789", "6081"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_overlay_from_worker" {
  name    = "${var.cluster_id}-master-in-overlay-from-worker"
  network = local.cluster_network

  # allow VXLAN and GENEVE
  allow {
    protocol = "udp"
    ports    = ["4789", "6081"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_internal" {
  name    = "${var.cluster_id}-master-in-internal"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_internal_from_worker" {
  name    = "${var.cluster_id}-master-in-internal-from-worker"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_internal_udp" {
  name    = "${var.cluster_id}-master-in-internal-udp"
  network = local.cluster_network

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_internal_from_worker_udp" {
  name    = "${var.cluster_id}-master-in-internal-from-worker-udp"
  network = local.cluster_network

  allow {
    protocol = "udp"
    ports    = ["9000-9999"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kube_scheduler" {
  name    = "${var.cluster_id}-master-in-kube-scheduler"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["10259"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kube_scheduler_from_worker" {
  name    = "${var.cluster_id}-master-in-kube-scheduler-fr-worker"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["10259"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kube_master_manager" {
  name    = "${var.cluster_id}-master-in-kube-master-manager"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["10257"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kube_master_manager_from_worker" {
  name    = "${var.cluster_id}-master-in-kube-master-mgr-fr-work"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["10257"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kubelet_secure" {
  name    = "${var.cluster_id}-master-in-kubelet-secure"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_kubelet_secure_from_worker" {
  name    = "${var.cluster_id}-master-in-kubelet-secure-fr-worker"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["10250"]
  }

  source_tags = ["${var.cluster_id}-worker"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_etcd" {
  name    = "${var.cluster_id}-master-in-etcd"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["2379-2380"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_services_tcp" {
  name    = "${var.cluster_id}-master-in-services-tcp"
  network = local.cluster_network

  allow {
    protocol = "tcp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}

resource "google_compute_firewall" "master_ingress_services_udp" {
  name    = "${var.cluster_id}-master-in-services-udp"
  network = local.cluster_network

  allow {
    protocol = "udp"
    ports    = ["30000-32767"]
  }

  source_tags = ["${var.cluster_id}-master"]
  target_tags = ["${var.cluster_id}-master"]
}
