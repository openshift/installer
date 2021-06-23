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
  name               = "${var.cluster_id}-bootstrap-ignition"
  location           = var.gcp_region
  bucket_policy_only = true
}

resource "google_storage_bucket_object" "ignition" {
  bucket  = google_storage_bucket.ignition.name
  name    = "bootstrap.ign"
  content = var.ignition_bootstrap
}

data "google_storage_object_signed_url" "ignition_url" {
  bucket   = google_storage_bucket.ignition.name
  path     = "bootstrap.ign"
  duration = "1h"
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
  name        = "${var.cluster_id}-bootstrap-in-ssh"
  network     = var.network
  description = local.description

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
