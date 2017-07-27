# Tectonic Changelog

## Current

### Console

* Kubectl-config cluster name is set to the Tectonic cluster name, defaulting to "tectonic" for backwards-compatibility

### Tectonic Installer

* Internet gateways and etcd node root volumes are tagged with default and user-supplied tags
* The bootkube kubeconfigs' cluster name is set to the Tectonic cluster name

## Tectonic 1.6.7-tectonic.2 (2017-07-28)

* Allow updates to the Tectonic v1.7.x releases.

## Tectonic 1.6.7-tectonic.1 (2017-07-11)

* Updates to Kubernetes v1.6.7.
* Update operators are available to all users to power automated operations
* Reduced flapping of node `NotReady`status
   - Increased controller manager health time out to be greater than the TTL of the load balancer DNS entry
   - Kubernetes default of 40s is below the minimum TTL of 60s for many platforms

### Console

* Updates to Console v1.7.4.
* All tables have sortable columns
* Removed broken Horizontal Pod Autoscalers UI
* Adds autocomplete for RBAC binding form dropdowns
* Adds ability to edit and duplicate RBAC bindings
* Adds RBAC edit binding roles dropdown filtering by namespace
* Improved support for valueless labels and annotations

### Tectonic Installer

* Installer will generate all TLS certificates for etcd
* Terraform tfvars are now pretty-printed

## Tectonic 1.6.4-tectonic.1 (2017-06-08)

* Updates to Terraform v0.9.6 (fixes some instances of `terraform destroy` not working).
* Updates to Kubernetes v1.6.4.
* Many components run as "nobody" instead of root.
* An option has been added to disable the creation of private zones.
* All resources are now tagged in AWS with the cluster id.
* A minimal IAM policy has been created.

### Console

* Updates to Console v1.6.3.
* CPU usage graphs now display usage instead of limits.
* Can now Create Role Bindings and many other supported resources.

### Tectonic Channel Operator

* Updates to Tectonic Channel Operator v0.3.4
* Requires signed payloads using the default CoreOS key.
* No longer creates components upon upgrade when they did not previously exist.

## Tectonic 1.6.2-tectonic.1 (2017-04-10)

Tectonic now uses Terraform for cluster installation. This supports greater customization of environments, enables scripted installs and generally makes it easier to manage the lifecycle of multiple clusters.

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

### Console

Role-based Access Control screens have been redesigned to make it easier to securely grant access to your clusters.

* Updates to Console v1.5.2.
* Adds binding name column to Role Bindings list pages
* Adds role binding name to fields searched by text filter
* Adds RBAC YAML editor
* Adds etcd cluster management pages

### Dex

* Updates to Dex v2.4.1.
* Adds support for login through SAML and GitHub Enterprise.

### Bug Fixes

* Fixes an issue where new nodes started automatically by auto-scalers would start with an outdated version of kubelet.
