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
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
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
	Filters       []Filter // filter(s) we will be searching for
	Logger        logrus.FieldLogger
	Region        string
	ClusterID     string
	ClusterDomain string

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
		Filters:       filters,
		Region:        region,
		Logger:        logger,
		ClusterID:     metadata.InfraID,
		ClusterDomain: metadata.AWS.ClusterDomain,
		Session:       session,
	}, nil
}

func (o *ClusterUninstaller) validate() error {
	if len(o.Filters) == 0 {
		return errors.Errorf("you must specify at least one tag filter")
	}
	switch r := o.Region; r {
	case "us-iso-east-1":
		return errors.Errorf("cannot destroy cluster in region %q", r)
	}
	return nil
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() error {
	_, err := o.RunWithContext(context.Background())
	return err
}

// RunWithContext runs the uninstall process with a context.
// The first return is the list of ARNs for resources that could not be destroyed.
func (o *ClusterUninstaller) RunWithContext(ctx context.Context) ([]string, error) {
	err := o.validate()
	if err != nil {
		return nil, err
	}

	awsSession := o.Session
	if awsSession == nil {
		// Relying on appropriate AWS ENV vars (eg AWS_PROFILE, AWS_ACCESS_KEY_ID, etc)
		awsSession, err = session.NewSession(aws.NewConfig().WithRegion(o.Region))
		if err != nil {
			return nil, err
		}
	}
	awsSession.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "openshiftInstaller.OpenshiftInstallerUserAgentHandler",
		Fn:   request.MakeAddToUserAgentHandler("OpenShift/4.x Destroyer", version.Raw),
	})

	tagClients := []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI{
		resourcegroupstaggingapi.New(awsSession),
	}

	switch o.Region {
	case endpoints.CnNorth1RegionID, endpoints.CnNorthwest1RegionID:
		if o.Region != endpoints.CnNorthwest1RegionID {
			tagClients = append(tagClients,
				resourcegroupstaggingapi.New(awsSession, aws.NewConfig().WithRegion(endpoints.CnNorthwest1RegionID)))
		}
	case endpoints.UsGovEast1RegionID, endpoints.UsGovWest1RegionID:
		if o.Region != endpoints.UsGovWest1RegionID {
			tagClients = append(tagClients,
				resourcegroupstaggingapi.New(awsSession, aws.NewConfig().WithRegion(endpoints.UsGovWest1RegionID)))
		}
	default:
		if o.Region != endpoints.UsEast1RegionID {
			tagClients = append(tagClients,
				resourcegroupstaggingapi.New(awsSession, aws.NewConfig().WithRegion(endpoints.UsEast1RegionID)))
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

	// Get the initial resources to delete, so that they can be returned if the context is canceled while terminating
	// instances.
	deleted := sets.NewString()
	resourcesToDelete, tagClientsWithResources, err := o.findResourcesToDelete(ctx, tagClients, iamClient, iamRoleSearch, iamUserSearch, deleted)
	if err != nil {
		o.Logger.WithError(err).Info("error while finding resources to delete")
		if err := ctx.Err(); err != nil {
			return resourcesToDelete.UnsortedList(), err
		}
	}

	tracker := new(errorTracker)

	// Terminate EC2 instances. The instances need to be terminated first so that we can ensure that there is nothing
	// running on the cluster creating new resources while we are attempting to delete resources, which could leak
	// the new resources.
	ec2Client := ec2.New(awsSession)
	lastTerminateTime := time.Now()
	err = wait.PollImmediateUntil(
		time.Second*10,
		func() (done bool, err error) {
			instancesRunning, instancesNotTerminated, err := o.findEC2Instances(ctx, ec2Client, deleted)
			if err != nil {
				o.Logger.WithError(err).Info("error while finding EC2 instances to delete")
				if err := ctx.Err(); err != nil {
					return false, err
				}
			}
			if len(instancesNotTerminated) == 0 && len(instancesRunning) == 0 && err == nil {
				return true, nil
			}
			instancesToDelete := instancesRunning
			if time.Since(lastTerminateTime) > 10*time.Minute {
				instancesToDelete = instancesNotTerminated
				lastTerminateTime = time.Now()
			}
			newlyDeleted, err := o.deleteResources(ctx, awsSession, instancesToDelete, tracker)
			// Delete from the resources-to-delete set so that the current state of the resources to delete can be
			// returned if the context is completed.
			resourcesToDelete = resourcesToDelete.Difference(newlyDeleted)
			deleted = deleted.Union(newlyDeleted)
			if err != nil {
				if err := ctx.Err(); err != nil {
					return false, err
				}
			}
			return false, nil
		},
		ctx.Done(),
	)
	if err != nil {
		return resourcesToDelete.UnsortedList(), err
	}

	// Delete the rest of the resources.
	err = wait.PollImmediateUntil(
		time.Second*10,
		func() (done bool, err error) {
			newlyDeleted, loopError := o.deleteResources(ctx, awsSession, resourcesToDelete.UnsortedList(), tracker)
			// Delete from the resources-to-delete set so that the current state of the resources to delete can be
			// returned if the context is completed.
			resourcesToDelete = resourcesToDelete.Difference(newlyDeleted)
			deleted = deleted.Union(newlyDeleted)
			if loopError != nil {
				if err := ctx.Err(); err != nil {
					return false, err
				}
			}
			// Store resources to delete in a temporary variable so that, in case the context is cancelled, the current
			// resources to delete are not lost.
			nextResourcesToDelete, nextTagClients, err := o.findResourcesToDelete(ctx, tagClientsWithResources, iamClient, iamRoleSearch, iamUserSearch, deleted)
			if err != nil {
				o.Logger.WithError(err).Info("error while finding resources to delete")
				if err := ctx.Err(); err != nil {
					return false, err
				}
				loopError = errors.Wrap(err, "error while finding resources to delete")
			}
			resourcesToDelete = nextResourcesToDelete
			tagClientsWithResources = nextTagClients
			return len(resourcesToDelete) == 0 && loopError == nil, nil
		},
		ctx.Done(),
	)
	if err != nil {
		return resourcesToDelete.UnsortedList(), err
	}

	err = o.removeSharedTags(ctx, awsSession, tagClients, tracker)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// findEC2Instances returns the EC2 instances with tags that satisfy the filters.
// returns two lists, first one is the list of all resources that are not terminated and are not in shutdown
// stage and the second list is the list of resources that are not terminated.
//   deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func (o *ClusterUninstaller) findEC2Instances(ctx context.Context, ec2Client *ec2.EC2, deleted sets.String) ([]string, []string, error) {
	if ec2Client.Config.Region == nil {
		return nil, nil, errors.New("EC2 client does not have region configured")
	}

	partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), *ec2Client.Config.Region)
	if !ok {
		return nil, nil, errors.Errorf("no partition found for region %q", *ec2Client.Config.Region)
	}

	var resourcesRunning []string
	var resourcesNotTerminated []string
	for _, filter := range o.Filters {
		o.Logger.Debugf("search for instances by tag matching %#+v", filter)
		instanceFilters := make([]*ec2.Filter, 0, len(filter))
		for key, value := range filter {
			instanceFilters = append(instanceFilters, &ec2.Filter{
				Name:   aws.String("tag:" + key),
				Values: []*string{aws.String(value)},
			})
		}
		err := ec2Client.DescribeInstancesPagesWithContext(
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

						instanceLogger := o.Logger.WithField("instance", *instance.InstanceId)
						arn := fmt.Sprintf("arn:%s:ec2:%s:%s:instance/%s", partition.ID(), *ec2Client.Config.Region, *reservation.OwnerId, *instance.InstanceId)
						if *instance.State.Name == "terminated" {
							if !deleted.Has(arn) {
								instanceLogger.Info("Terminated")
								deleted.Insert(arn)
							}
							continue
						}
						if *instance.State.Name != "shutting-down" {
							resourcesRunning = append(resourcesRunning, arn)
						}
						resourcesNotTerminated = append(resourcesNotTerminated, arn)
					}
				}
				return !lastPage
			},
		)
		if err != nil {
			err = errors.Wrap(err, "get ec2 instances")
			o.Logger.Info(err)
			return resourcesRunning, resourcesNotTerminated, err
		}
	}
	return resourcesRunning, resourcesNotTerminated, nil
}

// findResourcesToDelete returns the resources that should be deleted.
//   tagClients - clients of the tagging API to use to search for resources.
//   deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func (o *ClusterUninstaller) findResourcesToDelete(
	ctx context.Context,
	tagClients []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI,
	iamClient *iam.IAM,
	iamRoleSearch *iamRoleSearch,
	iamUserSearch *iamUserSearch,
	deleted sets.String,
) (sets.String, []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI, error) {
	resources := sets.NewString()
	var tagClientsWithResources []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI
	var errs []error

	// Find resources by tag
	for _, tagClient := range tagClients {
		resourcesInTagClient, err := o.findResourcesByTag(ctx, tagClient, deleted)
		if err != nil {
			errs = append(errs, err)
		}
		resources = resources.Union(resourcesInTagClient)
		// If there are still resources to be deleted for the tag client or if there was an error getting the resources
		// for the tag client, then retain the tag client for future queries.
		if len(resourcesInTagClient) > 0 || err != nil {
			tagClientsWithResources = append(tagClientsWithResources, tagClient)
		} else {
			o.Logger.Debugf("no deletions from %s, removing client", *tagClient.Config.Region)
		}
	}

	// Find IAM roles
	iamRoleResources, err := o.findIAMRoles(ctx, iamRoleSearch, deleted)
	if err != nil {
		errs = append(errs, err)
	}
	resources = resources.Union(iamRoleResources)

	// Find IAM users
	iamUserResources, err := o.findIAMUsers(ctx, iamUserSearch, deleted)
	if err != nil {
		errs = append(errs, err)
	}
	resources = resources.Union(iamUserResources)

	// Find untaggable resources
	untaggableResources, err := o.findUntaggableResources(ctx, iamClient, deleted)
	if err != nil {
		errs = append(errs, err)
	}
	resources = resources.Union(untaggableResources)

	return resources, tagClientsWithResources, utilerrors.NewAggregate(errs)
}

// findResourcesByTag returns the resources with tags that satisfy the filters.
//   tagClients - clients of the tagging API to use to search for resources.
//   deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func (o *ClusterUninstaller) findResourcesByTag(
	ctx context.Context,
	tagClient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI,
	deleted sets.String,
) (sets.String, error) {
	resources := sets.NewString()
	for _, filter := range o.Filters {
		o.Logger.Debugf("search for matching resources by tag in %s matching %#+v", *tagClient.Config.Region, filter)
		tagFilters := make([]*resourcegroupstaggingapi.TagFilter, 0, len(filter))
		for key, value := range filter {
			tagFilters = append(tagFilters, &resourcegroupstaggingapi.TagFilter{
				Key:    aws.String(key),
				Values: []*string{aws.String(value)},
			})
		}
		err := tagClient.GetResourcesPagesWithContext(
			ctx,
			&resourcegroupstaggingapi.GetResourcesInput{TagFilters: tagFilters},
			func(results *resourcegroupstaggingapi.GetResourcesOutput, lastPage bool) bool {
				for _, resource := range results.ResourceTagMappingList {
					arnString := *resource.ResourceARN
					if !deleted.Has(arnString) {
						resources.Insert(arnString)
					}
				}
				return !lastPage
			},
		)
		if err != nil {
			err = errors.Wrap(err, "get tagged resources")
			o.Logger.Info(err)
			return resources, err
		}
	}
	return resources, nil
}

// findIAMRoles returns the IAM roles for the cluster.
//   deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func (o *ClusterUninstaller) findIAMRoles(ctx context.Context, search *iamRoleSearch, deleted sets.String) (sets.String, error) {
	o.Logger.Debug("search for IAM roles")
	resources, _, err := search.find(ctx)
	if err != nil {
		o.Logger.Info(err)
		return nil, err
	}
	return sets.NewString(resources...).Difference(deleted), nil
}

// findIAMUsers returns the IAM users for the cluster.
//   deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func (o *ClusterUninstaller) findIAMUsers(ctx context.Context, search *iamUserSearch, deleted sets.String) (sets.String, error) {
	o.Logger.Debug("search for IAM users")
	resources, err := search.arns(ctx)
	if err != nil {
		o.Logger.Info(err)
		return nil, err
	}
	return sets.NewString(resources...).Difference(deleted), nil
}

// findUntaggableResources returns the resources for the cluster that cannot be tagged.
//   deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func (o *ClusterUninstaller) findUntaggableResources(ctx context.Context, iamClient *iam.IAM, deleted sets.String) (sets.String, error) {
	resources := sets.NewString()
	o.Logger.Debug("search for IAM instance profiles")
	for _, profileType := range []string{"master", "worker"} {
		profile := fmt.Sprintf("%s-%s-profile", o.ClusterID, profileType)
		response, err := iamClient.GetInstanceProfileWithContext(ctx, &iam.GetInstanceProfileInput{InstanceProfileName: &profile})
		if err != nil {
			if err.(awserr.Error).Code() == iam.ErrCodeNoSuchEntityException {
				continue
			}
			return resources, errors.Wrap(err, "failed to get IAM instance profile")
		}
		arnString := *response.InstanceProfile.Arn
		if !deleted.Has(arnString) {
			resources.Insert(arnString)
		}
	}
	return resources, nil
}

// deleteResources deletes the specified resources.
//   resources - the resources to be deleted.
// The first return is the ARNs of the resources that were successfully deleted
func (o *ClusterUninstaller) deleteResources(ctx context.Context, awsSession *session.Session, resources []string, tracker *errorTracker) (sets.String, error) {
	deleted := sets.NewString()
	for _, arnString := range resources {
		logger := o.Logger.WithField("arn", arnString)
		parsedARN, err := arn.Parse(arnString)
		if err != nil {
			logger.WithError(err).Debug("could not parse ARN")
			continue
		}
		if err := deleteARN(ctx, awsSession, parsedARN, o.Logger); err != nil {
			tracker.suppressWarning(arnString, err, logger)
			if err := ctx.Err(); err != nil {
				return deleted, err
			}
			continue
		}
		deleted.Insert(arnString)
	}
	return deleted, nil
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

type iamRoleSearch struct {
	client    *iam.IAM
	filters   []Filter
	logger    logrus.FieldLogger
	unmatched map[string]struct{}
}

func (search *iamRoleSearch) find(ctx context.Context) (arns []string, names []string, returnErr error) {
	if search.unmatched == nil {
		search.unmatched = map[string]struct{}{}
	}

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
						names = append(names, *role.RoleName)
					} else {
						search.unmatched[*role.Arn] = exists
					}
				}
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return arns, names, lastError
	}
	return arns, names, err
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

// getPublicHostedZone will find the ID of the non-Terraform-managed public route53 zone given the
// Terraform-managed zone's privateID.
func getPublicHostedZone(ctx context.Context, client *route53.Route53, privateID string, logger logrus.FieldLogger) (string, error) {
	response, err := client.GetHostedZoneWithContext(ctx, &route53.GetHostedZoneInput{
		Id: aws.String(privateID),
	})
	if err != nil {
		return "", err
	}

	privateName := *response.HostedZone.Name

	if response.HostedZone.Config != nil && response.HostedZone.Config.PrivateZone != nil {
		if !*response.HostedZone.Config.PrivateZone {
			return "", errors.Errorf("getPublicHostedZone requires a private ID, but was passed the public %s", privateID)
		}
	} else {
		logger.WithField("hosted zone", privateName).Warn("could not determine whether hosted zone is private")
	}

	return findAncestorPublicRoute53(ctx, client, privateName, logger)
}

// findAncestorPublicRoute53 finds a public route53 zone with the closest ancestor or match to dnsName.
// It returns "", when no public route53 zone could be found.
func findAncestorPublicRoute53(ctx context.Context, client *route53.Route53, dnsName string, logger logrus.FieldLogger) (string, error) {
	for len(dnsName) > 0 {
		sZone, err := findPublicRoute53(ctx, client, dnsName, logger)
		if err != nil {
			return "", err
		}
		if sZone != "" {
			return sZone, nil
		}

		idx := strings.Index(dnsName, ".")
		if idx == -1 {
			break
		}
		dnsName = dnsName[idx+1:]
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

func deleteARN(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	switch arn.Service {
	case "ec2":
		return deleteEC2(ctx, session, arn, logger)
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

func deleteEC2(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
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
		return deleteEC2Image(ctx, client, id, logger)
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

func deleteEC2Image(ctx context.Context, client *ec2.EC2, id string, logger logrus.FieldLogger) error {
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
	for _, image := range response.Images {
		var snapshots []*string
		for _, bdm := range image.BlockDeviceMappings {
			if bdm.Ebs != nil && bdm.Ebs.SnapshotId != nil {
				snapshots = append(snapshots, bdm.Ebs.SnapshotId)
			}
		}
		if len(snapshots) != 0 {
			_, err = client.CreateTagsWithContext(ctx, &ec2.CreateTagsInput{
				Resources: snapshots,
				Tags:      image.Tags,
			})
			if err != nil {
				return errors.Wrapf(err, "tagging snapshots for %s", id)
			}
		}
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
	// Ignore instances that are already terminated
	if instance.State == nil || *instance.State.Name == "terminated" {
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

	if group.GroupName != nil && *group.GroupName == "default" {
		logger.Debug("Skipping default security group")
		return nil
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

	err = client.ListAttachedRolePoliciesPagesWithContext(
		ctx,
		&iam.ListAttachedRolePoliciesInput{RoleName: &name},
		func(results *iam.ListAttachedRolePoliciesOutput, lastPage bool) bool {
			for _, policy := range results.AttachedPolicies {
				_, err := client.DetachRolePolicyWithContext(ctx, &iam.DetachRolePolicyInput{
					RoleName:  &name,
					PolicyArn: policy.PolicyArn,
				})
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "detaching IAM role policy %s", *policy.PolicyName)
				}
				logger.WithField("policy", *policy.PolicyName).Info("Detached")
			}

			return !lastPage
		},
	)

	if lastError != nil {
		return lastError
	}
	if err != nil {
		return errors.Wrap(err, "listing attached IAM role policies")
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

	publicZoneID, err := getPublicHostedZone(ctx, client, id, logger)
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

	publicEntries := map[string]*route53.ResourceRecordSet{}
	if len(publicZoneID) != 0 {
		err = client.ListResourceRecordSetsPagesWithContext(
			ctx,
			&route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(publicZoneID)},
			func(results *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
				for _, recordSet := range results.ResourceRecordSets {
					key := recordSetKey(recordSet)
					publicEntries[key] = recordSet
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
				if publicEntry, ok := publicEntries[key]; ok {
					err := deleteRoute53RecordSet(ctx, client, publicZoneID, publicEntry, logger.WithField("public zone", publicZoneID))
					if err != nil {
						if lastError != nil {
							logger.Debug(lastError)
						}
						lastError = errors.Wrapf(err, "deleting record set %#v from public zone %s", publicEntry, publicZoneID)
					}
					// do not delete the record set in the private zone if the delete failed in the public zone;
					// otherwise the record set in the public zone will get leaked
					continue
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
		if err.(awserr.Error).Code() == "NoSuchHostedZone" {
			return nil
		}
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
