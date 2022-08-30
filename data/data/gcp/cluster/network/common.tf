# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

data "google_compute_network" "preexisting_cluster_network" {
  count = var.preexisting_network ? 1 : 0

  name    = var.cluster_network
  project = var.network_project_id
}

data "google_compute_subnetwork" "preexisting_master_subnet" {
  count = var.preexisting_network ? 1 : 0

  name    = var.master_subnet
  project = var.network_project_id
}

data "google_compute_subnetwork" "preexisting_worker_subnet" {
  count = var.preexisting_network ? 1 : 0

  name    = var.worker_subnet
  project = var.network_project_id
}

locals {
  cluster_network = var.preexisting_network ? data.google_compute_network.preexisting_cluster_network[0].self_link : google_compute_network.cluster_network[0].self_link
  master_subnet   = var.preexisting_network ? data.google_compute_subnetwork.preexisting_master_subnet[0].self_link : google_compute_subnetwork.master_subnet[0].self_link
  description     = "Created By OpenShift Installer"
}
