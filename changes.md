# Tectonic Changelog

## Tectonic 1.7.3-tectonic.1 (2017-09-05)

### Core Components
 - Updates to Kubernetes v1.7.3

### Console
 - New Quick Start Guide for new users to Console.
 -  - Namespace selector shows only namespaces scoped to the user. Useful for using a restricted RBAC role.
 - Console redirects to the desired page after login rather than the cluster status page.
 - Improved ability to editing the YAML definition of Prometheus instances.
 - Improved automated operations UI for updates in progress

### Tectonic Installer
 - AWS accounts with a large amount of hosted zones are now paginated properly
 - UI enhancements around installation progress

### Known Issues
 - StatefulSet rolling updates must be executed manually. [More details](https://github.com/coreos/tectonic-docs/blob/master/Documentation/troubleshooting/tectonic-upgrade.md#upgrading-statefulsets).
 - Existing VPCs must be tagged manually for the AWS cloud provider to work correctly. [More details](https://github.com/coreos/tectonic-docs/blob/master/Documentation/install/aws/requirements.md#using-an-existing-vpc).

## Tectonic 1.7.1-tectonic.2 (2017-08-17)

* Makes the Container Linux instances on Azure start on the latest available version
* Fixes the tooltip preventing editing of CIDR inputs in the Installer
* Fixes a validation issue for STS tokens in the Installer
* Constrains the updater to go through every available versions

## Tectonic 1.6.8-tectonic.1 (2017-08-17)

* Updates to Kubernetes v1.6.8
* Constrains the updater to go through every available versions
 
## Tectonic 1.7.1-tectonic.1 (2017-08-09)

* Updates to Kubernetes v1.7.1.
* Support for Azure is Stable.

### Console

* Multiple update channels can be selected. See instructions below for additional details about updating from 1.6.x to 1.7.1.
  * 1.7-preproduction is available for testing and all non-production environments
  * 1.7-production should be used for all production environments
* Downloadable kubeconfigs now set the context name to the cluster name provided during installation, defaulting to "tectonic" for backwards-compatibility
* Added ability to view and configure Prometheus clusters run by the Prometheus operator
* Added ability to view Prometheus AlertManager configuration

### Tectonic Installer

* Container download and start up progress is output when booting a cluster
* Internet gateways and etcd node root volumes are tagged with default and user-supplied tags
* The bootkube kubeconfigs' cluster name is set to the Tectonic cluster name provided during installation

### Upgrade Notes - Requires 1.6.7-tectonic.2

To upgrade to Tectonic 1.7.1-tectonic.1, you must first update to `1.6.7-tectonic.2`. Once running `1.6.7-tectonic.2`, change the update channel to `1.7-preproduction` or `1.7-production` and click "Check for update". Update packages will be released to these channels in a rolling fashion, and will often be available on one channel but not the other.

If you encounter an error, confirm that you are running `1.6.7-tectonic.2` before reading the [troubleshooting guide](https://github.com/coreos/tectonic-installer/blob/master/Documentation/troubleshooting/tectonic-upgrade.md#upgrading-to-171-tectonic1).

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
