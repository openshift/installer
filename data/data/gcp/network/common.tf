# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

// Fetch a list of available AZs
data "google_compute_zones" "available" {}

// Only reference data sources which are guaranteed to exist at any time (above) in this locals{} block
locals {
  // List of possible AZs for each type of subnet
  zones = "${data.google_compute_zones.available.names}"
}
