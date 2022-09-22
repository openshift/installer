# Install: Installer provisioned infrastructure using Shared VPC (XPN)
The steps for performing a installer provisioned infrastructure installation
using Shared VPC (XPN) networking resources.

## Create an base install config
Create the `install-config.yaml` manually using the documentation
references/examples as you would for a normal IPI Installation.

### Configure passthrough credentialsMode
The cluster will need permissions in the Host Project in order to function
properly. This will be accomplished by configuring the service account used by
the installer to have these permissions and configuring the credentialsMode to
be passthrough. This configures the cluster to use the same service account
internally after installation. The exact permissions vary depending on choices
selected throughout this document as specified below.

```
credentialsMode: Passthrough
```

### Add pre-exsiting XPN network.
The following parameters are necessary to configure the installer to use
pre-exsting network resources, and specifcy that they exist in the
networkProjectID project, also known as the "Host Project".

```
platform:
  gcp:
    computeSubnet: <compute-subnet>
    controlPlaneSubnet: <control-plane-subnet>
    network: <network>
    networkProjectID: <networkProjectID>
```

The installer service account requires the following permissions on the network
in the Host Project in order to enable the cluster to use the network resources.

- Compute Network User

## Configure cloud firewall rules

Because firewall rules are attached to the network, they must exist in the Host
Project. As a result, additional permissions would be necessary to manage these
rules. Choose one of the options below based on if you wish to enable these
permissions:

### Option A (Default): Automatic firewall rules (requires permissions)

The default behavior is to automatically create and manage firewall rules. This
option requires the installer service account to have the following permissions
in the Host Project:

- Compute Network Admin
- Compute Security Admin

### Option B: Manually created, pre-exsiting firewall rules

When these permissions are not desired, it is possible to create firwall rules
in advance and configure the installer to skip their creation. It is important
to ensure the rules are sufficient to enable all cluster communications.
Creation of firewall rules can be disabled by adding the following to the
install-config.

```
platform:
  gcp:
    createFirewallRules: Disabled
```

The pre-existing rules can be cofigured specifically for the cluster instances
by using network tags in combination with the following tags paremeters in the
install-config:

```
compute:
- name: worker
  platform:
    gcp:
      tags:
      - mycomputetag
controlPlane:
  platform:
    gcp:
      tags:
      - mycontrolplanetag
platform:
  gcp:
    defaultMachinePlatform:
      tags:
      - mydefaulttag
```

## Configure Public DNS Zone

The public DNS zone can exist in either the Service Project or any other
project. Because the public zone is authorative, it is common for it to exist
in the Host Project.  As a result, when the public DNS Zone is not in the
Service Project, additional premissions would be necessary to manage these
resources. Choose one of the options below based on where your public DNS Zone
exists, and if you wish to enable these permissions:

### Option A (Default): Automatic in Service Project

The default behavior is to find the DNS Zone in the Service Project and use it
to create and manage DNS records. No special permissions are required.

### Option B: Automatic in Host Project (Requires Permissions)

Use the following install-config parameters to configure the cluster to use a
DNS Zone in the Host Project.

```
platform:
  gcp:
    publicDNSZone:
      id: my-public-zone-id
      project: my-project-id

```

This option requires the installer service account to have the following
permissions in the specified project:

- DNS Administrator

### Option C: Manual/External DNS management.

This option is not currently supported.

## Configure private DNS Zone

The private DNS zone has the same considerations as the public zone. However,
it is highly recommended to leave the private DNS zone in the service project
as configured by default. This should not generally be a problem as these zones
map 1 to 1 with clusters.

## Create Cluster

```console
[~]$ openshift-install create cluster

INFO Consuming Install Config from target directory
INFO Creating infrastructure resources...
INFO Waiting up to 30m0s for the Kubernetes API at https://api.mycluster.example.com:6443...
INFO API v1.14.0+37982ca up
INFO Waiting up to 30m0s for bootstrapping to complete...
INFO Destroying the bootstrap resources...
INFO Waiting up to 30m0s for the cluster at https://api.mycluster.example.com:6443 to initialize...
INFO Waiting up to 10m0s for the openshift-console route to be created...
INFO Install complete!
INFO To access the cluster as the system:admin user when using 'oc', run
    export KUBECONFIG=/home/user/auth/kubeconfig
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.mycluster.example.com
INFO Login to the console with user: kubeadmin, password: 5char-5char-5char-5char
```

## Running Cluster

In your GCP project, there will be a new private DNS zone (for internal lookups)

There will be six running VM instances in the Project.

The nodes within the Virtual Network utilize the internal DNS and use the Router and External API load balancers. External/Internet
access to the cluster use the Router and External API load balancers. Compute instances are spread equally across all running availability
zones for the region.

The OpenShift console is available via the kubeadmin login provided by the installer.
