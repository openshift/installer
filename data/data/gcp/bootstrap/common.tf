locals {
  external_ip = var.public_endpoints ? [google_compute_address.bootstrap.address] : []
}
