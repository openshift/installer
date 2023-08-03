package powervs

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"k8s.io/apimachinery/pkg/util/wait"
)

// listCloudConnections lists cloud connections in the cloud.
func (o *ClusterUninstaller) listCloudConnections() (cloudResources, error) {
	// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connections.go#L20-L25
	var cloudConnections *models.CloudConnections

	// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connection.go#L20-L71
	var cloudConnection *models.CloudConnection

	// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/job_reference.go#L18-L27
	var jobReference *models.JobReference

	var err error

	o.Logger.Debugf("Listing Cloud Connections")

	if o.cloudConnectionClient == nil {
		o.Logger.Infof("Skipping deleting Cloud Connections because no service instance was found")
		result := []cloudResource{}
		return cloudResources{}.insert(result...), nil
	}

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listCloudConnections: case <-ctx.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	result := []cloudResource{}

	cloudConnections, err = o.cloudConnectionClient.GetAll()
	if err != nil {
		message := err.Error()
		if strings.Contains(message, "Cloud connections are not currently available") {
			return cloudResources{}.insert(result...), nil
		}
		return nil, fmt.Errorf("failed to list cloud connections: %w", err)
	}

	var foundOne = false

	for _, cloudConnection = range cloudConnections.CloudConnections {
		if strings.Contains(*cloudConnection.Name, o.InfraID) {
			o.Logger.Debugf("listCloudConnections: FOUND: %s (%s)", *cloudConnection.Name, *cloudConnection.CloudConnectionID)
			foundOne = true

			jobReference, err = o.cloudConnectionClient.Delete(*cloudConnection.CloudConnectionID)
			if err != nil {
				return nil, fmt.Errorf("failed to delete cloud connection (%s): %w", *cloudConnection.CloudConnectionID, err)
			}

			o.Logger.Debugf("listCloudConnections: jobReference.ID = %s", *jobReference.ID)

			result = append(result, cloudResource{
				key:      *jobReference.ID,
				name:     *jobReference.ID,
				status:   "",
				typeName: jobTypeName,
				id:       *jobReference.ID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listCloudConnections: NO matching cloud connections against: %s", o.InfraID)
		for _, cloudConnection = range cloudConnections.CloudConnections {
			o.Logger.Debugf("listCloudConnections: only found cloud connection: %s", *cloudConnection.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

// destroyCloudConnections removes all cloud connections that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyCloudConnections() error {
	firstPassList, err := o.listCloudConnections()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(jobTypeName, firstPassList.list())

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyCloudConnections: case <-ctx.Done()")
			return o.Context.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			result, err2 := o.deleteJob(item)
			switch result {
			case DeleteJobSuccess:
				o.Logger.Debugf("destroyCloudConnections: deleteJob returns DeleteJobSuccess")
				return true, nil
			case DeleteJobRunning:
				o.Logger.Debugf("destroyCloudConnections: deleteJob returns DeleteJobRunning")
				return false, nil
			case DeleteJobError:
				o.Logger.Debugf("destroyCloudConnections: deleteJob returns DeleteJobError: %v", err2)
				return false, err2
			default:
				return false, fmt.Errorf("destroyCloudConnections: deleteJob unknown result enum %v", result)
			}
		})
		if err != nil {
			o.Logger.Fatal("destroyCloudConnections: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(jobTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyCloudConnections: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyCloudConnections: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listCloudConnections()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyCloudConnections: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyCloudConnections: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
