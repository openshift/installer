## Introduction

This module enables user provided Tectonic Identity certificates.
It actually doesn't contain any logic, but just passes user provided certificates from its input directly to its output.
This is to prevent changing existing references to the `tls/identity/self-signed` module, hence all `tls/identity/*` modules share
the same outputs.

## Usage

Comment out the existing self-signed identity TLS in your platform, i.e. `platforms/aws/tectonic.tf`:
```
/*
module "identity_certs" {
  source = "../../modules/tls/identity/self-signed"

  ca_cert_pem = "${module.kube_certs.ca_cert_pem}"
  ca_key_alg  = "${module.kube_certs.ca_key_alg}"
  ca_key_pem  = "${module.kube_certs.ca_key_pem}"
}
*/
```

Configure the user provided certificate paths in your platform, i.e. `platforms/aws/tectonic.tf`:
```
module "identity_certs" {
  source = "../../modules/tls/identity/user-provided"

  client_key_pem_path  = "/path/to/identity-client.key"
  client_cert_pem_path = "/path/to/identity-client.crt"
  server_key_pem_path  = "/path/to/identity-server.key"
  server_cert_pem_path = "/path/to/identity-server.crt"
}
```

The signed identity client certificate must have the following Key Usage associations:
```
$ openssl x509 -noout -text -in /path/to/identity-client.crt 
Certificate:
...
        X509v3 extensions:
            X509v3 Extended Key Usage: 
                TLS Web Client Authentication
```

The signed identity server certificate must have the following Key Usage associations:
```
$ openssl x509 -noout -text -in /path/to/identity-server.crt 
Certificate:
...
        X509v3 extensions:
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication
```

Note that the above identity certificates must use the same CA as the `modules/tls/kube` certificates.
