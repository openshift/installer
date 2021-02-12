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

    apiVersion <string>
      APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources

    baseDomain <string> -required-
      BaseDomain is the base domain to which the cluster should belong.

    bootstrapInPlace <object>
      BootstrapInPlace is the configuration for installing a single node with bootstrap in place installation.

    compute <[]object>
      Compute is the configuration for the machines that comprise the compute nodes.
      MachinePool is a pool of machines to be installed.

    controlPlane <object>
      ControlPlane is the configuration for the machines that comprise the control plane.

    credentialsMode <string>
      Valid Values: "","Mint","Passthrough","Manual"
      CredentialsMode is used to explicitly set the mode with which CredentialRequests are satisfied. 
 If this field is set, then the installer will not attempt to query the cloud permissions before attempting installation. If the field is not set or empty, then the installer will perform its normal verification that the credentials provided are sufficient to perform an installation. 
 There are three possible values for this field, but the valid values are dependent upon the platform being used. "Mint": create new credentials with a subset of the overall permissions for each CredentialsRequest "Passthrough": copy the credentials with all of the overall permissions for each CredentialsRequest "Manual": CredentialsRequests must be handled manually by the user 
 For each of the following platforms, the field can set to the specified values. For all other platforms, the field must not be set. AWS: "Mint", "Passthrough", "Manual" Azure: "Mint", "Passthrough", "Manual" GCP: "Mint", "Passthrough", "Manual"

    fips <boolean>
      Default: false
      FIPS configures https://www.nist.gov/itl/fips-general-information

    imageContentSources <[]object>
      ImageContentSources lists sources/repositories for the release-image content.
      ImageContentSource defines a list of sources/repositories that can be used to pull content.

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
    aws <object>
      AWS is the configuration used when installing on AWS.

    azure <object>
      Azure is the configuration used when installing on Azure.

    baremetal <object>
      BareMetal is the configuration used when installing on bare metal.

    gcp <object>
      GCP is the configuration used when installing on Google Cloud Platform.

    kubevirt <object>
      Kubevirt is the configuration used when installing on kubevirt.

    libvirt <object>
      Libvirt is the configuration used when installing on libvirt.

    none <object>
      None is the empty configuration used when installing on an unsupported platform.

    openstack <object>
      OpenStack is the configuration used when installing on OpenStack.

    ovirt <object>
      Ovirt is the configuration used when installing on oVirt.

    vsphere <object>
      VSphere is the configuration used when installing on vSphere.`,
	}, {
		path: []string{"platform", "aws"},
		desc: `FIELDS:
    amiID <string>
      AMIID is the AMI that should be used to boot machines for the cluster. If set, the AMI should belong to the same region as the cluster.

    defaultMachinePlatform <object>
      DefaultMachinePlatform is the default configuration used when installing on AWS for machine pools which do not define their own platform configuration.

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
    baseDomainResourceGroupName <string>
      BaseDomainResourceGroupName specifies the resource group where the Azure DNS zone for the base domain is found.

    cloudName <string>
      Valid Values: "","AzurePublicCloud","AzureUSGovernmentCloud","AzureChinaCloud","AzureGermanCloud"
      cloudName is the name of the Azure cloud environment which can be used to configure the Azure SDK with the appropriate Azure API endpoints. If empty, the value is equal to "AzurePublicCloud".

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
      Valid Values: "","Loadbalancer","UserDefinedRouting"
      OutboundType is a strategy for how egress from cluster is achieved. When not specified default is "Loadbalancer".

    region <string> -required-
      Region specifies the Azure region where the cluster will be created.

    resourceGroupName <string>
      ResourceGroupName is the name of an already existing resource group where the cluster should be installed. This resource group should only be used for this specific cluster and the cluster components will assume assume ownership of all resources in the resource group. Destroying the cluster using installer will delete this resource group. This resource group must be empty with no other resources when trying to use it for creating a cluster. If empty, a new resource group will created for the cluster.

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
