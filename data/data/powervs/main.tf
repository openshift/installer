provider "ibm" {
    ibmcloud_api_key = var.ibmcloud_api_key
    region           = var.ibmcloud_region
    zone             = var.ibmcloud_zone
}

module "bootstrap" {
  source = "./bootstrap"
  cloud_instance_id = var.cloud_instance_id
  cluster_id = var.cluster_id

  bootstrap = var.bootstrap
  sys_type = var.sys_type
  proc_type = var.proc_type
# TODO(mjturek): image and network IDs are not derived during terraform
#                for other providers. Need to investigate and follow how
#                other providers do this. cnorman's branch has some work
#                towards this already.

  image_name = var.image_name
  network_name = var.network_name
}
