locals {
  tcp_ports = concat([22623, 10258, 6443, 22], var.enable_snat ? [] : [5000])
}

resource "ibm_is_security_group" "ocp_security_group" {
  name           = "${var.cluster_id}-ocp-sec-group"
  resource_group = data.ibm_resource_group.resource_group.id
  vpc            = var.vpc_id
  tags           = [var.cluster_id]
}

resource "ibm_is_security_group_rule" "inbound_ports" {
  count     = length(local.tcp_ports)
  group     = ibm_is_security_group.ocp_security_group.id
  direction = "inbound"
  tcp {
    port_min = local.tcp_ports[count.index]
    port_max = local.tcp_ports[count.index]
  }
}

resource "ibm_is_security_group_rule" "outbound_any" {
  group     = ibm_is_security_group.ocp_security_group.id
  direction = "outbound"
}
