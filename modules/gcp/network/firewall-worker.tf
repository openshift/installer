resource "google_compute_firewall" "worker-ingress" {
  name    = "worker-ingress"
  network = "${google_compute_network.tectonic-network.name}"

  # ICMP
  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22", "80", "443"] # ssh, http, https
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["tectonic-workers"]
}

resource "google_compute_firewall" "worker-ingress-heapster" {
  name    = "worker-ingress-heapster"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["4194"]
  }

  source_tags = ["tectonic-masters", "tectonic-workers"]
  target_tags = ["tectonic-workers"]
}

resource "google_compute_firewall" "worker-ingress-flannel" {
  name    = "worker-ingress-flannel"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_tags = ["tectonic-masters", "tectonic-workers"]
  target_tags = ["tectonic-workers"]
}

resource "google_compute_firewall" "worker-ingress-node-exporter" {
  name    = "worker-ingress-node-exporter"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["9100"]
  }

  source_tags = ["tectonic-masters", "tectonic-workers"]
  target_tags = ["tectonic-workers"]
}

resource "google_compute_firewall" "worker-ingress-kubelet" {
  name    = "worker-ingress-kubelet"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["10250", "10255"] # insecure and secure ports
  }

  source_tags = ["tectonic-masters", "tectonic-workers"]
  target_tags = ["tectonic-workers"]
}

resource "google_compute_firewall" "worker-ingress-services" {
  name    = "worker-ingress-services"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["32000-32767"]
  }

  source_tags = ["tectonic-workers"]
  target_tags = ["tectonic-workers"]
}
