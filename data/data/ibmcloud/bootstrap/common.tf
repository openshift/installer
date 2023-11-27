locals {
  description = "Created By OpenShift Installer"
  # If any Service Endpoints are being overridden, set visibility to 'private'
  # for IBM Terraform Provider to use the endpoints JSON file.
  endpoint_visibility = var.ibmcloud_endpoints_json_file != "" ? "private" : "public"
  public_endpoints    = var.ibmcloud_publish_strategy == "External" ? true : false
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

  # Manage endpoints for IBM Cloud services
  visibility          = local.endpoint_visibility
  endpoints_file_path = var.ibmcloud_endpoints_json_file
}
