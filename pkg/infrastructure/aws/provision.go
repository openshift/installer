package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	awssession "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/tfvars"
	tfvarsaws "github.com/openshift/installer/pkg/tfvars/aws"
	typesaws "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/version"
)

func normalAWSProvision(a AWSInfraProvider, tfvarsFiles, fileList []*asset.File) (*asset.File, *asset.File, error) {
	// Unmarshall input from tf variables, so we can use it along with installConfig and other assets
	// as the contractual input regardless off the implementation.
	clusterConfig := &tfvars.Config{}
	clusterAWSConfig := &tfvarsaws.Config{}
	for _, file := range tfvarsFiles {
		if file.Filename == "terraform.tfvars.json" {
			if err := json.Unmarshal(file.Data, clusterConfig); err != nil {
				return nil, nil, err
			}
		}

		if file.Filename == "terraform.platform.auto.tfvars.json" {
			if err := json.Unmarshal(file.Data, clusterAWSConfig); err != nil {
				return nil, nil, err
			}
		}
	}

	eps := []typesaws.ServiceEndpoint{}
	for k, v := range clusterAWSConfig.CustomEndpoints {
		eps = append(eps, typesaws.ServiceEndpoint{Name: k, URL: v})
	}

	awsSession, err := awssession.GetSessionWithOptions(
		awssession.WithRegion(clusterAWSConfig.Region),
		awssession.WithServiceEndpoints(clusterAWSConfig.Region, eps),
	)
	if err != nil {
		return nil, nil, err
	}
	awsSession.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "openshiftInstaller.OpenshiftInstallerUserAgentHandler",
		Fn:   request.MakeAddToUserAgentHandler("OpenShift/4.x Creator", version.Raw),
	})

	// Create VPC resources.
	vpcInput := &CreateInfraOptions{
		Region:         clusterAWSConfig.Region,
		InfraID:        clusterConfig.ClusterID,
		Zones:          clusterAWSConfig.MasterAvailabilityZones,
		BaseDomain:     clusterConfig.BaseDomain,
		AdditionalTags: clusterAWSConfig.ExtraTags,
		public:         clusterAWSConfig.PublishStrategy == "External",
		cidrV4Blocks:   clusterConfig.MachineV4CIDRs,
		cidrV6Blocks:   clusterConfig.MachineV6CIDRs,
	}
	logger := logrus.StandardLogger()
	if err := createVPCResources(logger, awsSession, vpcInput); err != nil {
		return nil, nil, err
	}

	// Create DNS resources.
	dnsInput := &dnsInput{
		clusterID:                   clusterConfig.ClusterID,
		region:                      clusterAWSConfig.Region,
		baseDomain:                  clusterConfig.BaseDomain,
		clusterDomain:               clusterConfig.ClusterDomain,
		vpcID:                       vpcInput.vpcID,
		additionalTags:              clusterAWSConfig.ExtraTags,
		loadBalancerExternalZoneID:  vpcInput.LoadBalancers.External.ZoneID,
		loadBalancerExternalZoneDNS: vpcInput.LoadBalancers.External.DNSName,
		loadBalancerInternalZoneID:  vpcInput.LoadBalancers.Internal.ZoneID,
		loadBalancerInternalZoneDNS: vpcInput.LoadBalancers.Internal.DNSName,
	}
	if err := createDNSResources(context.TODO(), logger, awsSession, dnsInput); err != nil {
		return nil, nil, err
	}

	// Create Bootstrap resources.
	bootstrapInput := &bootstrapInput{
		clusterID:                clusterConfig.ClusterID,
		targetGroupARNs:          vpcInput.targetGroupARNs,
		ignitionBucket:           clusterAWSConfig.IgnitionBucket,
		ignitionContent:          clusterConfig.IgnitionBootstrap,
		userData:                 clusterAWSConfig.BootstrapIgnitionStub,
		amiID:                    clusterAWSConfig.AMI,
		instanceType:             clusterAWSConfig.BootstrapInstanceType,
		subnetID:                 vpcInput.publicSubnetIDs[0],
		securityGroupIDs:         []string{vpcInput.bootstrapSecurityGroupID, vpcInput.masterSecurityGroupID},
		associatePublicIPAddress: clusterAWSConfig.PublishStrategy == "External",
		additionalTags:           clusterAWSConfig.ExtraTags,
		volumeType:               "gp2",
		volumeSize:               30,
		volumeIOPS:               0,
		encrypted:                clusterAWSConfig.Encrypted,
		kmsKeyID:                 clusterAWSConfig.KMSKeyID,
	}
	if err := createBootstrapResources(logger, awsSession, bootstrapInput); err != nil {
		return nil, nil, err
	}

	// Create Control Plane resources.
	controlPlaneInput := &controlPlaneInput{
		clusterID:       clusterConfig.ClusterID,
		targetGroupARNs: vpcInput.targetGroupARNs,
		userData:        clusterConfig.IgnitionMaster,
		amiID:           clusterAWSConfig.AMI,
		instanceType:    clusterAWSConfig.MasterInstanceType,
		subnetIDs:       vpcInput.privateSubnetIDs,
		securityGroupID: vpcInput.masterSecurityGroupID,
		volumeType:      clusterAWSConfig.Type,
		volumeSize:      clusterAWSConfig.Size,
		volumeIOPS:      clusterAWSConfig.IOPS,
		additionalTags:  clusterAWSConfig.ExtraTags,
		encrypted:       clusterAWSConfig.Encrypted,
		kmsKeyID:        clusterAWSConfig.KMSKeyID,
		replicas:        clusterConfig.Masters,
	}
	if err := createControlPlaneResources(logger, awsSession, controlPlaneInput); err != nil {
		return nil, nil, err
	}

	// Create IAM resources.
	iamInput := &iamInput{
		clusterID: clusterConfig.ClusterID,
	}
	if err := createIAMResources(logger, awsSession, iamInput); err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}
