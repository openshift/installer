resource "google_compute_instance" "master" {
  count = var.instance_count

  name         = "${var.cluster_id}-master-${count.index}"
  machine_type = var.machine_type
  zone         = var.zones[count.index]

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
