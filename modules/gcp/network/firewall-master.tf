resource "google_compute_firewall" "master-ingress" {
  name    = "master-ingress"
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
  target_tags   = ["tectonic-masters"]
}

resource "google_compute_firewall" "master-ingress-heapster" {
  name    = "master-ingress-heapster"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["4194"]
  }

  source_tags = ["tectonic-masters", "tectonic-workers"]
  target_tags = ["tectonic-masters"]
}

resource "google_compute_firewall" "master-ingress-flannel" {
  name    = "master-ingress-flannel"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_tags = ["tectonic-masters", "tectonic-workers"]
  target_tags = ["tectonic-masters"]
}

resource "google_compute_firewall" "master-ingress-node-exporter" {
  name    = "master-ingress-node-exporter"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["9100"]
  }

  source_tags = ["tectonic-masters", "tectonic-workers"]
  target_tags = ["tectonic-masters"]
}

resource "google_compute_firewall" "master-ingress-kubelet" {
  name    = "master-ingress-kubelet"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["10250", "10255"] # insecure and secure ports
  }

  source_tags = ["tectonic-masters", "tectonic-workers"]
  target_tags = ["tectonic-masters"]
}

resource "google_compute_firewall" "master-ingress-etcd" {
  name    = "master-ingress-etcd"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["2379", "12379"] # etcd and bootstrap-etcd
  }

  source_tags = ["tectonic-masters"]
  target_tags = ["tectonic-masters", "tectonic-etcd"]
}

resource "google_compute_firewall" "master-ingress-services" {
  name    = "master-ingress-services"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["32000-32767"]
  }

  source_tags = ["tectonic-masters"]
  target_tags = ["tectonic-masters"]
}
