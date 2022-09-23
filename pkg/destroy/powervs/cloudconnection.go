package powervs

import (
	"context"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/pkg/errors"
	"log"
	"strings"
	"time"
)

// listCloudConnections lists cloud connections in the cloud.
func (o *ClusterUninstaller) listCloudConnections() (cloudResources, error) {
	var (
		ctx context.Context

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connections.go#L20-L25
		cloudConnections *models.CloudConnections

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connection.go#L20-L71
		cloudConnection *models.CloudConnection

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/job_reference.go#L18-L27
		jobReference *models.JobReference

		err error

		cloudConnectionID string

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connection_endpoint_v_p_c.go#L19-L26
		EndpointVpc *models.CloudConnectionEndpointVPC

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connection_v_p_c.go#L18-L26
		Vpc *models.CloudConnectionVPC

		foundOne       = false
		foundVpc       = false
		vpcStillExists = true
	)

	ctx, _ = o.contextWithTimeout()

	o.Logger.Debugf("Listing Cloud Connections")

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listCloudConnections: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	cloudConnections, err = o.cloudConnectionClient.GetAll()
	if err != nil {
		log.Fatalf("Failed to list cloud connections: %v", err)
	}

	result := []cloudResource{}
	for _, cloudConnection = range cloudConnections.CloudConnections {
		if !strings.Contains(*cloudConnection.Name, o.InfraID) {
			// Skip this one!
			continue
		}

		foundOne = true

		o.Logger.Debugf("listCloudConnections: FOUND: %s (%s)", *cloudConnection.Name, *cloudConnection.CloudConnectionID)

		cloudConnectionID = *cloudConnection.CloudConnectionID

		vpcStillExists = true

		for !o.timeout(ctx) {
			if !vpcStillExists {
				break
			}

			select {
			case <-o.Context.Done():
				o.Logger.Debugf("listCloudConnections: case <-o.Context.Done()")
				return nil, o.Context.Err() // we're cancelled, abort
			default:
			}

			cloudConnection, err = o.cloudConnectionClient.Get(cloudConnectionID)
			if err != nil {
				log.Fatalf("Failed to get cloud connection %s: %v", cloudConnectionID, err)
			}

			EndpointVpc = cloudConnection.Vpc
			o.Logger.Debugf("listCloudConnections: EndpointVpc = %+v", EndpointVpc)

			foundVpc = false
			for _, Vpc = range EndpointVpc.Vpcs {
				if Vpc != nil {
					foundVpc = true
				}
				o.Logger.Debugf("listCloudConnections: Vpc = %+v", Vpc)
			}
			o.Logger.Debugf("listCloudConnections: foundVpc = %v", foundVpc)
			if foundVpc {
				o.Logger.Debugf("listCloudConnections: This CC still has VPCs attached, waiting...")

				time.Sleep(15 * time.Second)
			} else {
				vpcStillExists = false
			}
		}

		// Finally delete the CloudConnection!
		jobReference, err = o.cloudConnectionClient.Delete(*cloudConnection.CloudConnectionID)
		if err != nil {
			errors.Errorf("Failed to delete cloud connection (%s): %v", *cloudConnection.CloudConnectionID, err)
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
	found, err := o.listCloudConnections()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(jobTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyCloudConnections: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted job %q", item.name)
				continue
			}
			err := o.deleteJob(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(jobTypeName)
		if len(items) == 0 {
			break
		}

		time.Sleep(15 * time.Second)
	}

	if items = o.getPendingItems(jobTypeName); len(items) > 0 {
		return errors.Errorf("destroyCloudConnections: %d undeleted items pending", len(items))
	}
	return nil
}
