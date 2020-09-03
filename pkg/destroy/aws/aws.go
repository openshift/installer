package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	awssession "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version"
)

var (
	exists = struct{}{}
)

// Filter holds the key/value pairs for the tags we will be matching against.
//
// A resource matches the filter if all of the key/value pairs are in its tags.
type Filter map[string]string

// ClusterUninstaller holds the various options for the cluster we want to delete
type ClusterUninstaller struct {

	// Filters is a slice of filters for matching resources.  A
	// resources matches the whole slice if it matches any of the
	// entries.  For example:
	//
	//   filter := []map[string]string{
	//     {
	//       "a": "b",
	//       "c": "d:,
	//     },
	//     {
	//       "d": "e",
	//     },
	//   }
	//
	// will match resources with (a:b and c:d) or d:e.
	Filters   []Filter // filter(s) we will be searching for
	Logger    logrus.FieldLogger
	Region    string
	ClusterID string

	// Session is the AWS session to be used for deletion.  If nil, a
	// new session will be created based on the usual credential
	// configuration (AWS_PROFILE, AWS_ACCESS_KEY_ID, etc.).
	Session *session.Session
}

// New returns an AWS destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	filters := make([]Filter, 0, len(metadata.ClusterPlatformMetadata.AWS.Identifier))
	for _, filter := range metadata.ClusterPlatformMetadata.AWS.Identifier {
		filters = append(filters, filter)
	}
	region := metadata.ClusterPlatformMetadata.AWS.Region
	session, err := awssession.GetSessionWithOptions(
		awssession.WithRegion(region),
		awssession.WithServiceEndpoints(region, metadata.ClusterPlatformMetadata.AWS.ServiceEndpoints),
	)
	if err != nil {
		return nil, err
	}

	return &ClusterUninstaller{
		Filters:   filters,
		Region:    region,
		Logger:    logger,
		ClusterID: metadata.InfraID,
		Session:   session,
	}, nil
}

func (o *ClusterUninstaller) validate() error {
	if len(o.Filters) == 0 {
		return errors.Errorf("you must specify at least one tag filter")
	}
	return nil
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() error {
	return o.RunWithContext(context.Background())
}

// RunWithContext runs the uninstall process with a context.
// The first return is the list of ARNs for resources that could not be destroyed.
func (o *ClusterUninstaller) RunWithContext(ctx context.Context) error {
	err := o.validate()
	if err != nil {
		return err
	}

	awsSession := o.Session
	if awsSession == nil {
		// Relying on appropriate AWS ENV vars (eg AWS_PROFILE, AWS_ACCESS_KEY_ID, etc)
		awsSession, err = session.NewSession(aws.NewConfig().WithRegion(o.Region))
		if err != nil {
			return err
		}
	}
	awsSession.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "openshiftInstaller.OpenshiftInstallerUserAgentHandler",
		Fn:   request.MakeAddToUserAgentHandler("OpenShift/4.x Destroyer", version.Raw),
	})

	tagClients := []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI{
		resourcegroupstaggingapi.New(awsSession),
	}
	tagClientNames := map[*resourcegroupstaggingapi.ResourceGroupsTaggingAPI]string{
		tagClients[0]: o.Region,
	}

	switch o.Region {
	case endpoints.CnNorth1RegionID, endpoints.CnNorthwest1RegionID:
		if o.Region != endpoints.CnNorthwest1RegionID {
			tagClient := resourcegroupstaggingapi.New(awsSession, aws.NewConfig().WithRegion(endpoints.CnNorthwest1RegionID))
			tagClients = append(tagClients, tagClient)
			tagClientNames[tagClient] = endpoints.CnNorthwest1RegionID
		}
	case endpoints.UsGovEast1RegionID, endpoints.UsGovWest1RegionID:
		if o.Region != endpoints.UsGovWest1RegionID {
			tagClient := resourcegroupstaggingapi.New(awsSession, aws.NewConfig().WithRegion(endpoints.UsGovWest1RegionID))
			tagClients = append(tagClients, tagClient)
			tagClientNames[tagClient] = endpoints.UsGovWest1RegionID
		}
	default:
		if o.Region != endpoints.UsEast1RegionID {
			tagClient := resourcegroupstaggingapi.New(awsSession, aws.NewConfig().WithRegion(endpoints.UsEast1RegionID))
			tagClients = append(tagClients, tagClient)
			tagClientNames[tagClient] = endpoints.UsEast1RegionID
		}
	}

	iamClient := iam.New(awsSession)
	iamRoleSearch := &iamRoleSearch{
		client:  iamClient,
		filters: o.Filters,
		logger:  o.Logger,
	}
	iamUserSearch := &iamUserSearch{
		client:  iamClient,
		filters: o.Filters,
		logger:  o.Logger,
	}

	deleted, err := terminateEC2InstancesByTags(ctx, ec2.New(awsSession), iamClient, o.Filters, o.Logger)
	if err != nil {
		return err
	}

	tracker := new(errorTracker)
	tagClientStack := append([]*resourcegroupstaggingapi.ResourceGroupsTaggingAPI(nil), tagClients...)
	err = wait.PollImmediateUntil(
		time.Second*10,
		func() (done bool, err error) {
			var loopError error
			nextTagClients := tagClients[:0]
			for _, tagClient := range tagClientStack {
				matched := false
				for _, filter := range o.Filters {
					o.Logger.Debugf("search for and delete matching resources by tag in %s matching %#+v", tagClientNames[tagClient], filter)
					tagFilters := make([]*resourcegroupstaggingapi.TagFilter, 0, len(filter))
					for key, value := range filter {
						tagFilters = append(tagFilters, &resourcegroupstaggingapi.TagFilter{
							Key:    aws.String(key),
							Values: []*string{aws.String(value)},
						})
					}
					err = tagClient.GetResourcesPagesWithContext(
						ctx,
						&resourcegroupstaggingapi.GetResourcesInput{TagFilters: tagFilters},
						func(results *resourcegroupstaggingapi.GetResourcesOutput, lastPage bool) bool {
							for _, resource := range results.ResourceTagMappingList {
								arnString := *resource.ResourceARN
								if _, ok := deleted[arnString]; !ok {
									arnLogger := o.Logger.WithField("arn", arnString)
									matched = true
									parsed, err := arn.Parse(arnString)
									if err != nil {
										arnLogger.Debug(err)
										continue
									}

									err = deleteARN(ctx, awsSession, parsed, filter, arnLogger)
									if err != nil {
										tracker.suppressWarning(arnString, err, arnLogger)
										err = errors.Wrapf(err, "deleting %s", arnString)
										continue
									}
									deleted[arnString] = exists
								}
							}

							return !lastPage
						},
					)
					if err != nil {
						err = errors.Wrap(err, "get tagged resources")
						o.Logger.Info(err)
						matched = true
						loopError = err
					}
				}

				if matched {
					nextTagClients = append(nextTagClients, tagClient)
				} else {
					o.Logger.Debugf("no deletions from %s, removing client", tagClientNames[tagClient])
				}
			}
			tagClientStack = nextTagClients

			o.Logger.Debug("search for IAM roles")
			arns, err := iamRoleSearch.arns(ctx)
			if err != nil {
				o.Logger.Info(err)
				loopError = err
			}

			o.Logger.Debug("search for IAM users")
			userARNs, err := iamUserSearch.arns(ctx)
			if err != nil {
				o.Logger.Info(err)
				loopError = err
			}
			arns = append(arns, userARNs...)

			if len(arns) > 0 {
				o.Logger.Debug("delete IAM roles and users")
			}
			for _, arnString := range arns {
				if _, ok := deleted[arnString]; !ok {
					arnLogger := o.Logger.WithField("arn", arnString)
					parsed, err := arn.Parse(arnString)
					if err != nil {
						arnLogger.Debug(err)
						loopError = err
						continue
					}

					err = deleteARN(ctx, awsSession, parsed, nil, arnLogger)
					if err != nil {
						tracker.suppressWarning(arnString, err, arnLogger)
						err = errors.Wrapf(err, "deleting %s", arnString)
						loopError = err
						continue
					}
					deleted[arnString] = exists
				}
			}

			return len(tagClientStack) == 0 && loopError == nil, nil
		},
		ctx.Done(),
	)
	if err != nil {
		return err
	}

	o.Logger.Debug("search for untaggable resources")
	if err := o.deleteUntaggedResources(ctx, awsSession); err != nil {
		o.Logger.Debug(err)
		return err
	}

	err = removeSharedTags(ctx, tagClients, tagClientNames, o.Filters, o.Logger)
	if err != nil {
		return err
	}

	return nil
}

func splitSlash(name string, input string) (base string, suffix string, err error) {
	segments := strings.SplitN(input, "/", 2)
	if len(segments) != 2 {
		return "", "", errors.Errorf("%s %q does not contain the expected slash", name, input)
	}
	return segments[0], segments[1], nil
}

func tagMatch(filters []Filter, tags map[string]string) bool {
	for _, filter := range filters {
		match := true
		for filterKey, filterValue := range filter {
			tagValue, ok := tags[filterKey]
			if !ok {
				match = false
				break
			}
			if tagValue != filterValue {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return len(filters) == 0
}

func tagsForFilter(filter Filter) []*ec2.Tag {
	tags := make([]*ec2.Tag, 0, len(filter))
	for key, value := range filter {
		tags = append(tags, &ec2.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}
	return tags
}

type iamRoleSearch struct {
	client    *iam.IAM
	filters   []Filter
	logger    logrus.FieldLogger
	unmatched map[string]struct{}
}

func (search *iamRoleSearch) arns(ctx context.Context) ([]string, error) {
	if search.unmatched == nil {
		search.unmatched = map[string]struct{}{}
	}

	arns := []string{}
	var lastError error
	err := search.client.ListRolesPagesWithContext(
		ctx,
		&iam.ListRolesInput{},
		func(results *iam.ListRolesOutput, lastPage bool) bool {
			for _, role := range results.Roles {
				if _, ok := search.unmatched[*role.Arn]; ok {
					continue
				}

				// Unfortunately role.Tags is empty from ListRoles, so we need to query each one
				response, err := search.client.GetRoleWithContext(ctx, &iam.GetRoleInput{RoleName: role.RoleName})
				if err != nil {
					if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
						search.unmatched[*role.Arn] = exists
					} else {
						if lastError != nil {
							search.logger.Debug(lastError)
						}
						lastError = errors.Wrapf(err, "get tags for %s", *role.Arn)
					}
				} else {
					role = response.Role
					tags := make(map[string]string, len(role.Tags))
					for _, tag := range role.Tags {
						tags[*tag.Key] = *tag.Value
					}
					if tagMatch(search.filters, tags) {
						arns = append(arns, *role.Arn)
					} else {
						search.unmatched[*role.Arn] = exists
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return arns, lastError
	}
	return arns, err
}

type iamUserSearch struct {
	client    *iam.IAM
	filters   []Filter
	logger    logrus.FieldLogger
	unmatched map[string]struct{}
}

func (search *iamUserSearch) arns(ctx context.Context) ([]string, error) {
	if search.unmatched == nil {
		search.unmatched = map[string]struct{}{}
	}

	arns := []string{}
	var lastError error
	err := search.client.ListUsersPagesWithContext(
		ctx,
		&iam.ListUsersInput{},
		func(results *iam.ListUsersOutput, lastPage bool) bool {
			for _, user := range results.Users {
				if _, ok := search.unmatched[*user.Arn]; ok {
					continue
				}

				// Unfortunately user.Tags is empty from ListUsers, so we need to query each one
				response, err := search.client.GetUserWithContext(ctx, &iam.GetUserInput{UserName: aws.String(*user.UserName)})
				if err != nil {
					if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
						search.unmatched[*user.Arn] = exists
					} else {
						if lastError != nil {
							search.logger.Debug(lastError)
						}
						lastError = errors.Wrapf(err, "get tags for %s", *user.Arn)
					}
				} else {
					user = response.User
					tags := make(map[string]string, len(user.Tags))
					for _, tag := range user.Tags {
						tags[*tag.Key] = *tag.Value
					}
					if tagMatch(search.filters, tags) {
						arns = append(arns, *user.Arn)
					} else {
						search.unmatched[*user.Arn] = exists
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return arns, lastError
	}
	return arns, err
}

// getSharedHostedZone will find the ID of the non-Terraform-managed public route53 zone given the
// Terraform-managed zone's privateID.
func getSharedHostedZone(ctx context.Context, client *route53.Route53, privateID string, logger logrus.FieldLogger) (string, error) {
	response, err := client.GetHostedZoneWithContext(ctx, &route53.GetHostedZoneInput{
		Id: aws.String(privateID),
	})
	if err != nil {
		return "", err
	}

	privateName := *response.HostedZone.Name

	if response.HostedZone.Config != nil && response.HostedZone.Config.PrivateZone != nil {
		if !*response.HostedZone.Config.PrivateZone {
			return "", errors.Errorf("getSharedHostedZone requires a private ID, but was passed the public %s", privateID)
		}
	} else {
		logger.WithField("hosted zone", privateName).Warn("could not determine whether hosted zone is private")
	}

	domain := privateName
	parents := []string{domain}
	for {
		idx := strings.Index(domain, ".")
		if idx == -1 {
			break
		}
		if len(domain[idx+1:]) > 0 {
			parents = append(parents, domain[idx+1:])
		}
		domain = domain[idx+1:]
	}

	for _, p := range parents {
		sZone, err := findPublicRoute53(ctx, client, p, logger)
		if err != nil {
			return "", err
		}
		if sZone != "" {
			return sZone, nil
		}
	}
	return "", nil
}

// findPublicRoute53 finds a public route53 zone matching the dnsName.
// It returns "", when no public route53 zone could be found.
func findPublicRoute53(ctx context.Context, client *route53.Route53, dnsName string, logger logrus.FieldLogger) (string, error) {
	request := &route53.ListHostedZonesByNameInput{
		DNSName: aws.String(dnsName),
	}
	for i := 0; true; i++ {
		logger.Debugf("listing AWS hosted zones %q (page %d)", dnsName, i)
		list, err := client.ListHostedZonesByNameWithContext(ctx, request)
		if err != nil {
			return "", err
		}

		for _, zone := range list.HostedZones {
			if *zone.Name != dnsName {
				// No name after this can match dnsName
				return "", nil
			}
			if zone.Config == nil || zone.Config.PrivateZone == nil {
				logger.WithField("hosted zone", *zone.Name).Warn("could not determine whether hosted zone is private")
				continue
			}
			if !*zone.Config.PrivateZone {
				return *zone.Id, nil
			}
		}

		if *list.IsTruncated && *list.NextDNSName == *request.DNSName {
			request.HostedZoneId = list.NextHostedZoneId
			continue
		}

		break
	}
	return "", nil
}

func deleteARN(ctx context.Context, session *session.Session, arn arn.ARN, filter Filter, logger logrus.FieldLogger) error {
	switch arn.Service {
	case "ec2":
		return deleteEC2(ctx, session, arn, filter, logger)
	case "elasticloadbalancing":
		return deleteElasticLoadBalancing(ctx, session, arn, logger)
	case "iam":
		return deleteIAM(ctx, session, arn, logger)
	case "route53":
		return deleteRoute53(ctx, session, arn, logger)
	case "s3":
		return deleteS3(ctx, session, arn, logger)
	default:
		return errors.Errorf("unrecognized ARN service %s (%s)", arn.Service, arn)
	}
}

func deleteEC2(ctx context.Context, session *session.Session, arn arn.ARN, filter Filter, logger logrus.FieldLogger) error {
	client := ec2.New(session)

	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "dhcp-options":
		return deleteEC2DHCPOptions(ctx, client, id, logger)
	case "elastic-ip":
		return deleteEC2ElasticIP(ctx, client, id, logger)
	case "image":
		return deleteEC2Image(ctx, client, id, filter, logger)
	case "instance":
		return terminateEC2Instance(ctx, client, iam.New(session), id, logger)
	case "internet-gateway":
		return deleteEC2InternetGateway(ctx, client, id, logger)
	case "natgateway":
		return deleteEC2NATGateway(ctx, client, id, logger)
	case "route-table":
		return deleteEC2RouteTable(ctx, client, id, logger)
	case "security-group":
		return deleteEC2SecurityGroup(ctx, client, id, logger)
	case "snapshot":
		return deleteEC2Snapshot(ctx, client, id, logger)
	case "network-interface":
		return deleteEC2NetworkInterface(ctx, client, id, logger)
	case "subnet":
		return deleteEC2Subnet(ctx, client, id, logger)
	case "volume":
		return deleteEC2Volume(ctx, client, id, logger)
	case "vpc":
		return deleteEC2VPC(ctx, client, elb.New(session), elbv2.New(session), id, logger)
	case "vpc-endpoint":
		return deleteEC2VPCEndpoint(ctx, client, id, logger)
	case "vpc-peering-connection":
		return deleteEC2VPCPeeringConnection(ctx, client, id, logger)
	default:
		return errors.Errorf("unrecognized EC2 resource type %s", resourceType)
	}
}

func deleteEC2DHCPOptions(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteDhcpOptionsWithContext(ctx, &ec2.DeleteDhcpOptionsInput{
		DhcpOptionsId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidDhcpOptionsID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2Image(ctx context.Context, client *ec2.EC2, id string, filter Filter, logger logrus.FieldLogger) error {
	// tag the snapshots used by the AMI so that the snapshots are matched
	// by the filter and deleted
	response, err := client.DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
		ImageIds: []*string{&id},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidAMIID.NotFound" {
			return nil
		}
		return err
	}
	snapshots := []*string{}
	for _, image := range response.Images {
		for _, bdm := range image.BlockDeviceMappings {
			if bdm.Ebs != nil && bdm.Ebs.SnapshotId != nil {
				snapshots = append(snapshots, bdm.Ebs.SnapshotId)
			}
		}
	}
	_, err = client.CreateTagsWithContext(ctx, &ec2.CreateTagsInput{
		Resources: snapshots,
		Tags:      tagsForFilter(filter),
	})
	if err != nil {
		err = errors.Wrapf(err, "tagging snapshots for %s", id)
	}

	_, err = client.DeregisterImageWithContext(ctx, &ec2.DeregisterImageInput{
		ImageId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidAMIID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2ElasticIP(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.ReleaseAddressWithContext(ctx, &ec2.ReleaseAddressInput{
		AllocationId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidAllocationID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Released")
	return nil
}

func terminateEC2Instance(ctx context.Context, ec2Client *ec2.EC2, iamClient *iam.IAM, id string, logger logrus.FieldLogger) error {
	response, err := ec2Client.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidInstanceID.NotFound" {
			return nil
		}
		return err
	}

	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			err = terminateEC2InstanceByInstance(ctx, ec2Client, iamClient, instance, logger)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func terminateEC2InstanceByInstance(ctx context.Context, ec2Client *ec2.EC2, iamClient *iam.IAM, instance *ec2.Instance, logger logrus.FieldLogger) error {
	// Skip 'shutting-down' and 'terminated' instances since they take a while to get cleaned up
	if instance.State == nil || *instance.State.Name == "shutting-down" || *instance.State.Name == "terminated" {
		return nil
	}

	if instance.IamInstanceProfile != nil {
		parsed, err := arn.Parse(*instance.IamInstanceProfile.Arn)
		if err != nil {
			return errors.Wrap(err, "parse ARN for IAM instance profile")
		}

		err = deleteIAMInstanceProfile(ctx, iamClient, parsed, logger)
		if err != nil {
			return errors.Wrapf(err, "deleting %s", parsed.String())
		}
	}

	_, err := ec2Client.TerminateInstancesWithContext(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []*string{instance.InstanceId},
	})
	if err != nil {
		return err
	}

	logger.Debug("Terminating")
	return nil
}

// terminateEC2InstancesByTags loops until there all instances which
// match the given tags are terminated.
func terminateEC2InstancesByTags(ctx context.Context, ec2Client *ec2.EC2, iamClient *iam.IAM, filters []Filter, logger logrus.FieldLogger) (map[string]struct{}, error) {
	if ec2Client.Config.Region == nil {
		return nil, errors.New("EC2 client does not have region configured")
	}

	partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), *ec2Client.Config.Region)
	if !ok {
		return nil, errors.Errorf("no partition found for region %q", *ec2Client.Config.Region)
	}

	terminated := map[string]struct{}{}
	err := wait.PollImmediateUntil(
		time.Second*10,
		func() (done bool, err error) {
			var loopError error
			matched := false
			for _, filter := range filters {
				logger.Debugf("search for and delete matching instances by tag matching %#+v", filter)
				instanceFilters := make([]*ec2.Filter, 0, len(filter))
				for key, value := range filter {
					instanceFilters = append(instanceFilters, &ec2.Filter{
						Name:   aws.String("tag:" + key),
						Values: []*string{aws.String(value)},
					})
				}
				err = ec2Client.DescribeInstancesPagesWithContext(
					ctx,
					&ec2.DescribeInstancesInput{Filters: instanceFilters},
					func(results *ec2.DescribeInstancesOutput, lastPage bool) bool {
						for _, reservation := range results.Reservations {
							if reservation.OwnerId == nil {
								continue
							}

							for _, instance := range reservation.Instances {
								if instance.InstanceId == nil || instance.State == nil {
									continue
								}

								instanceLogger := logger.WithField("instance", *instance.InstanceId)
								if *instance.State.Name == "terminated" {
									arn := fmt.Sprintf("arn:%s:ec2:%s:%s:instance/%s", partition.ID(), *ec2Client.Config.Region, *reservation.OwnerId, *instance.InstanceId)
									if _, ok := terminated[arn]; !ok {
										instanceLogger.Info("Terminated")
										terminated[arn] = exists
									}
									continue
								}

								matched = true
								err := terminateEC2InstanceByInstance(ctx, ec2Client, iamClient, instance, instanceLogger)
								if err != nil {
									instanceLogger.Debug(err)
									loopError = errors.Wrapf(err, "terminating %s", *instance.InstanceId)
									continue
								}
							}
						}

						return !lastPage
					},
				)
				if err != nil {
					err = errors.Wrap(err, "describe instances")
					logger.Info(err)
					matched = true
					loopError = err
				}
			}

			return !matched && loopError == nil, nil
		},
		ctx.Done(),
	)
	return terminated, err
}

// This is a bit of hack. Some objects, like Instance Profiles, can not be tagged in AWS.
// We "normally" find those objects by their relation to other objects. We have found,
// however, that people regularly delete all of their instances and roles outside of
// openshift-install destroy cluster. This means that we are unable to find the Instance
// Profiles.
//
// This code is a place to find specific objects like this which might be dangling.
func (o *ClusterUninstaller) deleteUntaggedResources(ctx context.Context, awsSession *session.Session) error {
	iamClient := iam.New(awsSession)
	masterProfile := fmt.Sprintf("%s-master-profile", o.ClusterID)
	if err := deleteIAMInstanceProfileByName(ctx, iamClient, &masterProfile, o.Logger); err != nil {
		return err
	}
	workerProfile := fmt.Sprintf("%s-worker-profile", o.ClusterID)
	if err := deleteIAMInstanceProfileByName(ctx, iamClient, &workerProfile, o.Logger); err != nil {
		return err
	}

	return nil
}

func deleteEC2InternetGateway(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeInternetGatewaysWithContext(ctx, &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{aws.String(id)},
	})
	if err != nil {
		return err
	}

	for _, gateway := range response.InternetGateways {
		for _, vpc := range gateway.Attachments {
			if vpc.VpcId == nil {
				logger.Warn("gateway does not have a VPC ID")
				continue
			}
			_, err := client.DetachInternetGatewayWithContext(ctx, &ec2.DetachInternetGatewayInput{
				InternetGatewayId: gateway.InternetGatewayId,
				VpcId:             vpc.VpcId,
			})
			if err == nil {
				logger.WithField("vpc", *vpc.VpcId).Debug("Detached")
			} else if err.(awserr.Error).Code() != "Gateway.NotAttached" {
				return errors.Wrapf(err, "detaching from %s", *vpc.VpcId)
			}
		}
	}

	_, err = client.DeleteInternetGatewayWithContext(ctx, &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: &id,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NATGateway(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteNatGatewayWithContext(ctx, &ec2.DeleteNatGatewayInput{
		NatGatewayId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "NatGatewayNotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NATGatewaysByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeNatGatewaysPagesWithContext(
		ctx,
		&ec2.DescribeNatGatewaysInput{
			Filter: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
		func(results *ec2.DescribeNatGatewaysOutput, lastPage bool) bool {
			for _, gateway := range results.NatGateways {
				err := deleteEC2NATGateway(ctx, client, *gateway.NatGatewayId, logger.WithField("NAT gateway", *gateway.NatGatewayId))
				if err != nil {
					if lastError != nil {
						logger.Debug(err)
					}
					lastError = errors.Wrapf(err, "deleting EC2 NAT gateway %s", *gateway.NatGatewayId)
					if failFast {
						return false
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteEC2RouteTable(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeRouteTablesWithContext(ctx, &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidRouteTableID.NotFound" {
			return nil
		}
		return err
	}

	for _, table := range response.RouteTables {
		err = deleteEC2RouteTableObject(ctx, client, table, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteEC2RouteTableObject(ctx context.Context, client *ec2.EC2, table *ec2.RouteTable, logger logrus.FieldLogger) error {
	hasMain := false
	for _, association := range table.Associations {
		if *association.Main {
			// can't remove the 'Main' association
			hasMain = true
			continue
		}
		_, err := client.DisassociateRouteTableWithContext(ctx, &ec2.DisassociateRouteTableInput{
			AssociationId: association.RouteTableAssociationId,
		})
		if err != nil {
			return errors.Wrapf(err, "dissociating %s", *association.RouteTableAssociationId)
		}
		logger.WithField("id", *association.RouteTableAssociationId).Info("Disassociated")
	}

	if hasMain {
		// can't delete route table with the 'Main' association
		// it will get cleaned up as part of deleting the VPC
		return nil
	}

	_, err := client.DeleteRouteTableWithContext(ctx, &ec2.DeleteRouteTableInput{
		RouteTableId: table.RouteTableId,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2RouteTablesByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeRouteTablesPagesWithContext(
		ctx,
		&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
		func(results *ec2.DescribeRouteTablesOutput, lastPage bool) bool {
			for _, table := range results.RouteTables {
				err := deleteEC2RouteTableObject(ctx, client, table, logger.WithField("table", *table.RouteTableId))
				if err != nil {
					if lastError != nil {
						logger.Debug(err)
					}
					lastError = errors.Wrapf(err, "deleting EC2 route table %s", *table.RouteTableId)
					if failFast {
						return false
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteEC2SecurityGroup(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	response, err := client.DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(id)},
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidGroup.NotFound" {
			return nil
		}
		return err
	}

	for _, group := range response.SecurityGroups {
		err = deleteEC2SecurityGroupObject(ctx, client, group, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteEC2SecurityGroupObject(ctx context.Context, client *ec2.EC2, group *ec2.SecurityGroup, logger logrus.FieldLogger) error {
	if group.GroupName != nil && *group.GroupName == "default" {
		logger.Debug("Skipping default security group")
		return nil
	}

	if len(group.IpPermissions) > 0 {
		_, err := client.RevokeSecurityGroupIngressWithContext(ctx, &ec2.RevokeSecurityGroupIngressInput{
			GroupId:       group.GroupId,
			IpPermissions: group.IpPermissions,
		})
		if err != nil {
			return errors.Wrap(err, "revoking ingress permissions")
		}
		logger.Debug("Revoked ingress permissions")
	}

	if len(group.IpPermissionsEgress) > 0 {
		_, err := client.RevokeSecurityGroupEgressWithContext(ctx, &ec2.RevokeSecurityGroupEgressInput{
			GroupId:       group.GroupId,
			IpPermissions: group.IpPermissionsEgress,
		})
		if err != nil {
			return errors.Wrap(err, "revoking egress permissions")
		}
		logger.Debug("Revoked egress permissions")
	}

	_, err := client.DeleteSecurityGroupWithContext(ctx, &ec2.DeleteSecurityGroupInput{
		GroupId: group.GroupId,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidGroup.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2SecurityGroupsByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeSecurityGroupsPagesWithContext(
		ctx,
		&ec2.DescribeSecurityGroupsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
		func(results *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			for _, group := range results.SecurityGroups {
				err := deleteEC2SecurityGroupObject(ctx, client, group, logger.WithField("security group", *group.GroupId))
				if err != nil {
					if lastError != nil {
						logger.Debug(err)
					}
					lastError = errors.Wrapf(err, "deleting EC2 security group %s", *group.GroupId)
					if failFast {
						return false
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteEC2Snapshot(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteSnapshotWithContext(ctx, &ec2.DeleteSnapshotInput{
		SnapshotId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidSnapshot.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NetworkInterface(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteNetworkInterfaceWithContext(ctx, &ec2.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidNetworkInterfaceID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2NetworkInterfaceByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeNetworkInterfacesPagesWithContext(
		ctx,
		&ec2.DescribeNetworkInterfacesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
		func(results *ec2.DescribeNetworkInterfacesOutput, lastPage bool) bool {
			for _, networkInterface := range results.NetworkInterfaces {
				err := deleteEC2NetworkInterface(ctx, client, *networkInterface.NetworkInterfaceId, logger.WithField("network interface", *networkInterface.NetworkInterfaceId))
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting EC2 network interface %s", *networkInterface.NetworkInterfaceId)
					if failFast {
						return false
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteEC2Subnet(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteSubnetWithContext(ctx, &ec2.DeleteSubnetInput{
		SubnetId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidSubnetID.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2SubnetsByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	// FIXME: port to DescribeSubnetsPages once we bump our vendored AWS package past v1.19.30
	results, err := client.DescribeSubnetsWithContext(
		ctx,
		&ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []*string{&vpc},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	var lastError error
	for _, subnet := range results.Subnets {
		err := deleteEC2Subnet(ctx, client, *subnet.SubnetId, logger.WithField("subnet", *subnet.SubnetId))
		if err != nil {
			err = errors.Wrapf(err, "deleting EC2 subnet %s", *subnet.SubnetId)
			if failFast {
				return err
			}
			if lastError != nil {
				logger.Debug(lastError)
			}
			lastError = err
		}
	}

	return lastError
}

func deleteEC2Volume(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVolumeWithContext(ctx, &ec2.DeleteVolumeInput{
		VolumeId: aws.String(id),
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidVolume.NotFound" {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPC(ctx context.Context, ec2Client *ec2.EC2, elbClient *elb.ELB, elbv2Client *elbv2.ELBV2, id string, logger logrus.FieldLogger) error {
	// first delete any Load Balancers under this VPC (not all of them are tagged)
	v1lbError := deleteElasticLoadBalancerClassicByVPC(ctx, elbClient, id, logger)
	v2lbError := deleteElasticLoadBalancerV2ByVPC(ctx, elbv2Client, id, logger)
	if v1lbError != nil {
		if v2lbError != nil {
			logger.Info(v2lbError)
		}
		return v1lbError
	} else if v2lbError != nil {
		return v2lbError
	}

	for _, child := range []struct {
		helper   func(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error
		failFast bool
	}{
		{helper: deleteEC2NATGatewaysByVPC, failFast: true},      // not always tagged
		{helper: deleteEC2NetworkInterfaceByVPC, failFast: true}, // not always tagged
		{helper: deleteEC2RouteTablesByVPC, failFast: true},      // not always tagged
		{helper: deleteEC2SecurityGroupsByVPC, failFast: false},  // not always tagged
		{helper: deleteEC2SubnetsByVPC, failFast: true},          // not always tagged
		{helper: deleteEC2VPCEndpointsByVPC, failFast: true},     // not taggable
	} {
		err := child.helper(ctx, ec2Client, id, child.failFast, logger)
		if err != nil {
			return err
		}
	}

	_, err := ec2Client.DeleteVpcWithContext(ctx, &ec2.DeleteVpcInput{
		VpcId: aws.String(id),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpoint(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVpcEndpointsWithContext(ctx, &ec2.DeleteVpcEndpointsInput{
		VpcEndpointIds: []*string{aws.String(id)},
	})
	if err != nil {
		return errors.Wrapf(err, "cannot delete VPC endpoint %s", id)
	}

	logger.Info("Deleted")
	return nil
}

func deleteEC2VPCEndpointsByVPC(ctx context.Context, client *ec2.EC2, vpc string, failFast bool, logger logrus.FieldLogger) error {
	response, err := client.DescribeVpcEndpointsWithContext(ctx, &ec2.DescribeVpcEndpointsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpc)},
			},
		},
	})

	if err != nil {
		return err
	}

	for _, endpoint := range response.VpcEndpoints {
		err := deleteEC2VPCEndpoint(ctx, client, *endpoint.VpcEndpointId, logger.WithField("VPC endpoint", *endpoint.VpcEndpointId))
		if err != nil {
			if err.(awserr.Error).Code() == "InvalidVpcID.NotFound" {
				return nil
			}
			return err
		}
	}

	return nil
}

func deleteEC2VPCPeeringConnection(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
	_, err := client.DeleteVpcPeeringConnectionWithContext(ctx, &ec2.DeleteVpcPeeringConnectionInput{
		VpcPeeringConnectionId: &id,
	})
	if err != nil {
		if err.(awserr.Error).Code() == "InvalidVpcPeeringConnectionID.NotFound" {
			return nil
		}
		return errors.Wrapf(err, "cannot delete VPC Peering Connection %s", id)
	}
	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancing(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "loadbalancer":
		segments := strings.SplitN(id, "/", 2)
		if len(segments) == 1 {
			return deleteElasticLoadBalancerClassic(ctx, elb.New(session), id, logger)
		} else if len(segments) != 2 {
			return errors.Errorf("cannot parse subresource %q into {subtype}/{id}", id)
		}
		subtype := segments[0]
		id = segments[1]
		switch subtype {
		case "net":
			return deleteElasticLoadBalancerV2(ctx, elbv2.New(session), arn, logger)
		default:
			return errors.Errorf("unrecognized elastic load balancing resource subtype %s", subtype)
		}
	case "targetgroup":
		return deleteElasticLoadBalancerTargetGroup(ctx, elbv2.New(session), arn, logger)
	default:
		return errors.Errorf("unrecognized elastic load balancing resource type %s", resourceType)
	}
}

func deleteElasticLoadBalancerClassic(ctx context.Context, client *elb.ELB, name string, logger logrus.FieldLogger) error {
	_, err := client.DeleteLoadBalancerWithContext(ctx, &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerClassicByVPC(ctx context.Context, client *elb.ELB, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeLoadBalancersPagesWithContext(
		ctx,
		&elb.DescribeLoadBalancersInput{},
		func(results *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, lb := range results.LoadBalancerDescriptions {
				lbLogger := logger.WithField("classic load balancer", *lb.LoadBalancerName)

				if lb.VPCId == nil {
					lbLogger.Warn("classic load balancer does not have a VPC ID so could not determine whether it should be deleted")
					continue
				}

				if *lb.VPCId != vpc {
					continue
				}

				err := deleteElasticLoadBalancerClassic(ctx, client, *lb.LoadBalancerName, lbLogger)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting classic load balancer %s", *lb.LoadBalancerName)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteElasticLoadBalancerTargetGroup(ctx context.Context, client *elbv2.ELBV2, arn arn.ARN, logger logrus.FieldLogger) error {
	_, err := client.DeleteTargetGroupWithContext(ctx, &elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(arn.String()),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerTargetGroupsByVPC(ctx context.Context, client *elbv2.ELBV2, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeTargetGroupsPagesWithContext(
		ctx,
		&elbv2.DescribeTargetGroupsInput{},
		func(results *elbv2.DescribeTargetGroupsOutput, lastPage bool) bool {
			for _, group := range results.TargetGroups {
				if group.VpcId == nil {
					logger.WithField("target group", *group.TargetGroupArn).Warn("load balancer target group does not have a VPC ID so could not determine whether it should be deleted")
					continue
				}

				if *group.VpcId != vpc {
					continue
				}

				parsed, err := arn.Parse(*group.TargetGroupArn)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrap(err, "parse ARN for target group")
					continue
				}

				err = deleteElasticLoadBalancerTargetGroup(ctx, client, parsed, logger.WithField("target group", parsed.Resource))
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting %s", parsed.String())
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteElasticLoadBalancerV2(ctx context.Context, client *elbv2.ELBV2, arn arn.ARN, logger logrus.FieldLogger) error {
	_, err := client.DeleteLoadBalancerWithContext(ctx, &elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn.String()),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticLoadBalancerV2ByVPC(ctx context.Context, client *elbv2.ELBV2, vpc string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.DescribeLoadBalancersPagesWithContext(
		ctx,
		&elbv2.DescribeLoadBalancersInput{},
		func(results *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, lb := range results.LoadBalancers {
				if lb.VpcId == nil {
					logger.WithField("load balancer", *lb.LoadBalancerArn).Warn("load balancer does not have a VPC ID so could not determine whether it should be deleted")
					continue
				}

				if *lb.VpcId != vpc {
					continue
				}

				parsed, err := arn.Parse(*lb.LoadBalancerArn)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrap(err, "parse ARN for load balancer")
					continue
				}

				err = deleteElasticLoadBalancerV2(ctx, client, parsed, logger.WithField("load balancer", parsed.Resource))
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting %s", parsed.String())
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	return err
}

func deleteIAM(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	client := iam.New(session)

	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "instance-profile":
		return deleteIAMInstanceProfile(ctx, client, arn, logger)
	case "role":
		return deleteIAMRole(ctx, client, arn, logger)
	case "user":
		return deleteIAMUser(ctx, client, id, logger)
	default:
		return errors.Errorf("unrecognized EC2 resource type %s", resourceType)
	}
}

func deleteIAMInstanceProfileByName(ctx context.Context, client *iam.IAM, name *string, logger logrus.FieldLogger) error {
	_, err := client.DeleteInstanceProfileWithContext(ctx, &iam.DeleteInstanceProfileInput{
		InstanceProfileName: name,
	})
	if err != nil {
		if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
			return nil
		}
		return err
	}
	logger.WithField("InstanceProfileName", *name).Info("Deleted")
	return err
}

func deleteIAMInstanceProfile(ctx context.Context, client *iam.IAM, profileARN arn.ARN, logger logrus.FieldLogger) error {
	resourceType, name, err := splitSlash("resource", profileARN.Resource)
	if err != nil {
		return err
	}

	if resourceType != "instance-profile" {
		return errors.Errorf("%s ARN passed to deleteIAMInstanceProfile: %s", resourceType, profileARN.String())
	}

	response, err := client.GetInstanceProfileWithContext(ctx, &iam.GetInstanceProfileInput{
		InstanceProfileName: &name,
	})
	if err != nil {
		if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
			return nil
		}
		return err
	}
	profile := response.InstanceProfile

	for _, role := range profile.Roles {
		_, err = client.RemoveRoleFromInstanceProfileWithContext(ctx, &iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: profile.InstanceProfileName,
			RoleName:            role.RoleName,
		})
		if err != nil {
			return errors.Wrapf(err, "dissociating %s", *role.RoleName)
		}
		logger.WithField("name", name).WithField("role", *role.RoleName).Info("Disassociated")
	}

	logger = logger.WithField("arn", profileARN.String())
	if err := deleteIAMInstanceProfileByName(ctx, client, profile.InstanceProfileName, logger); err != nil {
		return err
	}

	return nil
}

func deleteIAMRole(ctx context.Context, client *iam.IAM, roleARN arn.ARN, logger logrus.FieldLogger) error {
	resourceType, name, err := splitSlash("resource", roleARN.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("name", name)

	if resourceType != "role" {
		return errors.Errorf("%s ARN passed to deleteIAMRole: %s", resourceType, roleARN.String())
	}

	var lastError error
	err = client.ListRolePoliciesPagesWithContext(
		ctx,
		&iam.ListRolePoliciesInput{RoleName: &name},
		func(results *iam.ListRolePoliciesOutput, lastPage bool) bool {
			for _, policy := range results.PolicyNames {
				_, err := client.DeleteRolePolicyWithContext(ctx, &iam.DeleteRolePolicyInput{
					RoleName:   &name,
					PolicyName: policy,
				})
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting IAM role policy %s", *policy)
				}
				logger.WithField("policy", *policy).Info("Deleted")
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM role policies")
	}

	err = client.ListInstanceProfilesForRolePagesWithContext(
		ctx,
		&iam.ListInstanceProfilesForRoleInput{RoleName: &name},
		func(results *iam.ListInstanceProfilesForRoleOutput, lastPage bool) bool {
			for _, profile := range results.InstanceProfiles {
				parsed, err := arn.Parse(*profile.Arn)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrap(err, "parse ARN for IAM instance profile")
					continue
				}

				err = deleteIAMInstanceProfile(ctx, client, parsed, logger)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting %s", parsed.String())
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM instance profiles")
	}

	_, err = client.DeleteRoleWithContext(ctx, &iam.DeleteRoleInput{RoleName: &name})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteIAMUser(ctx context.Context, client *iam.IAM, id string, logger logrus.FieldLogger) error {
	var lastError error
	err := client.ListUserPoliciesPagesWithContext(
		ctx,
		&iam.ListUserPoliciesInput{UserName: &id},
		func(results *iam.ListUserPoliciesOutput, lastPage bool) bool {
			for _, policy := range results.PolicyNames {
				_, err := client.DeleteUserPolicyWithContext(ctx, &iam.DeleteUserPolicyInput{
					UserName:   &id,
					PolicyName: policy,
				})
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting IAM user policy %s", *policy)
				}
				logger.WithField("policy", *policy).Info("Deleted")
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM user policies")
	}

	err = client.ListAccessKeysPagesWithContext(
		ctx,
		&iam.ListAccessKeysInput{UserName: &id},
		func(results *iam.ListAccessKeysOutput, lastPage bool) bool {
			for _, key := range results.AccessKeyMetadata {
				_, err := client.DeleteAccessKeyWithContext(ctx, &iam.DeleteAccessKeyInput{
					UserName:    &id,
					AccessKeyId: key.AccessKeyId,
				})
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting IAM access key %s", *key.AccessKeyId)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing IAM access keys")
	}

	_, err = client.DeleteUserWithContext(ctx, &iam.DeleteUserInput{
		UserName: &id,
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteRoute53(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	if resourceType != "hostedzone" {
		return errors.Errorf("unrecognized Route 53 resource type %s", resourceType)
	}

	client := route53.New(session)

	sharedZoneID, err := getSharedHostedZone(ctx, client, id, logger)
	if err != nil {
		// In some cases AWS may return the zone in the list of tagged resources despite the fact
		// it no longer exists.
		if err.(awserr.Error).Code() == route53.ErrCodeNoSuchHostedZone {
			return nil
		}
		return err
	}

	recordSetKey := func(recordSet *route53.ResourceRecordSet) string {
		return fmt.Sprintf("%s %s", *recordSet.Type, *recordSet.Name)
	}

	sharedEntries := map[string]*route53.ResourceRecordSet{}
	if len(sharedZoneID) != 0 {
		err = client.ListResourceRecordSetsPagesWithContext(
			ctx,
			&route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(sharedZoneID)},
			func(results *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
				for _, recordSet := range results.ResourceRecordSets {
					key := recordSetKey(recordSet)
					sharedEntries[key] = recordSet
				}

				return !lastPage
			},
		)
		if err != nil {
			return err
		}
	} else {
		logger.Debug("shared public zone not found")
	}

	var lastError error
	err = client.ListResourceRecordSetsPagesWithContext(
		ctx,
		&route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(id)},
		func(results *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			for _, recordSet := range results.ResourceRecordSets {
				if *recordSet.Type == "SOA" || *recordSet.Type == "NS" {
					// can't delete SOA and NS types
					continue
				}
				key := recordSetKey(recordSet)
				if sharedEntry, ok := sharedEntries[key]; ok {
					err := deleteRoute53RecordSet(ctx, client, sharedZoneID, sharedEntry, logger.WithField("public zone", sharedZoneID))
					if err != nil {
						if lastError != nil {
							logger.Debug(lastError)
						}
						lastError = errors.Wrapf(err, "deleting public zone %s", sharedZoneID)
					}
				}

				err = deleteRoute53RecordSet(ctx, client, id, recordSet, logger)
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting record set %#+v from zone %s", recordSet, id)
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return err
	}

	_, err = client.DeleteHostedZoneWithContext(ctx, &route53.DeleteHostedZoneInput{
		Id: aws.String(id),
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteRoute53RecordSet(ctx context.Context, client *route53.Route53, zoneID string, recordSet *route53.ResourceRecordSet, logger logrus.FieldLogger) error {
	logger = logger.WithField("record set", fmt.Sprintf("%s %s", *recordSet.Type, *recordSet.Name))
	_, err := client.ChangeResourceRecordSetsWithContext(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action:            aws.String("DELETE"),
					ResourceRecordSet: recordSet,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteS3(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	client := s3.New(session)

	iter := s3manager.NewDeleteListIterator(client, &s3.ListObjectsInput{
		Bucket: aws.String(arn.Resource),
	})
	err := s3manager.NewBatchDeleteWithClient(client).Delete(ctx, iter)
	if err != nil && !isBucketNotFound(err) {
		return err
	}
	logger.Debug("Emptied")

	var lastError error
	err = client.ListObjectVersionsPagesWithContext(ctx, &s3.ListObjectVersionsInput{
		Bucket:  aws.String(arn.Resource),
		MaxKeys: aws.Int64(1000),
	}, func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
		var deleteObjects []*s3.ObjectIdentifier
		for _, deleteMarker := range page.DeleteMarkers {
			deleteObjects = append(deleteObjects, &s3.ObjectIdentifier{
				Key:       aws.String(*deleteMarker.Key),
				VersionId: aws.String(*deleteMarker.VersionId),
			})
		}
		for _, version := range page.Versions {
			deleteObjects = append(deleteObjects, &s3.ObjectIdentifier{
				Key:       aws.String(*version.Key),
				VersionId: aws.String(*version.VersionId),
			})
		}
		if len(deleteObjects) > 0 {
			_, err := client.DeleteObjectsWithContext(ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(arn.Resource),
				Delete: &s3.Delete{
					Objects: deleteObjects,
				},
			})
			if err != nil {
				lastError = errors.Wrapf(err, "delete object failed %v", err)
			}
		}
		return !lastPage
	})
	if lastError != nil {
		return lastError
	}
	if err != nil && !isBucketNotFound(err) {
		return err
	}
	logger.Debug("Versions Deleted")

	_, err = client.DeleteBucketWithContext(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(arn.Resource),
	})
	if err != nil && !isBucketNotFound(err) {
		return err
	}

	logger.Info("Deleted")
	return nil
}

func isBucketNotFound(err interface{}) bool {
	switch s3Err := err.(type) {
	case awserr.Error:
		if s3Err.Code() == "NoSuchBucket" {
			return true
		}
		origErr := s3Err.OrigErr()
		if origErr != nil {
			return isBucketNotFound(origErr)
		}
	case s3manager.Error:
		if s3Err.OrigErr != nil {
			return isBucketNotFound(s3Err.OrigErr)
		}
	case s3manager.Errors:
		if len(s3Err) == 1 {
			return isBucketNotFound(s3Err[0])
		}
	}
	return false
}

func removeSharedTags(ctx context.Context, tagClients []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI, tagClientNames map[*resourcegroupstaggingapi.ResourceGroupsTaggingAPI]string, filters []Filter, logger logrus.FieldLogger) error {
	for _, filter := range filters {
		for key, value := range filter {
			if strings.HasPrefix(key, "kubernetes.io/cluster/") {
				if value == "owned" {
					if err := removeSharedTag(ctx, tagClients, tagClientNames, key, logger); err != nil {
						return err
					}
				} else {
					logger.Warnf("Ignoring non-owned cluster key %s: %s for shared-tag removal", key, value)
				}
			}
		}
	}

	return nil
}

func removeSharedTag(ctx context.Context, tagClients []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI, tagClientNames map[*resourcegroupstaggingapi.ResourceGroupsTaggingAPI]string, key string, logger logrus.FieldLogger) error {
	request := &resourcegroupstaggingapi.UntagResourcesInput{
		TagKeys: []*string{aws.String(key)},
	}

	removed := map[string]struct{}{}
	tagClients = append([]*resourcegroupstaggingapi.ResourceGroupsTaggingAPI(nil), tagClients...)
	for len(tagClients) > 0 {
		nextTagClients := tagClients[:0]
		for _, tagClient := range tagClients {
			logger.Debugf("Search for and remove tags in %s matching %s: shared", tagClientNames[tagClient], key)
			arns := []string{}
			err := tagClient.GetResourcesPagesWithContext(
				ctx,
				&resourcegroupstaggingapi.GetResourcesInput{TagFilters: []*resourcegroupstaggingapi.TagFilter{{
					Key:    aws.String(key),
					Values: []*string{aws.String("shared")},
				}}},
				func(results *resourcegroupstaggingapi.GetResourcesOutput, lastPage bool) bool {
					for _, resource := range results.ResourceTagMappingList {
						arn := *resource.ResourceARN
						if _, ok := removed[arn]; !ok {
							arns = append(arns, arn)
						}
					}

					return !lastPage
				},
			)
			if err != nil {
				err = errors.Wrap(err, "get tagged resources")
				logger.Info(err)
				nextTagClients = append(nextTagClients, tagClient)
				continue
			}
			if len(arns) == 0 {
				logger.Debugf("No matches in %s for %s: shared, removing client", tagClientNames[tagClient], key)
				continue
			}
			nextTagClients = append(nextTagClients, tagClient)

			for i := 0; i < len(arns); i += 20 {
				request.ResourceARNList = make([]*string, 0, 20)
				for j := 0; i+j < len(arns) && j < 20; j++ {
					request.ResourceARNList = append(request.ResourceARNList, aws.String(arns[i+j]))
				}
				_, err = tagClient.UntagResourcesWithContext(ctx, request)
				if err != nil {
					err = errors.Wrap(err, "untag shared resources")
					logger.Info(err)
					continue
				}
				for j := 0; i+j < len(arns) && j < 20; j++ {
					arn := arns[i+j]
					logger.WithField("arn", arn).Infof("Removed tag %s: shared", key)
					removed[arn] = exists
				}
			}
		}
		tagClients = nextTagClients
	}

	return nil
}
