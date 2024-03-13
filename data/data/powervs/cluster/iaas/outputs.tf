output "si_name" {
  depends_on = [resource.time_sleep.wait_for_workspace]
  value      = var.service_instance_name == "" ? one(resource.ibm_resource_instance.created_service_instance[*].name) : one(data.ibm_resource_instance.existing_service_instance[*].name)
}

output "si_guid" {
  depends_on = [resource.time_sleep.wait_for_workspace]
  value      = var.service_instance_name == "" ? one(resource.ibm_resource_instance.created_service_instance[*].guid) : one(data.ibm_resource_instance.existing_service_instance[*].guid)
}

output "si_crn" {
  depends_on = [resource.time_sleep.wait_for_workspace]
  value      = var.service_instance_name == "" ? one(resource.ibm_resource_instance.created_service_instance[*].crn) : one(data.ibm_resource_instance.existing_service_instance[*].crn)
}
