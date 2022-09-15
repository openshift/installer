locals {
  description = "Created By OpenShift Installer"
}

resource "google_service_account" "master-node-sa" {
  count        = var.service_account == "" ? 1 : 0
  account_id   = "${var.cluster_id}-m"
  display_name = "${var.cluster_id}-master-node"
  description  = local.description
}

resource "google_project_iam_member" "master-compute-admin" {
  count   = var.service_account == "" ? 1 : 0
  project = var.project_id
  role    = "roles/compute.instanceAdmin"
  member  = "serviceAccount:${google_service_account.master-node-sa[0].email}"
}

resource "google_project_iam_member" "master-network-admin" {
  count   = var.service_account == "" ? 1 : 0
  project = var.project_id
  role    = "roles/compute.networkAdmin"
  member  = "serviceAccount:${google_service_account.master-node-sa[0].email}"
}

resource "google_project_iam_member" "master-compute-security" {
  count   = var.service_account == "" ? 1 : 0
  project = var.project_id
  role    = "roles/compute.securityAdmin"
  member  = "serviceAccount:${google_service_account.master-node-sa[0].email}"
}

resource "google_project_iam_member" "master-storage-admin" {
  count   = var.service_account == "" ? 1 : 0
  project = var.project_id
  role    = "roles/storage.admin"
  member  = "serviceAccount:${google_service_account.master-node-sa[0].email}"
}

resource "google_project_iam_member" "master-service-account-user" {
  count   = var.service_account == "" ? 1 : 0
  project = var.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.master-node-sa[0].email}"
}

resource "google_compute_instance" "master" {
  count       = var.instance_count
  description = local.description

  name         = "${var.cluster_id}-master-${count.index}"
  machine_type = var.machine_type
  zone         = element(var.zones, count.index)

  boot_disk {
    initialize_params {
      type  = var.root_volume_type
      size  = var.root_volume_size
      image = var.image
    }
    kms_key_self_link = var.root_volume_kms_key_link
  }

  network_interface {
    subnetwork = var.subnet
  }

  metadata = {
    user-data = var.ignition
  }

  tags = concat(
    ["${var.cluster_id}-master"],
    var.tags,
  )

  labels = var.labels

  service_account {
    email  = var.service_account != "" ? var.service_account : google_service_account.master-node-sa[0].email
    scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }

  lifecycle {
    # In GCP TF apply is run a second time to remove bootstrap node from LB.
    # If machine_type = n2-standard series, install will error as TF tries to
    # switch min_cpu_platform = "Intel Cascade Lake" -> null. BZ-1746119.
    # Also fails similarly with custom machine types: BZ-1908171.
    ignore_changes = [machine_type, min_cpu_platform]
  }
}

resource "google_compute_instance_group" "master" {
  count       = length(var.zones)
  description = local.description

  name = "${var.cluster_id}-master-${var.zones[count.index]}"
  #network = var.network
  zone = var.zones[count.index]

  named_port {
    name = "ignition"
    port = "22623"
  }

  named_port {
    name = "https"
    port = "6443"
  }

  instances = [for instance in google_compute_instance.master.* : instance.self_link if instance.zone == var.zones[count.index]]
}
