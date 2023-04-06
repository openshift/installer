locals {
  labels = merge(
    {
      "kubernetes-io-cluster-${var.cluster_id}" = "owned"
    },
    var.gcp_extra_labels,
  )
  description = "Created By OpenShift Installer"

  public_endpoints = var.gcp_publish_strategy == "External" ? true : false
  external_ip      = local.public_endpoints ? [google_compute_address.bootstrap.address] : []
}

provider "google" {
  credentials = var.gcp_service_account
  project     = var.gcp_project_id
  region      = var.gcp_region
}

resource "google_storage_bucket" "ignition" {
  name                        = "${var.cluster_id}-bootstrap-ignition"
  location                    = var.gcp_region
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "ignition" {
  bucket  = google_storage_bucket.ignition.name
  name    = "bootstrap.ign"
  content = var.ignition_bootstrap
}

resource "google_service_account" "bootstrap-node-sa" {
  count        = var.gcp_create_bootstrap_sa ? 1 : 0
  account_id   = "${var.cluster_id}-b"
  display_name = "${var.cluster_id}-bootstrap-node"
  description  = local.description
}

resource "google_service_account_key" "bootstrap" {
  count              = var.gcp_create_bootstrap_sa ? 1 : 0
  service_account_id = google_service_account.bootstrap-node-sa[0].name
}

resource "google_project_iam_member" "bootstrap-storage-admin" {
  count   = var.gcp_create_bootstrap_sa ? 1 : 0
  project = var.gcp_project_id
  role    = "roles/storage.admin"
  member  = "serviceAccount:${google_service_account.bootstrap-node-sa[0].email}"
}

data "google_storage_object_signed_url" "ignition_url" {
  bucket      = google_storage_bucket.ignition.name
  path        = "bootstrap.ign"
  duration    = "1h"
  credentials = var.gcp_create_bootstrap_sa ? base64decode(google_service_account_key.bootstrap[0].private_key) : null
}

data "ignition_config" "redirect" {
  replace {
    source = data.google_storage_object_signed_url.ignition_url.signed_url
  }
}

resource "google_compute_address" "bootstrap" {
  name        = "${var.cluster_id}-bootstrap-ip"
  description = local.description

  address_type = local.public_endpoints ? "EXTERNAL" : "INTERNAL"
  subnetwork   = local.public_endpoints ? null : var.master_subnet
}

resource "google_compute_firewall" "bootstrap_ingress_ssh" {
  count       = var.gcp_create_firewall_rules ? 1 : 0
  name        = "${var.cluster_id}-bootstrap-in-ssh"
  network     = var.network
  description = local.description
  project     = var.gcp_network_project_id

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [local.public_endpoints ? "0.0.0.0/0" : var.machine_v4_cidrs[0]]
  target_tags   = ["${var.cluster_id}-bootstrap"]
}

resource "google_compute_instance" "bootstrap" {
  name         = "${var.cluster_id}-bootstrap"
  description  = local.description
  machine_type = var.gcp_bootstrap_instance_type
  zone         = var.gcp_master_availability_zones[0]

  boot_disk {
    initialize_params {
      type  = var.gcp_master_root_volume_type
      size  = var.gcp_master_root_volume_size
      image = var.compute_image
    }
    kms_key_self_link = var.gcp_root_volume_kms_key_link
  }

  network_interface {
    subnetwork = var.master_subnet

    dynamic "access_config" {
      for_each = local.external_ip
      content {
        nat_ip = access_config.value
      }
    }

    network_ip = local.public_endpoints ? null : google_compute_address.bootstrap.address
  }

  dynamic "shielded_instance_config" {
    for_each = var.gcp_master_secure_boot != "" ? [1] : []
    content {
      enable_secure_boot = var.gcp_master_secure_boot == "Enabled"
    }
  }

  dynamic "confidential_instance_config" {
    for_each = var.gcp_master_confidential_compute != "" ? [1] : []
    content {
      enable_confidential_compute = var.gcp_master_confidential_compute == "Enabled"
    }
  }

  dynamic "scheduling" {
    for_each = var.gcp_master_on_host_maintenance != "" ? [1] : []
    content {
      on_host_maintenance = var.gcp_master_on_host_maintenance
    }
  }

  metadata = {
    user-data = data.ignition_config.redirect.rendered
  }

  tags = ["${var.cluster_id}-master", "${var.cluster_id}-bootstrap"]

  labels = local.labels

  lifecycle {
    # In GCP TF apply is run a second time to remove bootstrap node from LB.
    # If machine_type = n2-standard series, install will error as TF tries to
    # switch min_cpu_platform = "Intel Cascade Lake" -> null. BZ-1746119.
    # Also fails similarly with custom machine types: BZ-1908171.
    ignore_changes = [machine_type, min_cpu_platform]
  }
}

resource "google_compute_instance_group" "bootstrap" {

  name        = "${var.cluster_id}-bootstrap"
  description = local.description
  zone        = var.gcp_master_availability_zones[0]

  named_port {
    name = "ignition"
    port = "22623"
  }

  named_port {
    name = "https"
    port = "6443"
  }

  instances = google_compute_instance.bootstrap.*.self_link
}
