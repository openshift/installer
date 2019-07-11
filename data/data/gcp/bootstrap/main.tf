resource "google_compute_firewall" "bootstrap_ingress_ssh" {
  name    = "${var.cluster_id}-bootstrap-ingress-ssh"
  network = var.network

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["${var.cluster_id}-bootstrap"]
}

resource "google_compute_instance" "bootstrap" {
  name         = "${var.cluster_id}-bootstrap"
  machine_type = var.machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      type  = var.root_volume_type
      size  = var.root_volume_size
      image = var.image_name
    }
  }

  network_interface {
    subnetwork = var.subnet

    access_config {
    }
  }

  metadata = {
    user-data = var.ignition
  }

  tags = ["${var.cluster_id}-master", "${var.cluster_id}-bootstrap"]

  labels = var.labels
}
