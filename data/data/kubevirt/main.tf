provider "kubernetes" {
}

provider "kubevirt" {
}

module "datavolume" {
  source         = "./datavolume"
  pvc_name       = var.kubevirt_source_pvc_name
  namespace      = var.kubevirt_namespace
  labels         = var.kubevirt_labels
  pv_access_mode = var.kubevirt_pv_access_mode
  storage_class  = var.kubevirt_storage_class
  image_url      = var.kubevirt_image_url
}

module "masters" {
  source         = "./masters"
  master_count   = var.master_count
  cluster_id     = var.cluster_id
  ignition_data  = var.ignition_master
  namespace      = var.kubevirt_namespace
  storage        = var.kubevirt_master_storage
  memory         = var.kubevirt_master_memory
  cpu            = var.kubevirt_master_cpu
  storage_class  = var.kubevirt_storage_class
  network_name   = var.kubevirt_network_name
  pv_access_mode = var.kubevirt_pv_access_mode
  labels         = var.kubevirt_labels
  pvc_name       = module.datavolume.pvc_name
}

module "bootstrap" {
  source         = "./bootstrap"
  count          = var.bootstrapping ? 1 : 0
  cluster_id     = var.cluster_id
  ignition_data  = var.ignition_bootstrap
  namespace      = var.kubevirt_namespace
  storage_class  = var.kubevirt_storage_class
  network_name   = var.kubevirt_network_name
  pv_access_mode = var.kubevirt_pv_access_mode
  labels         = var.kubevirt_labels
  pvc_name       = module.datavolume.pvc_name
}
