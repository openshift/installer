provider "libvirt" {
  uri = var.libvirt_uri
}

provider "ironic" {
  url          = "http://${var.bootstrap_provisioning_ip}:6385/v1"
  inspector    = "http://${var.bootstrap_provisioning_ip}:5050/v1"
  microversion = "1.56"
  timeout      = 3600
}

module "bootstrap" {
  source = "./bootstrap"

  cluster_id          = var.cluster_id
  image               = var.bootstrap_os_image
  ignition            = var.ignition_bootstrap
  external_bridge     = var.external_bridge
  provisioning_bridge = var.provisioning_bridge
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
