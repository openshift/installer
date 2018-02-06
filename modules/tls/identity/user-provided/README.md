# Introduction

This module enables user provided Tectonic Identity certificates.
It does not contain any logic, but just passes user provided certificates from its input directly to its output. This is to prevent changing existing references to the `tls/identity/self-signed` module, hence all `tls/identity/*` modules share the same outputs.

For more information on using this module with Tectonic, please see [Enabling custom Tectonic Identity TLS certificates][tls-identity].


[tls-identity]: https://coreos.com/tectonic/docs/latest/tls/tls-identity.html
