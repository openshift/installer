## Introduction

This module enables user provided certificates to Tectonic.
It actually doesn't contain any logic, but just passing user provided certificates from its input directly to its output.
This is to prevent changing existing references to the `tls/ingress/self-signed` module, hence all `tls/ingress/*` modules share
the same outputs.

## Usage

Comment out the existing self-signed ingress TLS in your platform, i.e. `platforms/aws/tectonic.tf`:
```
/*
module "ingress_certs" {
  source = "../../modules/tls/ingress/self-signed"

  base_address = "${module.masters.ingress_internal_fqdn}"
  ca_cert_pem  = "${module.kube_certs.ca_cert_pem}"
  ca_key_alg   = "${module.kube_certs.ca_key_alg}"
  ca_key_pem   = "${module.kube_certs.ca_key_pem}"
}
*/
```

Configure the user provided certificates in your platform, i.e. `platforms/aws/tectonic.tf`:
```
module "ingress_certs" {
  source = "../../modules/tls/ingress/user-provided"

  ca_cert_pem = <<EOF
-----BEGIN CERTIFICATE-----
<contents of the public CA certificate in PEM format>
-----END CERTIFICATE-----
EOF

  cert_pem = <<EOF
-----BEGIN CERTIFICATE-----
<contents of the public ingress certificate signed by the above CA in PEM format>
-----END CERTIFICATE-----
EOF

  key_pem = <<EOF
-----BEGIN RSA PRIVATE KEY-----
<contents of the private ingress key used to generate the above certificate PEM format>
-----END RSA PRIVATE KEY-----
EOF
}
```
