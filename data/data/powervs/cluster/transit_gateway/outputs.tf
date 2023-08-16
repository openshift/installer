output "transit_gateway" {
  value = one(resource.ibm_tg_gateway.transit_gateway[*].id)
}
