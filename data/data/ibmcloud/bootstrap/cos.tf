locals {
  ignition_object_name = "bootstrap.ign"
}

############################################
# COS service instance
############################################

resource "ibm_resource_instance" "cos" {
  name              = "${local.prefix}-cos"
  service           = "cloud-object-storage"
  plan              = "standard"
  location          = "global"
  resource_group_id = var.resource_group_id
}

############################################
# COS bucket
############################################

resource "ibm_cos_bucket" "bootstrap_ignition" {
  bucket_name          = "${local.prefix}-bootstrap-ignition"
  resource_instance_id = ibm_resource_instance.cos.id
  region_location      = var.cos_bucket_region
  storage_class        = "smart"
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

############################################
# Upload ignition config to COS bucket
############################################

resource "null_resource" "upload_ignition" {
  triggers = {
    ignition_content = filesha512(var.ignition_file)
    cos_bucket_crn   = ibm_cos_bucket.bootstrap_ignition.crn
  }

  provisioner "local-exec" {
    command = <<-EOT
      curl -X PUT 'https://${ibm_cos_bucket.bootstrap_ignition.s3_endpoint_public}/${ibm_cos_bucket.bootstrap_ignition.bucket_name}/${local.ignition_object_name}' \
      -H 'Authorization: ${data.ibm_iam_auth_token.iam_token.iam_access_token}' \
      -H 'Content-Type: application/json' \
      -d @${var.ignition_file}
    EOT
  }
}
