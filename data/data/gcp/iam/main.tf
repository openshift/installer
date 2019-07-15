resource "google_service_account" "worker-node-sa" {
  account_id   = "${var.cluster_id}-w"
  display_name = "${var.cluster_id}-worker-node"
}

resource "google_project_iam_member" "worker-compute-viewer" {
  role   = "roles/compute.viewer"
  member = "serviceAccount:${google_service_account.worker-node-sa.email}"
}

resource "google_project_iam_member" "worker-storage-admin" {
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.worker-node-sa.email}"
}
