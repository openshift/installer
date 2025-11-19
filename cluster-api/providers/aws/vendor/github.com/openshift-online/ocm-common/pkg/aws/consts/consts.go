package consts

import (
	"os"
)

const (
	MaxAwsRoleLength = 64
)

var HomeDir, _ = os.UserHomeDir()
var QEFLAG string = os.Getenv("QE_FLAG")

const (
	HTTPConflict                   = 409
	AWSCredentialsFileRelativePath = "/.aws/credentials"
	AdminPolicyArn                 = "arn:aws:iam::aws:policy/AdministratorAccess"
	DefaultAWSCredentialUser       = "default"
	DefaultAWSRegion               = "us-east-2"
	DefaultAWSZone                 = "a"
)

var AWSCredentialsFilePath = HomeDir + AWSCredentialsFileRelativePath

const (
	DefaultVPCCIDR            = "10.0.0.0/16"
	DefaultCIDRPrefix         = 24
	RouteDestinationCidrBlock = "0.0.0.0/0"

	VpcDefaultName = "ocm-ci-vpc"

	CreationPrivateSelector = "private"
	CreationPublicSelector  = "public"
	CreationPairSelector    = "pair"
	CreationMultiSelector   = "multi"

	VpcDnsHostnamesAttribute = "DnsHostnames"
	VpcDnsSupportAttribute   = "DnsSupport"
	NetworkResourceFileName  = "resource.json"

	TCPProtocol = "tcp"
	UDPProtocol = "udp"

	BastionSecurityGroupName                  = "bastion-sg"
	ProxySecurityGroupName                    = "proxy-sg"
	AdditionalSecurityGroupName               = "ocm-additional-sg"
	ProxySecurityGroupDescription             = "security group for proxy"
	BastionSecurityGroupDescription           = "security group for bastion"
	DefaultAdditionalSecurityGroupDescription = "This security group is created for OCM testing"

	QEFlagKey = "ocm_ci_flag"

	// Proxy related
	ProxyName             = "ocm-proxy"
	InstanceKeyNamePrefix = "ocm-ci"
	AWSInstanceUser       = "ec2-user"
	BastionName           = "ocm-bastion"

	SSHPort = "22"

	// Squid related
	SquidConfigFilePath   = "/etc/squid/squid.conf"
	SquidPasswordFilePath = "/etc/squid/passwords"
	SquidProxyPort        = "3128"
)

var AmazonName = "amazon"

const (
	PublicSubNetTagKey   = "PublicSubnet"
	PublicSubNetTagValue = "true"
)

const (
	PrivateLBTag = "kubernetes.io/role/internal-elb"
	PublicLBTag  = "kubernetes.io/role/elb"
	// valid LB Tag is empty or 1
	LBTagValue = ""
)
