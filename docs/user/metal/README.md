# Support for Bare Metal Environments

OpenShift has support for bare metal deployments with either [User
provided infrastructure (UPI)](install_upi.md), or [Installer-provided
infrastructure (IPI)](install_ipi.md).

The following is a summary of key differences:

* UPI bare metal
  * Provisioning hosts is an external requirement
  * Requires extra DNS configuration
  * Requires setup of load balancers
  * Offers more control and choice over infrastructure

* IPI bare metal
  * Has built-in hardware provisioning components, will provision nodes with RHCOS automatically,
    and supports the Machine API for ongoing management of these hosts.
  * Automates internal DNS requirements
  * Automates setup of self-hosted load balancers
  * Supports “openshift-install create cluster” for bare metal environments
    using this infrastructure automation, but requires the use of compatible
    hardware, as described in [install_ipi.md](install_ipi.md).
