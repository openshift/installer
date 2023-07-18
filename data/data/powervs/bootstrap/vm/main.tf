data "ibm_resource_group" "cos_group" {
  name = var.resource_group
}

resource "ibm_resource_instance" "cos_instance" {
  name              = "${var.cluster_id}-cos"
  resource_group_id = data.ibm_resource_group.cos_group.id
  service           = "cloud-object-storage"
  plan              = "standard"
  location          = var.cos_instance_location
  tags              = [var.cluster_id]
}

# Create an IBM COS Bucket to store ignition
resource "ibm_cos_bucket" "ignition" {
  bucket_name          = "${var.cluster_id}-bootstrap-ign"
  resource_instance_id = ibm_resource_instance.cos_instance.id
  region_location      = var.cos_bucket_location
  storage_class        = var.cos_storage_class
}

resource "ibm_resource_key" "cos_service_cred" {
  name                 = "${var.cluster_id}-cred"
  role                 = "Reader"
  resource_instance_id = ibm_resource_instance.cos_instance.id
  parameters           = { HMAC = true }
}

# Place the bootstrap ignition file in the ignition COS bucket
resource "ibm_cos_bucket_object" "ignition" {
  bucket_crn      = ibm_cos_bucket.ignition.crn
  bucket_location = ibm_cos_bucket.ignition.region_location
  content         = var.ignition
  key             = "bootstrap.ign"
  etag            = md5(var.ignition)
}

data "ibm_iam_auth_token" "iam_token" {}

# Create the bootstrap instance
resource "ibm_pi_instance" "bootstrap" {
  pi_memory            = var.memory
  pi_processors        = var.processors
  pi_instance_name     = "${var.cluster_id}-bootstrap"
  pi_proc_type         = var.proc_type
  pi_image_id          = var.image_id
  pi_sys_type          = var.sys_type
  pi_cloud_instance_id = var.cloud_instance_id
  pi_network {
    network_id = var.dhcp_network_id
  }
  pi_user_data = base64encode(templatefile("${path.module}/templates/bootstrap.ign", {
    PROTOCOL    = var.enable_snat ? "https" : "http"
    HOSTNAME    = var.enable_snat ? ibm_cos_bucket.ignition.s3_endpoint_public : var.proxy_server_ip
    BUCKET_NAME = ibm_cos_bucket.ignition.bucket_name
    OBJECT_NAME = ibm_cos_bucket_object.ignition.key
    IAM_TOKEN   = data.ibm_iam_auth_token.iam_token.iam_access_token
  }))
  pi_key_pair_name = var.ssh_key_name
  pi_health_status = "WARNING"
}

resource "time_sleep" "wait_for_bootstrap_macs" {
  create_duration = "3m"

  depends_on = [ibm_pi_instance.bootstrap]
}

locals {
  bootstrap_ips = [for lease in data.ibm_pi_dhcp.dhcp_service_refresh.leases : lease.instance_ip if ibm_pi_instance.bootstrap.pi_network[0].mac_address == lease.instance_mac]
}

data "ibm_pi_dhcp" "dhcp_service_refresh" {
  depends_on           = [time_sleep.wait_for_bootstrap_macs]
  pi_cloud_instance_id = var.cloud_instance_id
  pi_dhcp_id           = var.dhcp_id
}
