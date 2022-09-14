locals {
  description       = "Created By OpenShift Installer"
  private_zone_name = var.create_private_zone ? google_dns_managed_zone.int[0].name : var.private_zone_name
}

resource "google_dns_managed_zone" "int" {
  count       = var.create_private_zone ? 1 : 0
  name        = "${var.cluster_id}-private-zone"
  description = local.description
  dns_name    = "${var.cluster_domain}."
  visibility  = "private"

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
  count        = var.create_private_zone_records ? 1 : 0
  name         = "api-int.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = local.private_zone_name
  project      = var.private_zone_project
  rrdatas      = [var.api_internal_lb_ip]
}

resource "google_dns_record_set" "api_external_internal_zone" {
  count        = var.create_private_zone_records ? 1 : 0
  name         = "api.${var.cluster_domain}."
  type         = "A"
  ttl          = "60"
  managed_zone = local.private_zone_name
  project      = var.private_zone_project
  rrdatas      = [var.api_internal_lb_ip]
}
