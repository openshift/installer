#######################################
# VPC module outputs
#######################################

output "control_plane_security_group_id" {
  # depends_on = [
  #   ibm_is_security_group_rule.control_plane_inbound,
  #   ibm_is_security_group_rule.control_plane_outbound,
  # ]
  value = ibm_is_security_group.control_plane.id
}

output "control_plane_subnet_id_list" {
  value = ibm_is_subnet.control_plane.*.id
}

output "control_plane_subnet_zone_list" {
  value = ibm_is_subnet.control_plane.*.zone
}

output "compute_security_group_id" {
  # depends_on = [
  #   ibm_is_security_group_rule.compute_inbound,
  #   ibm_is_security_group_rule.compute_outbound,
  # ]
  value = ibm_is_security_group.compute.id
}

output "compute_subnet_id_list" {
  value = ibm_is_subnet.compute.*.id
}

output "compute_subnet_zone_list" {
  value = ibm_is_subnet.compute.*.zone
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

output "vpc_crn" {
  value = ibm_is_vpc.vpc.crn
}

output "vpc_id" {
  value = ibm_is_vpc.vpc.id
}