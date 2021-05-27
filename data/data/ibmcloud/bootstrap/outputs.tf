#######################################
# Bootstrap module outputs
#######################################

output "name" {
  value = ibm_is_instance.bootstrap_node.name
}

output "primary_ipv4_address" {
  value = ibm_is_instance.bootstrap_node.primary_network_interface.0.primary_ipv4_address
}

output "ready" {
  depends_on = [
    ibm_is_lb_pool_member.kubernetes_api_public,
    ibm_is_lb_pool_member.kubernetes_api_private,
    ibm_is_lb_pool_member.machine_config,
  ]
  value = true
}