## Introduction

This module enables the use of RFC 2136 DNS updates with optional secret key based transaction authentication RFC 2845.

## Usage

### Specify the required Terraform variables
* `tectonic_base_domain`
* `tectonic_ddns_server`
* `tectonic_ddns_key_name`
* `tectonic_ddns_key_secret`
* `tectonic_ddns_key_algorithm`

### Comment out the existing DNS module in your platform e.g., `platforms/azure/main.tf`
```go
/*
module "dns" {
  source = "../../modules/dns/azure"

  etcd_count   = "${var.tectonic_etcd_count}"
  master_count = "${var.tectonic_master_count}"
  worker_count = "${var.tectonic_worker_count}"

  etcd_ip_addresses    = "${module.vnet.etcd_endpoints}"
  master_ip_addresses  = "${module.vnet.master_private_ip_addresses}"
  worker_ip_addresses  = "${module.vnet.worker_private_ip_addresses}"
  api_ip_addresses     = "${module.vnet.api_ip_addresses}"
  console_ip_addresses = "${module.vnet.console_ip_addresses}"

  base_domain  = "${var.tectonic_base_domain}"
  cluster_name = "${var.tectonic_cluster_name}"
  cluster_id   = "${module.tectonic.cluster_id}"

  location             = "${var.tectonic_azure_location}"
  external_dns_zone_id = "${var.tectonic_azure_external_dns_zone_id}"

  extra_tags = "${var.tectonic_azure_extra_tags}"
}
*/
```

### Configure the DDNS module in your platform e.g., `platforms/azure/main.tf`
```go
module "dns" {
  source = "../../modules/dns/ddns"

  etcd_count   = "${var.tectonic_etcd_count}"
  master_count = "${var.tectonic_master_count}"
  worker_count = "${var.tectonic_worker_count}"

  etcd_ip_addresses    = "${module.vnet.etcd_endpoints}"
  master_ip_addresses  = "${module.vnet.master_private_ip_addresses}"
  worker_ip_addresses  = "${module.vnet.worker_private_ip_addresses}"
  api_ip_addresses     = "${module.vnet.api_ip_addresses}"
  console_ip_addresses = "${module.vnet.console_ip_addresses}"

  base_domain  = "${var.tectonic_base_domain}"
  cluster_name = "${var.tectonic_cluster_name}"

  dns_server        = "${var.tectonic_ddns_server}"
  dns_key_name      = "${var.tectonic_ddns_key_name}"
  dns_key_secret    = "${var.tectonic_ddns_key_secret}"
  dns_key_algorithm = "${var.tectonic_ddns_key_algorithm}"
}
```
