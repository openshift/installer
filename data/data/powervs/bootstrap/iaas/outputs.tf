output "si_guid" {
  value = data.ibm_pi_workspace.powervs_service_instance.id
}

output "si_crn" {
  value = data.ibm_pi_workspace.powervs_service_instance.pi_workspace_details.crn
}
