output "si_name" {
  value = var.service_instance_name == "" ? one(resource.ibm_resource_instance.created_service_instance[*].name) : one(data.ibm_resource_instance.existing_service_instance[*].name)
}

output "si_guid" {
  value = var.service_instance_name == "" ? one(resource.ibm_resource_instance.created_service_instance[*].guid) : one(data.ibm_resource_instance.existing_service_instance[*].guid)
}

output "si_crn" {
  value = var.service_instance_name == "" ? one(resource.ibm_resource_instance.created_service_instance[*].crn) : one(data.ibm_resource_instance.existing_service_instance[*].crn)
}
