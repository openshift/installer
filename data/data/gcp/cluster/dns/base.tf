locals {
  description = "Created By OpenShift Installer"
}

resource "google_dns_managed_zone" "int" {
  count = var.private_zone_name != "" ? 0 : 1

  name        = "${var.cluster_id}-private-zone"
  description = local.description
  dns_name    = "${var.cluster_domain}."
  visibility  = "private"
  project     = var.project_id
  labels      = var.gcp_extra_labels

  private_visibility_config {
    networks {
      network_url = var.network
    }
  }

  depends_on = [google_dns_record_set.api_external]
}

resource "google_dns_record_set" "api_external" {
  count = var.public_endpoints ? 1 : 0

  name         = "api.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = var.public_zone_name
  rrdatas      = [var.api_external_lb_ip]
  project      = var.project_id
}

resource "google_dns_record_set" "api_internal" {
  name         = "api-int.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = var.private_zone_name != "" ? var.private_zone_name : google_dns_managed_zone.int[0].name
  rrdatas      = [var.api_internal_lb_ip]
  project      = var.project_id
}

resource "google_dns_record_set" "api_external_internal_zone" {
  name         = "api.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = var.private_zone_name != "" ? var.private_zone_name : google_dns_managed_zone.int[0].name
  rrdatas      = [var.api_internal_lb_ip]
  project      = var.project_id
}
