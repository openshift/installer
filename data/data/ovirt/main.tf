provider "ovirt" {
  url      = var.ovirt_url
  username = var.ovirt_username
  password = var.ovirt_password
}

module "template" {
  source                               = "./template"
  ovirt_cluster_id                     = var.ovirt_cluster_id
  ovirt_storage_domain_id              = var.ovirt_storage_domain_id
  ignition_bootstrap                   = var.ignition_bootstrap
  cluster_id                           = var.cluster_id
  openstack_base_image_name            = var.openstack_base_image_name
  openstack_base_image_local_file_path = var.openstack_base_image_local_file_path
  ovirt_network_name                   = var.ovirt_network_name
  ovirt_vnic_profile_id                = var.ovirt_vnic_profile_id
}

module "bootstrap" {
  source                               = "./bootstrap"
  ovirt_cluster_id                     = var.ovirt_cluster_id
  ovirt_template_id                    = module.template.releaseimage_template_id
  ignition_bootstrap                   = var.ignition_bootstrap
  cluster_id                           = var.cluster_id
  openstack_base_image_name            = var.openstack_base_image_name
  openstack_base_image_local_file_path = var.openstack_base_image_local_file_path
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
}
