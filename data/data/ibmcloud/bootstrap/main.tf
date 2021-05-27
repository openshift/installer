locals {
  prefix              = var.cluster_id
  port_kubernetes_api = 6443
  port_machine_config = 22623
}

############################################
# Bootstrap node
############################################

resource "ibm_is_instance" "bootstrap_node" {
  name           = "${local.prefix}-bootstrap"
  image          = var.vsi_image_id
  profile        = var.vsi_profile
  resource_group = var.resource_group_id

  primary_network_interface {
    name            = "eth0"
    subnet          = var.subnet_id
    security_groups = [ var.security_group_id ]
  }

  vpc  = var.vpc_id
  zone = var.zone
  keys = []

  # Use custom ignition config that pulls content from COS bucket
  user_data = templatefile("${path.module}/templates/bootstrap.ign", {
    HOSTNAME    = ibm_cos_bucket.bootstrap_ignition.s3_endpoint_public
    BUCKET_NAME = ibm_cos_bucket.bootstrap_ignition.bucket_name
    OBJECT_NAME = ibm_cos_bucket_object.bootstrap_ignition.key
    IAM_TOKEN   = data.ibm_iam_auth_token.iam_token.iam_access_token
  })
}

############################################
# Floating IP
############################################

resource "ibm_is_floating_ip" "bootstrap_floatingip" {
  name           = "${local.prefix}-bootstrap-node-ip"
  resource_group = var.resource_group_id
  target         = ibm_is_instance.bootstrap_node.primary_network_interface.0.id
}

############################################
# Load balancer backend pool members
############################################

resource "ibm_is_lb_pool_member" "kubernetes_api_public" {
  lb             = var.lb_kubernetes_api_public_id
  pool           = var.lb_pool_kubernetes_api_public_id
  port           = local.port_kubernetes_api
  target_address = ibm_is_instance.bootstrap_node.primary_network_interface.0.primary_ipv4_address
}

resource "ibm_is_lb_pool_member" "kubernetes_api_private" {
  lb             = var.lb_kubernetes_api_private_id
  pool           = var.lb_pool_kubernetes_api_private_id
  port           = local.port_kubernetes_api
  target_address = ibm_is_instance.bootstrap_node.primary_network_interface.0.primary_ipv4_address
}

resource "ibm_is_lb_pool_member" "machine_config" {
  lb             = var.lb_kubernetes_api_private_id
  pool           = var.lb_pool_machine_config_id
  port           = local.port_machine_config
  target_address = ibm_is_instance.bootstrap_node.primary_network_interface.0.primary_ipv4_address
}
