## Introduction

This module enables user provided Kubernetes certificates.
It actually doesn't contain any logic, but just passes user provided certificates from its input directly to its output.
This is to prevent changing existing references to the `tls/kube/self-signed` module, hence all `tls/kube/*` modules share
the same outputs.

## Usage

Comment out the existing self-signed kube TLS in your platform, i.e. `platforms/aws/tectonic.tf`:
```
/*
module "kube_certs" {
  source = "../../modules/tls/kube/self-signed"

  ca_cert_pem        = "${var.tectonic_ca_cert}"
  ca_key_alg         = "${var.tectonic_ca_key_alg}"
  ca_key_pem         = "${var.tectonic_ca_key}"
  kube_apiserver_url = "https://${module.masters.api_internal_fqdn}:443"
  service_cidr       = "${var.tectonic_service_cidr}"
}
*/
```

Configure the user provided certificate paths in your platform, i.e. `platforms/aws/tectonic.tf`:
```
module "kube_certs" {
  source = "../../modules/tls/kube/user-provided"

  ca_cert_pem_path        = "/path/to/ca.crt"
  kubelet_cert_pem_path   = "/path/to/kubelet.crt"
  kubelet_key_pem_path    = "/path/to/kubelet.key"
  apiserver_cert_pem_path = "/path/to/apiserver.crt"
  apiserver_key_pem_path  = "/path/to/apiserver.key"
}
```

The signed kubelet certificate must have the following key usage associations:
```
$ openssl x509 -noout -text -in /path/to/kubelet.crt
Certificate:
...
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication, TLS Web Client Authentication
...
```

The signed API certificate must have the following Subject Alternative Name (SAN) and Key Usage associations:
```
$ openssl x509 -noout -text -in /path/to/apiserver.crt 
Certificate:
...
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication, TLS Web Client Authentication
...
            X509v3 Subject Alternative Name: 
                DNS:<tectonic_cluster_name>-api.<tectonic_base_domain>,
                DNS:kubernetes,
                DNS:kubernetes.default,
                DNS:kubernetes.default.svc,
                DNS:kubernetes.default.svc.cluster.local,
                IP Address:10.3.0.1
```

The IP address `10.3.0.1` depends on the `tectonic_service_cidr` value. The API server always is assumed to be host number `1` of the configured service CIDR.
