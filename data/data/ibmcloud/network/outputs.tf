#######################################
# Network module outputs
#######################################

output "control_plane_dedicated_host_id_list" {
  value = module.dhost.control_plane_dedicated_host_id_list
}

output "control_plane_security_group_id_list" {
  value = module.vpc.control_plane_security_group_id_list
}

output "control_plane_subnet_id_list" {
  value = module.vpc.control_plane_subnet_id_list
}

output "control_plane_subnet_zone_list" {
  value = module.vpc.control_plane_subnet_zone_list
}

output "cos_resource_instance_crn" {
  value = ibm_resource_instance.cos.crn
}

output "lb_kubernetes_api_public_id" {
  value = module.vpc.lb_kubernetes_api_public_id
}

output "lb_kubernetes_api_private_id" {
  value = module.vpc.lb_kubernetes_api_private_id
}

output "lb_pool_kubernetes_api_public_id" {
  value = module.vpc.lb_pool_kubernetes_api_public_id
}

output "lb_pool_kubernetes_api_private_id" {
  value = module.vpc.lb_pool_kubernetes_api_private_id
}

output "lb_pool_machine_config_id" {
  value = module.vpc.lb_pool_machine_config_id
}

output "resource_group_id" {
  value = local.resource_group_id
}

output "vpc_id" {
  value = module.vpc.vpc_id
}

output "vsi_image_id" {
  value = module.image.vsi_image_id
}
