package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (o *ClusterUninstaller) removeSharedTags(
	ctx context.Context,
	session *session.Session,
	tagClients []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI,
	tracker *errorTracker,
) error {
	for _, key := range o.clusterOwnedKeys() {
		if err := o.removeSharedTag(ctx, session, tagClients, key, tracker); err != nil {
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

func (o *ClusterUninstaller) removeSharedTag(ctx context.Context, session *session.Session, tagClients []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI, key string, tracker *errorTracker) error {
	const sharedValue = "shared"

	request := &resourcegroupstaggingapi.UntagResourcesInput{
		TagKeys: []*string{aws.String(key)},
	}

	removed := map[string]struct{}{}
	tagClients = append([]*resourcegroupstaggingapi.ResourceGroupsTaggingAPI(nil), tagClients...)
	for len(tagClients) > 0 {
		nextTagClients := tagClients[:0]
		for _, tagClient := range tagClients {
			o.Logger.Debugf("Search for and remove tags in %s matching %s: shared", *tagClient.Config.Region, key)
			var arns []string
			err := tagClient.GetResourcesPagesWithContext(
				ctx,
				&resourcegroupstaggingapi.GetResourcesInput{TagFilters: []*resourcegroupstaggingapi.TagFilter{{
					Key:    aws.String(key),
					Values: []*string{aws.String(sharedValue)},
				}}},
				func(results *resourcegroupstaggingapi.GetResourcesOutput, lastPage bool) bool {
					for _, resource := range results.ResourceTagMappingList {
						arnString := aws.StringValue(resource.ResourceARN)
						logger := o.Logger.WithField("arn", arnString)
						parsedARN, err := arn.Parse(arnString)
						if err != nil {
							logger.WithError(err).Debug("could not parse ARN")
							continue
						}
						if _, ok := removed[arnString]; !ok {
							if err := o.cleanSharedARN(ctx, session, parsedARN, logger); err != nil {
								tracker.suppressWarning(arnString, err, logger)
								if err := ctx.Err(); err != nil {
									return false
								}
								continue
							}
							arns = append(arns, arnString)
						}
					}

					return !lastPage
				},
			)
			if err != nil {
				err2 := errors.Wrap(err, "get tagged resources")
				o.Logger.Info(err2)
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case resourcegroupstaggingapi.ErrorCodeInvalidParameterException:
						continue
					}
				}
				nextTagClients = append(nextTagClients, tagClient)
				continue
			}
			if len(arns) == 0 {
				o.Logger.Debugf("No matches in %s for %s: shared, removing client", *tagClient.Config.Region, key)
				continue
			}
			// appending the tag client here but it needs to be removed if there is a InvalidParameterException when trying to
			// untag below since that only leads to an infinite loop error.
			nextTagClients = append(nextTagClients, tagClient)

			for i := 0; i < len(arns); i += 20 {
				request.ResourceARNList = make([]*string, 0, 20)
				for j := 0; i+j < len(arns) && j < 20; j++ {
					request.ResourceARNList = append(request.ResourceARNList, aws.String(arns[i+j]))
				}
				result, err := tagClient.UntagResourcesWithContext(ctx, request)
				if err != nil {
					var awsErr awserr.Error
					ok := errors.As(err, &awsErr)
					if ok && awsErr.Code() == resourcegroupstaggingapi.ErrorCodeInvalidParameterException {
						nextTagClients = nextTagClients[:len(nextTagClients)-1]
					}
					err = errors.Wrap(err, "untag shared resources")
					o.Logger.Info(err)
					continue
				}
				for _, arn := range request.ResourceARNList {
					if info, failed := result.FailedResourcesMap[*arn]; failed {
						o.Logger.WithField("arn", *arn).Infof("Failed to remove tag %s: shared; error=%s", key, *info.ErrorMessage)
						continue
					}
					o.Logger.WithField("arn", *arn).Infof("Removed tag %s: shared", key)
					removed[*arn] = exists
				}
			}
		}
		tagClients = nextTagClients
	}

	iamClient := iam.New(session)
	iamRoleSearch := &iamRoleSearch{
		client:  iamClient,
		filters: []Filter{{key: sharedValue}},
		logger:  o.Logger,
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
					TagKeys:  []*string{&key},
				}
				if _, err := iamClient.UntagRoleWithContext(ctx, input); err != nil {
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

func (o *ClusterUninstaller) cleanSharedARN(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	switch service := arn.Service; service {
	case "route53":
		return o.cleanSharedRoute53(ctx, session, arn, logger)
	default:
		logger.Debugf("Nothing to clean for shared %s resource", service)
		return nil
	}
}

func (o *ClusterUninstaller) cleanSharedRoute53(ctx context.Context, session *session.Session, arn arn.ARN, logger logrus.FieldLogger) error {
	resourceType, id, err := splitSlash("resource", arn.Resource)
	if err != nil {
		return err
	}
	logger = logger.WithField("id", id)

	switch resourceType {
	case "hostedzone":
		return o.cleanSharedHostedZone(ctx, session, id, logger)
	default:
		logger.Debugf("Nothing to clean for shared %s resource", resourceType)
		return nil
	}
}

func (o *ClusterUninstaller) cleanSharedHostedZone(ctx context.Context, session *session.Session, id string, logger logrus.FieldLogger) error {
	// The private hosted zone (phz) may belong to a different account,
	// in which case we need a separate client.
	publicZoneClient := route53.New(session)
	privateZoneClient := route53.New(session)
	if o.HostedZoneRole != "" {
		creds := stscreds.NewCredentials(session, o.HostedZoneRole)
		privateZoneClient = route53.New(session, &aws.Config{Credentials: creds})
		logger.Infof("Assuming role %s to destroy records in private hosted zone", o.HostedZoneRole)
	}

	if o.ClusterDomain == "" {
		logger.Debug("No cluster domain specified in metadata; cannot clean the shared hosted zone")
		return nil
	}
	dottedClusterDomain := o.ClusterDomain + "."

	publicZoneID, err := findAncestorPublicRoute53(ctx, publicZoneClient, dottedClusterDomain, logger)
	if err != nil {
		return err
	}

	var lastError error
	err = privateZoneClient.ListResourceRecordSetsPagesWithContext(
		ctx,
		&route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(id)},
		func(results *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			for _, recordSet := range results.ResourceRecordSets {
				// skip record sets that are not part of the cluster
				name := aws.StringValue(recordSet.Name)
				if !strings.HasSuffix(name, dottedClusterDomain) {
					continue
				}
				if len(name) == len(dottedClusterDomain) {
					continue
				}
				recordsetFields := logrus.Fields{"recordset": fmt.Sprintf("%s (%s)", aws.StringValue(recordSet.Name), aws.StringValue(recordSet.Type))}
				// delete any matching record sets in the public hosted zone
				if publicZoneID != "" {
					publicZoneLogger := logger.WithField("id", publicZoneID)
					if err := deleteMatchingRecordSetInPublicZone(ctx, publicZoneClient, publicZoneID, recordSet, publicZoneLogger); err != nil {
						if lastError != nil {
							publicZoneLogger.Debug(lastError)
						}
						lastError = errors.Wrapf(err, "deleting record set matching %#v from public zone %s", recordSet, publicZoneID)
						// do not delete the record set in the private zone if the delete failed in the public zone;
						// otherwise the record set in the public zone will get leaked
						continue
					}
					publicZoneLogger.WithFields(recordsetFields).Debug("Deleted from public zone")
				}
				// delete the record set
				if err := deleteRoute53RecordSet(ctx, privateZoneClient, id, recordSet, logger); err != nil {
					if lastError != nil {
						logger.Debug(lastError)
					}
					lastError = errors.Wrapf(err, "deleting record set %#v from zone %s", recordSet, id)
				}
				logger.WithFields(recordsetFields).Debug("Deleted from public zone")
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

	logger.Info("Cleaned record sets from hosted zone")
	return nil
}

func deleteMatchingRecordSetInPublicZone(ctx context.Context, client *route53.Route53, zoneID string, recordSet *route53.ResourceRecordSet, logger logrus.FieldLogger) error {
	in := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID),
		MaxItems:        aws.String("1"),
		StartRecordName: recordSet.Name,
		StartRecordType: recordSet.Type,
	}
	out, err := client.ListResourceRecordSetsWithContext(ctx, in)
	if err != nil {
		return err
	}
	if len(out.ResourceRecordSets) == 0 {
		return nil
	}
	matchingRecordSet := out.ResourceRecordSets[0]
	if aws.StringValue(matchingRecordSet.Name) != aws.StringValue(recordSet.Name) ||
		aws.StringValue(matchingRecordSet.Type) != aws.StringValue(recordSet.Type) {
		return nil
	}
	return deleteRoute53RecordSet(ctx, client, zoneID, matchingRecordSet, logger)
}
