package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	tagtypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	awssession "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/version"
)

func (o *ClusterUninstaller) removeSharedTags(
	ctx context.Context,
	tagClients []*resourcegroupstaggingapi.Client,
	tracker *ErrorTracker,
) error {
	for _, key := range o.clusterOwnedKeys() {
		if err := o.removeSharedTag(ctx, tagClients, key, tracker); err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) clusterOwnedKeys() []string {
	var keys []string
	for _, filter := range o.Filters {
		for key, value := range filter {
			if !strings.HasPrefix(key, "kubernetes.io/cluster/") {
				continue
			}
			if value != "owned" {
				o.Logger.Warnf("Ignoring non-owned cluster key %s: %s for shared-tag removal", key, value)
			}
			keys = append(keys, key)
		}
	}
	return keys
}

func (o *ClusterUninstaller) removeSharedTag(ctx context.Context, tagClients []*resourcegroupstaggingapi.Client, key string, tracker *ErrorTracker) error {
	const sharedValue = "shared"

	request := &resourcegroupstaggingapi.UntagResourcesInput{
		TagKeys: []string{key},
	}

	removed := map[string]struct{}{}
	tagClients = append([]*resourcegroupstaggingapi.Client(nil), tagClients...)
	for len(tagClients) > 0 {
		nextTagClients := tagClients[:0]
		for _, tagClient := range tagClients {
			var lastErr error
			o.Logger.Debugf("Search for and remove tags in %s matching %s: shared", tagClient.Options().Region, key)
			var arns []string
			paginator := resourcegroupstaggingapi.NewGetResourcesPaginator(tagClient, &resourcegroupstaggingapi.GetResourcesInput{
				TagFilters: []tagtypes.TagFilter{{
					Key:    aws.String(key),
					Values: []string{sharedValue},
				}}})
			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil {
					o.Logger.Debugf("failed to get resources: %v", err)
					lastErr = err
					break
				}

				for _, resource := range page.ResourceTagMappingList {
					arnString := aws.ToString(resource.ResourceARN)
					logger := o.Logger.WithField("arn", arnString)
					parsedARN, err := arn.Parse(arnString)
					if err != nil {
						logger.WithError(err).Debug("could not parse ARN")
						continue
					}
					if _, ok := removed[arnString]; !ok {
						if err := o.cleanSharedARN(ctx, parsedARN, logger); err != nil {
							tracker.suppressWarning(arnString, err, logger)
							if err := ctx.Err(); err != nil {
								lastErr = fmt.Errorf("failed to remove tag %q: %w", key, err)
							}
							continue
						}
						arns = append(arns, arnString)
					}
				}
			}

			if lastErr != nil {
				o.Logger.Infof("failed to get tagged resources: %v", lastErr)
				var invalidParameter *tagtypes.InvalidParameterException
				if errors.As(lastErr, &invalidParameter) {
					continue
				}
				nextTagClients = append(nextTagClients, tagClient)
				continue
			}

			if len(arns) == 0 {
				o.Logger.Debugf("No matches in %s for %s: shared, removing client", tagClient.Options().Region, key)
				continue
			}
			// appending the tag client here but it needs to be removed if there is a InvalidParameterException when trying to
			// untag below since that only leads to an infinite loop error.
			nextTagClients = append(nextTagClients, tagClient)

			for i := 0; i < len(arns); i += 20 {
				request.ResourceARNList = make([]string, 0, 20)
				for j := 0; i+j < len(arns) && j < 20; j++ {
					request.ResourceARNList = append(request.ResourceARNList, arns[i+j])
				}
				result, err := tagClient.UntagResources(ctx, request)
				if err != nil {
					if awssession.IsInvalidParameter(err) {
						nextTagClients = nextTagClients[:len(nextTagClients)-1]
					}
					err = errors.Wrap(err, "untag shared resources")
					o.Logger.Info(err)
					continue
				}
				for _, arn := range request.ResourceARNList {
					if info, failed := result.FailedResourcesMap[arn]; failed {
						o.Logger.WithField("arn", arn).Infof("Failed to remove tag %s: shared; error=%s", key, *info.ErrorMessage)
						continue
					}
					o.Logger.WithField("arn", arn).Infof("Removed tag %s: shared", key)
					removed[arn] = exists
				}
			}
		}
		tagClients = nextTagClients
	}

	iamRoleSearch := &IamRoleSearch{
		Client:  o.IAMClient,
		Filters: []Filter{{key: sharedValue}},
		Logger:  o.Logger,
	}
	o.Logger.Debugf("Search for and remove shared tags for IAM roles matching %s: shared", key)
	if err := wait.PollImmediateUntil(
		time.Second*10,
		func() (bool, error) {
			_, sharedRoles, err := iamRoleSearch.find(ctx)
			if err != nil {
				o.Logger.Infof("Could not search for shared IAM roles: %v", err)
				return false, nil
			}
			done := true
			for _, role := range sharedRoles {
				o.Logger.Debugf("Removing the shared tag from the %q IAM role", role)
				input := &iam.UntagRoleInput{
					RoleName: &role,
					TagKeys:  []string{key},
				}
				if _, err := o.IAMClient.UntagRole(ctx, input); err != nil {
					done = false
					o.Logger.Infof("Could not remove the shared tag from the %q IAM role: %v", role, err)
				}
			}
			return done, nil
		},
		ctx.Done(),
	); err != nil {
		return errors.Wrap(err, "problem removing shared tags from IAM roles")
	}

	return nil
}

func (o *ClusterUninstaller) cleanSharedARN(ctx context.Context, arn arn.ARN, logger logrus.FieldLogger) error {
	switch service := arn.Service; service {
	case "route53":
		return o.cleanSharedRoute53(ctx, arn, logger)
	default:
		logger.Debugf("Nothing to clean for shared %s resource", service)
		return nil
	}
}

func (o *ClusterUninstaller) cleanSharedRoute53(ctx context.Context, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "hostedzone":
		return o.cleanSharedHostedZone(ctx, id, logger)
	default:
		logger.Debugf("Nothing to clean for shared %s resource", resourceType)
		return nil
	}
}

func (o *ClusterUninstaller) cleanSharedHostedZone(ctx context.Context, id string, logger logrus.FieldLogger) error {
	// The private hosted zone (phz) may belong to a different account,
	// in which case we need a separate client.
	// Note: the ClusterUninstaller has a basic Route53 client used for public zone resources.
	privateZoneClient, err := awssession.NewRoute53Client(ctx,
		awssession.EndpointOptions{Region: o.Region, Endpoints: o.endpoints}, o.HostedZoneRole,
		route53.WithAPIOptions(awsmiddleware.AddUserAgentKeyValue(awssession.OpenShiftInstallerDestroyerUserAgent, version.Raw)),
	)
	if err != nil {
		return fmt.Errorf("failed to create Route53 private zone client: %w", err)
	}

	if o.ClusterDomain == "" {
		logger.Debug("No cluster domain specified in metadata; cannot clean the shared hosted zone")
		return nil
	}
	dottedClusterDomain := o.ClusterDomain + "."

	publicZoneID, err := findAncestorPublicRoute53(ctx, o.Route53Client, dottedClusterDomain, logger)
	if err != nil {
		return err
	}

	var lastError error
	paginator := route53.NewListResourceRecordSetsPaginator(privateZoneClient, &route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(id)})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error listing shared hosted zone records: %w", err)
		}

		for _, recordSet := range page.ResourceRecordSets {
			// skip record sets that are not part of the cluster
			name := aws.ToString(recordSet.Name)
			if !strings.HasSuffix(name, dottedClusterDomain) {
				continue
			}
			if len(name) == len(dottedClusterDomain) {
				continue
			}
			recordsetFields := logrus.Fields{"recordset": fmt.Sprintf("%s (%s)", aws.ToString(recordSet.Name), recordSet.Type)}
			// delete any matching record sets in the public hosted zone
			if publicZoneID != "" {
				publicZoneLogger := logger.WithField("id", publicZoneID)
				if err := deleteMatchingRecordSetInPublicZone(ctx, o.Route53Client, publicZoneID, &recordSet, publicZoneLogger); err != nil {
					if lastError != nil {
						publicZoneLogger.Debug(lastError)
					}
					lastError = fmt.Errorf("deleting record set matching %#v from public zone %s: %w", recordSet, publicZoneID, err)
					// do not delete the record set in the private zone if the delete failed in the public zone;
					// otherwise the record set in the public zone will get leaked
					continue
				}
				publicZoneLogger.WithFields(recordsetFields).Debug("Deleted from public zone")
			}
			// delete the record set
			if err := deleteRoute53RecordSet(ctx, privateZoneClient, id, &recordSet, logger); err != nil {
				if lastError != nil {
					logger.Debug(lastError)
				}
				lastError = fmt.Errorf("deleting record set %#v from zone %s: %w", recordSet, id, err)
			}
			logger.WithFields(recordsetFields).Debug("Deleted from public zone")
		}
	}

	if lastError != nil {
		return lastError
	}
	logger.Info("Cleaned record sets from hosted zone")
	return nil
}

func deleteMatchingRecordSetInPublicZone(ctx context.Context, client *route53.Client, zoneID string, recordSet *route53types.ResourceRecordSet, logger logrus.FieldLogger) error {
	in := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID),
		MaxItems:        aws.Int32(1),
		StartRecordName: recordSet.Name,
		StartRecordType: recordSet.Type,
	}
	out, err := client.ListResourceRecordSets(ctx, in)
	if err != nil {
		return err
	}
	if len(out.ResourceRecordSets) == 0 {
		return nil
	}
	matchingRecordSet := out.ResourceRecordSets[0]
	if aws.ToString(matchingRecordSet.Name) != aws.ToString(recordSet.Name) || matchingRecordSet.Type != recordSet.Type {
		return nil
	}
	return deleteRoute53RecordSet(ctx, client, zoneID, &matchingRecordSet, logger)
}
