locals {
  description = "Created By OpenShift Installer"
}

resource "google_dns_managed_zone" "int" {
  name        = "${var.cluster_id}-private-zone"
  description = local.description
  dns_name    = "${var.cluster_domain}."
  visibility  = "private"
  project     = var.private_zone_project

  private_visibility_config {
    networks {
      network_url = var.network
    }
  }

  depends_on = [google_dns_record_set.api_external]
}

resource "google_dns_record_set" "api_external" {
  count = var.public_endpoints && var.create_public_zone_records ? 1 : 0

  name         = "api.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = var.public_zone_name
  project      = var.public_zone_project
  rrdatas      = [var.api_external_lb_ip]
}

resource "google_dns_record_set" "api_internal" {
  name         = "api-int.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = google_dns_managed_zone.int.name
  project      = var.private_zone_project
  rrdatas      = [var.api_internal_lb_ip]
}

resource "google_dns_record_set" "api_external_internal_zone" {
  name         = "api.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = google_dns_managed_zone.int.name
  project      = var.private_zone_project
  rrdatas      = [var.api_internal_lb_ip]
}
