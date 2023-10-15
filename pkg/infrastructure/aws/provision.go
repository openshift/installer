package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	awssession "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/tfvars"
	tfvarsaws "github.com/openshift/installer/pkg/tfvars/aws"
	typesaws "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/version"
	"k8s.io/apimachinery/pkg/util/sets"
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

	availabilityZones := sets.New(clusterAWSConfig.MasterAvailabilityZones...)
	availabilityZones.Insert(clusterAWSConfig.WorkerAvailabilityZones...)

	// Create VPC resources.
	vpcInput := &CreateInfraOptions{
		Region:         clusterAWSConfig.Region,
		InfraID:        clusterConfig.ClusterID,
		Zones:          sets.List(availabilityZones),
		BaseDomain:     clusterConfig.BaseDomain,
		AdditionalTags: clusterAWSConfig.ExtraTags,
		public:         clusterAWSConfig.PublishStrategy == "External",
		cidrV4Blocks:   clusterConfig.MachineV4CIDRs,
		cidrV6Blocks:   clusterConfig.MachineV6CIDRs,
		vpcID:          clusterAWSConfig.VPC,
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
		clusterID:         clusterConfig.ClusterID,
		targetGroupARNs:   vpcInput.targetGroupARNs,
		userData:          clusterConfig.IgnitionMaster,
		amiID:             clusterAWSConfig.AMI,
		instanceType:      clusterAWSConfig.MasterInstanceType,
		subnetIDs:         vpcInput.privateSubnetIDs,
		securityGroupIDs:  append(clusterAWSConfig.MasterSecurityGroups, vpcInput.masterSecurityGroupID),
		volumeType:        clusterAWSConfig.Type,
		volumeSize:        clusterAWSConfig.Size,
		volumeIOPS:        clusterAWSConfig.IOPS,
		additionalTags:    clusterAWSConfig.ExtraTags,
		encrypted:         clusterAWSConfig.Encrypted,
		kmsKeyID:          clusterAWSConfig.KMSKeyID,
		replicas:          clusterConfig.Masters,
		availabilityZones: clusterAWSConfig.MasterAvailabilityZones,
		zoneToSubnetIDMap: vpcInput.zoneToSubnetIDMap,
	}
	if err := createControlPlaneResources(logger, awsSession, controlPlaneInput); err != nil {
		return nil, nil, err
	}

	// Create IAM resources.
	if err := createIAMResources(logger, awsSession, clusterConfig.ClusterID, clusterAWSConfig.ExtraTags); err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}

func createInstanceProfile(logger *logrus.Logger, session *session.Session, namePrefix, assumeRolePolicy, policyDocument string, tags map[string]string) (string, error) {
	iamClient := iam.New(session)
	iamTags := iamCreateTags(tags)

	roleName := fmt.Sprintf("%s-role", namePrefix)
	role, err := iamGetRole(iamClient, roleName)
	if err != nil {
		return "", err
	}
	if role == nil {
		if _, err := iamCreateRole(iamClient, roleName, assumeRolePolicy, iamTags); err != nil {
			return "", err
		}
		logger.WithField("name", roleName).Infoln("Created role")
	} else {
		logger.WithField("name", roleName).Infoln("Found existing role")
	}

	profileName := fmt.Sprintf("%s-profile", namePrefix)
	instanceProfile, err := iamGetInstanceProfile(iamClient, profileName)
	if err != nil {
		return "", err
	}
	if instanceProfile == nil {
		instanceProfile, err = iamCreateInstanceProfile(iamClient, profileName, iamTags)
		if err != nil {
			return "", err
		}
		logger.WithField("name", profileName).Infoln("Created instance profile")
		//logger.WithField("name", profileName).Infoln("Instance profile was created and exists")
	} else {
		logger.WithField("name", profileName).Infoln("Found existing instance profile")
	}
	if err := iamAddRoleToProfile(iamClient, instanceProfile, roleName); err != nil {
		return "", err
	}
	logger.WithField("role", roleName).WithField("profile", profileName).Infoln("Added role to instance profile")

	rolePolicyName := fmt.Sprintf("%s-policy", profileName)
	policyName, err := iamGetRolePolicy(iamClient, roleName, rolePolicyName)
	if err != nil {
		return "", err
	}
	if policyName != rolePolicyName {
		if err := iamAddPolicyToRole(iamClient, roleName, rolePolicyName, policyDocument); err != nil {
			return "", err
		}
		logger.WithField("name", rolePolicyName).Infoln("Created role policy")
	}

	// We sleep here otherwise got an error when creating the ec2 instance referencing the profile.
	time.Sleep(10 * time.Second)

	return aws.StringValue(instanceProfile.Arn), nil
}

func createInstance(l *logrus.Logger, ec2Client ec2iface.EC2API, options instanceOptions) (*ec2.Instance, error) {
	// Check if an instance exists.
	instance, err := ec2GetInstance(ec2Client, []*ec2.Filter{
		ec2CreateFilter("tag:Name", options.name),
	})
	if err != nil {
		return nil, err
	}
	if instance != nil {
		l.WithField("id", aws.StringValue(instance.InstanceId)).Infoln("Instance already exists")
		return instance, nil
	}

	// Create a new EC2 instance.
	instance, err = ec2CreateInstance(ec2Client, options)
	if err != nil {
		return nil, err
	}
	l.WithField("id", aws.StringValue(instance.InstanceId)).Infoln("Created instance")

	return instance, nil
}
