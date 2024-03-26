package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"

	awsdestroy "github.com/openshift/installer/pkg/destroy/aws"
)

type destroyInputOptions struct {
	infraID          string
	region           string
	ignitionBucket   string
	preserveIgnition bool
}

func destroyBootstrapResources(ctx context.Context, logger logrus.FieldLogger, awsSession *session.Session, input *destroyInputOptions) error {
	iamClient := iam.New(awsSession)
	iamRoleSearch := &awsdestroy.IamRoleSearch{
		Client: iamClient,
		Filters: []awsdestroy.Filter{
			{
				"Name":                         fmt.Sprintf("%s-bootstrap-role", input.infraID),
				clusterOwnedTag(input.infraID): ownedTagValue,
			},
		},
		Logger: logger,
	}

	tagClients := []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI{
		resourcegroupstaggingapi.New(awsSession),
	}
	switch input.region {
	case endpoints.CnNorth1RegionID, endpoints.CnNorthwest1RegionID:
		break
	case endpoints.UsIsoEast1RegionID, endpoints.UsIsoWest1RegionID, endpoints.UsIsobEast1RegionID:
		break
	case endpoints.UsGovEast1RegionID:
		tagClients = append(tagClients, resourcegroupstaggingapi.New(awsSession, aws.NewConfig().WithRegion(endpoints.UsGovWest1RegionID)))
	case endpoints.UsGovWest1RegionID:
		break
	case endpoints.UsEast1RegionID:
		break
	default:
		tagClients = append(tagClients, resourcegroupstaggingapi.New(awsSession, aws.NewConfig().WithRegion(endpoints.UsEast1RegionID)))
	}

	entities := []string{"bootstrap", "bootstrap-profile", "bootstrap-role-policy", "bootstrap-sg"}
	filters := make([]awsdestroy.Filter, 0, len(entities))
	for _, entity := range entities {
		filters = append(filters, awsdestroy.Filter{
			"Name":                         fmt.Sprintf("%s-%s", input.infraID, entity),
			clusterOwnedTag(input.infraID): ownedTagValue,
		})
	}

	deleted := sets.New[string]()
	// Get the initial resources to delete, so that they can be returned if the context is cancelled while terminating instances.
	resourcesToDelete, tagClientsWithResources, err := awsdestroy.FindTaggedResourcesToDelete(ctx, logger, tagClients, filters, iamRoleSearch, nil, deleted)
	if err != nil {
		return fmt.Errorf("failed to collect bootstrap resources to delete: %w", err)
	}

	resourcesToPreserve := sets.New[string]()
	if input.preserveIgnition {
		logger.Debugln("Preserving ignition resources")
		var errs []error
		for _, resource := range resourcesToDelete.UnsortedList() {
			arn, err := arn.Parse(resource)
			if err != nil {
				// We don't care if we failed to parse some resource ARNs as
				// long as we are able to find the ignition bucket and object
				// we are looking for
				errs = append(errs, err)
				continue
			}
			if arn.Service != "s3" {
				continue
			}
			bucketName, objectName, objectFound := strings.Cut(arn.Resource, "/")
			if bucketName != input.ignitionBucket {
				continue
			}
			if !objectFound || objectName == ignitionKey {
				resourcesToPreserve.Insert(resource)
			}
		}
		// Should contain at least bucket
		if resourcesToPreserve.Len() < 1 {
			errMsg := "failed to find ignition resources to preserve"
			if len(errs) > 0 {
				return fmt.Errorf("%s: %w", errMsg, utilerrors.NewAggregate(errs))
			}
			return fmt.Errorf("%s", errMsg)
		}
		resourcesToDelete = resourcesToDelete.Difference(resourcesToPreserve)
		// Pretend the ignition objects have already been deleted to avoid them
		// being added again to the to-delete list
		deleted = deleted.Union(resourcesToPreserve)
	}

	tracker := new(awsdestroy.ErrorTracker)
	instanceFilters := []awsdestroy.Filter{
		{
			"Name":                         fmt.Sprintf("%s-bootstrap", input.infraID),
			clusterOwnedTag(input.infraID): ownedTagValue,
		},
	}
	err = awsdestroy.DeleteEC2Instances(ctx, logger, awsSession, instanceFilters, resourcesToDelete, deleted, tracker)
	if err != nil {
		logger.WithError(err).Infof("failed to delete the following resources: %v", resourcesToDelete.UnsortedList())
		return fmt.Errorf("failed to terminate bootstrap instance: %w", err)
	}

	// Delete the rest of the resources
	err = wait.PollUntilContextCancel(
		ctx,
		time.Second*10,
		true,
		func(ctx context.Context) (bool, error) {
			newlyDeleted, loopError := awsdestroy.DeleteResources(ctx, logger, awsSession, resourcesToDelete.UnsortedList(), tracker)
			// Delete from the resources-to-delete set so that the current
			// state of the resources to delete can be returned if the context
			// is completed.
			resourcesToDelete = resourcesToDelete.Difference(newlyDeleted)
			deleted = deleted.Union(newlyDeleted)
			if loopError != nil {
				if err := ctx.Err(); err != nil {
					return false, err
				}
			}
			// Store resources to delete in a temporary variable so that, in
			// case the context is cancelled, the current resources to delete
			// are not lost.
			nextResourcesToDelete, nextTagClients, err := awsdestroy.FindTaggedResourcesToDelete(ctx, logger, tagClientsWithResources, filters, iamRoleSearch, nil, deleted)
			if err != nil {
				logger.WithError(err).Info("error while finding resources to delete")
				if err := ctx.Err(); err != nil {
					return false, err
				}
				loopError = fmt.Errorf("error while finding resources to delete: %w", err)
			}
			resourcesToDelete = nextResourcesToDelete
			tagClientsWithResources = nextTagClients
			return len(resourcesToDelete) == 0 && loopError == nil, nil
		},
	)
	if err != nil {
		logger.WithError(err).Infof("failed to delete the following resources: %v", resourcesToDelete.UnsortedList())
		return fmt.Errorf("failed to delete bootstrap resources: %w", err)
	}

	logger.Debugf("Preserving the following resources: %v", resourcesToPreserve.UnsortedList())

	return nil
}
