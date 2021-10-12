package alibabacloud

import (
	"encoding/json"
	"fmt"
	"strings"
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

	Region          string
	InfraID         string
	ClusterID       string
	ClusterDomain   string
	ResourceGroupID string
	TagKey          string
	TagValue        string
	TagResources    struct {
		ecsInstances   []ResourceArn
		securityGroups []ResourceArn
		vpcs           []ResourceArn
		vSwitchs       []ResourceArn
		eips           []ResourceArn
		natgateways    []ResourceArn
		slbs           []ResourceArn
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
		ClusterID:       metadata.ClusterID,
		ClusterDomain:   metadata.AlibabaCloud.ClusterDomain,
		TagKey:          fmt.Sprintf("kubernetes.io/cluster/%s", metadata.InfraID),
		TagValue:        "owned",
	}, nil
}

func (o *ClusterUninstaller) executeDeleteFunction(execute func() error, resourceName string) (err error) {
	err = wait.PollImmediateInfinite(
		time.Second*10,
		func() (bool, error) {
			ferr := execute()
			if ferr != nil {
				o.Logger.Debugf("failed to delete %s: %v", resourceName, ferr)
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	var err error
	deletedFuncs := []struct {
		resourceName string
		executeFunc  func() error
	}{
		{resourceName: "DNS records", executeFunc: o.deleteDNSRecords},
		{resourceName: "OSS buckets", executeFunc: o.deleteBucket},
		{resourceName: "private zones", executeFunc: o.deletePrivateZones},
		{resourceName: "RAM roles", executeFunc: o.deleteRAMRoles},
		{resourceName: "ECS instances", executeFunc: o.deleteEcsInstances},
		{resourceName: "ECS security groups", executeFunc: o.deleteSecurityGroups},
		{resourceName: "Nat gateways", executeFunc: o.deleteNatGateways},
		{resourceName: "EIPs", executeFunc: o.deleteEips},
		{resourceName: "SLBs", executeFunc: o.deleteSlbs},
		{resourceName: "VSwitchs", executeFunc: o.deleteVSwitchs},
		{resourceName: "VPCs", executeFunc: o.deleteVpcs},
	}

	err = o.configureClients()
	if err != nil {
		return nil, err
	}

	err = o.findResources()
	if err != nil {
		return nil, err
	}

	// TODO: more appropriate to use asynchronous. It is advisable to optimise in the future
	for _, execute := range deletedFuncs {
		err = o.executeDeleteFunction(execute.executeFunc, execute.resourceName)
		if err != nil {
			return nil, err
		}
	}

	err = o.waitComplete()
	if err != nil {
		return nil, err
	}

	return nil, nil
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
				return err
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
		default:
			o.TagResources.others = append(o.TagResources.others, resourceArn)
		}
	}
}

func (o *ClusterUninstaller) waitComplete() (err error) {
	var isSuccessful bool
	err = wait.Poll(
		2*time.Second,
		30*time.Second,
		func() (bool, error) {
			tagResources, err := o.findResourcesByTag()
			if err != nil {
				return false, err
			}
			if len(tagResources) == 0 {
				isSuccessful = true
				return true, nil
			}
			return false, nil
		},
	)

	if !isSuccessful {
		notDeletedResources := []string{}
		tagResources, err := o.findResourcesByTag()
		if err != nil {
			return err
		}
		for _, arn := range tagResources {
			notDeletedResources = append(notDeletedResources, arn.ResourceARN)
		}
		return errors.New(fmt.Sprintf("There are undeleted cloud resources %q", notDeletedResources))
	}
	return
}

func (o *ClusterUninstaller) deleteBucket() (err error) {
	bucketName := fmt.Sprintf("%s-bootstrap", o.InfraID)
	result, err := o.ossClient.ListBuckets(oss.Prefix(bucketName))
	if err != nil || len(result.Buckets) == 0 {
		return
	}

	o.Logger.Debugf("Start to delete buckets %q", bucketName)

	keys := []string{o.TagKey}
	arns := []string{fmt.Sprintf("arn:acs:oss:%s:*:bucket/%s", o.Region, bucketName)}
	err = o.unTagResource(&keys, &arns)
	if err != nil {
		return err
	}

	bucket, err := o.ossClient.Bucket(bucketName)
	if err != nil {
		return err
	}

	err = o.deleteObjects(bucket)
	if err != nil {
		return err
	}

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

	return
}

func (o *ClusterUninstaller) deleteObjects(bucket *oss.Bucket) (err error) {
	result, err := bucket.ListObjectsV2()
	if err != nil {
		return err
	}
	if len(result.Objects) == 0 {
		return
	}

	o.Logger.Debugf("Start to delete objects of buckets %s", bucket.BucketName)
	var keys []string

	for _, object := range result.Objects {
		keys = append(keys, object.Key)
	}
	bucket.DeleteObjects(keys)

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
	return
}

func (o *ClusterUninstaller) deleteSlbs() (err error) {
	if len(o.TagResources.slbs) <= 0 {
		return nil
	}

	var slbIDs []string
	for _, slbArn := range o.TagResources.slbs {
		slbIDs = append(slbIDs, slbArn.ResourceID)
	}

	o.Logger.Debugf("Start to delete SLBs %q", slbIDs)
	for _, slbID := range slbIDs {
		err = o.deleteSlb(slbID)
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
	return
}

func (o *ClusterUninstaller) listSlb(slbIDs []string) (response *slb.DescribeLoadBalancersResponse, err error) {
	request := slb.CreateDescribeLoadBalancersRequest()
	request.LoadBalancerId = strings.Join(slbIDs, ",")
	response, err = o.slbClient.DescribeLoadBalancers(request)
	return
}

func (o *ClusterUninstaller) deleteSlb(slbID string) (err error) {
	o.Logger.Debugf("Start to delete SLB %q", slbID)
	request := slb.CreateDeleteLoadBalancerRequest()
	request.LoadBalancerId = slbID
	_, err = o.slbClient.DeleteLoadBalancer(request)
	return
}

func (o *ClusterUninstaller) deleteVSwitchs() (err error) {
	if len(o.TagResources.vSwitchs) <= 0 {
		return nil
	}

	var vSwitchIDs []string
	for _, vSwitchArn := range o.TagResources.vSwitchs {
		vSwitchIDs = append(vSwitchIDs, vSwitchArn.ResourceID)
	}

	o.Logger.Debugf("Start to delete VSwitchs %q", vSwitchIDs)
	for _, vSwitchID := range vSwitchIDs {
		err = wait.Poll(
			5*time.Second,
			30*time.Second,
			func() (bool, error) {
				err = o.deleteVSwitch(vSwitchID)
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
	return
}

func (o *ClusterUninstaller) listVSwitch(vSwitchIDs []string) (response *vpc.DescribeVSwitchesResponse, err error) {
	request := vpc.CreateDescribeVSwitchesRequest()
	request.VSwitchId = strings.Join(vSwitchIDs, ",")
	response, err = o.vpcClient.DescribeVSwitches(request)
	return
}

func (o *ClusterUninstaller) deleteVSwitch(vSwitchID string) (err error) {
	o.Logger.Debugf("Start to delete VSwitch %q", vSwitchID)
	request := vpc.CreateDeleteVSwitchRequest()
	request.VSwitchId = vSwitchID
	_, err = o.vpcClient.DeleteVSwitch(request)
	return
}

func (o *ClusterUninstaller) deleteVpcs() (err error) {
	if len(o.TagResources.vpcs) <= 0 {
		return nil
	}

	var vpcIDs []string
	for _, vpcArn := range o.TagResources.vpcs {
		vpcIDs = append(vpcIDs, vpcArn.ResourceID)
	}

	o.Logger.Debugf("Start to delete VPCs %q", vpcIDs)
	for _, vpcID := range vpcIDs {
		err = wait.Poll(
			5*time.Second,
			30*time.Second,
			func() (bool, error) {
				err = o.deleteVpc(vpcID)
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

	return
}

func (o *ClusterUninstaller) deleteVpc(vpcID string) (err error) {
	o.Logger.Debugf("Start to delete VPC %q", vpcID)
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

func (o *ClusterUninstaller) deleteEips() (err error) {
	if len(o.TagResources.eips) <= 0 {
		return nil
	}

	var eipIDs []string
	for _, eipArn := range o.TagResources.eips {
		eipIDs = append(eipIDs, eipArn.ResourceID)
	}

	o.Logger.Debugf("Start to delete EIPs %q", eipIDs)
	for _, eipID := range eipIDs {
		err = o.deleteEip(eipID)
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
	return err
}

func (o *ClusterUninstaller) listEip(eipIDs []string) (response *vpc.DescribeEipAddressesResponse, err error) {
	request := vpc.CreateDescribeEipAddressesRequest()
	request.AllocationId = strings.Join(eipIDs, ",")
	response, err = o.vpcClient.DescribeEipAddresses(request)
	return response, err
}

func (o *ClusterUninstaller) deleteEip(eipID string) (err error) {
	o.Logger.Debugf("Start to delete EIP %q", eipID)
	request := vpc.CreateReleaseEipAddressRequest()
	request.AllocationId = eipID
	_, err = o.vpcClient.ReleaseEipAddress(request)
	return
}

func (o *ClusterUninstaller) deleteNatGateways() (err error) {
	if len(o.TagResources.natgateways) <= 0 {
		return nil
	}

	var natGatewayIDs []string
	for _, natGatewayArn := range o.TagResources.natgateways {
		natGatewayIDs = append(natGatewayIDs, natGatewayArn.ResourceID)
	}

	o.Logger.Debugf("Start to delete NAT gateways %q", natGatewayIDs)
	for _, natGatewayID := range natGatewayIDs {
		err = o.deleteNatGateway(natGatewayID)
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
	}
	return
}

func (o *ClusterUninstaller) listNatGateways(natGatewayID string) (response *vpc.DescribeNatGatewaysResponse, err error) {
	request := vpc.CreateDescribeNatGatewaysRequest()
	request.NatGatewayId = natGatewayID
	response, err = o.vpcClient.DescribeNatGateways(request)
	return
}

func (o *ClusterUninstaller) deleteNatGateway(natGatewayID string) (err error) {
	o.Logger.Debugf("Start to delete NAT gateway %q", natGatewayID)
	request := vpc.CreateDeleteNatGatewayRequest()
	request.NatGatewayId = natGatewayID
	request.Force = "true"
	_, err = o.vpcClient.DeleteNatGateway(request)
	return
}

func (o *ClusterUninstaller) deleteSecurityGroups() (err error) {
	if len(o.TagResources.securityGroups) <= 0 {
		return nil
	}

	var securityGroupIDs []string
	for _, securityGroupArn := range o.TagResources.securityGroups {
		securityGroupIDs = append(securityGroupIDs, securityGroupArn.ResourceID)
	}

	o.Logger.Debugf("Start to delete security groups %q", securityGroupIDs)

	for _, securityGroupID := range securityGroupIDs {
		err = o.deleteSecurityGroupRules(securityGroupID)
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

	for _, securityGroupID := range securityGroupIDs {
		err = wait.Poll(
			5*time.Second,
			30*time.Second,
			func() (bool, error) {
				err = o.deleteSecurityGroup(securityGroupID)
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

	return
}

func (o *ClusterUninstaller) deleteSecurityGroup(securityGroupID string) (err error) {
	o.Logger.Debugf("Start to delete security group %q", securityGroupID)
	request := ecs.CreateDeleteSecurityGroupRequest()
	request.SecurityGroupId = securityGroupID
	_, err = o.ecsClient.DeleteSecurityGroup(request)
	return
}

func (o *ClusterUninstaller) deleteSecurityGroupRules(securityGroupID string) (err error) {
	o.Logger.Debugf("Start to delete security group %q rules ", securityGroupID)
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

func (o *ClusterUninstaller) deleteEcsInstances() (err error) {
	if len(o.TagResources.ecsInstances) <= 0 {
		return nil
	}

	var instanceIDs []string
	for _, instanceArn := range o.TagResources.ecsInstances {
		instanceIDs = append(instanceIDs, instanceArn.ResourceID)
	}

	o.Logger.Debugf("Start to delete ECS instances %q", instanceIDs)

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
	return
}

func (o *ClusterUninstaller) findResourcesByTag() (tagResources []tag.TagResource, err error) {
	tags := map[string]string{o.TagKey: o.TagValue}
	tagsString, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}

	o.Logger.Debugf("Retrieving cloud resources by tag %s", tagsString)

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
	o.Logger.Debugf("Untag cloud resources %q with tags %q", arns, keys)
	request := tag.CreateUntagResourcesRequest()
	request.TagKey = keys
	request.ResourceARN = arns
	_, err = o.tagClient.UntagResources(request)
	return
}

func (o *ClusterUninstaller) deleteRAMRoles() (err error) {
	roles := []string{"bootstrap", "master", "worker"}

	for _, role := range roles {
		roleName := fmt.Sprintf("%s-role-%s", o.InfraID, role)
		policyName := fmt.Sprintf("%s-policy-%s", o.InfraID, role)

		err = o.detachRAMPolicy(policyName)
		if err != nil {
			return err
		}
		err = o.deletePolicyByName(policyName)
		if err != nil {
			return err
		}
		err = o.deleteRAMRole(roleName)
		if err != nil && !strings.Contains(err.Error(), "EntityNotExist.Role") {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) deleteRAMRole(roleName string) (err error) {
	o.Logger.Debugf("Start to search and delete RAM role %q", roleName)
	request := ram.CreateDeleteRoleRequest()
	request.Scheme = "https"
	request.RoleName = roleName
	_, err = o.ramClient.DeleteRole(request)
	return
}

func (o *ClusterUninstaller) deletePolicyByName(policyName string) (err error) {
	err = o.deletePolicy(policyName)
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

func (o *ClusterUninstaller) deletePolicy(policyName string) (err error) {
	o.Logger.Debugf("Start to search and delete RAM policy %q", policyName)
	request := ram.CreateDeletePolicyRequest()
	request.Scheme = "https"
	request.PolicyName = policyName
	_, err = o.ramClient.DeletePolicy(request)
	return
}

func (o *ClusterUninstaller) detachRAMPolicy(policyName string) (err error) {
	o.Logger.Debugf("Start to search RAM policy %q attachments", policyName)
	attachmentsResponse, err := o.listPolicyAttachments(policyName)
	if err != nil {
		return err
	}
	if attachmentsResponse.TotalCount == 0 {
		return nil
	}

	o.Logger.Debugf("Start to detach RAM policy %q", policyName)
	for _, a := range attachmentsResponse.PolicyAttachments.PolicyAttachment {
		err = o.detachPolicy(a.PolicyName, a.PolicyType, a.PrincipalName, a.PrincipalType, a.ResourceGroupId)
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

func (o *ClusterUninstaller) detachPolicy(policyName string, policyType string, principalName string, principalType string, resourceGroupID string) (err error) {
	o.Logger.Debugf("Start to detach RAM policy %q with %q", policyName, principalName)
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

func (o *ClusterUninstaller) deletePrivateZones() (err error) {
	clusterDomain := o.ClusterDomain
	o.Logger.Debug("Start to search private zones")
	zones, err := o.listPrivateZone(clusterDomain)
	if err != nil {
		return err
	}
	if len(zones) == 0 {
		return nil
	}
	if len(zones) > 1 {
		return errors.Wrap(err, fmt.Sprintf("matched to multiple private zones by clustedomain %q", clusterDomain))
	}

	zoneID := zones[0].ZoneId
	err = o.bindZoneVpc(zoneID)
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

	o.Logger.Debug("Start to delete private zones")
	// Delete a private zone does not require delete the record in advance
	err = o.deletePrivateZone(zoneID)
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

	return nil
}

func (o *ClusterUninstaller) deletePrivateZone(zoneID string) (err error) {
	o.Logger.Debugf("Start to delete private zone %q", zoneID)
	request := pvtz.CreateDeleteZoneRequest()
	request.ZoneId = zoneID
	_, err = o.pvtzClient.DeleteZone(request)
	return
}

func (o *ClusterUninstaller) bindZoneVpc(zoneID string) (err error) {
	o.Logger.Debugf("Start to unbind/bind private zone %q with vpc", zoneID)
	request := pvtz.CreateBindZoneVpcRequest()
	request.ZoneId = zoneID
	_, err = o.pvtzClient.BindZoneVpc(request)
	return
}

func (o *ClusterUninstaller) listPrivateZone(clusterDomain string) ([]pvtz.Zone, error) {
	request := pvtz.CreateDescribeZonesRequest()
	request.Lang = "en"
	request.Keyword = clusterDomain

	response, err := o.pvtzClient.DescribeZones(request)
	if err != nil {
		return nil, err
	}
	return response.Zones.Zone, nil
}

func (o *ClusterUninstaller) deleteDNSRecords() (err error) {
	o.Logger.Debug("Start to search DNS records")

	baseDomain := strings.Join(strings.Split(o.ClusterDomain, ".")[1:], ".")
	domains, err := o.listDomain(baseDomain)
	if err != nil {
		return
	}
	if len(domains) == 0 {
		return
	}

	records, err := o.listRecord(baseDomain)
	if err != nil {
		return
	}
	if len(records) == 0 {
		return
	}

	o.Logger.Debug("Start to delete DNS records")
	for _, record := range records {
		err = o.deleteRecord(record.RecordId)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("DNS record %q", record.RecordId))
			o.Logger.Info(err)
			return
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

			if len(records) == 0 {
				return true, nil
			}
			return false, nil
		},
	)
	if err != nil {
		return
	}

	return nil
}

func (o *ClusterUninstaller) deleteRecord(recordID string) error {
	o.Logger.Debugf("Start to delete DNS record %q", recordID)
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
