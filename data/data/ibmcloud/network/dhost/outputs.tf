#######################################
# Dedicated Host module outputs
#######################################

output "control_plane_dedicated_host_id_list" {
  value = local.dhosts_master_merged
}
