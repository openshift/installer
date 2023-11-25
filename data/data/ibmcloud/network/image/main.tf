locals {
  # Use the direct COS endpoint if IBM Cloud Service Endpoints are being overridden,
  # as public and private may not be available. The direct endpoint requires
  # additional IBM Cloud Account configuration, which must be configured when using
  # Service Endpoint overrides.
  cos_endpoint_type = var.endpoint_visibility == "private" ? "direct" : "public"
  prefix            = var.cluster_id
}

resource "ibm_cos_bucket" "images" {
  bucket_name = "${local.prefix}-vsi-image"
  # Use the direct COS endpoint if IBM Cloud Service endpoints are being overridden,
  # as public and private may not be available. Direct requires additional IBM Cloud
  # Account configuration
  endpoint_type        = local.cos_endpoint_type
  resource_instance_id = var.cos_resource_instance_crn
  region_location      = var.region
  storage_class        = "smart"
}

resource "ibm_cos_bucket_object" "file" {
  bucket_crn      = ibm_cos_bucket.images.crn
  bucket_location = ibm_cos_bucket.images.region_location
  content_file    = var.image_filepath
  endpoint_type   = local.cos_endpoint_type
  key             = basename(var.image_filepath)
}

resource "ibm_iam_authorization_policy" "policy" {
  source_service_name         = "is"
  source_resource_type        = "image"
  target_service_name         = "cloud-object-storage"
  target_resource_instance_id = element(split(":", var.cos_resource_instance_crn), 7)
  roles                       = ["Reader"]
}

resource "ibm_is_image" "image" {
  depends_on = [
    ibm_iam_authorization_policy.policy
  ]

  name             = var.name
  href             = "cos://${ibm_cos_bucket.images.region_location}/${ibm_cos_bucket.images.bucket_name}/${ibm_cos_bucket_object.file.key}"
  operating_system = "rhel-coreos-stable-amd64"
  resource_group   = var.resource_group_id
  tags             = var.tags
}
