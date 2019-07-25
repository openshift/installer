resource "google_service_account" "master-node-sa" {
  account_id   = "${var.cluster_id}-m"
  display_name = "${var.cluster_id}-master-node"
}

resource "google_project_iam_member" "master-compute-admin" {
  role   = "roles/compute.instanceAdmin"
  member = "serviceAccount:${google_service_account.master-node-sa.email}"
}

resource "google_project_iam_member" "master-network-admin" {
  role   = "roles/compute.networkAdmin"
  member = "serviceAccount:${google_service_account.master-node-sa.email}"
}

resource "google_project_iam_member" "master-storage-admin" {
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.master-node-sa.email}"
}

resource "google_project_iam_member" "master-object-storage-admin" {
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.master-node-sa.email}"
}

resource "google_compute_instance" "master" {
  count = var.instance_count

  name         = "${var.cluster_id}-m-${count.index}"
  machine_type = var.machine_type
  zone         = element(var.zones, count.index)

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

  service_account {
    email  = google_service_account.master-node-sa.email
    scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }
}
