resource "google_compute_firewall" "global" {
  name = "${var.cluster_name}-global"
  network = "${google_compute_network.default.self_link}"

  allow {
    protocol = "icmp"
  }
}

resource "google_compute_firewall" "internal" {
  name = "${var.cluster_name}-internal"
  network = "${google_compute_network.default.self_link}"

  allow {
    protocol = "tcp"
    ports = ["80", "6443-6445", "4789", "9000-9990", "10250-10255", "30000-32767"]
  }

  source_tags = ["ocp"]
  target_tags = ["ocp"]
}

resource "google_compute_firewall" "master-ssh" {
  name = "${var.cluster_name}-master-ssh"
  network = "${google_compute_network.default.self_link}"

  allow {
    protocol = "tcp"
    ports = ["22"]
  }

  target_tags = ["ocp-master"]
}

resource "google_compute_firewall" "master-internal" {
  name = "${var.cluster_name}-tcp"
  network = "${google_compute_network.default.self_link}"

  allow {
    protocol = "tcp"
    ports = ["2379-2380", "12379-12380"]
  }

  source_tags = ["ocp-master"]
  target_tags = ["ocp-master"]
}
