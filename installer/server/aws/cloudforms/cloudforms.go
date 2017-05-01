package cloudforms

import (
	"net"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/session"
)

// Config: External configuration interface
type Config struct {
	// Name of cloudformation stack
	ClusterName string `json:"clusterName"`

	// Region to deploy cluster in
	Region string `json:"region"`

	// CoreOS Channel - alpha/beta/stable
	Channel string `json:"channel"`

	// DNS name for Kubernetes Controller Load Balancer
	// Must be contained within hosted zone
	ControllerDomain string `json:"controllerDomain"`

	// DNS name for Tectonic Load Balancer
	// Must be contained within hosted zone
	TectonicDomain string `json:"tectonicDomain"`

	// ELBs and controllers should be 'internet-facing' or 'internal'
	ELBScheme string `json:"elbScheme"`

	// Hosted zone ID to add DNS records to
	HostedZoneID string `json:"hostedZoneID"`

	// CIDR for new VPC
	VPCCIDR string `json:"vpcCIDR"`

	// Existing VPC ID (leave blank to create new VPC)
	VPCID string `json:"vpcID,omitempty"`

	// OPTIONAL: Existing VPC route table to attach subnets to.
	// (Leave blank to use main route table in existing VPC)
	RouteTableID string `json:"routeTableID,omitempty"`

	// List of subnets in VPC (new or existing) to spread controllers across.
	ControllerSubnets []VPCSubnet `json:"controllerSubnets"`

	// List of subnets in VPC (new or existing) to spread workers across.
	WorkerSubnets []VPCSubnet `json:"workerSubnets"`

	// ARN of KMS key used to encrypt secrets
	KMSKeyARN string `json:"kmsKeyARN"`

	// EC2 ssh key for instances (controller and worker)
	KeyName string `json:"keyName"`

	// EC2 etcd instance settings
	ETCDCount          int    `json:"etcdCount"`
	ETCDInstanceType   string `json:"etcdInstanceType"`
	ETCDRootVolumeType string `json:"etcdRootVolumeType"`
	ETCDRootVolumeIOPS int    `json:"etcdRootVolumeIOPS"`
	ETCDRootVolumeSize int    `json:"etcdRootVolumeSize"`

	// External etcd client endpoint, e.g. etcd.example.com:2379
	ExternalETCDClient string `json:"externalETCDClient"`

	// EC2 controller instances
	ControllerCount          int    `json:"controllerCount"`
	ControllerInstanceType   string `json:"controllerInstanceType"`
	ControllerRootVolumeType string `json:"controllerRootVolumeType"`
	ControllerRootVolumeIOPS int    `json:"controllerRootVolumeIOPS"`
	ControllerRootVolumeSize int    `json:"controllerRootVolumeSize"`

	// EC2 worker instances
	WorkerCount          int    `json:"workerCount"`
	WorkerInstanceType   string `json:"workerInstanceType"`
	WorkerRootVolumeType string `json:"workerRootVolumeType"`
	WorkerRootVolumeIOPS int    `json:"workerRootVolumeIOPS"`
	WorkerRootVolumeSize int    `json:"workerRootVolumeSize"`

	PodCIDR     string `json:"podCIDR"`
	ServiceCIDR string `json:"serviceCIDR"`

	// Cloudformation tags
	Tags []Tag `json:"tags"`

	// Userdata templates
	ControllerTemplate *template.Template `json:"-"`
	WorkerTemplate     *template.Template `json:"-"`
	EtcdTemplate       *template.Template `json:"-"`

	// Cloudformation stack template
	StackTemplate *template.Template `json:"-"`

	// Computed IPs for self-hosted Kubernetes
	APIServiceIP net.IP
	DNSServiceIP net.IP

	// computed fields (set during initialize(), and/or overwritable after)
	ETCDInstances           []ETCDInstance
	ETCDEndpoints           string
	ETCDInitialCluster      string
	APIServers              string
	SecureAPIServers        string
	AMI                     string
	CreateControllerSubnets bool
	CreateWorkerSubnets     bool
	HostedZoneName          string

	// Encoded assets
	EncodedSecrets *compactSecretAssets

	// Logical names of dynamic resources
	VPCLogicalName string

	// Reference strings for dynamic resources
	VPCRef string

	// Logical name for the VPC internet gateway
	InternetGatewayLogicalName string

	// Reference to an existing VPC internet gateway
	InternetGatewayRef string

	// Asset S3 location information
	AssetsS3File   string
	AssetsS3Bucket string
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VPCSubnet struct {
	// Identifier of the subnet if already existing
	ID string `json:"id"`
	// Logical name for this subnet
	// ignored if existing
	Name string `json:"name"`
	// Availability zone for this subnet
	// Max one subnet per availability zone
	AvailabilityZone string `json:"availabilityZone"`
	// CIDR for this subnet
	// must be disjoint from other subnets
	// must be contained by VPC CIDR
	InstanceCIDR string `json:"instanceCIDR"`
}

type ETCDInstance struct {
	// Nam of the ETCD instance
	Name string
	// DNS name addressing the EC2 Instance
	DomainName string
	// Subnet for this EC2 Instance
	Subnet VPCSubnet
}

func NewCloudFormation(config *Config, sess *session.Session, secrets *SecretAssets) (*Cluster, error) {
	// ensure the config is valid
	if err := config.Valid(); err != nil {
		return nil, err
	}

	// create a stackVars for rendering a cloud-formation 'stack'
	stackVars := newStackVars(config)

	// encode and compress secret assets into stackVars
	if err := stackVars.encodeSecretAssets(sess, secrets); err != nil {
		return nil, err
	}

	if err := stackVars.validateUserData(); err != nil {
		return nil, err
	}

	// render the controller, worker user-data and stack template
	if err := stackVars.render(); err != nil {
		return nil, err
	}

	// upload stack to S3
	if err := stackVars.upload(sess); err != nil {
		return nil, err
	}

	// make AWS validate the stack
	if err := stackVars.validateAll(sess); err != nil {
		stackVars.remove(sess)
		return nil, err
	}

	bundle := &Cluster{
		ClusterName:      config.ClusterName,
		ControllerDomain: config.ControllerDomain,
		Region:           config.Region,
		StackBody:        stackVars.StackBody,
		StackURL:         stackVars.StackURL,
	}

	return bundle, nil
}
