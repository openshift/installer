locals {
  # NOTE: Defined in ./vpc.tf
  # prefix = var.cluster_id
}

# NOTE: Security group quota
# 5 per network interface (NIC) on a virtual server instance

############################################
# Security group (Cluster-wide)
############################################

resource "ibm_is_security_group" "cluster_wide" {
  name           = "${local.prefix}-security-group-cluster-wide"
  resource_group = var.resource_group_id
  tags           = var.tags
  vpc            = ibm_is_vpc.vpc.id
}

# SSH
resource "ibm_is_security_group_rule" "cluster_wide_ssh_inbound" {
  group     = ibm_is_security_group.cluster_wide.id
  direction = "inbound"
  remote    = ibm_is_security_group.cluster_wide.id
  tcp {
    port_min = 22
    port_max = 22
  }
}

# ICMP
resource "ibm_is_security_group_rule" "cluster_wide_icmp_inbound" {
  group     = ibm_is_security_group.cluster_wide.id
  direction = "inbound"
  remote    = ibm_is_security_group.cluster_wide.id
  icmp {}
}

# VXLAN and Geneve - port 4789
resource "ibm_is_security_group_rule" "cluster_wide_vxlan_geneve_4789_inbound" {
  group     = ibm_is_security_group.cluster_wide.id
  direction = "inbound"
  remote    = ibm_is_security_group.cluster_wide.id
  udp {
    port_min = 4789
    port_max = 4789
  }
}

# VXLAN and Geneve - port 6081
resource "ibm_is_security_group_rule" "cluster_wide_vxlan_geneve_6081_inbound" {
  group     = ibm_is_security_group.cluster_wide.id
  direction = "inbound"
  remote    = ibm_is_security_group.cluster_wide.id
  udp {
    port_min = 6081
    port_max = 6081
  }
}

# Outbound
resource "ibm_is_security_group_rule" "cluster_wide_outbound" {
  group     = ibm_is_security_group.cluster_wide.id
  direction = "outbound"
  remote    = "0.0.0.0/0"
}

############################################
# Security group (OpenShift network)
############################################

resource "ibm_is_security_group" "openshift_network" {
  name           = "${local.prefix}-security-group-openshift-network"
  resource_group = var.resource_group_id
  tags           = var.tags
  vpc            = ibm_is_vpc.vpc.id
}

# Host level services - TCP
resource "ibm_is_security_group_rule" "openshift_network_host_services_tcp_inbound" {
  group     = ibm_is_security_group.openshift_network.id
  direction = "inbound"
  remote    = ibm_is_security_group.openshift_network.id
  tcp {
    port_min = 9000
    port_max = 9999
  }
}

# Host level services - UDP
resource "ibm_is_security_group_rule" "openshift_network_host_services_udp_inbound" {
  group     = ibm_is_security_group.openshift_network.id
  direction = "inbound"
  remote    = ibm_is_security_group.openshift_network.id
  udp {
    port_min = 9000
    port_max = 9999
  }
}

# Kubernetes default ports
resource "ibm_is_security_group_rule" "openshift_network_kube_default_ports_inbound" {
  group     = ibm_is_security_group.openshift_network.id
  direction = "inbound"
  remote    = ibm_is_security_group.openshift_network.id
  tcp {
    port_min = 10250
    port_max = 10250
  }
}

# Kubernetes node ports - TCP
resource "ibm_is_security_group_rule" "openshift_network_node_ports_tcp_inbound" {
  group     = ibm_is_security_group.openshift_network.id
  direction = "inbound"
  remote    = ibm_is_security_group.openshift_network.id
  tcp {
    port_min = 30000
    port_max = 32767
  }
}

# Kubernetes node ports - UDP
resource "ibm_is_security_group_rule" "openshift_network_node_ports_udp_inbound" {
  group     = ibm_is_security_group.openshift_network.id
  direction = "inbound"
  remote    = ibm_is_security_group.openshift_network.id
  udp {
    port_min = 30000
    port_max = 32767
  }
}

############################################
# Security group (Kubernetes API LB)
############################################

resource "ibm_is_security_group" "kubernetes_api_lb" {
  name           = "${local.prefix}-security-group-kubernetes-api-lb"
  resource_group = var.resource_group_id
  tags           = var.tags
  vpc            = ibm_is_vpc.vpc.id
}

# Kubernetes API LB - inbound
resource "ibm_is_security_group_rule" "kubernetes_api_lb_inbound" {
  group     = ibm_is_security_group.kubernetes_api_lb.id
  direction = "inbound"
  remote    = var.public_endpoints ? "0.0.0.0/0" : ibm_is_security_group.cluster_wide.id
  tcp {
    port_min = 6443
    port_max = 6443
  }
}

# Kubernetes API LB - outbound
resource "ibm_is_security_group_rule" "kubernetes_api_lb_outbound" {
  group     = ibm_is_security_group.kubernetes_api_lb.id
  direction = "outbound"
  remote    = ibm_is_security_group.control_plane.id
  tcp {
    port_min = 6443
    port_max = 6443
  }
}

# Machine config server LB - inbound
resource "ibm_is_security_group_rule" "kubernetes_api_lb_machine_config_inbound" {
  group     = ibm_is_security_group.kubernetes_api_lb.id
  direction = "inbound"
  remote    = ibm_is_security_group.cluster_wide.id
  tcp {
    port_min = 22623
    port_max = 22623
  }
}

# Machine config server LB - outbound
resource "ibm_is_security_group_rule" "kubernetes_api_lb_machine_config_outbound" {
  group     = ibm_is_security_group.kubernetes_api_lb.id
  direction = "outbound"
  remote    = ibm_is_security_group.control_plane.id
  tcp {
    port_min = 22623
    port_max = 22623
  }
}

############################################
# Security group (Control plane)
############################################

resource "ibm_is_security_group" "control_plane" {
  name           = "${local.prefix}-security-group-control-plane"
  resource_group = var.resource_group_id
  tags           = var.tags
  vpc            = ibm_is_vpc.vpc.id
}

# etcd
resource "ibm_is_security_group_rule" "control_plane_etcd_inbound" {
  group     = ibm_is_security_group.control_plane.id
  direction = "inbound"
  remote    = ibm_is_security_group.control_plane.id
  tcp {
    port_min = 2379
    port_max = 2380
  }
}

# Kubernetes default ports
resource "ibm_is_security_group_rule" "control_plane_kube_default_ports_inbound" {
  group     = ibm_is_security_group.control_plane.id
  direction = "inbound"
  remote    = ibm_is_security_group.control_plane.id
  tcp {
    port_min = 10257
    port_max = 10259
  }
}

# Kubernetes API - inbound
resource "ibm_is_security_group_rule" "control_plane_kubernetes_api_inbound" {
  group     = ibm_is_security_group.control_plane.id
  direction = "inbound"
  remote    = ibm_is_security_group.cluster_wide.id
  tcp {
    port_min = 6443
    port_max = 6443
  }
}

# Kubernetes API - inbound via LB
resource "ibm_is_security_group_rule" "control_plane_kubernetes_api_lb_inbound" {
  group     = ibm_is_security_group.control_plane.id
  direction = "inbound"
  remote    = ibm_is_security_group.kubernetes_api_lb.id
  tcp {
    port_min = 6443
    port_max = 6443
  }
}

# Machine config server - inbound via LB
resource "ibm_is_security_group_rule" "control_plane_machine_config_lb_inbound" {
  group     = ibm_is_security_group.control_plane.id
  direction = "inbound"
  remote    = ibm_is_security_group.kubernetes_api_lb.id
  tcp {
    port_min = 22623
    port_max = 22623
  }
}
