locals {
  # Mint Mode is default so assume Mint unless explicitly set otherwise.
  mint_mode = var.credentials_mode != "Passthrough" && var.credentials_mode != "Manual"
}

resource "google_service_account" "worker-node-sa" {
  count        = local.mint_mode ? 1 : 0
  account_id   = split("@", var.worker_sa_email)[0]
  display_name = "${var.cluster_id}-worker-node"
}

resource "google_project_iam_member" "worker-compute-viewer" {
  count  = local.mint_mode ? 1 : 0
  role   = "roles/compute.viewer"
  member = "serviceAccount:${google_service_account.worker-node-sa[0].email}"
}

resource "google_project_iam_member" "worker-storage-admin" {
  count  = local.mint_mode ? 1 : 0
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.worker-node-sa[0].email}"
}

resource "google_service_account" "master-node-sa" {
  count        = local.mint_mode ? 1 : 0
  account_id   = split("@", var.master_sa_email)[0]
  display_name = "${var.cluster_id}-master-node"
}

resource "google_project_iam_member" "master-compute-admin" {
  count  = local.mint_mode ? 1 : 0
  role   = "roles/compute.instanceAdmin"
  member = "serviceAccount:${google_service_account.master-node-sa[0].email}"
}

resource "google_project_iam_member" "master-network-admin" {
  count  = local.mint_mode ? 1 : 0
  role   = "roles/compute.networkAdmin"
  member = "serviceAccount:${google_service_account.master-node-sa[0].email}"
}

resource "google_project_iam_member" "master-compute-security" {
  count  = local.mint_mode ? 1 : 0
  role   = "roles/compute.securityAdmin"
  member = "serviceAccount:${google_service_account.master-node-sa[0].email}"
}

resource "google_project_iam_member" "master-storage-admin" {
  count  = local.mint_mode ? 1 : 0
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.master-node-sa[0].email}"
}

resource "google_project_iam_member" "master-service-account-user" {
  count  = local.mint_mode ? 1 : 0
  role   = "roles/iam.serviceAccountUser"
  member = "serviceAccount:${google_service_account.master-node-sa[0].email}"
}
