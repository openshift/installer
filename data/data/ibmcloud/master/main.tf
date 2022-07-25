locals {
  prefix              = var.cluster_id
  port_kubernetes_api = 6443
  port_machine_config = 22623
  subnet_count        = length(var.control_plane_subnet_id_list)
  zone_count          = length(var.control_plane_subnet_zone_list)
}

############################################
# Master nodes
############################################

resource "ibm_is_instance" "master_node" {
  count = var.master_count

  name           = "${local.prefix}-master-${count.index}"
  image          = var.vsi_image_id
  profile        = var.ibmcloud_master_instance_type
  resource_group = var.resource_group_id
  tags           = local.tags

  primary_network_interface {
    name            = "eth0"
    subnet          = var.control_plane_subnet_id_list[count.index % local.subnet_count]
    security_groups = var.control_plane_security_group_id_list
  }

  dedicated_host = length(var.control_plane_dedicated_host_id_list) > 0 ? var.control_plane_dedicated_host_id_list[count.index % local.zone_count] : null

  vpc  = var.vpc_id
  zone = var.control_plane_subnet_zone_list[count.index % local.zone_count]
  keys = []

  user_data = var.ignition_master
}

############################################
# Load balancer backend pool members
############################################

resource "ibm_is_lb_pool_member" "kubernetes_api_public" {
  count = local.public_endpoints ? var.master_count : 0

  lb             = var.lb_kubernetes_api_public_id
  pool           = var.lb_pool_kubernetes_api_public_id
  port           = local.port_kubernetes_api
  target_address = ibm_is_instance.master_node[count.index].primary_network_interface.0.primary_ipv4_address
}

resource "ibm_is_lb_pool_member" "kubernetes_api_private" {
  count = var.master_count

  lb             = var.lb_kubernetes_api_private_id
  pool           = var.lb_pool_kubernetes_api_private_id
  port           = local.port_kubernetes_api
  target_address = ibm_is_instance.master_node[count.index].primary_network_interface.0.primary_ipv4_address
}

resource "ibm_is_lb_pool_member" "machine_config" {
  count = var.master_count

  lb             = var.lb_kubernetes_api_private_id
  pool           = var.lb_pool_machine_config_id
  port           = local.port_machine_config
  target_address = ibm_is_instance.master_node[count.index].primary_network_interface.0.primary_ipv4_address
}