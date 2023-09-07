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
	"github.com/aws/aws-sdk-go/service/efs"
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
		Filters:        filters,
		Region:         region,
		Logger:         logger,
		ClusterID:      metadata.InfraID,
		ClusterDomain:  metadata.AWS.ClusterDomain,
		Session:        session,
		HostedZoneRole: metadata.AWS.HostedZoneRole,
	}, nil
}

func (o *ClusterUninstaller) validate() error {
	if len(o.Filters) == 0 {
		return errors.Errorf("you must specify at least one tag filter")
	}
	return nil
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	_, err := o.RunWithContext(context.Background())
	return nil, err
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

	if o.HostedZoneRole != "" {
		cfg := awssession.GetR53ClientCfg(awsSession, o.HostedZoneRole)
		// This client is specifically for finding route53 zones,
		// so it needs to use the global us-east-1 region.
		cfg.Region = aws.String(endpoints.UsEast1RegionID)
		tagClients = append(tagClients, resourcegroupstaggingapi.New(awsSession, cfg))
	}

	switch o.Region {
	case endpoints.CnNorth1RegionID, endpoints.CnNorthwest1RegionID:
		break
	case endpoints.UsIsoEast1RegionID, endpoints.UsIsoWest1RegionID, endpoints.UsIsobEast1RegionID:
		break
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
			instancesRunning, instancesNotTerminated, err := findEC2Instances(ctx, ec2Client, deleted, o.Filters, o.Logger)
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

// findResourcesToDelete returns the resources that should be deleted.
//
//	tagClients - clients of the tagging API to use to search for resources.
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
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
	iamRoleResources, err := findIAMRoles(ctx, iamRoleSearch, deleted, o.Logger)
	if err != nil {
		errs = append(errs, err)
	}
	resources = resources.Union(iamRoleResources)

	// Find IAM users
	iamUserResources, err := findIAMUsers(ctx, iamUserSearch, deleted, o.Logger)
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
//
//	tagClients - clients of the tagging API to use to search for resources.
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
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

// findUntaggableResources returns the resources for the cluster that cannot be tagged.
//
//	deleted - the resources that have already been deleted. Any resources specified in this set will be ignored.
func (o *ClusterUninstaller) findUntaggableResources(ctx context.Context, iamClient *iam.IAM, deleted sets.String) (sets.String, error) {
	resources := sets.NewString()
	o.Logger.Debug("search for IAM instance profiles")
	for _, profileType := range []string{"master", "worker", "bootstrap"} {
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
//
//	resources - the resources to be deleted.
//
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
	case "elasticfilesystem":
		return deleteElasticFileSystem(ctx, session, arn, logger)
	default:
		return errors.Errorf("unrecognized ARN service %s (%s)", arn.Service, arn)
	}
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

func deleteElasticFileSystem(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	client := efs.New(session)

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

func deleteFileSystem(ctx context.Context, client *efs.EFS, fsid string, logger logrus.FieldLogger) error {
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

	_, err = client.DeleteFileSystemWithContext(ctx, &efs.DeleteFileSystemInput{FileSystemId: aws.String(fsid)})
	if err != nil {
		if err.(awserr.Error).Code() == efs.ErrCodeFileSystemNotFound {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}

func getAccessPoints(ctx context.Context, client *efs.EFS, apID string) ([]string, error) {
	var accessPointIDs []string
	err := client.DescribeAccessPointsPagesWithContext(
		ctx,
		&efs.DescribeAccessPointsInput{FileSystemId: aws.String(apID)},
		func(page *efs.DescribeAccessPointsOutput, lastPage bool) bool {
			for _, ap := range page.AccessPoints {
				apName := ap.AccessPointId
				if apName == nil {
					continue
				}
				accessPointIDs = append(accessPointIDs, *apName)
			}
			return !lastPage

		},
	)
	if err != nil {
		return nil, err
	}
	return accessPointIDs, nil
}

func getMountTargets(ctx context.Context, client *efs.EFS, fsid string) ([]string, error) {
	var mountTargetIDs []string
	// There is no DescribeMountTargetsPagesWithContext.
	// Number of Mount Targets should be equal to nr. of subnets that can access the volume, i.e. relatively small.
	rsp, err := client.DescribeMountTargetsWithContext(
		ctx,
		&efs.DescribeMountTargetsInput{FileSystemId: aws.String(fsid)},
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

func deleteAccessPoint(ctx context.Context, client *efs.EFS, id string, logger logrus.FieldLogger) error {
	logger = logger.WithField("AccessPoint ID", id)
	_, err := client.DeleteAccessPointWithContext(ctx, &efs.DeleteAccessPointInput{AccessPointId: aws.String(id)})
	if err != nil {
		if err.(awserr.Error).Code() == efs.ErrCodeAccessPointNotFound {
			return nil
		}

		return err
	}

	logger.Info("Deleted")
	return nil
}

func deleteMountTarget(ctx context.Context, client *efs.EFS, id string, logger logrus.FieldLogger) error {
	logger = logger.WithField("Mount Target ID", id)
	_, err := client.DeleteMountTargetWithContext(ctx, &efs.DeleteMountTargetInput{MountTargetId: aws.String(id)})
	if err != nil {
		if err.(awserr.Error).Code() == efs.ErrCodeMountTargetNotFound {
			return nil
		}
		return err
	}

	logger.Info("Deleted")
	return nil
}
