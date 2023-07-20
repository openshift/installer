package explain

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PrintFields(t *testing.T) {
	schema, err := loadSchema(loadCRD(t))
	assert.NoError(t, err)

	cases := []struct {
		path []string

		desc string
		err  string
	}{{
		desc: `FIELDS:
    additionalTrustBundle <string>
      AdditionalTrustBundle is a PEM-encoded X.509 certificate bundle that will be added to the nodes' trusted certificate store.

    additionalTrustBundlePolicy <string>
      Valid Values: "","Proxyonly","Always"
      AdditionalTrustBundlePolicy determines when to add the AdditionalTrustBundle to the nodes' trusted certificate store. "Proxyonly" is the default. The field can be set to following specified values. "Proxyonly" : adds the AdditionalTrustBundle to nodes when http/https proxy is configured. "Always" : always adds AdditionalTrustBundle.

    apiVersion <string>
      APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources

    baseDomain <string> -required-
      BaseDomain is the base domain to which the cluster should belong.

    bootstrapInPlace <object>
      BootstrapInPlace is the configuration for installing a single node with bootstrap in place installation.

    capabilities <object>
      Capabilities configures the installation of optional core cluster components.

    compute <[]object>
      Compute is the configuration for the machines that comprise the compute nodes.
      MachinePool is a pool of machines to be installed.

    controlPlane <object>
      ControlPlane is the configuration for the machines that comprise the control plane.

    cpuPartitioningMode <string>
      Default: "None"
      Valid Values: "None","AllNodes"
      CPUPartitioning determines if a cluster should be setup for CPU workload partitioning at install time. When this field is set the cluster will be flagged for CPU Partitioning allowing users to segregate workloads to specific CPU Sets. This does not make any decisions on workloads it only configures the nodes to allow CPU Partitioning. The "AllNodes" value will setup all nodes for CPU Partitioning, the default is "None". This feature is currently in TechPreview.

    credentialsMode <string>
      Valid Values: "","Mint","Passthrough","Manual"
      CredentialsMode is used to explicitly set the mode with which CredentialRequests are satisfied. 
 If this field is set, then the installer will not attempt to query the cloud permissions before attempting installation. If the field is not set or empty, then the installer will perform its normal verification that the credentials provided are sufficient to perform an installation. 
 There are three possible values for this field, but the valid values are dependent upon the platform being used. "Mint": create new credentials with a subset of the overall permissions for each CredentialsRequest "Passthrough": copy the credentials with all of the overall permissions for each CredentialsRequest "Manual": CredentialsRequests must be handled manually by the user 
 For each of the following platforms, the field can set to the specified values. For all other platforms, the field must not be set. AWS: "Mint", "Passthrough", "Manual" Azure: "Passthrough", "Manual" AzureStack: "Manual" GCP: "Mint", "Passthrough", "Manual" IBMCloud: "Manual" AlibabaCloud: "Manual" PowerVS: "Manual" Nutanix: "Manual"

    featureGates <[]string>
      FeatureGates enables a set of custom feature gates. May only be used in conjunction with FeatureSet "CustomNoUpgrade". Features may be enabled or disabled by providing a true or false value for the feature gate. E.g. "featureGates": ["FeatureGate1=true", "FeatureGate2=false"].

    featureSet <string>
      FeatureSet enables features that are not part of the default feature set. Valid values are "Default", "TechPreviewNoUpgrade" and "CustomNoUpgrade". When omitted, the "Default" feature set is used.

    fips <boolean>
      Default: false
      FIPS configures https://www.nist.gov/itl/fips-general-information

    imageContentSources <[]object>
      ImageContentSources lists sources/repositories for the release-image content. The field is deprecated. Please use imageDigestSources.
      ImageContentSource defines a list of sources/repositories that can be used to pull content. The field is deprecated. Please use imageDigestSources.

    imageDigestSources <[]object>
      ImageDigestSources lists sources/repositories for the release-image content.
      ImageDigestSource defines a list of sources/repositories that can be used to pull content.

    kind <string>
      Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds

    metadata <object> -required-
      <empty>

    networking <object>
      Networking is the configuration for the pod network provider in the cluster.

    platform <object> -required-
      Platform is the configuration for the specific platform upon which to perform the installation.

    proxy <object>
      Proxy defines the proxy settings for the cluster. If unset, the cluster will not be configured to use a proxy.

    publish <string>
      Default: "External"
      Valid Values: "","External","Internal"
      Publish controls how the user facing endpoints of the cluster like the Kubernetes API, OpenShift routes etc. are exposed. When no strategy is specified, the strategy is "External".

    pullSecret <string> -required-
      PullSecret is the secret to use when pulling images.

    sshKey <string>
      SSHKey is the public Secure Shell (SSH) key to provide access to instances.`,
	}, {
		path: []string{"publish"},
		desc: ``,
	}, {
		path: []string{"platform"},
		desc: `FIELDS:
    alibabacloud <object>
      AlibabaCloud is the configuration used when installing on Alibaba Cloud.

    aws <object>
      AWS is the configuration used when installing on AWS.

    azure <object>
      Azure is the configuration used when installing on Azure.

    baremetal <object>
      BareMetal is the configuration used when installing on bare metal.

    external <object>
      External is the configuration used when installing on an external cloud provider.

    gcp <object>
      GCP is the configuration used when installing on Google Cloud Platform.

    ibmcloud <object>
      IBMCloud is the configuration used when installing on IBM Cloud.

    libvirt <object>
      Libvirt is the configuration used when installing on libvirt.

    none <object>
      None is the empty configuration used when installing on an unsupported platform.

    nutanix <object>
      Nutanix is the configuration used when installing on Nutanix.

    openstack <object>
      OpenStack is the configuration used when installing on OpenStack.

    ovirt <object>
      Ovirt is the configuration used when installing on oVirt.

    powervs <object>
      PowerVS is the configuration used when installing on Power VS.

    vsphere <object>
      VSphere is the configuration used when installing on vSphere.`,
	}, {
		path: []string{"platform", "aws"},
		desc: `FIELDS:
    amiID <string>
      AMIID is the AMI that should be used to boot machines for the cluster. If set, the AMI should belong to the same region as the cluster.

    defaultMachinePlatform <object>
      DefaultMachinePlatform is the default configuration used when installing on AWS for machine pools which do not define their own platform configuration.

    experimentalPropagateUserTags <boolean>
      The field is deprecated. ExperimentalPropagateUserTags is an experimental flag that directs in-cluster operators to include the specified user tags in the tags of the AWS resources that the operators create.

    hostedZone <string>
      HostedZone is the ID of an existing hosted zone into which to add DNS records for the cluster's internal API. An existing hosted zone can only be used when also using existing subnets. The hosted zone must be associated with the VPC containing the subnets. Leave the hosted zone unset to have the installer create the hosted zone on your behalf.

    hostedZoneRole <string>
      HostedZoneRole is the ARN of an IAM role to be assumed when performing operations on the provided HostedZone. HostedZoneRole can be used in a shared VPC scenario when the private hosted zone belongs to a different account than the rest of the cluster resources. If HostedZoneRole is set, HostedZone must also be set.

    lbType <string>
      LBType is an optional field to specify a load balancer type. 
 When this field is specified, the default ingresscontroller will be created using the specified load-balancer type. 
 Following are the accepted values: 
 * "Classic": A Classic Load Balancer that makes routing decisions at either the transport layer (TCP/SSL) or the application layer (HTTP/HTTPS). See the following for additional details: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/load-balancer-types.html#clb 
 * "NLB": A Network Load Balancer that makes routing decisions at the transport layer (TCP/SSL). See the following for additional details: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/load-balancer-types.html#nlb 
 If this field is not set explicitly, it defaults to "Classic".  This default is subject to change over time.

    preserveBootstrapIgnition <boolean>
      PreserveBootstrapIgnition is an optional field that can be used to make the S3 deletion optional during bootstrap destroy.

    propagateUserTags <boolean>
      PropagateUserTags is a flag that directs in-cluster operators to include the specified user tags in the tags of the AWS resources that the operators create.

    region <string> -required-
      Region specifies the AWS region where the cluster will be created.

    serviceEndpoints <[]object>
      ServiceEndpoints list contains custom endpoints which will override default service endpoint of AWS Services. There must be only one ServiceEndpoint for a service.
      ServiceEndpoint store the configuration for services to override existing defaults of AWS Services.

    subnets <[]string>
      Subnets specifies existing subnets (by ID) where cluster resources will be created.  Leave unset to have the installer create subnets in a new VPC on your behalf.

    userTags <object>
      UserTags additional keys and values that the installer will add as tags to all resources that it creates. Resources created by the cluster itself may not include these tags.`,
	}, {
		path: []string{"platform", "azure"},
		desc: `FIELDS:
    armEndpoint <string>
      ARMEndpoint is the endpoint for the Azure API when installing on Azure Stack.

    baseDomainResourceGroupName <string>
      BaseDomainResourceGroupName specifies the resource group where the Azure DNS zone for the base domain is found. This field is optional when creating a private cluster, otherwise required.

    cloudName <string>
      Valid Values: "","AzurePublicCloud","AzureUSGovernmentCloud","AzureChinaCloud","AzureGermanCloud","AzureStackCloud"
      cloudName is the name of the Azure cloud environment which can be used to configure the Azure SDK with the appropriate Azure API endpoints. If empty, the value is equal to "AzurePublicCloud".

    clusterOSImage <string>
      ClusterOSImage is the url of a storage blob in the Azure Stack environment containing an RHCOS VHD. This field is required for Azure Stack and not applicable to Azure.

    computeSubnet <string>
      ComputeSubnet specifies an existing subnet for use by compute nodes

    controlPlaneSubnet <string>
      ControlPlaneSubnet specifies an existing subnet for use by the control plane nodes

    defaultMachinePlatform <object>
      DefaultMachinePlatform is the default configuration used when installing on Azure for machine pools which do not define their own platform configuration.

    networkResourceGroupName <string>
      NetworkResourceGroupName specifies the network resource group that contains an existing VNet

    outboundType <string>
      Default: "Loadbalancer"
      Valid Values: "","Loadbalancer","NatGateway","UserDefinedRouting"
      OutboundType is a strategy for how egress from cluster is achieved. When not specified default is "Loadbalancer".

    region <string> -required-
      Region specifies the Azure region where the cluster will be created.

    resourceGroupName <string>
      ResourceGroupName is the name of an already existing resource group where the cluster should be installed. This resource group should only be used for this specific cluster and the cluster components will assume ownership of all resources in the resource group. Destroying the cluster using installer will delete this resource group. This resource group must be empty with no other resources when trying to use it for creating a cluster. If empty, a new resource group will created for the cluster.

    userTags <object>
      UserTags has additional keys and values that the installer will add as tags to all resources that it creates on AzurePublicCloud alone. Resources created by the cluster itself may not include these tags.

    virtualNetwork <string>
      VirtualNetwork specifies the name of an existing VNet for the installer to use`,
	}, {
		path: []string{"platform", "aws", "region"},
		desc: ``,
	}, {
		path: []string{"platform", "aws", "subnets"},
		desc: ``,
	}, {
		path: []string{"platform", "aws", "userTags"},
		desc: ``,
	}, {
		path: []string{"platform", "aws", "serviceEndpoints"},
		desc: `FIELDS:
    name <string> -required-
      Name is the name of the AWS service. This must be provided and cannot be empty.

    url <string> -required-
      URL is fully qualified URI with scheme https, that overrides the default generated endpoint for a client. This must be provided and cannot be empty.`,
	}, {
		path: []string{"platform", "aws", "serviceEndpoints", "url"},
		desc: ``,
	}}
	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			got, err := lookup(schema, test.path)
			assert.NoError(t, err)
			buf := &bytes.Buffer{}
			(printer{Writer: buf}).PrintFields(got)
			assert.Equal(t, test.desc, strings.TrimSpace(buf.String()))
		})
	}
}

func Test_PrintNonFields(t *testing.T) {
	schema, err := loadSchema(loadCRD(t))
	assert.NoError(t, err)

	cases := []struct {
		path []string

		desc string
		err  string
	}{{
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <object>
  InstallConfig is the configuration for an OpenShift install.
		`,
	}, {
		path: []string{"publish"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <string>
  Default: "External"
  Valid Values: "","External","Internal"
  Publish controls how the user facing endpoints of the cluster like the Kubernetes API, OpenShift routes etc. are exposed. When no strategy is specified, the strategy is "External".
		`,
	}, {
		path: []string{"platform"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <object>
  Platform is the configuration for the specific platform upon which to perform the installation.
		`,
	}, {
		path: []string{"platform", "aws"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <object>
  AWS is the configuration used when installing on AWS.
		`,
	}, {
		path: []string{"platform", "azure"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <object>
  Azure is the configuration used when installing on Azure.
		`,
	}, {
		path: []string{"platform", "aws", "region"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <string>
  Region specifies the AWS region where the cluster will be created.
		`,
	}, {
		path: []string{"platform", "aws", "subnets"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <[]string>
  Subnets specifies existing subnets (by ID) where cluster resources will be created.  Leave unset to have the installer create subnets in a new VPC on your behalf.
		`,
	}, {
		path: []string{"platform", "aws", "userTags"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <object>
  UserTags additional keys and values that the installer will add as tags to all resources that it creates. Resources created by the cluster itself may not include these tags.
		`,
	}, {
		path: []string{"platform", "aws", "serviceEndpoints"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <[]object>
  ServiceEndpoints list contains custom endpoints which will override default service endpoint of AWS Services. There must be only one ServiceEndpoint for a service.
		`,
	}, {
		path: []string{"platform", "aws", "serviceEndpoints", "url"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <string>
  URL is fully qualified URI with scheme https, that overrides the default generated endpoint for a client. This must be provided and cannot be empty.
		`,
	}, {
		path: []string{"compute", "platform", "aws", "iamRole"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <string>
  IAMRole is the name of the IAM Role to use for the instance profile of the machine. Leave unset to have the installer create the IAM Role on your behalf.
	`,
	}, {
		path: []string{"controlPlane", "platform", "aws", "iamRole"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <string>
  IAMRole is the name of the IAM Role to use for the instance profile of the machine. Leave unset to have the installer create the IAM Role on your behalf.
	`,
	}, {
		path: []string{"platform", "aws", "defaultMachinePlatform", "iamRole"},
		desc: `
KIND:     InstallConfig
VERSION:  v1

RESOURCE: <string>
  IAMRole is the name of the IAM Role to use for the instance profile of the machine. Leave unset to have the installer create the IAM Role on your behalf.
	`,
	}}
	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			got, err := lookup(schema, test.path)
			assert.NoError(t, err)
			buf := &bytes.Buffer{}
			p := printer{Writer: buf}
			p.PrintKindAndVersion()
			p.PrintResource(got)
			assert.Equal(t, strings.TrimSpace(test.desc), strings.TrimSpace(buf.String()))
		})
	}
}
