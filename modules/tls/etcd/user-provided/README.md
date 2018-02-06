# Introduction

This module enables user provided etcd certificates.
It does not contain any logic, but just passes user provided certificates from its input directly to its output. This is to prevent changing existing references to the `tls/etcd/signed` module, hence all `tls/etcd/*` modules share the same outputs.

For more information on using this module with Tectonic, please see [Enabling custom etcd TLS certificates][tls-etcd].


[tls-etcd]: https://coreos.com/tectonic/docs/latest/tls/tls-etcd.html
