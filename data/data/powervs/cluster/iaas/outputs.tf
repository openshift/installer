output "service_instance_guid" {
  depends_on = [resource.time_sleep.wait_for_workspace]
  value = var.service_instance_guid == "" ? one(data.ibm_pi_workspace.created_workspace[*].id) : one(data.ibm_pi_workspace.existing_workspace[*].id)
}

output "pi_workspace_crn" {
  depends_on = [resource.time_sleep.wait_for_workspace]
  value = var.service_instance_guid == "" ? one(data.ibm_pi_workspace.created_workspace[*].pi_workspace_details.crn) : one(data.ibm_pi_workspace.existing_workspace[*].pi_workspace_details.crn)
}
