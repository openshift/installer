# InstallConfig

## Goals

1. Define InstallConfig object that will serve as the configuration for the 4.0 installer.
2. Must support multiple platforms. For example, AWS, Azure, local(libvirt) etc...
3. Must be a kubernetes Custom Resource Complaint object. This allows easy versioning and generation of the object.

## Overview

The installer requires configurations of mostly 2 types:

1. Cloud-infra specific options

    a. Certain options are used in cloud resource launch phase only but are not required to be persisted in cluster for future management.
    **For example VPC-id, Route53 zone-id, IAM roles in AWS are configurations that is required only at install time.**

    b. While some option especially related to the machines need to be persisted in cluster for cluster scaling management.
    **For example master/worker machine types, ELBs in AWS need to be persisted for scaling machinesets later on.**

2. Kubernetes cluster-level configuration options. **For example, API server URL address, cluster networking etc.**

## 4.0 Installer

### InstallConfig object

The installconfig object provides only necessary configurations options that are valuable for most users. Any extra user customization needs to happen in `render` and `prepare` phases of installer.

```go
// InstallConfig is the configuration for an OpenShift install.
type InstallConfig struct {
    // +optional
    metav1.TypeMeta `json:",inline"`

    metav1.ObjectMeta `json:"metadata"`

    // ClusterID is the ID of the cluster.
    ClusterID string `json:"clusterID"`

    // Admin is the configuration for the admin user.
    Admin Admin `json:"admin"`

    // BaseDomain is the base domain to which the cluster should belong.
    BaseDomain string `json:"baseDomain"`

    // Networking defines the pod network provider in the cluster.
    Networking `json:"networking"`

    // Machines is the list of MachinePools that need to be installed.
    Machines []MachinePool `json:"machines"`

    // Platform is the configuration for the specific platform upon which to
    // perform the installation.
    Platform `json:"platform"`

    // License is an OpenShift license needed to install a cluster.
    License string `json:"license"`

    // PullSecret is the secret to use when pulling images.
    PullSecret string `json:"pullSecret"`
}

// Admin is the configuration for the admin user.
type Admin struct {
    // Email is the email address of the admin user.
    Email string `json:"email"`
    // Password is the password of the admin user.
    Password string `json:"password"`
    // SSHKey to use for the access to compute instances.
    SSHKey string `json:"sshKey,omitempty"`
}

// Platform is the configuration for the specific platform upon which to perform
// the installation. Only one of the platform configuration should be set.
type Platform struct {
    // AWS is the configuration used when installing on AWS.
    AWS *AWSPlatform `json:"aws,omitempty"`
    // Libvirt is the configuration used when installing on libvirt.
    Libvirt *LibvirtPlatform `json:"libvirt,omitempty"`
}

// Networking defines the pod network provider in the cluster.
type Networking struct {
    Type        NetworkType `json:"type"`
    ServiceCIDR net.IPNet   `json:"serviceCIDR"`
    PodCIDR     net.IPNet   `json:"podCIDR"`
}

// NetworkType defines the pod network provider in the cluster.
type NetworkType string

const (
    // NetworkTypeOpenshiftSDN is used to install with SDN.
    NetworkTypeOpenshiftSDN NetworkType = "openshift-sdn"
    // NetworkTypeOpenshiftOVN is used to install with OVN.
    NetworkTypeOpenshiftOVN NetworkType = "openshift-ovn"
)

// AWSPlatform stores all the global configuration that
// all machinesets use.
type AWSPlatform struct {
    // Region specifies the AWS region where the cluster will be created.
    Region string `json:"region"`

    // VPCID specifies the vpc to associate with the cluster.
    // If empty, new vpc will be created.
    // +optional
    VPCID string `json:"vpcID"`

    // VPCCIDRBlock
    // +optional
    VPCCIDRBlock string `json:"vpcCIDRBlock"`
}

// LibvirtPlatform stores all the global configuration that
// all machinesets use.
type LibvirtPlatform struct {
    // URI
    URI string `json:"URI"`

    // Network
    Network LibvirtNetwork `json:"network"`

    // MasterIPs
    MasterIPs []net.IP `json:"masterIPs"`
}

// LibvirtNetwork is the configuration of the libvirt network.
type LibvirtNetwork struct {
    // Name is the name of the nework.
    Name string `json:"name"`
    // IfName is the name of the network interface.
    IfName string `json:"if"`
    // DNSServer is the name of the DNS server.
    DNSServer string `json:"resolver"`
    // IPRange is the range of IPs to use.
    IPRange string `json:"ipRange"`
}
```

### MachinePools

MachinePools are InstallConfig's absraction that only allows user to choose replica count and OS channel. This doesn't map to MachineSets.
Any user customization on `MachineSets` needs to happen at render phase and controlled by Machine Operator.

```go
type MachinePools struct {
    // Name is the name of the machine pool.
    Name string

    // This is count of machines for this machine pool.
    // Default is 1.
    Replicas *int64 `json:"replicas"`

    // PlatformConfig is configuration for machine pool specfic to the platfrom.
    PlatformConfig MachinePoolPlatformConfig `json: platformConfig`
}

type MachinePoolPlatformConfig struct {
    AWS *AWSMachinePoolPlatformConfig `json: "aws,omitempty"`
    Libvirt *LibvirtMachinePoolPlatformConfig `json:"libvirt,omitempty"`
}
```

AWS specific config options for machine pool.

```go
type AWSMachinePoolPlatformConfig struct {
    // InstanceType defines the ec2 instance type.
    // eg. m4-large
    InstanceType    string `json:"type"`

    // IAMRoleName defines the IAM role associated
    // with the ec2 instance.
    IAMRoleName     string `json:"iamRoleName"`

    // RootVolume defines the storage for ec2 instance.
    EC2RootVolume      `json:"rootVolume"`

}

type EC2RootVolume struct {
    IOPS int    `json:"iops"`
    Size int    `json:"size"`
    Type string `json:"type"`
}
```

Libvirt specific config options for machine pool.

```go
type LibvirtMachinePoolPlatformConfig struct {
    // QCOWImagePath
    QCOWImagePath string `json:"qcowImagePath"`
}
```

### Why not embed `MachineSets` in `InstallConfig`?

This installer is opinionated, and the `ClusterMachineSet` properties do not need to be configured by users.
