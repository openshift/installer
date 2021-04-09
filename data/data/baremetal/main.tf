provider "libvirt" {
  uri = var.libvirt_uri
}

provider "ironic" {
  url                = var.ironic_uri
  inspector          = var.inspector_uri
  microversion       = "1.56"
  timeout            = 3600
  auth_strategy      = "http_basic"
  ironic_username    = var.ironic_username
  ironic_password    = var.ironic_password
  inspector_username = var.ironic_username
  inspector_password = var.ironic_password
}

module "bootstrap" {
  source = "./bootstrap"
  count  = var.bootstrapping ? 1 : 0

  cluster_id = var.cluster_id
  image      = var.bootstrap_os_image
  ignition   = var.ignition_bootstrap
  bridges    = var.bridges
}

module "masters" {
  source = "./masters"

  master_count   = var.master_count
  ignition       = var.ignition_master
  hosts          = var.hosts
  properties     = var.properties
  root_devices   = var.root_devices
  driver_infos   = var.driver_infos
  instance_infos = var.instance_infos
}
