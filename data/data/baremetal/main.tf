provider "libvirt" {
  uri = var.libvirt_uri
}

provider "ironic" {
  url          = var.ironic_uri
  microversion = "1.52"
}

module "bootstrap" {
  source = "./bootstrap"

  cluster_id          = var.cluster_id
  image               = var.os_image
  ignition            = var.ignition_bootstrap
  external_bridge     = var.external_bridge
  provisioning_bridge = var.provisioning_bridge
}

module "masters" {
  source = "./masters"

  ironic_uri     = var.ironic_uri
  master_count   = var.master_count
  ignition       = var.ignition_master
  hosts          = var.hosts
  properties     = var.properties
  root_devices   = var.root_devices
  driver_infos   = var.driver_infos
  instance_infos = var.instance_infos
}
