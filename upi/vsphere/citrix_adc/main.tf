provider "citrixadc" {
    username = "${var.citrix_adc_username}"
    password = "${var.citrix_adc_password}"
    endpoint = format("http://%s/", var.citrix_adc_ip)
    insecure_skip_verify = true
}

locals {
  api_server_sg_members = formatlist("%s:6443", var.api_backend_addresses)
  machine_config_sg_members = formatlist("%s:22623", var.api_backend_addresses)
  ingress_http_sg_members = formatlist("%s:80", var.ingress_backend_addresses)
  ingress_https_sg_members = formatlist("%s:443", var.ingress_backend_addresses)
}

// API Server Load-Balancing
resource "citrixadc_lbvserver" "openshift_api_server" {
  name = "openshift_lb_api_server"
  ipv46 = "${var.lb_ip_address}"
  port = "${var.api_server["lb_port"]}"
  servicetype = "${var.api_server["lb_protocol"]}"
}
resource "citrixadc_servicegroup" "openshift_api_server" {
  servicegroupname = "openshift_sg_api_server"
  servicetype = "${var.api_server["lb_protocol"]}"
  lbvservers = ["${citrixadc_lbvserver.openshift_api_server.name}"]
  servicegroupmembers = local.api_server_sg_members
}

// Machine Config Server Load-Balancing
resource "citrixadc_lbvserver" "openshift_machine_config_server" {
  name = "openshift_lb_machine_config_server"
  ipv46 = "${var.lb_ip_address}"
  port = "${var.machine_config_server["lb_port"]}"
  servicetype = "${var.machine_config_server["lb_protocol"]}"
}
resource "citrixadc_servicegroup" "openshift_machine_config_server" {
  servicegroupname = "openshift_sg_machine_config_server"
  servicetype = "${var.machine_config_server["lb_protocol"]}"
  lbvservers = ["${citrixadc_lbvserver.openshift_machine_config_server.name}"]
  servicegroupmembers = local.machine_config_sg_members
}

// HTTP Ingress Load-Balancing
resource "citrixadc_lbvserver" "openshift_ingress_http" {
  name = "openshift_lb_ingress_http"
  ipv46 = "${var.lb_ip_address}"
  port = "${var.ingress_http["lb_port"]}"
  servicetype = "${var.ingress_http["lb_protocol"]}"
}
resource "citrixadc_servicegroup" "openshift_ingress_http" {
  servicegroupname = "openshift_sg_ingress_http"
  servicetype = "${var.ingress_http["lb_protocol"]}"
  lbvservers = ["${citrixadc_lbvserver.openshift_ingress_http.name}"]
  servicegroupmembers = local.ingress_http_sg_members
}

// HTTPS Ingress Load-Balancing
resource "citrixadc_lbvserver" "openshift_ingress_https" {
  name = "openshift_lb_ingress_https"
  ipv46 = "${var.lb_ip_address}"
  port = "${var.ingress_https["lb_port"]}"
  servicetype = "${var.ingress_https["lb_protocol"]}"
}
resource "citrixadc_servicegroup" "openshift_ingress_https" {
  servicegroupname = "openshift_sg_ingress_https"
  servicetype = "${var.ingress_https["lb_protocol"]}"
  lbvservers = ["${citrixadc_lbvserver.openshift_ingress_https.name}"]
  servicegroupmembers = local.ingress_https_sg_members
}
