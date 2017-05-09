# Tectonic Changelog

## Tectonic 1.6.2-tectonic.1 (2017-04-10)

* Switches provisioning methods on AWS & Bare-Metal to Terraform exclusively.
* Adds support for customizing the Tectonic infrastructure via Terraform.
* Introduces experimental support for self-hosted etcd using its operator, and associated UI.
* Adds Container Linux Update Operator(CLUO).
* Updates to Kubernetes v1.6.2.
* Updates to bootkube v0.4.2.
* GUI Installer with Terraform on AWS and bare-metal.
* Segregates control-plane / user workloads to master / worker nodes respectively.
* API server-to-etcd communication is secured over TLS.
* Removes locksmithd, etcd-gateway.
* Enables audit-logs for the API Server.
* Removes final manual installation step of copying over assets folder.

## Console

* Updates to Console v1.5.2.
* Adds binding name column to Role Bindings list pages
* Adds role binding name to fields searched by text filter
* Adds RBAC YAML editor
* Adds etcd cluster management pages

## Dex

* Updates to Dex v2.4.1.
* Adds support for login through SAML and GitHub Enterprise.

## Bug Fixes

* Fixes an issue where new nodes started automatically by auto-scalers would start with an outdated version of kubelet.
