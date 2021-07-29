locals {
  description      = "Created By OpenShift Installer"
  public_endpoints = var.ibmcloud_publish_strategy == "External" ? true : false
  tags = concat(
    ["kubernetes.io_cluster_${var.cluster_id}:owned"],
    var.ibmcloud_extra_tags
  )
}

############################################
# IBM Cloud provider
############################################

provider "ibm" {
  ibmcloud_api_key = var.ibmcloud_api_key
  region           = var.ibmcloud_region
}