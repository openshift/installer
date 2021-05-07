output "pvc_name" {
  value = kubevirt_data_volume.data_volume.metadata[0].name
}
