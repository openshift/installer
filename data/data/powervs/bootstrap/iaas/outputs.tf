output "si_guid" {
  value = data.ibm_resource_instance.powervs_service_instance.guid
}

output "si_crn" {
  value = data.ibm_resource_instance.powervs_service_instance.crn
}
