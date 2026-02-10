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
  project = "${var.project}"
  region = "${var.region}"
}

variable "infra_id" {
  type        = string
  description = "OpenShift Installer Infrastructure ID"
}

variable "project" {
  type        = string
  description = "Project ID"
}

variable "region" {
  type        = string
  description = "GCP Region where the resources will be created."
  default     = "us-central1"
}

variable "master_subnet_cidr" {
  type        = string
  description = "CIDR for the control plane subnet."
}

variable "worker_subnet_cidr" {
  type        = string
  description = "CIDR for the compute subnet."
}

resource "google_compute_network" "cluster_network" {
  provider = google-beta

  name = "${var.infra_id}-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "master_subnet" {
  provider = google-beta

  name = "${var.infra_id}-master-subnet"
  ip_cidr_range = "${var.master_subnet_cidr}"
  region = "${var.region}"
  network = google_compute_network.cluster_network.self_link
}

resource "google_compute_subnetwork" "worker_subnet" {
  provider = google-beta

  name = "${var.infra_id}-worker-subnet"
  ip_cidr_range = "${var.worker_subnet_cidr}"
  region = "${var.region}"
  network = google_compute_network.cluster_network.self_link
}

#tfimport-terraform import google_compute_router._router  __project__//-router
resource "google_compute_router" "router" {
  provider = google-beta

  name = "${var.infra_id}-router"
  network = google_compute_network.cluster_network.self_link
  region = "${var.region}"
}
resource "google_compute_router_nat" "master_nat" {
  provider = google-beta

  name = "${var.infra_id}-nat-master"
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  nat_ip_allocate_option = "AUTO_ONLY"
  min_ports_per_vm = 7168
  subnetwork {
    name = google_compute_subnetwork.master_subnet.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  router = google_compute_router.router.name
  region = "${var.region}"

  depends_on = [
    google_compute_router.router
  ]
}
resource "google_compute_router_nat" "worker_nat" {
  provider = google-beta

  name = "${var.infra_id}-nat-worker"
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  nat_ip_allocate_option = "AUTO_ONLY"
  min_ports_per_vm = 512
  subnetwork {
    name = google_compute_subnetwork.worker_subnet.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  router = google_compute_router.router.name
  region = "${var.region}"

  depends_on = [
    google_compute_router.router
  ]
}
