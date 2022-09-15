#######################################
# VPC module outputs
#######################################

output "control_plane_security_group_id_list" {
  value = [
    ibm_is_security_group.cluster_wide.id,
    ibm_is_security_group.openshift_network.id,
    ibm_is_security_group.control_plane.id,
    ibm_is_security_group.control_plane_internal.id,
  ]
}

output "control_plane_subnet_id_list" {
  value = local.control_plane_subnets[*].id
}

output "control_plane_subnet_zone_list" {
  value = local.control_plane_subnets[*].zone
}

output "lb_kubernetes_api_public_hostname" {
  value = var.public_endpoints ? ibm_is_lb.kubernetes_api_public.0.hostname : ""
}

output "lb_kubernetes_api_public_id" {
  # Wait for frontend listeners to be ready before use
  depends_on = [
    ibm_is_lb_listener.kubernetes_api_public
  ]
  value = var.public_endpoints ? ibm_is_lb.kubernetes_api_public.0.id : ""
}

output "lb_kubernetes_api_private_hostname" {
  value = ibm_is_lb.kubernetes_api_private.hostname
}

output "lb_kubernetes_api_private_id" {
  # Wait for frontend listeners to be ready before use
  depends_on = [
    ibm_is_lb_listener.kubernetes_api_private,
    ibm_is_lb_listener.machine_config,
  ]
  value = ibm_is_lb.kubernetes_api_private.id
}

output "lb_pool_kubernetes_api_public_id" {
  value = var.public_endpoints ? ibm_is_lb_pool.kubernetes_api_public.0.id : ""
}

output "lb_pool_kubernetes_api_private_id" {
  value = ibm_is_lb_pool.kubernetes_api_private.id
}

output "lb_pool_machine_config_id" {
  value = ibm_is_lb_pool.machine_config.id
}

output "vpc_id" {
  value = local.vpc_id
}

output "vpc_crn" {
  value = local.vpc_crn
}
