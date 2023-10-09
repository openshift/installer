package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"

	awssession "github.com/openshift/installer/pkg/asset/installconfig/aws"
	awsdestroy "github.com/openshift/installer/pkg/destroy/aws"
	"github.com/openshift/installer/pkg/tfvars"
	tfvarsaws "github.com/openshift/installer/pkg/tfvars/aws"
	typesaws "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/version"
)

func normalAWSDestroy(a AWSInfraProvider, directory string, varFiles []string) error {
	logger := logrus.StandardLogger()

	// Unmarshall input from tf variables, so we can use it along with installConfig and other assets
	// as the contractual input regardless off the implementation.
	clusterConfig := &tfvars.Config{}
	clusterAWSConfig := &tfvarsaws.Config{}
	for _, file := range varFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if strings.HasSuffix(file, "terraform.tfvars.json") {
			//if file.Filename == "terraform.tfvars.json" {
			if err := json.Unmarshal(data, clusterConfig); err != nil {
				return err
			}
		}

		if strings.HasSuffix(file, "terraform.platform.auto.tfvars.json") {
			// if file.Filename == "terraform.platform.auto.tfvars.json" {
			if err := json.Unmarshal(data, clusterAWSConfig); err != nil {
				return err
			}
		}
	}

	// FIXME: we are assuming the above files were found. If they were not, we
	// will trigger segfaults
	eps := []typesaws.ServiceEndpoint{}
	for k, v := range clusterAWSConfig.CustomEndpoints {
		eps = append(eps, typesaws.ServiceEndpoint{Name: k, URL: v})
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
		Fn:   request.MakeAddToUserAgentHandler("OpenShift/4.x Creator", version.Raw),
	})

	tracker := new(awsdestroy.ErrorTracker)

	resourcesToDelete := sets.NewString()

	entities := []string{"bootstrap", "bootstrap-profile", "bootstrap-role", "bootstrap-role-policy", "bootstrap-sg"}
	filters := []awsdestroy.Filter{}
	for _, entity := range entities {
		filters = append(filters, awsdestroy.Filter{
			"Name":                              fmt.Sprintf("%s-%s", clusterConfig.ClusterID, entity),
			clusterTag(clusterConfig.ClusterID): clusterTagValue,
		})
	}

	iamClient := iam.New(awsSession)
	iamRoleSearch := &awsdestroy.IamRoleSearch{
		Client:  iamClient,
		Filters: filters,
		Logger:  logger,
	}

	tagClients := []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI{
		resourcegroupstaggingapi.New(awsSession),
	}
	switch clusterAWSConfig.Region {
	case endpoints.CnNorth1RegionID, endpoints.CnNorthwest1RegionID:
		break
	case endpoints.UsIsoEast1RegionID, endpoints.UsIsoWest1RegionID, endpoints.UsIsobEast1RegionID:
		break
	case endpoints.UsGovEast1RegionID, endpoints.UsGovWest1RegionID:
		if clusterAWSConfig.Region != endpoints.UsGovWest1RegionID {
			tagClients = append(tagClients,
				resourcegroupstaggingapi.New(awsSession, aws.NewConfig().WithRegion(endpoints.UsGovWest1RegionID)))
		}
	default:
		if clusterAWSConfig.Region != endpoints.UsEast1RegionID {
			tagClients = append(tagClients,
				resourcegroupstaggingapi.New(awsSession, aws.NewConfig().WithRegion(endpoints.UsEast1RegionID)))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Get the initial resources to delete, so that they can be returned if the context is canceled while terminating
	// instances.
	deleted := sets.NewString()
	resourcesToDelete, tagClientsWithResources, err := findResourcesToDelete(ctx, tagClients, iamClient, iamRoleSearch, deleted, filters, logger, clusterConfig.ClusterID)
	if err != nil {
		logger.WithError(err).Info("error while finding resources to delete")
		if err := ctx.Err(); err != nil {
			return err
		}
	}

	instFilters := []awsdestroy.Filter{
		{
			"Name":                              fmt.Sprintf("%s-bootstrap", clusterConfig.ClusterID),
			clusterTag(clusterConfig.ClusterID): clusterTagValue,
		},
	}
	// Terminate EC2 instances. The instances need to be terminated first so that we can ensure that there is nothing
	// running on the cluster creating new resources while we are attempting to delete resources, which could leak
	// the new resources.
	ec2Client := ec2.New(awsSession)
	lastTerminateTime := time.Now()
	err = wait.PollUntilContextCancel(
		ctx,
		time.Second*10,
		true,
		func(ctx context.Context) (done bool, err error) {
			instancesRunning, instancesNotTerminated, err := awsdestroy.FindEC2Instances(ctx, ec2Client, deleted, instFilters, logger)
			if err != nil {
				logger.WithError(err).Info("error while finding EC2 instances to delete")
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
			newlyDeleted, err := awsdestroy.DeleteResources(ctx, awsSession, instancesToDelete, tracker, logger)
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
	)
	if err != nil {
		logger.WithError(err).Infof("Leaking resources: %v", resourcesToDelete.UnsortedList())
		return err
	}

	// Delete the rest of the resources.
	err = wait.PollUntilContextCancel(
		ctx,
		time.Second*10,
		true,
		func(ctx context.Context) (done bool, err error) {
			newlyDeleted, loopError := awsdestroy.DeleteResources(ctx, awsSession, resourcesToDelete.UnsortedList(), tracker, logger)
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
			nextResourcesToDelete, nextTagClients, err := findResourcesToDelete(ctx, tagClientsWithResources, iamClient, iamRoleSearch, deleted, filters, logger, clusterConfig.ClusterID)
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
		logger.WithError(err).Infof("Leaking resources: %v", resourcesToDelete.UnsortedList())
		return err
	}
	return nil
}

func findResourcesToDelete(ctx context.Context, tagClients []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI, iamClient iamiface.IAMAPI, iamRoleSearch *awsdestroy.IamRoleSearch, deleted sets.String, filters []awsdestroy.Filter, logger logrus.FieldLogger, clusterID string) (sets.String, []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI, error) {
	// ignition s3 bucket
	// ignition s3 object
	// bootstrap role policy
	// bootstrap instance
	// bootstrap lb target group attach
	// bootstrap security group
	// bootstrap ssh security group rule
	// bootstrap journald gw security group rule
	var errs []error
	resources := sets.NewString()
	var tagClientsWithResources []*resourcegroupstaggingapi.ResourceGroupsTaggingAPI
	// Find resources by tag
	for _, tagClient := range tagClients {
		resourcesInTagClient, err := awsdestroy.FindResourcesByTag(ctx, tagClient, deleted, filters, logger)
		if err != nil {
			errs = append(errs, err)
		}
		resources = resources.Union(resourcesInTagClient)
		// If there are still resources to be deleted for the tag client or if there was an error getting the resources
		// for the tag client, then retain the tag client for future queries.
		if len(resourcesInTagClient) > 0 || err != nil {
			tagClientsWithResources = append(tagClientsWithResources, tagClient)
		} else {
			logger.Debugf("no deletions from %s, removing client", *tagClient.Config.Region)
		}
	}

	// bootstrap iam role
	// Find IAM roles
	iamRoleResources, err := awsdestroy.FindIAMRoles(ctx, iamRoleSearch, deleted, logger)
	if err != nil {
		errs = append(errs, err)
	}
	resources = resources.Union(iamRoleResources)

	// bootstrap instance profile
	profile := fmt.Sprintf("%s-bootstrap-profile", clusterID)
	response, err := iamClient.GetInstanceProfileWithContext(ctx, &iam.GetInstanceProfileInput{InstanceProfileName: &profile})
	if err != nil {
		var awsErr awserr.Error
		if !errors.As(err, &awsErr) || awsErr.Code() != iam.ErrCodeNoSuchEntityException {
			errs = append(errs, fmt.Errorf("failed to get IAM instance profile: %w", err))
		}
	} else {
		arnString := *response.InstanceProfile.Arn
		if !deleted.Has(arnString) {
			resources.Insert(arnString)
		}
	}

	return resources, tagClientsWithResources, utilerrors.NewAggregate(errs)
}
