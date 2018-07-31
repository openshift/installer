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

## Previous work

### Tectonic Installer

Tectonic-installer in `track-2` used a `Cluster` object to make install time configurations.

#### `Cluster` object

```go
// Cluster defines the config for a cluster.
type Cluster struct {
    Admin           `json:",inline" yaml:"admin,omitempty"`
    aws.AWS         `json:",inline" yaml:"aws,omitempty"`
    BaseDomain      string `json:"tectonic_base_domain,omitempty" yaml:"baseDomain,omitempty"`
    CA              `json:",inline" yaml:"CA,omitempty"`
    ContainerLinux  `json:",inline" yaml:"containerLinux,omitempty"`
    Etcd            `json:",inline" yaml:"etcd,omitempty"`
    IgnitionEtcd    string `json:"tectonic_ignition_etcd,omitempty" yaml:"-"`
    IgnitionMaster  string `json:"tectonic_ignition_master,omitempty" yaml:"-"`
    IgnitionWorker  string `json:"tectonic_ignition_worker,omitempty" yaml:"-"`
    Internal        `json:",inline" yaml:"-"`
    libvirt.Libvirt `json:",inline" yaml:"libvirt,omitempty"`
    LicensePath     string `json:"tectonic_license_path,omitempty" yaml:"licensePath,omitempty"`
    Master          `json:",inline" yaml:"master,omitempty"`
    Name            string `json:"tectonic_cluster_name,omitempty" yaml:"name,omitempty"`
    Networking      `json:",inline" yaml:"networking,omitempty"`
    NodePools       `json:"-" yaml:"nodePools"`
    Platform        Platform `json:"tectonic_platform" yaml:"platform,omitempty"`
    PullSecretPath  string   `json:"tectonic_pull_secret_path,omitempty" yaml:"pullSecretPath,omitempty"`
    Worker          `json:",inline" yaml:"worker,omitempty"`
}
```

#### AWS specfic configuration

```go
// AWS converts AWS related config.
type AWS struct {
    AutoScalingGroupExtraTags []map[string]string `json:"tectonic_autoscaling_group_extra_tags,omitempty" yaml:"autoScalingGroupExtraTags,omitempty"`
    EC2AMIOverride            string              `json:"tectonic_aws_ec2_ami_override,omitempty" yaml:"ec2AMIOverride,omitempty"`
    Endpoints                 Endpoints           `json:"tectonic_aws_endpoints,omitempty" yaml:"endpoints,omitempty"`
    Etcd                      `json:",inline" yaml:"etcd,omitempty"`
    External                  `json:",inline" yaml:"external,omitempty"`
    ExtraTags                 map[string]string `json:"tectonic_aws_extra_tags,omitempty" yaml:"extraTags,omitempty"`
    InstallerRole             string            `json:"tectonic_aws_installer_role,omitempty" yaml:"installerRole,omitempty"`
    Master                    `json:",inline" yaml:"master,omitempty"`
    Profile                   string `json:"tectonic_aws_profile,omitempty" yaml:"profile,omitempty"`
    Region                    string `json:"tectonic_aws_region,omitempty" yaml:"region,omitempty"`
    SSHKey                    string `json:"tectonic_aws_ssh_key,omitempty" yaml:"sshKey,omitempty"`
    VPCCIDRBlock              string `json:"tectonic_aws_vpc_cidr_block,omitempty" yaml:"vpcCIDRBlock,omitempty"`
    Worker                    `json:",inline" yaml:"worker,omitempty"`
}

// External converts external related config.
type External struct {
    MasterSubnetIDs []string `json:"tectonic_aws_external_master_subnet_ids,omitempty" yaml:"masterSubnetIDs,omitempty"`
    PrivateZone     string   `json:"tectonic_aws_external_private_zone,omitempty" yaml:"privateZone,omitempty"`
    VPCID           string   `json:"tectonic_aws_external_vpc_id,omitempty" yaml:"vpcID,omitempty"`
    WorkerSubnetIDs []string `json:"tectonic_aws_external_worker_subnet_ids,omitempty" yaml:"workerSubnetIDs,omitempty"`
}

// Master converts master related config.
type Master struct {
    CustomSubnets    map[string]string `json:"tectonic_aws_master_custom_subnets,omitempty" yaml:"customSubnets,omitempty"`
    EC2Type          string            `json:"tectonic_aws_master_ec2_type,omitempty" yaml:"ec2Type,omitempty"`
    ExtraSGIDs       []string          `json:"tectonic_aws_master_extra_sg_ids,omitempty" yaml:"extraSGIDs,omitempty"`
    IAMRoleName      string            `json:"tectonic_aws_master_iam_role_name,omitempty" yaml:"iamRoleName,omitempty"`
    MasterRootVolume `json:",inline" yaml:"rootVolume,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
    CustomSubnets    map[string]string `json:"tectonic_aws_worker_custom_subnets,omitempty" yaml:"customSubnets,omitempty"`
    EC2Type          string            `json:"tectonic_aws_worker_ec2_type,omitempty" yaml:"ec2Type,omitempty"`
    ExtraSGIDs       []string          `json:"tectonic_aws_worker_extra_sg_ids,omitempty" yaml:"extraSGIDs,omitempty"`
    IAMRoleName      string            `json:"tectonic_aws_worker_iam_role_name,omitempty" yaml:"iamRoleName,omitempty"`
    LoadBalancers    []string          `json:"tectonic_aws_worker_load_balancers,omitempty" yaml:"loadBalancers,omitempty"`
    WorkerRootVolume `json:",inline" yaml:"rootVolume,omitempty"`
}
```

#### libvirt specfic configuration

```go
type Libvirt struct {
    URI           string `json:"tectonic_libvirt_uri,omitempty" yaml:"uri"`
    SSHKey        string `json:"tectonic_libvirt_ssh_key,omitempty" yaml:"sshKey"`
    QCOWImagePath string `json:"tectonic_coreos_qcow_path,omitempty" yaml:"imagePath"`
    Network       `json:",inline" yaml:"network"`
    MasterIPs     []string `json:"tectonic_libvirt_master_ips,omitempty" yaml:"masterIPs"`
}

// Network describes a libvirt network configuration.
type Network struct {
    Name      string `json:"tectonic_libvirt_network_name,omitempty" yaml:"name"`
    IfName    string `json:"tectonic_libvirt_network_if,omitempty" yaml:"ifName"`
    DNSServer string `json:"tectonic_libvirt_resolver,omitempty" yaml:"dnsServer"`
    IPRange   string `json:"tectonic_libvirt_ip_range,omitempty" yaml:"ipRange"`
}
```

### Cluster Operator

Cluster Operator also defines a `ClusterDeployment` object to represent a cluster managed by clusteroperator.

#### `ClusterDeployment` object

```go
type ClusterDeployment struct {
    metav1.TypeMeta
    metav1.ObjectMeta

    Spec ClusterDeploymentSpec
    Status ClusterDeploymentStatus
}

type ClusterDeploymentSpec struct {
    ClusterID string

    // Hardware specifies the hardware that the cluster will run on
    Hardware ClusterHardwareSpec

    // Config specifies cluster-wide OpenShift configuration
    Config ClusterConfigSpec
    DefaultHardwareSpec *MachineSetHardwareSpec

    // MachineSets specifies the configuration of all machine sets for the cluster
    MachineSets []ClusterMachineSet

    ClusterVersionRef ClusterVersionReference
}
```

`ClusterHardwareSpec` contains the global configurations options for platforms.

```go
type ClusterHardwareSpec struct {
    // AWS specifies cluster hardware configuration on AWS
    // +optional
    AWS *AWSClusterSpec

    // TODO: Add other cloud-specific Specs as needed
}
```

#### AWS specific configurations

```go
type AWSClusterSpec struct {
    Defaults *MachineSetAWSHardwareSpec
    AccountSecret corev1.LocalObjectReference
    SSHSecret corev1.LocalObjectReference
    SSHUser string
    SSLSecret corev1.LocalObjectReference

    KeyPairName string

    Region string
    VPCName string
    VPCSubnet string
}
```

#### Machine configuration using `ClusterMachineSet`

Machines in clusteroperator are defined using `ClusterMachineSet`.

```go
type ClusterMachineSet struct {
    ShortName string
    // MachineSetConfig is the configuration for the MachineSet
    MachineSetConfig
}

type MachineSetConfig struct {
    // NodeType is the type of nodes that comprise the MachineSet
    // TODO: remove in favor of upstream MachineTemplateSpec roles.
    NodeType NodeType

    // Infra indicates whether this machine set should contain infrastructure
    // pods
    // TODO: remove in favor of upstream MachineTemplateSpec roles.
    Infra bool

    // Size is the number of nodes that the node group should contain
    // TODO: remove in favor of upstream MachineSet and MachineDeployment replicas.
    Size int

    // Hardware defines what the hardware should look like for this
    // MachineSet. The specification will vary based on the cloud provider.
    // +optional
    Hardware *MachineSetHardwareSpec

    // NodeLabels specifies the labels that will be applied to nodes in this
    // MachineSet
    NodeLabels map[string]string
}
```

`MachineSetHardwareSpec` is used to define the machines on various platforms.

```go
type MachineSetHardwareSpec struct {
    AWS *MachineSetAWSHardwareSpec
}

type MachineSetAWSHardwareSpec struct {
    InstanceType string
}
```

## 4.0 Installer

### InstallConfig object

The installconfig object provides only necessary configurations options that are valuable for most users. Any extra user customization needs to happen in `render` and `prepare` phases of installer.

```go
type InstallConfig struct {
    // +optional
    metav1.TypeMeta   `json:",inline"`

    metav1.ObjectMeta `json:"metadata"`

    // ClusterID is the ID of the cluster.
    ClusterID       uuid.UUID `json:"clusterID"`

    // Networking defines the pod network provider in the cluster.
    Networking      `json:"networking"`

    // Machines is the list of MachinePools that need to be installed.
    Machines        []MachinePools `json:"machines"`

    // only one of the platform configuration should be set
    Platform    `json:"platform"`
}

type Platform struct {
    AWS     *AWSPlatform           `json:"aws,omitempty"`
    Libvirt *LibvirtPlatform       `json:"libvirt,omitempty"`
}

type Networking struct {
    Type        NetworkType `json:"type"`
    ServiceCIDR net.IPNet `json:"serviceCIDR"`
    PodCIDR     net.IPNet `json:"podCIDR"`
}

// NetworkType defines the pod network provider in the cluster.
type NetworkType string

const (
    // NetworkTypeOpenshiftSDN
    NetworkTypeOpenshiftSDN NetworkType = "openshift-sdn"
    // NetworkTypeOpenshiftOVN
    NetworkTypeOpenshiftOVN NetworkType = "openshift-ovn"
)

// AWS stores all the global configuration that
// all machinesets use.
type AWS struct {
    // Region specifies the AWS region where the cluster will be created.
    Region              string `json:"region"`

    // KeyPairName is the name of the AWS key pair to use for SSH access to EC2
    // instances in this cluster.
    KeyPairName         string `json:"keyPairName"`

    // VPCID specifies the vpc to associate with the cluster.
    // If empty, new vpc will be created.
    // +optional
    VPCID               string `json:"vpcID"`

    // VPCCIDRBlock
    // +optional
    VPCCIDRBlock        string `json:"vpcCIDRBlock"`
}

// Libvirt stores all the global configuration that
// all machinesets use.
type Libvirt struct {
    // URI
    URI           string `json:"URI"`

    // SSHKey
    SSHKey         string `json:"sshKey"`

    // Network
    Network       `json:"network"`

    // MasterIPs
    MasterIPs     []net.IP `json:"masterIPs"`
}

type LibvirtNetwork struct {
    Name      string `json:"name"`
    IfName    string `json:"if"`
    DNSServer string `json:"resolver"`
    IPRange   string `json:"ipRange"`
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
