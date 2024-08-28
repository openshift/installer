output "transit_gateway" {
  value = var.tg_name == "" ? one(resource.ibm_tg_gateway.transit_gateway[*].id) : one(data.ibm_tg_gateway.existing_transit_gateway[*].id)
}
