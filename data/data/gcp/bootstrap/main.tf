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
  }

  metadata = {
    user-data = var.ignition
  }

  tags = ["${var.cluster_id}-master"]

  labels = var.labels
}

resource "google_compute_instance_group" "bootstrap" {
  name    = "${var.cluster_id}-bootstrap"
  network = var.network
  zone    = var.zone

  named_port {
    name = "ignition"
    port = "22623"
  }

  named_port {
    name = "https"
    port = "6443"
  }

  instances = [google_compute_instance.bootstrap.self_link]
}
