locals {
  # Use the direct COS endpoint if IBM Cloud Service Endpoints are being overridden,
  # as public and private may not be available. The direct endpoint requires
  # additional IBM Cloud Account configuration, which must be configured when using
  # Service Endpoint overrides.
  cos_endpoint_type = local.endpoint_visibility == "private" ? "direct" : "public"
}

############################################
# COS bucket
############################################

resource "ibm_cos_bucket" "bootstrap_ignition" {
  bucket_name          = "${local.prefix}-bootstrap-ignition"
  endpoint_type        = local.cos_endpoint_type
  resource_instance_id = var.cos_resource_instance_crn
  region_location      = var.ibmcloud_region
  storage_class        = "smart"
}

############################################
# COS object
############################################

resource "ibm_cos_bucket_object" "bootstrap_ignition" {
  bucket_crn      = ibm_cos_bucket.bootstrap_ignition.crn
  bucket_location = ibm_cos_bucket.bootstrap_ignition.region_location
  content_file    = var.ignition_bootstrap_file
  endpoint_type   = local.cos_endpoint_type
  etag            = filemd5(var.ignition_bootstrap_file)
  key             = "bootstrap.ign"
}

############################################
# IAM service credentials
############################################

# NOTE/TODO: Get IAM token for created Service ID, not supported in provider
data "ibm_iam_auth_token" "iam_token" {}

# NOTE: Not used at the moment
# resource "ibm_iam_service_id" "cos" {
#   name = "${local.prefix}-cos-service-id"
# }

# NOTE: Not used at the moment
# resource "ibm_resource_key" "cos_reader" {
#   name                 = "${local.prefix}-cos-reader"
#   role                 = "Reader"
#   resource_instance_id = ibm_resource_instance.cos.id
#   parameters           = {
#     HMAC          = true
#     serviceid_crn = ibm_iam_service_id.cos.crn
#   }
# }

# NOTE: Not used at the moment
# resource "ibm_resource_key" "cos_writer" {
#   name                 = "${local.prefix}-cos-writer"
#   role                 = "Writer"
#   resource_instance_id = ibm_resource_instance.cos.id
#   parameters           = {
#     HMAC          = true
#     serviceid_crn = ibm_iam_service_id.cos.crn
#   }
# }
