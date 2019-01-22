How to add a new platform to OpenShift
======================================

This document describes the process for adding a new platform to OpenShift in stages. Because there are many components to an automated platform, the process is defined in terms of delivering specific levels of functionality over time.

Terminology:

* `User Provided Infrastructure (UPI)` - The documentation and enablement that describes how to launch OpenShift on a particular platform following the supported best practices.
* `Install Provided Infrastructure (IPI)` - Infrastructure created as part of the install process following the best practices on a platform. IPI may support options that allow portions of the infrastructure to be user-provided.
* `Cloud Provider` - The set of controllers in OpenShift that automatically manage storage, networking, and host failure detection by invoking infrastructure APIs.
* `Dynamic Compute` - The 4.x cloud provider feature that allows OpenShift to automatically manage creating, deleting, and scaling nodes by invoking infrastructure APIs. Exposed via the Machine API (`Machine`, `MachineSet`, and `MachineDeployment`) and implemented per platform via an `actuator` controller.

The general order of adding a new platform for OpenShift is:

### Enable core platform

1. **Boot** - Ensure RH CoreOS boots on the desired platform, that Ignition works, and that you have VM / machine images to test with
2. **Arch** - Identify the correct opinionated configuration for a desired platform supporting the default features.
3. **CI** - Identify credentials and setup for a CI environment, ensure those credentials exist and can be used in the CI enviroment
4. **Name** - Identify and get approved the correct naming for adding a new platform to the core API objects (specifically the [infrastructure config](https://github.com/openshift/api/blob/master/config/v1/types_infrastructure.go) and the installer config (https://github.com/openshift/installer/blob/master/pkg/types/aws/doc.go)) so that we are consistent
5. **Enable Provisioning** Add a hidden installer option to this repo for the desired platform as a PR and implement the minimal features for bootstrap as well as a reliable teardown
6. **Enable Platform** Ensure all operators treat your platform as a no-op
7. **CI Job** Add a new CI job to the installer that uses the credentials above to run the installer against the platform and correctly tear down resources
8. **Publish Images** Ensure RH CoreOS images on the platform are being published to a location CI can test
9. **Merge** Merge the platform PR to the installer with a passing platform specific CI job

At this point the platform is said to be an `unsupported IPI` (installer provided infrastructure) install - work can begin enabling in other repositories. Once these steps have been completed and official documentation is available, the platform can said to be `supported UPI without cloud provider` (user provided infrastructure) for the set of options in **Arch**.

### Enable component infrastructure management

Once the platform can be launched and tested, system features must be implemented. The sections below are roughly independent:

* General requirements:
    * Replace the installer terraform destroy with one that doesn't rely on terraform state
    * Ensure the installer IPI support is consistent with other platform features (private config, etc)
        * Not all platforms will support all features, so IPI is taken to be a spectrum of support
    * Enable a CI job that verifies the e2e suite for the given platform runs successfully
* Requirements for dynamic storage and dynamic load balancing
    * Ensure the cloud provider is enabled in Kubernetes for your platform (this is required for `supported UPI with cloud provider`)
    * Enable cluster-storage-operator to set the correct default storage class
* Requirements for dynamic compute:
    * Enable the cloud credential operator for the platform to subdivide access for individual operators
    * Enable dynamic compute (MachineSets) by adding a cloud actuator for that platform
* Requirements for dynamic ingress and images:
    * Enable cluster-image-registry-operator to provision a storage bucket (if your platform supports object storage)
    * Enable cluster-ingress-operator to provision the wildcard domain names

At this point the platform is said to be a `supported IPI with Dynamic Compute` if the platform supports
MachineSets, or `supported IPI without Dynamic Compute` if it does not.


OpenShift Architectural Design
------------------------------

OpenShift 4 combines Kubernetes, fully-managed host lifecycle, dynamic infrastructure management, and a comprehensive set of fully-automated platform extensions that can be upgraded and managed uniformly.  The foundation of the platform is Red Hat CoreOS, an immutable operating system based on RHEL that is capable of acting as a fully-integrated part of the cluster. The Kubernetes control plane is hosted on the cluster, along with a number of other fundamental extensions like cluster networking, ingress, storage, and application development tooling. Each of those extensions is fully managed on-cluster via a cluster operator that reacts to top level global configuration APIs and can automatically reconfigure the affected components. The Operator Lifecycle Manager (OLM) allows additional ecosystem operators to be installed, upgraded and managed. All of these components - from the operating system kernel to the web console - are part of a unified update lifecycle under the top level Cluster Version Operator which offers worry free rolling updates of the entire infrastructure.

### Core configuration

An OpenShift cluster programs the infrastructure it runs on to provide operational simplicity. For every platform, the minimum requirements are:

1. The control plane nodes:
   1. Run RH CoreOS, allowing in-place updates
   2. Are fronted by a load balancer that allows raw TCP connections to port 6443 and exposes port 443
   3. Have low latency interconnections connections (<5ms RTT) and persistent disks that survive reboot and are provisoned for at least 300 IOPS
   4. Have cloud or infrastructure firewall rules that at minimum allow the standard ports to be opened (see AWS provider)
   5. Do *not* have automatic cloud provider permissions to perform infrastructure API calls
   6. Have a domain name pointing to the load balancer IP(s) that is `api.<BASE_DOMAIN>`
   7. Has an internal DNS CNAME pointing to each master called `etcd-N.<BASE_DOMAIN>` that 
   8. Has an optional internal load balancer that TCP load balances all master nodes, with a DNS name `internal-api.<BASE_DOMAIN>` pointing to the load balancer.
2. The bootstrap node:
   1. Runs RH CoreOS
   2. Is reachable by control plane nodes over the network
   3. Is part of the control plane load balancer until it is removed
   4. Can reach a network endpoint that hosts the bootstrap ignition file securely, or has the bootstrap ignition injected
3. All other compute nodes:
   1. Must be able to reach the internal IPs reported by the master nodes directly
   2. Have cloud or infrastructure firewall rules that at minimum allow ports 4789, 6443, 9000-10000, and 10250-10255 to be reachable

The following clarifications to configurations are noted:

1. The control plane load balancer does not need to be exposed to the public internet, but the DNS entry must be visible from the location the installer is run.
2. Master nodes are not required to expose external IPs for SSH access, but can instead allow SSH from a bastion inside a protected network.
3. Compute nodes do not require external IPs

For dynamic infrastructure, the following permissions are required to be provided as part of the install:

1. Service LoadBalancer - Load balancers can be created and removed, infastructure nodes can be queried
2. Dynamic Storage - New volumes can be created, deleted, attached, and detached from nodes. Snapshot creation is optional if the platform supports snapshotting
3. Dynamic Compute - New instances can be created, deleted, and restarted inside of the cluster's network / infrastructure, plus any platform specific constructs like programming instance groups for master load balancing on GCP.


Booting RH CoreOS
-----------------

Red Hat CoreOS uses ignition to receive initial configuration from a remote source. Ignition has platform specific behavior to read that configuration that is determined by the `oemID` embedded in the VM image.

To boot RHCoS to a new platform, you must:

1. Ensure [ignition](https://github.com/coreos/ignition) supports that platform
2. Ensure that RHCoS has any necessary platform specific code to communicate with the platform (for instance, on Azure the instance must periodically health check) - see [cloud support tracker on Fedora CoreOS for more info](https://github.com/coreos/fedora-coreos-tracker/issues/95).
3. Have a RHCoS image with the appropriate oemID tag set.

There is a script that assists you in converting the generic VM image to have a specific oemID set in the [coreos-assembler repo as gf-oemid](https://github.com/coreos/coreos-assembler/blob/master/src/gf-oemid). See the instructions there to create an image with the appropriate ID.

Once you have uploaded the image to your platform, and the machine stays up, you can begin porting the installer to have a minimal IPI.


Continuous Integration
----------------------

To enable a new platform, require a core continuous integration testing loop that verifies that new changes do not regress our support for the platform. The minimum steps required are:

1. Have an infrastructure that can receive API calls from the OpenShift CI system to provision/destroy instances
2. Support at minimum 3 concurrent clusters on that infrastructure as "per release image" testing (https://origin-release.svc.ci.openshift.org) that verify a release supports that platform
3. Also support a per-PR target that can be selectively run on the installer, core, and operator repositories in OpenShift in order to allow developers to test incremental changes to those components

No PR will be merged to openshift/installer for platform support that cannot satisfy the above steps.


Naming
------

The platform name will be part of our public API and must go through standard API review. The name
should be consistent with common usage of the platform and be recognizable to a consumer.

The following names for platforms are good examples of what is expected:

* Amazon Web Services -> `aws` or `AWS`
* Google Cloud Platform -> `gcp` or `GCP`
* Azure -> `azure` or `Azure`
* Libvirt -> `libvirt` or `Libvirt`
* OpenStack -> `openstack` or `OpenStack`


Enable Provisioning
-------------------

Since CI testing requires the ability to provision via an API, we define the basic path for supporting a platform as having a minimal provisioning path in the OpenShift installer. Not all platforms we support will have full infrastructure provisioning supported, but the basic path must be invokable via Go code in openshift-install before a platform can be certified. This ensures we have at least one path to installation.

The OpenShift installer has normal and hidden provisioners. The hidden provisioners are explicitly unsupported for production use but are supported for testing. 

1. Add a new hidden provisioner
2. Define the minimal platform parameters that the provisioner must support
3. Use Terraform or direct Go code to provision that platform via the credentials provided to the installer.

A minimal provisioner must be able to launch the control plane and bootstrap node via an API call and accept any "environmental" settings like network or region as inputs. The installer should use the Route53 DNS provisioning code to set up round robin to the bootstrap and control plane nodes if necessary.


Enable Platform
---------------

OpenShift handles platform functionality as a set of operators running on the platform that interface with users, admins, and infrastructure. Because operators handle day 2 reconfiguration of the cluster, many "installation" related duties are delegated to the operators.

Operators derive their configuration from top level API objects called `global configuration`. One such object is the `Infrastructure` global config, which reports which platform the cluster is running on.

All operators that react to infrastructure must support a `None` option, and any unrecognized infrastructure platform **MUST** be treated as `None`. When an operator starts, it should log a single warning if the infrastructure provider is not recognized and then fall back to `None`.

When adding a new platform to the installer, the infrastructure setting should happen automatically during bootstrapping, and if a component does not correctly treat your new platform as `None` it should be fixed immediately.


CI Job
------

The initial CI job for a new platform PR to `openshift/installer` must use the `cluster-installer-e2e` template but with an alternate profile, and the CI infrastructure should be configured with the credentials for your infrastructure in a `cluster-secrets-PLATFORM` secret. Talk to the testplatform team.
This CI job will then be reused whenever a repo wants to test, or when we add new release tests.

A new platform should pass many of the kubernetes conformance tests, so the default job would run the e2e suite `kubernetes/conformance`.  We may define a more scoped job if the platform cannot pass.

The teardown behavior of the cluster is the hardest part of this process - because we run so many tests a day, it must be 100% reliable from the beginning. You should implement a reliable teardown mechanism in your `destroy` method, leveraging the OpenStack and AWS examples.

We **will not** merge a new job if it does not have reliable cleanup in the face of failures, rate limits, etc, because it blocks other work.


Publishing Red Hat CoreOS Images
--------------------------------

RHCoS nodes can be upgraded to newer versions of kernel, userspace, and Kubelet post-creation. For this reason, the installer launches a recent version of RHCoS that is then upgraded at boot time to the version of RHCoS content that is included in the OpenShift release payload.

Once a version of RHCoS supports the desired platform, an image with that OEM ID embedded must be published to either to the cloud or a publicly available download location on a regular schedule. The installer may then embed logic to identify the most recent location for the payload and automatically provide that to the installer provisioning steps.


Merge the initial platform support
----------------------------------

After all of the steps above have been completed, the pull request enabling the platform may be merged with documentation updated to indicate the platform is in an unsupported pre-release configuration. Other components may now begin their integration work.


Integration to individual operators
-----------------------------------

1. Machine API Operator
2. Machine Config Operator
3. Cluster Storage Operator
4. Cloud Credential Operator
5. Cluster Ingress Operator
6. Cluster Image Registry Operator

TODO: add details
