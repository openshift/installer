# Introduction

This module enables user provided certificates to Tectonic Console (ingress).

It does not contain any logic, but just passes user provided certificates from its input directly to its output. This is to prevent changing existing references to the `tls/ingress/self-signed` module, hence all `tls/ingress/*` modules share the same outputs.

For more information on using this module with Tectonic, please see [Enabling custom Tectonic Console TLS certificates][tls-ingress].


[tls-ingress]: https://coreos.com/tectonic/docs/latest/tls/tls-ingress.html
