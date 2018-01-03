## Introduction

This module enables user provided etcd certificates.
It actually doesn't contain any logic, but just passes user provided certificates from its input directly to its output.
This is to prevent changing existing references to the `tls/etcd/signed` module, hence all `tls/etcd/*` modules share
the same outputs.

## Usage

Comment out the existing signed etcd TLS in your platform, i.e. `platforms/aws/tectonic.tf`:
```
/*
module "etcd_certs" {
  source = "../../modules/tls/etcd"

  etcd_ca_cert_path     = "${var.tectonic_etcd_ca_cert_path}"
  etcd_cert_dns_names   = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_client_cert_path = "${var.tectonic_etcd_client_cert_path}"
  etcd_client_key_path  = "${var.tectonic_etcd_client_key_path}"
  self_signed           = "${var.tectonic_etcd_tls_enabled}"
  service_cidr          = "${var.tectonic_service_cidr}"
}
*/
```

Configure the user provided certificate paths in your platform, i.e. `platforms/aws/tectonic.tf`:
```
module "etcd_certs" {
  source = "../../modules/tls/etcd/user-provided"

  etcd_ca_crt_pem_path     = "/path/to/etcd-ca.crt"
  etcd_server_crt_pem_path = "/path/to/etcd-server.crt"
  etcd_server_key_pem_path = "/path/to/etcd-server.key"
  etcd_peer_crt_pem_path   = "/path/to/etcd-peer.crt"
  etcd_peer_key_pem_path   = "/path/to/etcd-peer.key"
  etcd_client_crt_pem_path = "/path/to/etcd-client.crt"
  etcd_client_key_pem_path = "/path/to/etcd-client.key"
}
```

The signed etcd server certificate must have the following Subject Alternative Name (SAN) and Key Usage associations:
```
$ openssl x509 -noout -text -in /path/to/etcd-server.crt 
Certificate:
...
        X509v3 extensions:
            X509v3 Key Usage: critical
                Key Encipherment
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication
...
            X509v3 Subject Alternative Name: 
                DNS:<tectonic_cluster_name>-etcd-0.<tectonic_base_domain>,
                DNS:<tectonic_cluster_name>-etcd-1.<tectonic_base_domain>,
                ...
                DNS:<tectonic_cluster_name>-etcd-<tectonic_etcd_count-1>.<tectonic_base_domain>,
```

The signed etcd peer certificate must have the following Subject Alternative Name (SAN) and Key Usage associations:
```
$ openssl x509 -noout -text -in /path/to/etcd-peer.crt 
Certificate:
        X509v3 extensions:
            X509v3 Key Usage: critical
                Key Encipherment
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication, TLS Web Client Authentication
...
            X509v3 Subject Alternative Name: 
                DNS:<tectonic_cluster_name>-etcd-0.<tectonic_base_domain>,
                DNS:<tectonic_cluster_name>-etcd-1.<tectonic_base_domain>,
                ...
                DNS:<tectonic_cluster_name>-etcd-<tectonic_etcd_count-1>.<tectonic_base_domain>,
```

The signed etcd client certificate must have the following Key Usage associations:
```
$ openssl x509 -noout -text -in /path/to/etcd-client.crt 
Certificate:
...
        X509v3 extensions:
            X509v3 Extended Key Usage: 
                TLS Web Client Authentication
```
