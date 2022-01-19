package alibabacloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/tag"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	icalibabacloud "github.com/openshift/installer/pkg/asset/installconfig/alibabacloud"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
)

// ClusterUninstaller holds the various options for the cluster we want to delete
type ClusterUninstaller struct {
	Logger          logrus.FieldLogger
	AccessKeyID     string
	AccessKeySecret string
	Auth            auth.Credential

	Region        string
	InfraID       string
	ClusterDomain string
	Tags          []map[string]string
	TagResources  struct {
		ecsInstances   []ResourceArn
		securityGroups []ResourceArn
		vpcs           []ResourceArn
		vSwitchs       []ResourceArn
		eips           []ResourceArn
		natgateways    []ResourceArn
		slbs           []ResourceArn
		buckets        []ResourceArn
		others         []ResourceArn
	}

	ecsClient      *ecs.Client
	dnsClient      *alidns.Client
	pvtzClient     *pvtz.Client
	vpcClient      *vpc.Client
	ramClient      *ram.Client
	tagClient      *tag.Client
	slbClient      *slb.Client
	ossClient      *oss.Client
	rmanagerClient *resourcemanager.Client
}

// ResourceArn holds the information contained in the cloud resource Arn string
type ResourceArn struct {
	Service      string
	Region       string
	Account      string
	ResourceType string
	ResourceID   string
	Arn          string
}

func (o *ClusterUninstaller) configureClients() error {
	var err error
	config := sdk.NewConfig()
	config.AutoRetry = true
	config.MaxRetryTime = 3
	ossEndpoint := fmt.Sprintf("http://oss-%s.aliyuncs.com", o.Region)

	o.ecsClient, err = ecs.NewClientWithOptions(o.Region, config, o.Auth)
	if err != nil {
		return err
	}

	o.dnsClient, err = alidns.NewClientWithOptions(o.Region, config, o.Auth)
	if err != nil {
		return err
	}

	o.pvtzClient, err = pvtz.NewClientWithOptions(o.Region, config, o.Auth)
	if err != nil {
		return err
	}

	o.ramClient, err = ram.NewClientWithOptions(o.Region, config, o.Auth)
	if err != nil {
		return err
	}

	o.vpcClient, err = vpc.NewClientWithOptions(o.Region, config, o.Auth)
	if err != nil {
		return err
	}

	o.tagClient, err = tag.NewClientWithOptions(o.Region, config, o.Auth)
	if err != nil {
		return err
	}

	o.slbClient, err = slb.NewClientWithOptions(o.Region, config, o.Auth)
	if err != nil {
		return err
	}

	o.ossClient, err = oss.New(ossEndpoint, o.AccessKeyID, o.AccessKeySecret)
	if err != nil {
		return err
	}

	o.rmanagerClient, err = resourcemanager.NewClientWithOptions(o.Region, config, o.Auth)
	if err != nil {
		return err
	}

	return nil
}

// New returns an Alibaba Cloud destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	region := metadata.ClusterPlatformMetadata.AlibabaCloud.Region
	client, err := icalibabacloud.NewClient(region)
	if err != nil {
		return nil, err
	}

	auth := credentials.NewAccessKeyCredential(client.AccessKeyID, client.AccessKeySecret)

	return &ClusterUninstaller{
		Logger:          logger,
		Auth:            auth,
		AccessKeyID:     client.AccessKeyID,
		AccessKeySecret: client.AccessKeySecret,
		Region:          region,
		InfraID:         metadata.InfraID,
		ClusterDomain:   metadata.AlibabaCloud.ClusterDomain,
		Tags: []map[string]string{
			{
				fmt.Sprintf("kubernetes.io/cluster/%s", metadata.InfraID): "owned",
			},
			{
				"ack.aliyun.com": metadata.InfraID,
			},
		},
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	var err error

	err = o.configureClients()
	if err != nil {
		return nil, err
	}

	err = o.findResources()
	if err != nil {
		return nil, err
	}

	err = o.destroyCluster()
	if err != nil {
		return nil, errors.Wrap(err, "failed to destroy cluster")
	}

	return nil, nil
}

func (o *ClusterUninstaller) destroyCluster() error {
	stagedFuncs := [][]struct {
		name    string
		execute func(logrus.FieldLogger) error
	}{
		{
			{name: "DNS records", execute: o.deleteDNSRecords},
			{name: "OSS buckets", execute: o.deleteBuckets},
			{name: "RAM roles", execute: o.deleteRAMRoles},
			{name: "ECS instances", execute: o.deleteEcsInstances},
		},
		{
			{name: "private zones", execute: o.deletePrivateZones},
			{name: "ECS security groups", execute: o.deleteSecurityGroups},
			{name: "Nat gateways", execute: o.deleteNatGateways},
			{name: "SLBs", execute: o.deleteSlbs},
		},
		{
			{name: "EIPs", execute: o.deleteEips},
		},
		{
			{name: "VSwitchs", execute: o.deleteVSwitchs},
		},
		{
			{name: "VPCs", execute: o.deleteVpcs},
		},
		{
			{name: "resource groups", execute: o.deleteResourceGroup},
		},
	}

	for _, stage := range stagedFuncs {
		var wg sync.WaitGroup
		errCh := make(chan error)
		wgDone := make(chan bool)

		for _, f := range stage {
			wg.Add(1)
			go o.executeStageFunction(f, errCh, &wg)
		}

		go func() {
			wg.Wait()
			close(wgDone)
		}()

		select {
		case <-wgDone:
			// On to the next stage
			continue
		case err := <-errCh:
			return err
		}
	}

	return nil
}

func (o *ClusterUninstaller) executeStageFunction(f struct {
	name    string
	execute func(logrus.FieldLogger) error
}, errCh chan error, wg *sync.WaitGroup) error {
	defer wg.Done()

	err := wait.PollImmediateInfinite(
		time.Second*10,
		func() (bool, error) {
			stageLogger := o.Logger.WithField("stage", f.name)
			ferr := f.execute(stageLogger)
			if ferr != nil {
				stageLogger.WithError(ferr).Debugf("Error executing stage")
				return false, nil
			}
			return true, nil
		},
	)

	if err != nil {
		errCh <- err
	}
	return nil
}

func (o *ClusterUninstaller) findResources() (err error) {
	tagResources, err := o.findResourcesByTag()
	if err != nil {
		return errors.Wrap(err, "failed to find resource by tag")
	}
	if len(tagResources) > 0 {
		var deletedResources []ResourceArn
		for _, resource := range tagResources {
			arn, err := convertResourceArn(resource.ResourceARN)
			if err != nil {
				continue
			}
			deletedResources = append(deletedResources, arn)
		}

		o.tidyResourceByArn(deletedResources)
	}
	return nil
}

func (o *ClusterUninstaller) tidyResourceByArn(resourceArns []ResourceArn) {
	for _, resourceArn := range resourceArns {
		switch resourceArn.Service {
		case "ecs":
			switch resourceArn.ResourceType {
			case "instance":
				o.TagResources.ecsInstances = append(o.TagResources.ecsInstances, resourceArn)
			case "securitygroup":
				o.TagResources.securityGroups = append(o.TagResources.securityGroups, resourceArn)
			default:
				o.TagResources.others = append(o.TagResources.others, resourceArn)
			}
		case "vpc":
			switch resourceArn.ResourceType {
			case "vpc":
				o.TagResources.vpcs = append(o.TagResources.vpcs, resourceArn)
			case "vswitch":
				o.TagResources.vSwitchs = append(o.TagResources.vSwitchs, resourceArn)
			case "eip":
				o.TagResources.eips = append(o.TagResources.eips, resourceArn)
			case "natgateway":
				o.TagResources.natgateways = append(o.TagResources.natgateways, resourceArn)
			default:
				o.TagResources.others = append(o.TagResources.others, resourceArn)
			}
		case "slb":
			switch resourceArn.ResourceType {
			case "instance":
				o.TagResources.slbs = append(o.TagResources.slbs, resourceArn)
			default:
				o.TagResources.others = append(o.TagResources.others, resourceArn)
			}
		case "oss":
			switch resourceArn.ResourceType {
			case "bucket":
				o.TagResources.buckets = append(o.TagResources.buckets, resourceArn)
			default:
				o.TagResources.others = append(o.TagResources.others, resourceArn)
			}
		default:
			o.TagResources.others = append(o.TagResources.others, resourceArn)
		}
	}
}

func (o *ClusterUninstaller) deleteResourceGroup(logger logrus.FieldLogger) (err error) {
	resourceGroupName := fmt.Sprintf("%s-rg", o.InfraID)
	logger = logger.WithField("name", resourceGroupName)
	logger.Debug("Searching resource groups")

	response, err := o.listResourceGroups(resourceGroupName)
	if err != nil {
		return err
	}

	resourceGroupID := ""
	for _, resourceGroup := range response.ResourceGroups.ResourceGroup {
		if resourceGroup.Name == resourceGroupName {
			resourceGroupID = resourceGroup.Id
		}
	}
	if resourceGroupID == "" {
		return
	}

	err = o.deleteResourceGroupByID(resourceGroupID, logger)
	if err != nil {
		return err
	}

	err = wait.Poll(
		2*time.Second,
		1*time.Minute,
		func() (bool, error) {
			response, err := o.listResourceGroups(resourceGroupName)
			if err != nil {
				return false, err
			}
			for _, resourceGroup := range response.ResourceGroups.ResourceGroup {
				if resourceGroup.Name == resourceGroupName {
					return false, nil
				}
			}
			return true, nil
		},
	)

	logger.Info("Resource group deleted")
	return
}

func (o *ClusterUninstaller) deleteResourceGroupByID(resourceGroupID string, logger logrus.FieldLogger) (err error) {
	logger.WithField("resourceGroupID", resourceGroupID).Debug("Deleting")
	request := resourcemanager.CreateDeleteResourceGroupRequest()
	request.Scheme = "https"
	request.ResourceGroupId = resourceGroupID
	_, err = o.rmanagerClient.DeleteResourceGroup(request)
	return
}

func (o *ClusterUninstaller) listResourceGroups(resourceGroupName string) (response *resourcemanager.ListResourceGroupsResponse, err error) {
	request := resourcemanager.CreateListResourceGroupsRequest()
	request.Scheme = "https"
	request.QueryParams["Name"] = resourceGroupName
	response, err = o.rmanagerClient.ListResourceGroups(request)
	return
}

func (o *ClusterUninstaller) deleteBuckets(logger logrus.FieldLogger) (err error) {
	if len(o.TagResources.buckets) <= 0 {
		return nil
	}

	var bucketNames []string
	for _, bucketArn := range o.TagResources.buckets {
		bucketNames = append(bucketNames, bucketArn.ResourceID)
	}

	for _, bucketName := range bucketNames {
		ossLogger := logger.WithField("bucketName", bucketName)
		err = o.deleteBucket(bucketName, ossLogger)
		if err != nil {
			return err
		}
	}
	logger.Info("OSS buckets deleted")
	return
}

func (o *ClusterUninstaller) deleteBucket(bucketName string, logger logrus.FieldLogger) (err error) {
	logger.Debug("Searching OSS bucket")
	result, err := o.ossClient.ListBuckets(oss.Prefix(bucketName))
	if err != nil || len(result.Buckets) == 0 {
		return
	}

	keys := []string{fmt.Sprintf("kubernetes.io/cluster/%s", o.InfraID)}
	arns := []string{fmt.Sprintf("arn:acs:oss:%s:*:bucket/%s", o.Region, bucketName)}
	logger.WithField("tags", keys).Debug("Unbinding tags for OSS bucket")
	err = o.unTagResource(&keys, &arns)
	if err != nil {
		return err
	}

	bucket, err := o.ossClient.Bucket(bucketName)
	if err != nil {
		return err
	}

	err = o.deleteObjects(bucket, logger)
	if err != nil {
		return err
	}

	logger.Debug("Deleting OSS bucket")
	err = o.ossClient.DeleteBucket(bucketName)
	if err != nil {
		return err
	}

	err = wait.Poll(
		2*time.Second,
		2*time.Minute,
		func() (bool, error) {
			result, err := o.ossClient.ListBuckets(oss.Prefix(bucketName))
			if err != nil {
				return false, err
			}
			if len(result.Buckets) == 0 {
				return true, nil
			}
			return false, nil
		},
	)

	logger.Info("OSS bucket deleted")
	return
}

func (o *ClusterUninstaller) deleteObjects(bucket *oss.Bucket, logger logrus.FieldLogger) (err error) {
	logger.Debug("Searching OSS bucket objects")
	result, err := bucket.ListObjectsV2()
	if err != nil {
		return err
	}
	if len(result.Objects) == 0 {
		return
	}

	var keys []string
	for _, object := range result.Objects {
		keys = append(keys, object.Key)
	}
	logger = logger.WithField("objects", keys)
	logger.Debug("Deleting bucket objects")
	_, err = bucket.DeleteObjects(keys)
	if err != nil {
		return err
	}

	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			result, err := bucket.ListObjectsV2()
			if err != nil {
				return false, err
			}
			if len(result.Objects) == 0 {
				return true, nil
			}
			return false, nil
		},
	)
	logger.Debug("OSS bucket objects deleted")
	return
}

func (o *ClusterUninstaller) deleteSlbs(logger logrus.FieldLogger) (err error) {
	if len(o.TagResources.slbs) <= 0 {
		return nil
	}

	var slbIDs []string
	for _, slbArn := range o.TagResources.slbs {
		slbIDs = append(slbIDs, slbArn.ResourceID)
	}

	logger.WithField("SLBs", slbIDs).Debug("Deleting")
	for _, slbID := range slbIDs {
		slbLogger := logger.WithField("slbID", slbID)
		err = o.setSlbModificationProtection(slbID, slbLogger)
		if err != nil {
			return err
		}

		err = o.setSlbDeleteProtection(slbID, slbLogger)
		if err != nil {
			return err
		}

		err = o.deleteSlb(slbID, slbLogger)
		if err != nil {
			return err
		}
	}

	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			response, err := o.listSlb(slbIDs)
			if err != nil {
				return false, err
			}
			if response.TotalCount == 0 {
				return true, nil
			}
			return false, nil
		},
	)

	logger.Info("SLB instances deleted")
	return
}

func (o *ClusterUninstaller) listSlb(slbIDs []string) (response *slb.DescribeLoadBalancersResponse, err error) {
	request := slb.CreateDescribeLoadBalancersRequest()
	request.LoadBalancerId = strings.Join(slbIDs, ",")
	response, err = o.slbClient.DescribeLoadBalancers(request)
	return
}

func (o *ClusterUninstaller) setSlbModificationProtection(slbID string, logger logrus.FieldLogger) (err error) {
	logger.Debug("Turn off the modification protection")
	request := slb.CreateSetLoadBalancerModificationProtectionRequest()
	request.LoadBalancerId = slbID
	request.ModificationProtectionStatus = "NonProtection"
	_, err = o.slbClient.SetLoadBalancerModificationProtection(request)
	return
}

func (o *ClusterUninstaller) setSlbDeleteProtection(slbID string, logger logrus.FieldLogger) (err error) {
	logger.Debug("Turn off the deletion protection")
	request := slb.CreateSetLoadBalancerDeleteProtectionRequest()
	request.LoadBalancerId = slbID
	request.DeleteProtection = "off"
	_, err = o.slbClient.SetLoadBalancerDeleteProtection(request)
	return
}

func (o *ClusterUninstaller) deleteSlb(slbID string, logger logrus.FieldLogger) (err error) {
	logger.Debug("Deleting")
	request := slb.CreateDeleteLoadBalancerRequest()
	request.LoadBalancerId = slbID
	_, err = o.slbClient.DeleteLoadBalancer(request)
	return
}

func (o *ClusterUninstaller) deleteVSwitchs(logger logrus.FieldLogger) (err error) {
	if len(o.TagResources.vSwitchs) <= 0 {
		return nil
	}

	var vSwitchIDs []string
	for _, vSwitchArn := range o.TagResources.vSwitchs {
		vSwitchIDs = append(vSwitchIDs, vSwitchArn.ResourceID)
	}

	logger.WithField("vSwitchIDs", vSwitchIDs).Debug("Deleting VSwitches")
	for _, vSwitchID := range vSwitchIDs {
		err = wait.Poll(
			5*time.Second,
			30*time.Second,
			func() (bool, error) {
				err = o.deleteVSwitch(vSwitchID, logger)
				if err == nil {
					return true, nil
				}
				if strings.Contains(err.Error(), "DependencyViolation") {
					return false, nil
				}
				return false, err
			},
		)
		if err != nil {
			return err
		}
		logger.WithField("vSwitchID", vSwitchID).Debug("Deleted")
	}

	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			response, err := o.listVSwitch(vSwitchIDs)
			if err != nil {
				return false, err
			}
			if response.TotalCount == 0 {
				return true, nil
			}
			return false, nil
		},
	)

	logger.Info("VSwitches deleted")
	return
}

func (o *ClusterUninstaller) listVSwitch(vSwitchIDs []string) (response *vpc.DescribeVSwitchesResponse, err error) {
	request := vpc.CreateDescribeVSwitchesRequest()
	request.VSwitchId = strings.Join(vSwitchIDs, ",")
	response, err = o.vpcClient.DescribeVSwitches(request)
	return
}

func (o *ClusterUninstaller) deleteVSwitch(vSwitchID string, logger logrus.FieldLogger) (err error) {
	logger.WithField("vSwitchID", vSwitchID).Debug("Deleting")
	request := vpc.CreateDeleteVSwitchRequest()
	request.VSwitchId = vSwitchID
	_, err = o.vpcClient.DeleteVSwitch(request)
	return
}

func (o *ClusterUninstaller) deleteVpcs(logger logrus.FieldLogger) (err error) {
	if len(o.TagResources.vpcs) <= 0 {
		return nil
	}

	var vpcIDs []string
	for _, vpcArn := range o.TagResources.vpcs {
		vpcIDs = append(vpcIDs, vpcArn.ResourceID)
	}

	logger.WithField("vpcIDs", vpcIDs).Debug("Deleting VPCs")
	for _, vpcID := range vpcIDs {
		err = wait.Poll(
			5*time.Second,
			30*time.Second,
			func() (bool, error) {
				err = o.deleteVpc(vpcID, logger)
				if err == nil {
					return true, nil
				}
				if strings.Contains(err.Error(), "DependencyViolation") {
					return false, nil
				}
				return false, err
			},
		)
		if err != nil {
			return err
		}
		logger.WithField("vpcID", vpcID).Debug("Deleted")
	}

	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			response, err := o.listVpc(vpcIDs)
			if err != nil {
				return false, err
			}
			if response.TotalCount == 0 {
				return true, nil
			}
			return false, nil
		},
	)

	logger.Info("VPCs deleted")
	return
}

func (o *ClusterUninstaller) deleteVpc(vpcID string, logger logrus.FieldLogger) (err error) {
	logger.WithField("vpcID", vpcID).Debug("Deleting")
	request := vpc.CreateDeleteVpcRequest()
	request.VpcId = vpcID
	_, err = o.vpcClient.DeleteVpc(request)
	return
}

func (o *ClusterUninstaller) listVpc(vpcIDs []string) (response *vpc.DescribeVpcsResponse, err error) {
	request := vpc.CreateDescribeVpcsRequest()
	request.VpcId = strings.Join(vpcIDs, ",")
	response, err = o.vpcClient.DescribeVpcs(request)
	return
}

func (o *ClusterUninstaller) deleteEips(logger logrus.FieldLogger) (err error) {
	if len(o.TagResources.eips) <= 0 {
		return nil
	}

	var eipIDs []string
	for _, eipArn := range o.TagResources.eips {
		eipIDs = append(eipIDs, eipArn.ResourceID)
	}

	logger.WithField("eipIDs", eipIDs).Debug("Deleting EIPs")
	for _, eipID := range eipIDs {
		err = o.deleteEip(eipID, logger)
		if err != nil {
			return err
		}
	}
	err = wait.Poll(
		2*time.Second,
		2*time.Minute,
		func() (bool, error) {
			response, err := o.listEip(eipIDs)
			if err != nil {
				return false, err
			}
			if response.TotalCount == 0 {
				return true, nil
			}
			return false, nil
		},
	)
	logger.Info("EIPs deleted")
	return err
}

func (o *ClusterUninstaller) listEip(eipIDs []string) (response *vpc.DescribeEipAddressesResponse, err error) {
	request := vpc.CreateDescribeEipAddressesRequest()
	request.AllocationId = strings.Join(eipIDs, ",")
	response, err = o.vpcClient.DescribeEipAddresses(request)
	return response, err
}

func (o *ClusterUninstaller) deleteEip(eipID string, logger logrus.FieldLogger) (err error) {
	logger.WithField("eipID", eipID).Debug("Deleting")
	request := vpc.CreateReleaseEipAddressRequest()
	request.AllocationId = eipID
	_, err = o.vpcClient.ReleaseEipAddress(request)
	return
}

func (o *ClusterUninstaller) deleteNatGateways(logger logrus.FieldLogger) (err error) {
	if len(o.TagResources.natgateways) <= 0 {
		return nil
	}

	var natGatewayIDs []string
	for _, natGatewayArn := range o.TagResources.natgateways {
		natGatewayIDs = append(natGatewayIDs, natGatewayArn.ResourceID)
	}

	logger.WithField("natGatewayIDs", natGatewayIDs).Debug("Deleting NAT gateways")
	for _, natGatewayID := range natGatewayIDs {
		err = o.deleteNatGateway(natGatewayID, logger)
		if err != nil {
			return err
		}
		err = wait.Poll(
			3*time.Second,
			3*time.Minute,
			func() (bool, error) {
				response, err := o.listNatGateways(natGatewayID)
				if err != nil {
					return false, err
				}
				if response.TotalCount == 0 {
					return true, nil
				}
				return false, nil
			},
		)
		if err != nil {
			return err
		}
		logger.WithField("natGatewayID", natGatewayID).Debug("Deleted")
	}
	logger.Info("NAT gateways deleted")
	return
}

func (o *ClusterUninstaller) listNatGateways(natGatewayID string) (response *vpc.DescribeNatGatewaysResponse, err error) {
	request := vpc.CreateDescribeNatGatewaysRequest()
	request.NatGatewayId = natGatewayID
	response, err = o.vpcClient.DescribeNatGateways(request)
	return
}

func (o *ClusterUninstaller) deleteNatGateway(natGatewayID string, logger logrus.FieldLogger) (err error) {
	logger.WithField("natGatewayID", natGatewayID).Debug("Deleting")
	request := vpc.CreateDeleteNatGatewayRequest()
	request.NatGatewayId = natGatewayID
	request.Force = "true"
	_, err = o.vpcClient.DeleteNatGateway(request)
	return
}

func (o *ClusterUninstaller) deleteSecurityGroups(logger logrus.FieldLogger) (err error) {
	if len(o.TagResources.securityGroups) <= 0 {
		return nil
	}

	var securityGroupIDs []string
	for _, securityGroupArn := range o.TagResources.securityGroups {
		securityGroupIDs = append(securityGroupIDs, securityGroupArn.ResourceID)
	}

	logger.WithField("securityGroupIDs", securityGroupIDs).Debug("Revoking dependency for security groups")
	for _, securityGroupID := range securityGroupIDs {
		err = o.deleteSecurityGroupRules(securityGroupID, logger)
		if err != nil {
			return err
		}
	}

	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			response, err := o.listSecurityGroupReferences(securityGroupIDs)
			if err != nil {
				return false, err
			}
			if len(response.SecurityGroupReferences.SecurityGroupReference) == 0 {
				return true, nil
			}
			return false, nil
		},
	)
	if err != nil {
		return err
	}

	logger.WithField("securityGroupIDs", securityGroupIDs).Debug("Deleting security groups")
	for _, securityGroupID := range securityGroupIDs {
		err = wait.Poll(
			5*time.Second,
			30*time.Second,
			func() (bool, error) {
				err = o.deleteSecurityGroup(securityGroupID, logger)
				if err == nil {
					return true, nil
				}
				if strings.Contains(err.Error(), "DependencyViolation") {
					return false, nil
				}
				return false, err
			},
		)
		if err != nil {
			return err
		}
		logger.WithField("securityGroupID", securityGroupID).Debug("Deleted")
	}

	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			response, err := o.listSecurityGroup(securityGroupIDs)
			if err != nil {
				return false, err
			}
			if response.TotalCount == 0 {
				return true, nil
			}
			return false, nil
		},
	)

	logger.Info("Security groups deleted")
	return
}

func (o *ClusterUninstaller) deleteSecurityGroup(securityGroupID string, logger logrus.FieldLogger) (err error) {
	logger.WithField("securityGroupID", securityGroupID).Debug("Deleting")
	request := ecs.CreateDeleteSecurityGroupRequest()
	request.SecurityGroupId = securityGroupID
	_, err = o.ecsClient.DeleteSecurityGroup(request)
	return
}

func (o *ClusterUninstaller) deleteSecurityGroupRules(securityGroupID string, logger logrus.FieldLogger) (err error) {
	logger.WithField("securityGroupID", securityGroupID).Debug("Revoking")
	response, err := o.getSecurityGroup(securityGroupID)
	if err != nil {
		return err
	}
	for _, permission := range response.Permissions.Permission {
		if permission.SourceGroupId != "" {
			err = o.revokeSecurityGroup(securityGroupID, permission.SourceGroupId, permission.IpProtocol, permission.PortRange, permission.NicType)
		}
	}
	return
}

func (o *ClusterUninstaller) revokeSecurityGroup(securityGroupID string, sourceGroupID string, ipProtocol string, portRange string, nicType string) (err error) {
	request := ecs.CreateRevokeSecurityGroupRequest()
	request.SecurityGroupId = securityGroupID
	request.SourceGroupId = sourceGroupID
	request.IpProtocol = ipProtocol
	request.PortRange = portRange
	request.NicType = nicType

	_, err = o.ecsClient.RevokeSecurityGroup(request)
	return
}

func (o *ClusterUninstaller) getSecurityGroup(securityGroupID string) (response *ecs.DescribeSecurityGroupAttributeResponse, err error) {
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	request.SecurityGroupId = securityGroupID
	response, err = o.ecsClient.DescribeSecurityGroupAttribute(request)
	return
}

func (o *ClusterUninstaller) listSecurityGroupReferences(securityGroupIDs []string) (response *ecs.DescribeSecurityGroupReferencesResponse, err error) {
	request := ecs.CreateDescribeSecurityGroupReferencesRequest()
	request.SecurityGroupId = &securityGroupIDs
	response, err = o.ecsClient.DescribeSecurityGroupReferences(request)
	return
}

func (o *ClusterUninstaller) listSecurityGroup(securityGroupIDs []string) (response *ecs.DescribeSecurityGroupsResponse, err error) {
	request := ecs.CreateDescribeSecurityGroupsRequest()
	securityGroupIDsString, err := json.Marshal(securityGroupIDs)
	if err != nil {
		return nil, err
	}
	request.SecurityGroupIds = string(securityGroupIDsString)
	response, err = o.ecsClient.DescribeSecurityGroups(request)
	return
}

func (o *ClusterUninstaller) listEcsInstance(instanceIDs []string) (response *ecs.DescribeInstancesResponse, err error) {
	request := ecs.CreateDescribeInstancesRequest()
	instanceIDsString, err := json.Marshal(instanceIDs)
	if err != nil {
		return nil, err
	}
	request.InstanceIds = string(instanceIDsString)
	response, err = o.ecsClient.DescribeInstances(request)
	return
}

func (o *ClusterUninstaller) modifyDeletionProtection(instanceID string, logger logrus.FieldLogger) (err error) {
	logger.WithField("ecsID", instanceID).Debug("Turn off the deletion protection")
	request := ecs.CreateModifyInstanceAttributeRequest()
	request.InstanceId = instanceID
	request.DeletionProtection = "false"
	_, err = o.ecsClient.ModifyInstanceAttribute(request)
	return
}

func (o *ClusterUninstaller) modifyECSInstancesDeletionProtection(instanceIDs []string, logger logrus.FieldLogger) (err error) {
	response, err := o.listEcsInstance(instanceIDs)
	if err != nil {
		return err
	}
	for _, instance := range response.Instances.Instance {
		if instance.DeletionProtection {
			err := o.modifyDeletionProtection(instance.InstanceId, logger)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (o *ClusterUninstaller) deleteEcsInstances(logger logrus.FieldLogger) (err error) {
	if len(o.TagResources.ecsInstances) <= 0 {
		return nil
	}

	var instanceIDs []string
	for _, instanceArn := range o.TagResources.ecsInstances {
		instanceIDs = append(instanceIDs, instanceArn.ResourceID)
	}

	err = o.modifyECSInstancesDeletionProtection(instanceIDs, logger)
	if err != nil {
		return err
	}

	logger.WithField("ecsIDs", instanceIDs).Debug("Deleting ECS instances")
	request := ecs.CreateDeleteInstancesRequest()
	request.InstanceId = &instanceIDs
	request.Force = "true"
	_, err = o.ecsClient.DeleteInstances(request)
	if err != nil {
		return err
	}

	err = wait.Poll(
		5*time.Second,
		5*time.Minute,
		func() (bool, error) {
			response, err := o.listEcsInstance(instanceIDs)
			if err != nil {
				return false, err
			}
			if response.TotalCount == 0 {
				return true, nil
			}
			return false, nil
		},
	)

	logger.Info("ECS instances deleted")
	return
}

func (o *ClusterUninstaller) findResourcesByTag() (tagResources []tag.TagResource, err error) {
	for _, tags := range o.Tags {
		resources, err := o.listTagResources(tags)
		if err != nil {
			return nil, err
		}
		tagResources = append(tagResources, resources...)
	}
	return tagResources, nil
}

func (o *ClusterUninstaller) listTagResources(tags map[string]string) (tagResources []tag.TagResource, err error) {
	tagsString, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}

	o.Logger.WithField("tags", string(tagsString)).Debug("Retrieving cloud resources")

	request := tag.CreateListTagResourcesRequest()
	request.PageSize = "1000"
	request.Tags = string(tagsString)
	request.Category = "Custom"
	response, err := o.tagClient.ListTagResources(request)
	if err != nil {
		return nil, err
	}
	return response.TagResources, nil
}

func (o *ClusterUninstaller) unTagResource(keys *[]string, arns *[]string) (err error) {
	request := tag.CreateUntagResourcesRequest()
	request.TagKey = keys
	request.ResourceARN = arns
	_, err = o.tagClient.UntagResources(request)
	return
}

func (o *ClusterUninstaller) deleteRAMRoles(logger logrus.FieldLogger) (err error) {
	roles := []string{"bootstrap", "master", "worker"}

	for _, role := range roles {
		roleName := fmt.Sprintf("%s-role-%s", o.InfraID, role)
		policyName := fmt.Sprintf("%s-policy-%s", o.InfraID, role)

		err = o.detachRAMPolicy(policyName, logger)
		if err != nil {
			return err
		}
		err = o.deletePolicyByName(policyName, logger)
		if err != nil {
			return err
		}
		err = o.deleteRAMRole(roleName, logger)
		if err != nil && !strings.Contains(err.Error(), "EntityNotExist.Role") {
			return err
		}
	}

	logger.Info("RAM roles deleted")
	return nil
}

func (o *ClusterUninstaller) deleteRAMRole(roleName string, logger logrus.FieldLogger) (err error) {
	logger.WithField("roleName", roleName).Debugf("Deleting")
	request := ram.CreateDeleteRoleRequest()
	request.Scheme = "https"
	request.RoleName = roleName
	_, err = o.ramClient.DeleteRole(request)
	return
}

func (o *ClusterUninstaller) deletePolicyByName(policyName string, logger logrus.FieldLogger) (err error) {
	err = o.deletePolicy(policyName, logger)
	if err != nil && !strings.Contains(err.Error(), "EntityNotExist.Policy") {
		return err
	}

	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			_, err := o.getPolicy(policyName)
			if err != nil {
				if strings.Contains(err.Error(), "EntityNotExist.Policy") {
					return true, nil
				}
				return false, err
			}
			return false, nil
		},
	)
	return
}

func (o *ClusterUninstaller) getPolicy(policyName string) (response *ram.GetPolicyResponse, err error) {
	request := ram.CreateGetPolicyRequest()
	request.Scheme = "https"
	request.PolicyName = policyName
	request.PolicyType = "Custom"
	response, err = o.ramClient.GetPolicy(request)
	return
}

func (o *ClusterUninstaller) deletePolicy(policyName string, logger logrus.FieldLogger) (err error) {
	logger.WithField("policyName", policyName).Debug("Deleting")
	request := ram.CreateDeletePolicyRequest()
	request.Scheme = "https"
	request.PolicyName = policyName
	_, err = o.ramClient.DeletePolicy(request)
	return
}

func (o *ClusterUninstaller) detachRAMPolicy(policyName string, logger logrus.FieldLogger) (err error) {
	logger.WithField("policyName", policyName).Debug("Searching RAM policy")
	attachmentsResponse, err := o.listPolicyAttachments(policyName)
	if err != nil {
		return err
	}
	if attachmentsResponse.TotalCount == 0 {
		return nil
	}

	for _, a := range attachmentsResponse.PolicyAttachments.PolicyAttachment {
		err = o.detachPolicy(a.PolicyName, a.PolicyType, a.PrincipalName, a.PrincipalType, a.ResourceGroupId, logger)
		if err != nil {
			return err
		}
	}

	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			attachmentsResponse, err = o.listPolicyAttachments(policyName)
			if err != nil {
				return false, err
			}
			if attachmentsResponse.TotalCount == 0 {
				return true, nil
			}
			return false, nil
		},
	)
	logger.WithField("policyName", policyName).Debug("Policy detached")
	return
}

func (o *ClusterUninstaller) listPolicyAttachments(policyName string) (response *resourcemanager.ListPolicyAttachmentsResponse, err error) {
	request := resourcemanager.CreateListPolicyAttachmentsRequest()
	request.Scheme = "https"
	request.PolicyName = policyName
	response, err = o.rmanagerClient.ListPolicyAttachments(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (o *ClusterUninstaller) detachPolicy(policyName string, policyType string, principalName string, principalType string, resourceGroupID string, logger logrus.FieldLogger) (err error) {
	logger.WithFields(logrus.Fields{"policyName": policyName, "principalName": principalName}).Debug("Detaching policy for RAM role")
	request := resourcemanager.CreateDetachPolicyRequest()
	request.Scheme = "https"
	request.PolicyName = policyName
	request.PolicyType = policyType
	request.PrincipalName = principalName
	request.PrincipalType = principalType
	request.ResourceGroupId = resourceGroupID
	_, err = o.rmanagerClient.DetachPolicy(request)
	return
}

func (o *ClusterUninstaller) deletePrivateZones(logger logrus.FieldLogger) (err error) {
	clusterDomain := o.ClusterDomain
	logger.WithField("clusterDomain", clusterDomain).Debug("Searching private zone")
	zoneID, err := o.getPrivateZoneID()
	if err != nil {
		return err
	}
	if zoneID == "" {
		return nil
	}

	err = o.bindZoneVpc(zoneID, logger)
	if err != nil {
		return err
	}

	// Wait for unbind vpc to complete
	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			zones, err := o.listPrivateZone(clusterDomain)
			if err != nil {
				return false, err
			}

			if len(zones[0].Vpcs.Vpc) == 0 {
				return true, nil
			}
			return false, nil
		},
	)
	if err != nil {
		return
	}

	// Delete a private zone does not require delete the record in advance
	err = o.deletePrivateZone(zoneID, logger)
	if err != nil {
		return err
	}

	// Wait for deletion private zone to complete
	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			zones, err := o.listPrivateZone(clusterDomain)
			if err != nil {
				return false, err
			}

			if len(zones) == 0 {
				return true, nil
			}
			return false, nil
		},
	)
	if err != nil {
		return err
	}

	logger.Info("Private zones deleted")
	return nil
}

func (o *ClusterUninstaller) getPrivateZoneID() (zoneID string, err error) {
	clusterDomain := o.ClusterDomain
	zones, err := o.listPrivateZone(clusterDomain)
	if err != nil {
		return "", err
	}
	if len(zones) == 0 {
		return "", nil
	}
	if len(zones) > 1 {
		return "", errors.Wrap(err, fmt.Sprintf("matched to multiple private zones by clusterdomain %q", clusterDomain))
	}

	return zones[0].ZoneId, nil
}

func (o *ClusterUninstaller) deletePrivateZone(zoneID string, logger logrus.FieldLogger) (err error) {
	logger.WithField("zoneID", zoneID).Debug("Deleting private zone")
	request := pvtz.CreateDeleteZoneRequest()
	request.ZoneId = zoneID
	_, err = o.pvtzClient.DeleteZone(request)
	return
}

func (o *ClusterUninstaller) bindZoneVpc(zoneID string, logger logrus.FieldLogger) (err error) {
	logger.WithField("zoneID", zoneID).Debug("Unbinding private zone with vpc")
	request := pvtz.CreateBindZoneVpcRequest()
	request.ZoneId = zoneID
	_, err = o.pvtzClient.BindZoneVpc(request)
	return
}

func (o *ClusterUninstaller) listPrivateZone(clusterDomain string) ([]pvtz.Zone, error) {
	request := pvtz.CreateDescribeZonesRequest()
	request.Lang = "en"
	request.Keyword = clusterDomain
	request.SearchMode = "EXACT"

	response, err := o.pvtzClient.DescribeZones(request)
	if err != nil {
		return nil, err
	}
	return response.Zones.Zone, nil
}

func (o *ClusterUninstaller) listPrivateZoneRecords(zoneID string) ([]pvtz.Record, error) {
	request := pvtz.CreateDescribeZoneRecordsRequest()
	request.Lang = "en"
	request.ZoneId = zoneID

	response, err := o.pvtzClient.DescribeZoneRecords(request)
	if err != nil {
		return nil, err
	}
	return response.Records.Record, nil
}

func (o *ClusterUninstaller) deleteDNSRecords(logger logrus.FieldLogger) (err error) {
	logger.Debug("Searching DNS records")

	// Get the base domain from the cluster domain. the format of cluster domain is '<cluster name>.<base domain>'.
	domainParts := strings.Split(o.ClusterDomain, ".")
	if len(domainParts) < 2 {
		return errors.New("could not determine cluster name from cluster domain")
	}
	clusterName := domainParts[0]
	baseDomain := strings.Join(domainParts[1:], ".")

	domains, err := o.listDomain(baseDomain)
	if err != nil {
		return
	}
	if len(domains) == 0 {
		return
	}

	recordSetKey := func(recordType string, rr string) string {
		return fmt.Sprintf("%s %s", recordType, rr)
	}

	// Get the parsing record of privatezone and delete the record in publiczone
	// When the user manually deletes the private zone and records, it may cause the public records to be leaked.
	privateZoneID, err := o.getPrivateZoneID()
	if err != nil {
		return
	}
	if privateZoneID == "" {
		o.Logger.Info("The private zone ID is not obtained, and the DNS public records cannot be deleted")
		return nil
	}

	privateRecords := map[string]bool{}
	records, err := o.listPrivateZoneRecords(privateZoneID)
	if err != nil {
		return
	}
	for _, record := range records {
		key := recordSetKey(record.Type, fmt.Sprintf("%s.%s", record.Rr, clusterName))
		privateRecords[key] = true
	}

	publicRecords, err := o.listRecord(baseDomain)
	if err != nil {
		return
	}
	if len(publicRecords) == 0 {
		return
	}

	var lastErr error
	for _, record := range publicRecords {
		recordLogger := logger.WithFields(logrus.Fields{"recordID": record.RecordId, "domain": baseDomain, "rr": record.RR})
		key := recordSetKey(record.Type, record.RR)
		if privateRecords[key] {
			err = o.deleteRecord(record.RecordId, recordLogger)
			if err != nil {
				privateRecords[key] = false
				lastErr = errors.Wrap(err, fmt.Sprintf("DNS record %q", record.RecordId))
				o.Logger.Info(lastErr)
			}
		}
	}

	// Wait for deletion to complete
	err = wait.Poll(
		1*time.Second,
		1*time.Minute,
		func() (bool, error) {
			records, err := o.listRecord(baseDomain)
			if err != nil {
				return false, err
			}
			for _, record := range records {
				key := recordSetKey(record.Type, record.RR)
				if privateRecords[key] {
					return false, nil
				}
			}
			return true, nil
		},
	)
	if err != nil {
		return err
	}

	logger.Debug("Public DNS records deleted")
	return lastErr
}

func (o *ClusterUninstaller) deleteRecord(recordID string, logger logrus.FieldLogger) error {
	logger.Debug("Deleting")
	request := alidns.CreateDeleteDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = recordID
	_, err := o.dnsClient.DeleteDomainRecord(request)
	if err != nil {
		return err
	}
	return nil
}

func (o *ClusterUninstaller) listDomain(baseDomain string) ([]alidns.DomainInDescribeDomains, error) {
	request := alidns.CreateDescribeDomainsRequest()
	request.Scheme = "https"
	request.KeyWord = baseDomain
	request.SearchMode = "EXACT"
	response, err := o.dnsClient.DescribeDomains(request)
	if err != nil {
		return nil, err
	}
	return response.Domains.Domain, nil
}

func (o *ClusterUninstaller) listRecord(baseDomain string) ([]alidns.Record, error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = baseDomain
	response, err := o.dnsClient.DescribeDomainRecords(request)
	if err != nil {
		return nil, err
	}
	return response.DomainRecords.Record, nil
}

func convertResourceArn(arn string) (resourceArn ResourceArn, err error) {
	_arn := strings.Split(arn, "/")
	serviceInfos := strings.Split(_arn[0], ":")

	resourceArn.Service = serviceInfos[2]
	resourceArn.Region = serviceInfos[3]
	resourceArn.Account = serviceInfos[4]
	resourceArn.ResourceType = serviceInfos[5]
	resourceArn.ResourceID = _arn[1]
	resourceArn.Arn = arn
	return resourceArn, nil
}
