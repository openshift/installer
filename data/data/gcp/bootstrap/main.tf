resource "google_storage_bucket" "ignition" {
  name = "${var.cluster_id}-bootstrap-ignition"
}

resource "google_storage_bucket_object" "ignition" {
  bucket  = "${google_storage_bucket.ignition.name}"
  name    = "bootstrap.ign"
  content = var.ignition
}

data "google_storage_object_signed_url" "ignition_url" {
  bucket   = "${google_storage_bucket.ignition.name}"
  path     = "bootstrap.ign"
  duration = "1h"
}

data "ignition_config" "redirect" {
  replace {
    source = "${data.google_storage_object_signed_url.ignition_url.signed_url}"
  }
}

resource "google_compute_address" "bootstrap" {
  name = "${var.cluster_id}-bootstrap-public-ip"
}

resource "google_compute_firewall" "bootstrap_ingress_ssh" {
  name    = "${var.cluster_id}-bootstrap-in-ssh"
  network = var.network

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["${var.cluster_id}-bootstrap"]
}

resource "google_compute_instance" "bootstrap" {
  count = var.bootstrap_enabled ? 1 : 0

  name         = "${var.cluster_id}-b"
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
      nat_ip = "${google_compute_address.bootstrap.address}"
    }
  }

  metadata = {
    user-data = data.ignition_config.redirect.rendered
  }

  tags = ["${var.cluster_id}-master", "${var.cluster_id}-bootstrap"]

  labels = var.labels
}
