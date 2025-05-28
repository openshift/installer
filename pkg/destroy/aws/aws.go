package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	ec2v2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2v2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	iamv2 "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	tagtypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"

	awssession "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/version"
)

var exists = struct{}{}

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
	Filters        []Filter // filter(s) we will be searching for
	Logger         logrus.FieldLogger
	Region         string
	ClusterID      string
	ClusterDomain  string
	HostedZoneRole string
	endpoints      []awstypes.ServiceEndpoint

	EC2Client   *ec2v2.Client
	EFSClient   *efs.Client
	ELBClient   *elb.Client
	ELBV2Client *elbv2.Client

	IAMClient     *iamv2.Client
	Route53Client *route53.Client
	S3Client      *s3.Client
}

const (
	endpointUSEast1      = "us-east-1"
	endpointCNNorth1     = "cn-north-1"
	endpointCNNorthWest1 = "cn-northwest-1"
	endpointISOEast1     = "us-iso-east-1"
	endpointISOWest1     = "us-iso-west-1"
	endpointISOBEast1    = "us-isob-east-1"
	endpointUSGovEast1   = "us-gov-east-1"
	endpointUSGovWest1   = "us-gov-west-1"
)

// New returns an AWS destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	filters := make([]Filter, 0, len(metadata.ClusterPlatformMetadata.AWS.Identifier))
	for _, filter := range metadata.ClusterPlatformMetadata.AWS.Identifier {
		filters = append(filters, filter)
	}
	region := metadata.ClusterPlatformMetadata.AWS.Region

	ctx := context.Background()
	ec2Client, err := awssession.NewEC2Client(ctx, awssession.EndpointOptions{
		Region:    region,
		Endpoints: metadata.AWS.ServiceEndpoints,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create EC2 client: %w", err)
	}

	iamClient, err := awssession.NewIAMClient(ctx, awssession.EndpointOptions{
		Region:    region,
		Endpoints: metadata.AWS.ServiceEndpoints,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM client: %w", err)
	}

	// FIXME: remove this code when the elb and elbv2 clients are "fixed" or figured out
	elbCfg, err := awssession.GetConfigWithOptions(ctx, configv2.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config for elb client: %w", err)
	}
	elbclient := elb.NewFromConfig(elbCfg, func(options *elb.Options) {
		options.Region = region
		for _, endpoint := range metadata.AWS.ServiceEndpoints {
			if strings.EqualFold(endpoint.Name, "elb") {
				options.BaseEndpoint = awsv2.String(endpoint.URL)
			}
		}
	})

	// FIXME: remove this code when the elb and elbv2 clients are "fixed" or figured out
	elbv2Cfg, err := awssession.GetConfigWithOptions(ctx, configv2.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config for elbv2 client: %w", err)
	}
	elbv2client := elbv2.NewFromConfig(elbv2Cfg, func(options *elbv2.Options) {
		options.Region = region
		for _, endpoint := range metadata.AWS.ServiceEndpoints {
			if strings.EqualFold(endpoint.Name, "elbv2") {
				options.BaseEndpoint = awsv2.String(endpoint.URL)
			}
		}
	})

	// FIXME: remove this code when the s3client is made
	s3Cfg, err := awssession.GetConfigWithOptions(ctx, configv2.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config for S3 client: %w", err)
	}
	s3Client := s3.NewFromConfig(s3Cfg, func(options *s3.Options) {
		options.Region = region
		for _, endpoint := range metadata.AWS.ServiceEndpoints {
			if strings.EqualFold(endpoint.Name, "s3") {
				options.BaseEndpoint = awsv2.String(endpoint.URL)
			}
		}
	})

	// FIXME: remove this code when the EFS client is made
	efsCfg, err := awssession.GetConfigWithOptions(ctx, configv2.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config for EFS client: %w", err)
	}
	efsClient := efs.NewFromConfig(efsCfg, func(options *efs.Options) {
		options.Region = region
		for _, endpoint := range metadata.AWS.ServiceEndpoints {
			if strings.EqualFold(endpoint.Name, "efs") {
				options.BaseEndpoint = awsv2.String(endpoint.URL)
			}
		}
	})

	route53Client, err := awssession.NewRoute53Client(ctx, awssession.EndpointOptions{
		Region:    region,
		Endpoints: metadata.AWS.ServiceEndpoints,
	}, "") // FIXME: Do we need an ARN here?
	if err != nil {
		return nil, fmt.Errorf("failed to create Route53 client: %w", err)
	}

	return &ClusterUninstaller{
		Filters:        filters,
		Region:         region,
		Logger:         logger,
		ClusterID:      metadata.InfraID,
		ClusterDomain:  metadata.AWS.ClusterDomain,
		HostedZoneRole: metadata.AWS.HostedZoneRole,
		endpoints:      metadata.AWS.ServiceEndpoints,
		EC2Client:      ec2Client,
		IAMClient:      iamClient,
		ELBClient:      elbclient,
		ELBV2Client:    elbv2client,
		S3Client:       s3Client,
		EFSClient:      efsClient,
		Route53Client:  route53Client,
	}, nil
}

// validate runs before the uninstall process to ensure that
// all prerequisites are met for a safe destroy.
func (o *ClusterUninstaller) validate(ctx context.Context) error {
	if len(o.Filters) == 0 {
		return errors.Errorf("you must specify at least one tag filter")
	}

	return o.ValidateOwnedSubnets(ctx)
}

// ValidateOwnedSubnets validates whether the subnets owned by the cluster are safe to destroy. That is, the subnets are not currently in use (shared) by other clusters.
// This scenario is a misconfiguration and should not happen, but in practice it did: https://issues.redhat.com//browse/OCPBUGS-60071
// Thus, we add a preflight check to abort the uninstall process in this case to avoid disruptions to other clusters.
func (o *ClusterUninstaller) ValidateOwnedSubnets(ctx context.Context) error {
	o.Logger.Debug("Checking owned subnets for shared tags...")

	subnets := make(map[string]awssession.Subnet, 0)

	// Retrieve the subnet(s) to be destroyed during the uninstall process.
	for _, tags := range o.Filters {
		subnetFilters := make([]ec2v2types.Filter, 0, len(tags))
		for tagKey, tagValue := range tags {
			subnetFilters = append(subnetFilters, ec2v2types.Filter{
				Name:   aws.String(fmt.Sprintf("tag:%s", tagKey)),
				Values: []string{tagValue},
			})
		}

		input := &ec2v2.DescribeSubnetsInput{Filters: subnetFilters}
		paginator := ec2v2.NewDescribeSubnetsPaginator(o.EC2Client, input)
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return fmt.Errorf("failed to describe subnets by tags: %w", err)
			}

			for _, subnet := range page.Subnets {
				id := aws.StringValue(subnet.SubnetId)
				if id == "" {
					continue
				}
				// If the results return the subnet, skip adding.
				if _, ok := subnets[id]; ok {
					continue
				}

				subnets[id] = awssession.Subnet{
					ID:   id,
					Tags: awssession.FromAWSTags(subnet.Tags),
				}
			}
		}
	}

	// The cluster does not own any subnets (i.e. BYO VPC/Subnet use case)
	// so we can skip the check.
	if len(subnets) == 0 {
		o.Logger.Debug("No owned subnets found, skipping validation")
		return nil
	}

	for _, subnet := range subnets {
		// The subnet is marked for deletion but has a shared tag.
		// We abort the uninstall process.
		if subnet.Tags.HasClusterSharedTag() {
			sharedClusterIDs := subnet.Tags.GetClusterIDs(awssession.TagValueShared)

			errMsg := fmt.Sprintf("shared tags found from clusters %v on subnet %s, owned by cluster %s. Destroying cluster %s will delete resources depended on by other clusters, resulting in an outage",
				sharedClusterIDs, subnet.ID, o.ClusterID, o.ClusterID)
			resolveMsg := fmt.Sprintf("To destroy cluster %s, first destroy clusters %v", o.ClusterID, sharedClusterIDs)
			return fmt.Errorf("%s. %s", errMsg, resolveMsg)
		}
	}

	return nil
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	_, err := o.RunWithContext(context.Background())
	return nil, err
}

func createResourceTaggingClientWithConfig(cfg awsv2.Config, region string, endpoints []awstypes.ServiceEndpoint) *resourcegroupstaggingapi.Client {
	return resourcegroupstaggingapi.NewFromConfig(cfg, func(options *resourcegroupstaggingapi.Options) {
		options.Region = region
		for _, endpoint := range endpoints {
			if strings.EqualFold(endpoint.Name, "resourcegroupstaggingapi") {
				options.BaseEndpoint = awsv2.String(endpoint.URL)
			}
		}
	})
}

func createResourceTaggingClient(region string, endpoints []awstypes.ServiceEndpoint) (*resourcegroupstaggingapi.Client, error) {
	cfg, err := awssession.GetConfigWithOptions(context.Background(), configv2.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config for resource tagging client: %w", err)
	}
	return createResourceTaggingClientWithConfig(cfg, region, endpoints), nil
}

// RunWithContext runs the uninstall process with a context.
// The first return is the list of ARNs for resources that could not be destroyed.
func (o *ClusterUninstaller) RunWithContext(ctx context.Context) ([]string, error) {
	err := o.validate(ctx)
	if err != nil {
		return nil, err
	}

	baseTaggingClient, err := createResourceTaggingClient(o.Region, o.endpoints)
	if err != nil {
		return nil, err
	}
	tagClients := []*resourcegroupstaggingapi.Client{baseTaggingClient}

	if o.HostedZoneRole != "" {
		cfg, err := awssession.GetConfigWithOptions(ctx, configv2.WithRegion(endpointUSEast1))
		if err != nil {
			return nil, fmt.Errorf("failed to create AWS config for resource tagging client: %w", err)
		}
		stsSvc, err := awssession.NewSTSClient(ctx, awssession.EndpointOptions{
			Region:    endpointUSEast1,
			Endpoints: o.endpoints,
		}, sts.WithAPIOptions(middleware.AddUserAgentKeyValue("OpenShift/4.x Destroyer", version.Raw)))
		if err != nil {
			return nil, fmt.Errorf("failed to create STS client: %w", err)
		}

		creds := stscreds.NewAssumeRoleProvider(stsSvc, o.HostedZoneRole)
		cfg.Credentials = awsv2.NewCredentialsCache(creds)
		// This client is specifically for finding route53 zones,
		// so it needs to use the global us-east-1 region.

		tagClients = append(tagClients, createResourceTaggingClientWithConfig(cfg, endpointUSEast1, o.endpoints))
	}

	switch o.Region {
	case endpointCNNorth1, endpointCNNorthWest1:
		break
	case endpointISOEast1, endpointISOWest1, endpointISOBEast1:
		break
	case endpointUSGovEast1, endpointUSGovWest1:
		if o.Region != endpointUSGovWest1 {
			tagClient, err := createResourceTaggingClient(endpointUSGovWest1, o.endpoints)
			if err != nil {
				return nil, fmt.Errorf("failed to create resource tagging client for us-gov-west-1: %w", err)
			}
			tagClients = append(tagClients, tagClient)
		}
	default:
		if o.Region != endpointUSEast1 {
			tagClient, err := createResourceTaggingClient(endpointUSEast1, o.endpoints)
			if err != nil {
				return nil, fmt.Errorf("failed to create resource tagging client for default us-east-1: %w", err)
			}
			tagClients = append(tagClients, tagClient)
		}
	}

	iamRoleSearch := &IamRoleSearch{
		Client:  o.IAMClient,
		Filters: o.Filters,
		Logger:  o.Logger,
	}
	iamUserSearch := &IamUserSearch{
		client:  o.IAMClient,
		filters: o.Filters,
		logger:  o.Logger,
	}

	// Get the initial resources to delete, so that they can be returned if the context is canceled while terminating
	// instances.
	deleted := sets.New[string]()
	resourcesToDelete, tagClientsWithResources, err := o.findResourcesToDelete(ctx, tagClients, iamRoleSearch, iamUserSearch, deleted)
	if err != nil {
		o.Logger.WithError(err).Info("error while finding resources to delete")
		if err := ctx.Err(); err != nil {
			return resourcesToDelete.UnsortedList(), err
		}
	}

	tracker := new(ErrorTracker)

	// Terminate EC2 instances. The instances need to be terminated first so that we can ensure that there is nothing
	// running on the cluster creating new resources while we are attempting to delete resources, which could leak
	// the new resources.
	err = o.DeleteEC2Instances(ctx, resourcesToDelete, deleted, tracker)
	if err != nil {
		return resourcesToDelete.UnsortedList(), err
	}

	// Delete the rest of the resources.
	err = wait.PollImmediateUntil(
		time.Second*10,
		func() (done bool, err error) {
			newlyDeleted, loopError := o.DeleteResources(ctx, resourcesToDelete.UnsortedList(), tracker)
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
			nextResourcesToDelete, nextTagClients, err := o.findResourcesToDelete(ctx, tagClientsWithResources, iamRoleSearch, iamUserSearch, deleted)
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

	err = o.removeSharedTags(ctx, tagClients, tracker)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// findUntaggableResources returns the resources for the cluster that cannot be tagged. Any resource that contains
// a shared tag will be ignored.
//
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func (o *ClusterUninstaller) findUntaggableResources(ctx context.Context, deleted sets.Set[string]) (sets.Set[string], error) {
	resources := sets.New[string]()
	o.Logger.Debug("search for IAM instance profiles")
	for _, profileType := range []string{"master", "worker", "bootstrap"} {
		profile := fmt.Sprintf("%s-%s-profile", o.ClusterID, profileType)
		response, err := o.IAMClient.GetInstanceProfile(ctx, &iamv2.GetInstanceProfileInput{InstanceProfileName: &profile})
		if err != nil {
			if strings.Contains(HandleErrorCode(err), "NoSuchEntity") {
				continue
			}
			return resources, fmt.Errorf("failed to get IAM instance profile: %w", err)
		}
		arnString := *response.InstanceProfile.Arn
		if !deleted.Has(arnString) {
			resources.Insert(arnString)
		}
	}
	return resources, nil
}

// findResourcesToDelete returns the resources that should be deleted.
//
// tagClients - clients of the tagging API to use to search for resources.
// deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func (o *ClusterUninstaller) findResourcesToDelete(
	ctx context.Context,
	tagClients []*resourcegroupstaggingapi.Client,
	iamRoleSearch *IamRoleSearch,
	iamUserSearch *IamUserSearch,
	deleted sets.Set[string],
) (sets.Set[string], []*resourcegroupstaggingapi.Client, error) {
	var errs []error
	resources, tagClients, err := FindTaggedResourcesToDelete(ctx, o.Logger, tagClients, o.Filters, iamRoleSearch, iamUserSearch, deleted)
	if err != nil {
		errs = append(errs, err)
	}

	// Find untaggable resources
	untaggableResources, err := o.findUntaggableResources(ctx, deleted)
	if err != nil {
		errs = append(errs, err)
	}
	resources = resources.Union(untaggableResources)

	return resources, tagClients, utilerrors.NewAggregate(errs)
}

// FindTaggedResourcesToDelete returns the tagged resources that should be deleted.
//
//	tagClients - clients of the tagging API to use to search for resources.
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func FindTaggedResourcesToDelete(
	ctx context.Context,
	logger logrus.FieldLogger,
	tagClients []*resourcegroupstaggingapi.Client,
	filters []Filter,
	iamRoleSearch *IamRoleSearch,
	iamUserSearch *IamUserSearch,
	deleted sets.Set[string],
) (sets.Set[string], []*resourcegroupstaggingapi.Client, error) {
	resources := sets.New[string]()
	var tagClientsWithResources []*resourcegroupstaggingapi.Client
	var errs []error

	// Find resources by tag
	for _, tagClient := range tagClients {
		resourcesInTagClient, err := findResourcesByTag(ctx, logger, tagClient, filters, deleted)
		if err != nil {
			errs = append(errs, err)
		}
		resources = resources.Union(resourcesInTagClient)
		// If there are still resources to be deleted for the tag client or if there was an error getting the resources
		// for the tag client, then retain the tag client for future queries.
		if len(resourcesInTagClient) > 0 || err != nil {
			tagClientsWithResources = append(tagClientsWithResources, tagClient)
		} else {
			logger.Debugf("no deletions from %s, removing client", tagClient.Options().Region)
		}
	}

	// Find IAM roles
	if iamRoleSearch != nil {
		iamRoleResources, err := findIAMRoles(ctx, iamRoleSearch, deleted, logger)
		if err != nil {
			errs = append(errs, err)
		}
		resources = resources.Union(iamRoleResources)
	}

	// Find IAM users
	if iamUserSearch != nil {
		iamUserResources, err := findIAMUsers(ctx, iamUserSearch, deleted, logger)
		if err != nil {
			errs = append(errs, err)
		}
		resources = resources.Union(iamUserResources)
	}

	return resources, tagClientsWithResources, utilerrors.NewAggregate(errs)
}

// findResourcesByTag returns the resources with tags that satisfy the filters.
//
//	tagClients - clients of the tagging API to use to search for resources.
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func findResourcesByTag(
	ctx context.Context,
	logger logrus.FieldLogger,
	tagClient *resourcegroupstaggingapi.Client,
	filters []Filter,
	deleted sets.Set[string],
) (sets.Set[string], error) {
	resources := sets.New[string]()
	for _, filter := range filters {
		logger.Debugf("search for matching resources by tag in %s matching %#+v", tagClient.Options().Region, filter)
		tagFilters := make([]tagtypes.TagFilter, 0, len(filter))
		for key, value := range filter {
			tagFilters = append(tagFilters, tagtypes.TagFilter{
				Key:    awsv2.String(key),
				Values: []string{value},
			})
		}

		paginator := resourcegroupstaggingapi.NewGetResourcesPaginator(tagClient, &resourcegroupstaggingapi.GetResourcesInput{TagFilters: tagFilters})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, fmt.Errorf("failed to fetch resources by tag: %w", err)
			}

			for _, resource := range page.ResourceTagMappingList {
				arnString := *resource.ResourceARN
				if !deleted.Has(arnString) {
					resources.Insert(arnString)
				}
			}
		}
	}

	return resources, nil
}

// DeleteResources deletes the specified resources.
//
//	resources - the resources to be deleted.
//
// The first return is the ARNs of the resources that were successfully deleted.
func (o *ClusterUninstaller) DeleteResources(ctx context.Context, resources []string, tracker *ErrorTracker) (sets.Set[string], error) {
	deleted := sets.New[string]()
	for _, arnString := range resources {
		l := o.Logger.WithField("arn", arnString)
		parsedARN, err := arn.Parse(arnString)
		if err != nil {
			l.WithError(err).Debug("could not parse ARN")
			continue
		}
		if err := o.deleteARN(ctx, parsedARN, o.Logger); err != nil {
			tracker.suppressWarning(arnString, err, l)
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

// getPublicHostedZone will find the ID of the non-Terraform-managed public route53 zone given the
// Terraform-managed zone's privateID.
func getPublicHostedZone(ctx context.Context, client *route53.Client, privateID string, logger logrus.FieldLogger) (string, error) {
	response, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{
		Id: awsv2.String(privateID),
	})
	if err != nil {
		return "", err
	}

	privateName := *response.HostedZone.Name

	if response.HostedZone.Config != nil && !response.HostedZone.Config.PrivateZone {
		return "", fmt.Errorf("getPublicHostedZone requires a private ID, but was passed the public %s", privateID)
	} else {
		logger.WithField("hosted zone", privateName).Warn("could not determine whether hosted zone is private")
	}

	return findAncestorPublicRoute53(ctx, client, privateName, logger)
}

// findAncestorPublicRoute53 finds a public route53 zone with the closest ancestor or match to dnsName.
// It returns "", when no public route53 zone could be found.
func findAncestorPublicRoute53(ctx context.Context, client *route53.Client, dnsName string, logger logrus.FieldLogger) (string, error) {
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
func findPublicRoute53(ctx context.Context, client *route53.Client, dnsName string, logger logrus.FieldLogger) (string, error) {
	request := &route53.ListHostedZonesByNameInput{
		DNSName: awsv2.String(dnsName),
	}
	for i := 0; true; i++ {
		logger.Debugf("listing AWS hosted zones %q (page %d)", dnsName, i)
		list, err := client.ListHostedZonesByName(ctx, request)
		if err != nil {
			return "", err
		}

		for _, zone := range list.HostedZones {
			if *zone.Name != dnsName {
				// No name after this can match dnsName
				return "", nil
			}
			if zone.Config == nil {
				logger.WithField("hosted zone", *zone.Name).Warn("could not determine whether hosted zone is private")
				continue
			}
			if !zone.Config.PrivateZone {
				return *zone.Id, nil
			}
		}

		if list.IsTruncated && *list.NextDNSName == *request.DNSName {
			request.HostedZoneId = list.NextHostedZoneId
			continue
		}

		break
	}
	return "", nil
}

func (o *ClusterUninstaller) deleteARN(ctx context.Context, arn arn.ARN, logger logrus.FieldLogger) error {
	switch arn.Service {
	case "ec2":
		return o.deleteEC2(ctx, arn, logger)
	case "elasticloadbalancing":
		return o.deleteElasticLoadBalancing(ctx, arn, logger)
	case "iam":
		return o.deleteIAM(ctx, o.IAMClient, arn, logger)
	case "route53":
		return deleteRoute53(ctx, o.Route53Client, arn, logger)
	case "s3":
		return deleteS3(ctx, o.S3Client, arn, logger)
	case "elasticfilesystem":
		return deleteElasticFileSystem(ctx, o.EFSClient, arn, logger)
	default:
		return errors.Errorf("unrecognized ARN service %s (%s)", arn.Service, arn)
	}
}

func deleteRoute53(ctx context.Context, client *route53.Client, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	if resourceType != "hostedzone" {
		return errors.Errorf("unrecognized Route 53 resource type %s", resourceType)
	}

	publicZoneID, err := getPublicHostedZone(ctx, client, id, logger)
	if err != nil {
		// In some cases AWS may return the zone in the list of tagged resources despite the fact
		// it no longer exists.
		if strings.Contains(HandleErrorCode(err), "NoSuchHostedZone") {
			return nil
		}
		return err
	}

	recordSetKey := func(recordSet route53types.ResourceRecordSet) string {
		return fmt.Sprintf("%s %s", recordSet.Type, *recordSet.Name)
	}

	publicEntries := map[string]route53types.ResourceRecordSet{}
	if len(publicZoneID) != 0 {

		paginator := route53.NewListResourceRecordSetsPaginator(client, &route53.ListResourceRecordSetsInput{HostedZoneId: awsv2.String(publicZoneID)})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return fmt.Errorf("listing record sets for public zone: %w", err)
			}
			for _, recordSet := range page.ResourceRecordSets {
				key := recordSetKey(recordSet)
				publicEntries[key] = recordSet
			}
		}
	} else {
		logger.Debug("shared public zone not found")
	}

	var lastError error
	paginator := route53.NewListResourceRecordSetsPaginator(client, &route53.ListResourceRecordSetsInput{HostedZoneId: awsv2.String(id)})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("listing record sets: %w", err)
		}

		for _, recordSet := range page.ResourceRecordSets {
			if recordSet.Type == "SOA" || recordSet.Type == "NS" {
				// can't delete SOA and NS types
				continue
			}
			key := recordSetKey(recordSet)
			if publicEntry, ok := publicEntries[key]; ok {
				err := deleteRoute53RecordSet(ctx, client, publicZoneID, &publicEntry, logger.WithField("public zone", publicZoneID))
				if err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = fmt.Errorf("deleting record set %#v from public zone %s: %w", publicEntry, publicZoneID, err)
				}
				// do not delete the record set in the private zone if the delete failed in the public zone;
				// otherwise the record set in the public zone will get leaked
				continue
			}

			err = deleteRoute53RecordSet(ctx, client, id, &recordSet, logger)
			if err != nil {
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = fmt.Errorf("deleting record set %#+v from zone %s: %w", recordSet, id, err)
			}
		}
	}

	if lastError != nil {
		return lastError
	}

	_, err = client.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{
		Id: awsv2.String(id),
	})
	if err != nil {
		if strings.Contains(HandleErrorCode(err), "NoSuchHostedZone") {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteRoute53RecordSet(ctx context.Context, client *route53.Client, zoneID string, recordSet *route53types.ResourceRecordSet, logger logrus.FieldLogger) error {
	logger = logger.WithField("record set", fmt.Sprintf("%s %s", recordSet.Type, *recordSet.Name))
	_, err := client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: awsv2.String(zoneID),
		ChangeBatch: &route53types.ChangeBatch{
			Changes: []route53types.Change{
				{
					Action:            route53types.ChangeActionDelete,
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

func deleteS3(ctx context.Context, client *s3.Client, arn arn.ARN, logger logrus.FieldLogger) error {
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket:  awsv2.String(arn.Resource),
		MaxKeys: awsv2.Int32(1000),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list objects in bucket %s: %w", arn.Resource, err)
		}

		var objects []s3types.ObjectIdentifier
		for _, object := range page.Contents {
			objects = append(objects, s3types.ObjectIdentifier{
				Key: object.Key,
			})
		}

		if len(objects) > 0 {
			if _, err = client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: awsv2.String(arn.Resource),
				Delete: &s3types.Delete{
					Objects: objects,
				},
			}); err != nil {
				return fmt.Errorf("failed to delete objects in bucket %s: %w", arn.Resource, err)
			}
		}
	}

	if _, err := client.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: awsv2.String(arn.Resource)}); err != nil {
		var noSuckBucket *s3types.NoSuchBucket
		if errors.As(err, &noSuckBucket) {
			logrus.Debugf("bucket %s already deleted", arn.Resource)
			return nil
		}
		return fmt.Errorf("failed to delete bucket %s: %w", arn.Resource, err)
	}

	logger.Info("Deleted")
	return nil
}

func deleteElasticFileSystem(ctx context.Context, client *efs.Client, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}

	switch resourceType {
	case "file-system":
		return deleteFileSystem(ctx, client, id, logger)
	case "access-point":
		return deleteAccessPoint(ctx, client, id, logger)
	default:
		return errors.Errorf("unrecognized elastic file system resource type %s", resourceType)
	}
}

func deleteFileSystem(ctx context.Context, client *efs.Client, fsid string, logger logrus.FieldLogger) error {
	logger = logger.WithField("Elastic FileSystem ID", fsid)

	// Delete all MountTargets + AccessPoints under given FS ID
	mountTargetIDs, err := getMountTargets(ctx, client, fsid)
	if err != nil {
		return err
	}
	for _, mt := range mountTargetIDs {
		err := deleteMountTarget(ctx, client, mt, logger)
		if err != nil {
			return err
		}
	}
	accessPointIDs, err := getAccessPoints(ctx, client, fsid)
	if err != nil {
		return err
	}
	for _, ap := range accessPointIDs {
		err := deleteAccessPoint(ctx, client, ap, logger)
		if err != nil {
			return err
		}
	}

	_, err = client.DeleteFileSystem(ctx, &efs.DeleteFileSystemInput{FileSystemId: awsv2.String(fsid)})
	if err != nil {
		if strings.Contains(HandleErrorCode(err), "FileSystemNotFound") {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func getAccessPoints(ctx context.Context, client *efs.Client, apID string) ([]string, error) {
	var accessPointIDs []string
	paginator := efs.NewDescribeAccessPointsPaginator(client, &efs.DescribeAccessPointsInput{FileSystemId: awsv2.String(apID)})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("describing access points: %w", err)
		}

		for _, ap := range page.AccessPoints {
			apName := ap.AccessPointId
			if apName == nil {
				continue
			}
			accessPointIDs = append(accessPointIDs, *apName)
		}
	}

	return accessPointIDs, nil
}

func getMountTargets(ctx context.Context, client *efs.Client, fsid string) ([]string, error) {
	var mountTargetIDs []string
	// There is no DescribeMountTargetsPagesWithContext.
	// Number of Mount Targets should be equal to nr. of subnets that can access the volume, i.e. relatively small.
	rsp, err := client.DescribeMountTargets(
		ctx,
		&efs.DescribeMountTargetsInput{FileSystemId: awsv2.String(fsid)},
	)
	if err != nil {
		return nil, err
	}

	for _, mt := range rsp.MountTargets {
		mtName := mt.MountTargetId
		if mtName == nil {
			continue
		}
		mountTargetIDs = append(mountTargetIDs, *mtName)
	}

	return mountTargetIDs, nil
}

func deleteAccessPoint(ctx context.Context, client *efs.Client, id string, logger logrus.FieldLogger) error {
	logger = logger.WithField("AccessPoint ID", id)
	_, err := client.DeleteAccessPoint(ctx, &efs.DeleteAccessPointInput{AccessPointId: awsv2.String(id)})
	if err != nil {
		if strings.Contains(HandleErrorCode(err), "AccessPointNotFound") {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteMountTarget(ctx context.Context, client *efs.Client, id string, logger logrus.FieldLogger) error {
	logger = logger.WithField("Mount Target ID", id)
	_, err := client.DeleteMountTarget(ctx, &efs.DeleteMountTargetInput{MountTargetId: awsv2.String(id)})
	if err != nil {
		if strings.Contains(HandleErrorCode(err), "MountTargetNotFound") {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}
