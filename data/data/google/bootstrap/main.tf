resource "google_service_account" "cluster" {
  account_id   = "${var.cluster_name}-cluster"
  display_name = "Cluster service account"
}

resource "google_project_iam_member" "cluster" {
  role    = "roles/editor"
  member  = "serviceAccount:${google_service_account.cluster.email}"
}

resource "google_compute_instance" "bootstrap" {
  name         = "${var.cluster_name}-bootstrap"
  machine_type = "${var.instance_type}"
  zone         = "${var.zone}"

  metadata = {
    user-data = "${var.ignition}"
  }

  tags = ["ocp", "ocp-master"]

  service_account = {
    email = "${google_service_account.cluster.email}"
    scopes = ["compute-rw"]
  }

  network_interface = {
    network    = "${var.subnetwork != "" ? "" : var.network}"
    subnetwork = "${var.subnetwork}"

    access_config = {    
    }
  }

  boot_disk {
    initialize_params {
      type  = "${var.root_volume_type}"
      size  = "${var.root_volume_size}"
      image = "${var.image_name}"
    }
  }

  labels = "${merge(map(
    "cluster-kubernetes-io", "${var.cluster_name}",
    ), var.extra_labels)}"
}

resource "google_compute_instance_group" "bootstrap" {
  name        = "${var.cluster_name}-bootstrap"
  zone        = "${var.zone}"
  network     = "${var.network}"

  named_port {
    name = "https"
    port = "6443"
  }

  instances = ["${google_compute_instance.bootstrap.self_link}"]
}
