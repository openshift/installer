resource "google_compute_firewall" "etcd-ingress" {
  name    = "etcd-ingress"
  network = "${google_compute_network.tectonic-network.name}"

  # ICMP
  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["22"] # ssh
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["tectonic-etcd"]
}

resource "google_compute_firewall" "etcd" {
  name    = "etcd"
  network = "${google_compute_network.tectonic-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["2379", "2380", "12379"] # etcd and bootstrap-etcd
  }

  source_tags = ["tectonic-etcd"]
  target_tags = ["tectonic-etcd"]
}
