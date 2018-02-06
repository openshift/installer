# Introduction

This module enables user provided Kubernetes certificates.

It does not contain any logic, but just passes user provided certificates from its input directly to its output. This is to prevent changing existing references to the `tls/kube/self-signed` module, hence all `tls/kube/*` modules share the same outputs.

For more information on using this module with Tectonic, please see [Enabling custom Kubernetes TLS certificates][tls-kube].


[tls-kube]: https://coreos.com/tectonic/docs/latest/tls/tls-kube.html
