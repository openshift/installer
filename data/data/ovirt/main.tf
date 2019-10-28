provider "ovirt" {
  url      = var.ovirt_url
  username = var.ovirt_username
  password = var.ovirt_password
}

module "bootstrap" {
  source             = "./bootstrap"
  ovirt_cluster_id   = var.ovirt_cluster_id
  ovirt_template_id  = var.ovirt_template_id
  ignition_bootstrap = var.ignition_bootstrap
  cluster_id         = var.cluster_id
}

module "masters" {
  source            = "./masters"
  master_count      = var.master_count
  ovirt_cluster_id  = var.ovirt_cluster_id
  ovirt_template_id = var.ovirt_template_id
  ignition_master   = var.ignition_master
  cluster_domain    = var.cluster_domain
  cluster_id        = var.cluster_id
}
