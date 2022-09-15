locals {
  description = "Created By OpenShift Installer"
}

resource "google_service_account" "worker-node-sa" {
  count        = var.service_account == "" ? 1 : 0
  account_id   = "${var.cluster_id}-w"
  display_name = "${var.cluster_id}-worker-node"
  description  = local.description
}

resource "google_project_iam_member" "worker-compute-viewer" {
  count   = var.service_account == "" ? 1 : 0
  project = var.project_id
  role    = "roles/compute.viewer"
  member  = "serviceAccount:${google_service_account.worker-node-sa[0].email}"
}

resource "google_project_iam_member" "worker-storage-admin" {
  count   = var.service_account == "" ? 1 : 0
  project = var.project_id
  role    = "roles/storage.admin"
  member  = "serviceAccount:${google_service_account.worker-node-sa[0].email}"
}
