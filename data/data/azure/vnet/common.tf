# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

// Only reference data sources which are gauranteed to exist at any time (above) in this locals{} block
locals {
  // The VPC ID to use to build the rest of the vpc data sources
  vnet_id = "${azurerm_virtual_network.new_vnet.id}"

  subnet_ids    = "${azurerm_subnet.public_subnet.id}"

  lb_fqdn = "${azurerm_lb.external_lb.id}"
  
  lb_backend_pool_id ="${azurerm_lb_backend_address_pool.master_elb_pool.id}"

  external_lb_id = "${azurerm_lb.external_lb.id}"
}
