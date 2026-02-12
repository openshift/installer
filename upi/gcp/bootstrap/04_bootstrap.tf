terraform {
  # Infra manager supports specific Terraform veresions; ensure compatibility
  required_version = ">=1.2.3"
  required_providers {
    google = {
      source = "hashicorp/google"
      version = ">= 4.0.0"
    }
  }
}

provider "google-beta" {
  project = "openshift-dev-installer"
  region = "${var.region}"
}

variable "infra_id" {
  type        = string
  description = "OpenShift Installer Infrastructure ID"
}

variable "region" {
  type        = string
  description = "GCP Region where the resources will be created."
  default     = "us-central1"
}

variable "zone" {
  type        = string
  description = "Zone inside of the region where the bootstrap node is created."
}

variable "cluster_network" {
  type        = string
  description = "Full link to the cluster network."
}

variable "subnet" {
  type        = string
  description = "Control plane subnet."
}

variable "image" {
  type        = string
  description = "Cluster Image."
}

variable "machine_type" {
  type        = string
  description = "Machine type for the bootstrap machine."
  default     = "n1-standard-4"
}

variable "root_volume_size" {
  type        = string
  description = "Size in GB for the root volume."
  default     = "128"
}

variable "bootstrap_ign" {
  type        = string
  description = "Bootstrap ignition data."
}

resource "google_compute_address" "bootstrap_public_ip" {
  provider = google-beta

  name = "${var.infra_id}-bootstrap-public-ip"
  region = "${var.region}"
}

resource "google_compute_instance" "bootstrap" {
  provider = google-beta

  name = "${var.infra_id}-bootstrap"
  zone = "${var.zone}"
  machine_type = "${var.machine_type}"
  tags = [
    "${var.infra_id}-master",
    "${var.infra_id}-bootstrap"
  ]
  boot_disk {
    auto_delete = true
    initialize_params {
      size = "${var.root_volume_size}"
      image = "${var.image}"
    }
  }
  network_interface {
    subnetwork = "${var.subnet}"
    access_config {
      nat_ip = google_compute_address.bootstrap_public_ip.address
    }
  }
  metadata = {
    user-data = "{\"ignition\":{\"config\":{\"replace\":{\"source\":\"${var.bootstrap_ign}\"}},\"version\":\"3.2.0\"}}"
  }
}

resource "google_compute_instance_group" "bootstrap_ig" {
  provider = google-beta

  name = "${var.infra_id}-bootstrap-ig"
  network = "${var.cluster_network}"
  zone = "${var.zone}"
  named_port {
    name = "ignition"
    port = 22623
  }
  named_port {
    name = "https"
    port = 6443
  }
}
