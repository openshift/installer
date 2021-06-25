provider "ovirt" {
  url       = var.ovirt_url
  username  = var.ovirt_username
  password  = var.ovirt_password
  cafile    = var.ovirt_cafile
  ca_bundle = var.ovirt_ca_bundle
  insecure  = var.ovirt_insecure
}

module "template" {
  source                    = "./template"
  ovirt_cluster_id          = var.ovirt_cluster_id
  cluster_id                = var.cluster_id
  openstack_base_image_name = var.ovirt_base_image_name
  tmp_import_vm_id          = var.tmp_import_vm_id
}

module "affinity_group" {
  source                = "./affinity_group"
  ovirt_cluster_id      = var.ovirt_cluster_id
  ovirt_affinity_groups = var.ovirt_affinity_groups
}

module "masters" {
  source                        = "./masters"
  master_count                  = var.master_count
  ovirt_cluster_id              = var.ovirt_cluster_id
  ovirt_template_id             = module.template.releaseimage_template_id
  ignition_master               = var.ignition_master
  cluster_domain                = var.cluster_domain
  cluster_id                    = var.cluster_id
  ovirt_master_instance_type_id = var.ovirt_master_instance_type_id
  ovirt_master_cores            = var.ovirt_master_cores
  ovirt_master_sockets          = var.ovirt_master_sockets
  ovirt_master_memory           = var.ovirt_master_memory
  ovirt_master_vm_type          = var.ovirt_master_vm_type
  ovirt_master_os_disk_size_gb  = var.ovirt_master_os_disk_gb
  ovirt_master_affinity_groups  = var.ovirt_master_affinity_groups
  ovirt_affinity_group_count    = module.affinity_group.ovirt_affinity_group_count
}

