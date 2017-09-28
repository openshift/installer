/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
resource "google_storage_bucket" "tectonic" {
  name          = "${var.tectonic_cluster_name}-${var.tectonic_gcp_region}-${var.tectonic_gcp_ext_google_managedzone_name}"
  location      = "${upper(var.tectonic_gcp_region)}"
  storage_class = "REGIONAL"
}

resource "google_storage_bucket_acl" "private_acl" {
  bucket         = "${google_storage_bucket.tectonic.name}"
  predefined_acl = "projectprivate"
}

# Bootkube / Tectonic assets
resource "google_storage_bucket_object" "tectonic-assets" {
  name   = "assets.zip"
  bucket = "${google_storage_bucket.tectonic.name}"
  source = "${data.archive_file.assets.output_path}"
}

# kubeconfig
resource "google_storage_bucket_object" "kubeconfig" {
  name    = "kubeconfig"
  bucket  = "${google_storage_bucket.tectonic.name}"
  content = "${module.bootkube.kubeconfig}"
}
