package workflow

type configTemplateData struct {
	BaseDomain     string
	LicensePath    string
	Name           string
	PullSecretPath string
}

const configTemplate = `admin:
  email: "a@b.c"
  password: "verysecure"

aws:
  # (optional) Unique name under which the Amazon S3 bucket will be created. Bucket name must start with a lower case name and is limited to 63 characters.
  # The Tectonic Installer uses the bucket to store tectonic assets and kubeconfig.
  # If name is not provided the installer will construct the name using "name", current AWS region and "baseDomain"
  # assetsS3BucketName:

  # (optional) Extra AWS tags to be applied to created autoscaling group resources.
  # This is a list of maps having the keys ` + "`key`" + `, ` + "`value`" + ` and ` + "`propagate_at_launch`" + `.
  #
  # Example: ` + "`[ { key = \"foo\", value = \"bar\", propagate_at_launch = true } ]`" + `
  # autoScalingGroupExtraTags:

  # (optional) AMI override for all nodes. Example: ` + "`ami-foobar123`" + `.
  # ec2AMIOverride:

  etcd:
    # Instance size for the etcd node(s). Example: ` + "`t2.medium`" + `. Read the [etcd recommended hardware](https:#coreos.com/etcd/docs/latest/op-guide/hardware.html) guide for best performance
    ec2Type: t2.medium

    # (optional) List of additional security group IDs for etcd nodes.
    #
    # Example: ` + "`[\"sg-51530134\", \"sg-b253d7cc\"]`" + `
    # extraSGIDs:

    # (optional) Name of IAM role to use for the instance profiles of etcd nodes.
    # The name is also the last part of a role's ARN.
    #
    # Example:
    #  * Role ARN  = arn:aws:iam::123456789012:role/tectonic-installer
    #  * Role Name = tectonic-installer
    # iamRoleName:

    rootVolume:
      # The amount of provisioned IOPS for the root block device of etcd nodes.
      # Ignored if the volume type is not io1.
      iops: 100

      # The size of the volume in gigabytes for the root block device of etcd nodes.  size: 30 
      # The type of volume for the root block device of etcd nodes.
      type: gp2

  external:
    # (optional) List of subnet IDs within an existing VPC to deploy master nodes into.
    # Required to use an existing VPC and the list must match the AZ count.
    #
    # masterSubnetIDs:

    # (optional) If set, the given Route53 zone ID will be used as the internal (private) zone.
    # This zone will be used to create etcd DNS records as well as internal API and internal Ingress records.
    # If set, no additional private zone will be created.
    #
    # Example: ` + "`\"Z1ILINNUJGTAO1\"`" + `
    # privateZone:

    # (optional) ID of an existing VPC to launch nodes into.
    # If unset a new VPC is created.
    # Example: ` + "`vpc-123456`" + `
    # vpcID:

    # (optional) List of subnet IDs within an existing VPC to deploy worker nodes into.
    # Required to use an existing VPC and the list must match the AZ count.
    #
    # Example: ` + "`[\"subnet-111111\", \"subnet-222222\", \"subnet-333333\"]`" + `
    # workerSubnetIDs:

  # (optional) Extra AWS tags to be applied to created resources.
  #
  # Example: ` + "`{ \"key\" = \"value\", \"foo\" = \"bar\" }`" + `
  # extraTags:

  # (optional) Name of IAM role to use to access AWS in order to deploy the Tectonic Cluster.
  # The name is also the full role's ARN.
  #
  # Example:
  #  * Role ARN  = arn:aws:iam::123456789012:role/tectonic-installer
  # installerRole:

  master:
    # (optional) This configures master availability zones and their corresponding subnet CIDRs directly.
    #
    # Example:
    # ` + "`{ eu-west-1a = \"10.0.0.0/20\", eu-west-1b = \"10.0.16.0/20\" }`" + `
    # customSubnets:

    # Instance size for the master node(s). Example: ` + "`t2.medium`" + `.
    ec2Type: t2.medium

    # (optional) List of additional security group IDs for master nodes.
    #
    # Example: ` + "`[\"sg-51530134\", \"sg-b253d7cc\"]`" + `
    # extraSGIDs:

    # (optional) Name of IAM role to use for the instance profiles of master nodes.
    # The name is also the last part of a role's ARN.
    #
    # Example:
    #  * Role ARN  = arn:aws:iam::123456789012:role/tectonic-installer
    #  * Role Name = tectonic-installer
    # iamRoleName:

    rootVolume:
      # The amount of provisioned IOPS for the root block device of master nodes.
      # Ignored if the volume type is not io1.
      iops: 100

      # The size of the volume in gigabytes for the root block device of master nodes.
      size: 30

      # The type of volume for the root block device of master nodes.
      type: gp2

  # (optional) If set to true, create private-facing ingress resources (ELB, A-records).
  # If set to false, no private-facing ingress resources will be provisioned and all DNS records will be created in the public Route53 zone.
  # privateEndpoints: true

  # (optional) This declares the AWS credentials profile to use.
  # profile: default

  # (optional) If set to true, create public-facing ingress resources (ELB, A-records).
  # If set to false, no public-facing ingress resources will be created.
  # publicEndpoints: true

  # The target AWS region for the cluster.
  region: eu-west-1

  # Name of an SSH key located within the AWS region. Example: coreos-user.
  sshKey:

  # Block of IP addresses used by the VPC.
  # This should not overlap with any other networks, such as a private datacenter connected via Direct Connect.
  vpcCIDRBlock: 10.0.0.0/16

  worker:
    # (optional) This configures worker availability zones and their corresponding subnet CIDRs directly.
    #
    # Example: ` + "`{ eu-west-1a = \"10.0.64.0/20\", eu-west-1b = \"10.0.80.0/20\" }`" + `
    # customSubnets:

    # Instance size for the worker node(s). Example: ` + "`t2.medium`" + `.
    ec2Type: t2.medium

    # (optional) List of additional security group IDs for worker nodes.
    #
    # Example: ` + "`[\"sg-51530134\", \"sg-b253d7cc\"]`" + `
    # extraSGIDs:

    # (optional) Name of IAM role to use for the instance profiles of worker nodes.
    # The name is also the last part of a role's ARN.
    #
    # Example:
    #  * Role ARN  = arn:aws:iam::123456789012:role/tectonic-installer
    #  * Role Name = tectonic-installer
    # iamRoleName:

    # (optional) List of ELBs to attach all worker instances to.
    # This is useful for exposing NodePort services via load-balancers managed separately from the cluster.
    #
    # Example:
    #  * ` + "`[\"ingress-nginx\"]`" + `
    # loadBalancers:

    rootVolume:
      # The amount of provisioned IOPS for the root block device of worker nodes.
      # Ignored if the volume type is not io1.
      iops: 100

      # The size of the volume in gigabytes for the root block device of worker nodes.
      size: 30

      # The type of volume for the root block device of worker nodes.
      type: gp2

# The base DNS domain of the cluster. It must NOT contain a trailing period. Some
# DNS providers will automatically add this if necessary.
#
# Example: ` + "`openshift.example.com`" + `.
#
# Note: This field MUST be set manually prior to creating the cluster.
# This applies only to cloud platforms.
baseDomain: {{.BaseDomain}}

ca:
  # (optional) The content of the PEM-encoded CA certificate, used to generate Tectonic Console's server certificate.
  # If left blank, a CA certificate will be automatically generated.
  # cert:

  # (optional) The content of the PEM-encoded CA key, used to generate Tectonic Console's server certificate.
  # This field is mandatory if ` + "`ca_cert`" + ` is set.
  # key:

  # (optional) The algorithm used to generate ca_key.
  # The default value is currently recommended.
  # This field is mandatory if ` + "`ca_cert`" + ` is set.
  # keyAlg: RSA

containerLinux:
  # (optional) The Container Linux update channel.
  #
  # Examples: ` + "`stable`" + `, ` + "`beta`" + `, ` + "`alpha`" + `
  # channel: stable

  # The Container Linux version to use. Set to ` + "`latest`" + ` to select the latest available version for the selected update channel.
  #
  # Examples: ` + "`latest`" + `, ` + "`1465.6.0`" + `
  version: latest

  # (optional) A list of PEM encoded CA files that will be installed in /etc/ssl/certs on etcd, master, and worker nodes.
  # customCAPEMList:

etcd:
  # The name of the node pool(s) to use for etcd nodes
  nodePools:
    - etcd

iscsi:
  # (optional) Start iscsid.service to enable iscsi volume attachment.
  # enabled: false

# The path to the tectonic licence file.
# You can download the Tectonic license file from your Account overview page at [1].
#
# [1] https://account.coreos.com/overview
licensePath: {{.LicensePath}}

master:
  # The name of the node pool(s) to use for master nodes
  nodePools:
    - master

# The name of the cluster.
# If used in a cloud-environment, this will be prepended to ` + "`baseDomain`" + ` resulting in the URL to the Tectonic console.
#
# Note: This field MUST be set manually prior to creating the cluster.
# Warning: Special characters in the name like '.' may cause errors on OpenStack platforms due to resource name constraints.
name: {{.Name}}

networking:
  # (optional) This declares the MTU used by Calico.
  # mtu:

  # This declares the IP range to assign Kubernetes pod IPs in CIDR notation.
  podCIDR: 10.2.0.0/16

  # This declares the IP range to assign Kubernetes service cluster IPs in CIDR notation.
  # The maximum size of this IP range is /12
  serviceCIDR: 10.3.0.0/16

  # (optional) Configures the network to be used in Tectonic. One of the following values can be used:
  #
  # - "flannel": enables overlay networking only. This is implemented by flannel using VXLAN.
  #
  # - "canal": enables overlay networking including network policy. Overlay is implemented by flannel using VXLAN. Network policy is implemented by Calico.
  #
  # - "calico-ipip": [ALPHA] enables BGP based networking. Routing and network policy is implemented by Calico. Note this has been tested on baremetal installations only.
  #
  # - "none": disables the installation of any Pod level networking layer provided by Tectonic. By setting this value, users are expected to deploy their own solution to enable network connectivity for Pods and Services.
  # type: canal

nodePools:
    # The number of etcd nodes to be created.
    # If set to zero, the count of etcd nodes will be determined automatically.
  - count: 3
    name: etcd

    # The number of master nodes to be created.
    # This applies only to cloud platforms.
  - count: 1
    name: master

    # The number of worker nodes to be created.
    # This applies only to cloud platforms.
  - count: 3
    name: worker

# The platform used for deploying.
platform: aws

# The path the pull secret file in JSON format.
# This is known to be a "Docker pull secret" as produced by the docker login [1] command.
# A sample JSON content is shown in [2].
# You can download the pull secret from your Account overview page at [3].
#
# [1] https://docs.docker.com/engine/reference/commandline/login/
#
# [2] https://coreos.com/os/docs/latest/registry-authentication.html#manual-registry-auth-setup
#
# [3] https://account.coreos.com/overview
pullSecretPath: {{.PullSecretPath}}

worker:
  # The name of the node pool(s) to use for workers
  nodePools:
    - worker
`
