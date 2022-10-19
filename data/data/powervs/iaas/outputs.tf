output "cloud_instance_id" {
  value = var.powervs_cloud_instance_id == "" ? ibm_resource_instance.powervs_service_instance[0].guid : var.powervs_cloud_instance_id
}

