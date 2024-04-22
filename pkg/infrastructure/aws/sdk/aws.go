package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset"
	tfvarsAsset "github.com/openshift/installer/pkg/asset/cluster/tfvars"
	awssession "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/tfvars"
	awstfvars "github.com/openshift/installer/pkg/tfvars/aws"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/version"
)

const (
	tfVarsFileName         = "terraform.tfvars.json"
	tfPlatformVarsFileName = "terraform.platform.auto.tfvars.json"
	clusterOutputFileName  = "cluster.awssdk.vars.json"

	defaultDescription = "Created by Openshift Installer"
	ownedTagKey        = "kubernetes.io/cluster/%s"
	ownedTagValue      = "owned"
)

// InfraProvider is the AWS SDK infra provider.
type InfraProvider struct{}

// InitializeProvider initializes the AWS SDK provider.
func InitializeProvider() infrastructure.Provider {
	return InfraProvider{}
}

type output struct {
	VpcID            string   `json:"vpc_id,omitempty"`
	MasterSGID       string   `json:"master_sg_id,omitempty"`
	WorkerSGID       string   `json:"worker_sg_id,omitempty"`
	BootstrapIP      string   `json:"bootstrap_ip,omitempty"`
	ControlPlaneIPs  []string `json:"master_ips,omitempty"`
	TargetGroupARNs  []string `json:"lb_target_group_arns,omitempty"`
	PublicSubnetIDs  []string `json:"public_subnet_ids,omitempty"`
	PrivateSubnetIDs []string `json:"private_subnet_ids,omitempty"`
}

// Provision creates cluster infrastructure using AWS SDK calls.
func (a InfraProvider) Provision(ctx context.Context, dir string, parents asset.Parents) ([]*asset.File, error) {
	terraformVariables := &tfvarsAsset.TerraformVariables{}
	parents.Get(terraformVariables)
	// Unmarshall input from tf variables, so we can use it along with
	// installConfig and other assets as the contractual input regardless of
	// the implementation.
	clusterConfig := &tfvars.Config{}
	clusterAWSConfig := &awstfvars.Config{}
	for _, file := range terraformVariables.Files() {
		switch file.Filename {
		case tfVarsFileName:
			if err := json.Unmarshal(file.Data, clusterConfig); err != nil {
				return nil, err
			}
		case tfPlatformVarsFileName:
			if err := json.Unmarshal(file.Data, clusterAWSConfig); err != nil {
				return nil, err
			}
		}
	}

	if clusterConfig == (&tfvars.Config{}) || clusterAWSConfig == (&awstfvars.Config{}) {
		return nil, fmt.Errorf("could not find terraform config files")
	}

	eps := []awstypes.ServiceEndpoint{}
	for k, v := range clusterAWSConfig.CustomEndpoints {
		eps = append(eps, awstypes.ServiceEndpoint{Name: k, URL: v})
	}

	awsSession, err := awssession.GetSessionWithOptions(
		awssession.WithRegion(clusterAWSConfig.Region),
		awssession.WithServiceEndpoints(clusterAWSConfig.Region, eps),
	)
	if err != nil {
		return nil, err
	}

	availabilityZones := sets.New(clusterAWSConfig.MasterSecurityGroups...)
	availabilityZones.Insert(clusterAWSConfig.WorkerAvailabilityZones...)

	tags := mergeTags(clusterAWSConfig.ExtraTags, map[string]string{
		clusterOwnedTag(clusterConfig.ClusterID): ownedTagValue,
	})

	usePublicEndpoints := clusterAWSConfig.PublishStrategy == "External"

	logger := logrus.StandardLogger()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	ec2Client := ec2.New(awsSession)
	amiID := clusterAWSConfig.AMI
	if clusterAWSConfig.Region != clusterAWSConfig.AMIRegion {
		logger.Infof("Copying AMI to region %s", clusterAWSConfig.Region)
		amiID, err = copyAMIToRegion(ctx, ec2Client, clusterAWSConfig.AMI, clusterAWSConfig.AMIRegion, clusterAWSConfig.Region, clusterConfig.ClusterID, tags)
		if err != nil {
			return nil, fmt.Errorf("failed to copy AMI to region (%s): %w", clusterAWSConfig.Region, err)
		}
	}

	logger.Infoln("Creating VPC resources")
	vpcInput := vpcInputOptions{
		infraID:          clusterConfig.ClusterID,
		region:           clusterAWSConfig.Region,
		vpcID:            clusterAWSConfig.VPC,
		cidrV4Block:      clusterConfig.MachineV4CIDRs[0],
		zones:            sets.List(availabilityZones),
		tags:             tags,
		privateSubnetIDs: clusterAWSConfig.PrivateSubnets,
		edgeZones:        clusterAWSConfig.EdgeLocalZones,
		edgeParentMap:    clusterAWSConfig.EdgeZonesGatewayIndex,
	}
	if clusterAWSConfig.PublicSubnets != nil {
		vpcInput.publicSubnetIDs = *clusterAWSConfig.PublicSubnets
	}

	vpcOutput, err := createVPCResources(ctx, logger, ec2Client, &vpcInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create VPC resources: %w", err)
	}

	logger.Infoln("Creating Load Balancer resources")
	elbClient := elbv2.New(awsSession)
	lbInput := lbInputOptions{
		infraID:          clusterConfig.ClusterID,
		vpcID:            vpcOutput.vpcID,
		privateSubnetIDs: vpcOutput.privateSubnetIDs,
		publicSubnetIDs:  vpcOutput.publicSubnetIDs,
		tags:             tags,
		isPrivateCluster: !usePublicEndpoints,
	}
	lbOutput, err := createLoadBalancers(ctx, logger, elbClient, &lbInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancers: %w", err)
	}

	logger.Infoln("Creating DNS resources")
	r53Config := awssession.GetR53ClientCfg(awsSession, clusterAWSConfig.InternalZoneRole)
	if len(clusterAWSConfig.InternalZoneRole) > 0 {
		logger.WithField("role", clusterAWSConfig.InternalZoneRole).Debugln("Assuming role for private hosted zone")
	}
	r53Client := route53.New(awsSession)
	assumedRoleClient := route53.New(awsSession, r53Config)
	dnsInput := dnsInputOptions{
		infraID:           clusterConfig.ClusterID,
		region:            clusterAWSConfig.Region,
		baseDomain:        clusterConfig.BaseDomain,
		clusterDomain:     clusterConfig.ClusterDomain,
		vpcID:             vpcOutput.vpcID,
		tags:              tags,
		lbExternalZoneID:  lbOutput.external.zoneID,
		lbExternalZoneDNS: lbOutput.external.dnsName,
		lbInternalZoneID:  lbOutput.internal.zoneID,
		lbInternalZoneDNS: lbOutput.internal.dnsName,
		isPrivateCluster:  !usePublicEndpoints,
		internalZone:      clusterAWSConfig.InternalZone,
	}
	err = createDNSResources(ctx, logger, r53Client, assumedRoleClient, &dnsInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS rsources: %w", err)
	}

	logger.Infoln("Creating security groups")
	sgInput := sgInputOptions{
		infraID:          clusterConfig.ClusterID,
		vpcID:            vpcOutput.vpcID,
		cidrV4Blocks:     clusterConfig.MachineV4CIDRs,
		isPrivateCluster: !usePublicEndpoints,
		tags:             tags,
	}
	sgOutput, err := createSecurityGroups(ctx, logger, ec2Client, &sgInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create security groups: %w", err)
	}

	partitionDNSSuffix := "amazonaws.com"
	if ps, found := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), clusterAWSConfig.Region); found {
		partitionDNSSuffix = ps.DNSSuffix()
	}
	logger.Debugf("Using partition DNS suffix: %s", partitionDNSSuffix)

	logger.Infoln("Creating bootstrap resources")
	bootstrapSubnet := vpcOutput.privateSubnetIDs[0]
	if usePublicEndpoints {
		bootstrapSubnet = vpcOutput.publicSubnetIDs[0]
	}
	bootstrapInput := bootstrapInputOptions{
		instanceInputOptions: instanceInputOptions{
			infraID:            clusterConfig.ClusterID,
			amiID:              amiID,
			instanceType:       clusterAWSConfig.MasterInstanceType,
			iamRole:            clusterAWSConfig.MasterIAMRoleName,
			volumeType:         "gp2",
			volumeSize:         30,
			volumeIOPS:         0,
			isEncrypted:        true,
			metadataAuth:       clusterAWSConfig.BootstrapMetadataAuthentication,
			kmsKeyID:           clusterAWSConfig.KMSKeyID,
			securityGroupIds:   []string{sgOutput.bootstrap, sgOutput.controlPlane},
			targetGroupARNs:    lbOutput.targetGroupArns,
			subnetID:           bootstrapSubnet,
			associatePublicIP:  usePublicEndpoints,
			userData:           clusterAWSConfig.BootstrapIgnitionStub,
			partitionDNSSuffix: partitionDNSSuffix,
			tags:               tags,
		},
		ignitionBucket:  clusterAWSConfig.IgnitionBucket,
		ignitionContent: clusterConfig.IgnitionBootstrap,
	}
	iamClient := iam.New(awsSession)
	s3Client := s3.New(awsSession)
	bootstrapOut, err := createBootstrapResources(ctx, logger, ec2Client, iamClient, s3Client, elbClient, &bootstrapInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create bootstrap resources: %w", err)
	}

	logger.Infoln("Creating control plane resources")
	controlPlaneInput := controlPlaneInputOptions{
		instanceInputOptions: instanceInputOptions{
			infraID:            clusterConfig.ClusterID,
			amiID:              amiID,
			instanceType:       clusterAWSConfig.MasterInstanceType,
			iamRole:            clusterAWSConfig.MasterIAMRoleName,
			volumeType:         clusterAWSConfig.Type,
			volumeSize:         clusterAWSConfig.Size,
			volumeIOPS:         clusterAWSConfig.IOPS,
			isEncrypted:        clusterAWSConfig.Encrypted,
			kmsKeyID:           clusterAWSConfig.KMSKeyID,
			metadataAuth:       clusterAWSConfig.MasterMetadataAuthentication,
			securityGroupIds:   append(clusterAWSConfig.MasterSecurityGroups, sgOutput.controlPlane),
			targetGroupARNs:    lbOutput.targetGroupArns,
			associatePublicIP:  len(os.Getenv("OPENSHIFT_INSTALL_AWS_PUBLIC_ONLY")) > 0,
			userData:           clusterConfig.IgnitionMaster,
			partitionDNSSuffix: partitionDNSSuffix,
			tags:               tags,
		},
		nReplicas:         clusterConfig.Masters,
		privateSubnetIDs:  vpcOutput.privateSubnetIDs,
		zoneToSubnetMap:   vpcOutput.zoneToSubnetMap,
		availabilityZones: clusterAWSConfig.MasterAvailabilityZones,
	}
	controlPlaneOut, err := createControlPlaneResources(ctx, logger, ec2Client, iamClient, elbClient, &controlPlaneInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create control plane resources: %w", err)
	}

	logger.Infoln("Creating compute resources")
	computeInput := computeInputOptions{
		infraID:            clusterConfig.ClusterID,
		partitionDNSSuffix: partitionDNSSuffix,
		tags:               tags,
	}
	err = createComputeResources(ctx, logger, iamClient, &computeInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute resources: %w", err)
	}

	bootstrapIP := bootstrapOut.privateIP
	if usePublicEndpoints {
		bootstrapIP = bootstrapOut.publicIP
	}
	out := &output{
		BootstrapIP:      bootstrapIP,
		VpcID:            vpcOutput.vpcID,
		TargetGroupARNs:  lbOutput.targetGroupArns,
		PublicSubnetIDs:  vpcOutput.publicSubnetIDs,
		PrivateSubnetIDs: vpcOutput.privateSubnetIDs,
		MasterSGID:       sgOutput.controlPlane,
		WorkerSGID:       sgOutput.compute,
		ControlPlaneIPs:  controlPlaneOut.controlPlaneIPs,
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to write cluster output: %w", err)
	}
	return []*asset.File{
		{Filename: clusterOutputFileName, Data: data},
	}, nil
}

// DestroyBootstrap destroys the temporary bootstrap resources.
func (a InfraProvider) DestroyBootstrap(dir string, _ *types.ClusterMetadata) error {
	// Unmarshall input from tf variables, so we can use it along with
	// installConfig and other assets as the contractual input regardless of
	// the implementation.
	clusterConfig := &tfvars.Config{}
	data, err := os.ReadFile(filepath.Join(dir, tfVarsFileName))
	if err == nil {
		err = json.Unmarshal(data, clusterConfig)
	}
	if err != nil {
		return fmt.Errorf("failed to load cluster terraform variables: %w", err)
	}
	clusterAWSConfig := &awstfvars.Config{}
	data, err = os.ReadFile(filepath.Join(dir, tfPlatformVarsFileName))
	if err == nil {
		err = json.Unmarshal(data, clusterAWSConfig)
	}
	if err != nil {
		return fmt.Errorf("failed to load AWS terraform variables: %w", err)
	}

	eps := []awstypes.ServiceEndpoint{}
	for k, v := range clusterAWSConfig.CustomEndpoints {
		eps = append(eps, awstypes.ServiceEndpoint{Name: k, URL: v})
	}

	awsSession, err := awssession.GetSessionWithOptions(
		awssession.WithRegion(clusterAWSConfig.Region),
		awssession.WithServiceEndpoints(clusterAWSConfig.Region, eps),
	)
	if err != nil {
		return err
	}
	awsSession.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "openshiftInstaller.OpenshiftInstallerUserAgentHandler",
		Fn:   request.MakeAddToUserAgentHandler("OpenShift/4.x Destroyer", version.Raw),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	logger := logrus.StandardLogger()
	input := &destroyInputOptions{
		infraID:          clusterConfig.ClusterID,
		region:           clusterAWSConfig.Region,
		ignitionBucket:   clusterAWSConfig.IgnitionBucket,
		preserveIgnition: clusterAWSConfig.PreserveBootstrapIgnition,
	}
	err = destroyBootstrapResources(ctx, logger, awsSession, input)
	if err != nil {
		return fmt.Errorf("failed to delete bootstrap resources: %w", err)
	}

	return nil
}

// ExtractHostAddresses extracts the IPs of the bootstrap and control plane machines.
func (a InfraProvider) ExtractHostAddresses(dir string, ic *types.InstallConfig, ha *infrastructure.HostAddresses) error {
	clusterOutput := &output{}
	data, err := os.ReadFile(filepath.Join(dir, clusterOutputFileName))
	if err == nil {
		err = json.Unmarshal(data, clusterOutput)
	}
	if err != nil {
		return fmt.Errorf("failed to load cluster terraform variables: %w", err)
	}

	ha.Bootstrap = clusterOutput.BootstrapIP
	ha.Masters = append(ha.Masters, clusterOutput.ControlPlaneIPs...)

	return nil
}

func mergeTags(lhsTags, rhsTags map[string]string) map[string]string {
	merged := make(map[string]string, len(lhsTags)+len(rhsTags))
	for k, v := range lhsTags {
		merged[k] = v
	}
	for k, v := range rhsTags {
		merged[k] = v
	}
	return merged
}

func clusterOwnedTag(infraID string) string {
	return fmt.Sprintf(ownedTagKey, infraID)
}

func copyAMIToRegion(ctx context.Context, client ec2iface.EC2API, sourceAMI string, sourceRegion string, targetRegion string, infraID string, tags map[string]string) (string, error) {
	name := fmt.Sprintf("%s-ami-%s", infraID, targetRegion)
	amiTags := mergeTags(tags, map[string]string{
		"Name":         name,
		"sourceAMI":    sourceAMI,
		"sourceRegion": sourceRegion,
	})
	res, err := client.CopyImageWithContext(ctx, &ec2.CopyImageInput{
		Name:          aws.String(fmt.Sprintf("%s-master", infraID)),
		ClientToken:   aws.String(infraID),
		Description:   aws.String(defaultDescription),
		SourceImageId: aws.String(sourceAMI),
		SourceRegion:  aws.String(sourceRegion),
		Encrypted:     aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	_, err = client.CreateTagsWithContext(ctx, &ec2.CreateTagsInput{
		Resources: []*string{res.ImageId},
		Tags:      ec2Tags(amiTags),
	})
	if err != nil {
		return "", fmt.Errorf("failed to tag AMI copy (%s): %w", name, err)
	}

	return aws.StringValue(res.ImageId), nil
}
