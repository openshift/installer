package powervs

import (
	"context"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/pkg/errors"
	"log"
	"strings"
	"time"
)

// listVPCInCloudConnections removes VPCs attached to CloudConnections and returs a list of jobs.
func (o *ClusterUninstaller) listVPCInCloudConnections() (cloudResources, error) {
	var (
		ctx context.Context

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connections.go#L20-L25
		cloudConnections *models.CloudConnections

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connection.go#L20-L71
		cloudConnection          *models.CloudConnection
		cloudConnectionUpdateNew *models.CloudConnection

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/job_reference.go#L18-L27
		jobReference *models.JobReference

		err error

		cloudConnectionID string

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connection_endpoint_v_p_c.go#L19-L26
		endpointVpc       *models.CloudConnectionEndpointVPC
		endpointUpdateVpc models.CloudConnectionEndpointVPC

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connection_v_p_c.go#L18-L26
		Vpc *models.CloudConnectionVPC

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connection_update.go#L20
		cloudConnectionUpdate models.CloudConnectionUpdate

		foundOne bool = false
		foundVpc bool = false
	)

	ctx, _ = o.contextWithTimeout()

	o.Logger.Debugf("Listing VPCs in Cloud Connections")

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listVPCInCloudConnections: case <-ctx.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	cloudConnections, err = o.cloudConnectionClient.GetAll()
	if err != nil {
		log.Fatalf("Failed to list cloud connections: %v", err)
	}

	result := []cloudResource{}
	for _, cloudConnection = range cloudConnections.CloudConnections {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listVPCInCloudConnections: case <-ctx.Done()")
			return nil, o.Context.Err() // we're cancelled, abort
		default:
		}

		if !strings.Contains(*cloudConnection.Name, o.InfraID) {
			// Skip this one!
			continue
		}

		foundOne = true

		o.Logger.Debugf("listVPCInCloudConnections: FOUND: %s (%s)", *cloudConnection.Name, *cloudConnection.CloudConnectionID)

		cloudConnectionID = *cloudConnection.CloudConnectionID

		cloudConnection, err = o.cloudConnectionClient.Get(cloudConnectionID)
		if err != nil {
			log.Fatalf("Failed to get cloud connection %s: %v", cloudConnectionID, err)
		}

		endpointVpc = cloudConnection.Vpc

		o.Logger.Debugf("listVPCInCloudConnections: endpointVpc = %+v", endpointVpc)

		foundVpc = false
		for _, Vpc = range endpointVpc.Vpcs {
			o.Logger.Debugf("listVPCInCloudConnections: Vpc = %+v", Vpc)
			o.Logger.Debugf("listVPCInCloudConnections: Vpc.Name = %v, o.InfraID = %v", Vpc.Name, o.InfraID)
			if strings.Contains(Vpc.Name, o.InfraID) {
				foundVpc = true
			}
		}
		o.Logger.Debugf("listVPCInCloudConnections: foundVpc = %v", foundVpc)
		if !foundVpc {
			continue
		}

		// https://github.com/IBM-Cloud/power-go-client/blob/v1.0.88/power/models/cloud_connection_v_p_c.go#L18
		var vpcsUpdate []*models.CloudConnectionVPC

		for _, Vpc = range endpointVpc.Vpcs {
			if !strings.Contains(Vpc.Name, o.InfraID) {
				vpcsUpdate = append(vpcsUpdate, Vpc)
			}
		}

		endpointUpdateVpc.Enabled = len(vpcsUpdate) > 0
		endpointUpdateVpc.Vpcs = vpcsUpdate

		cloudConnectionUpdate.Vpc = &endpointUpdateVpc

		var vpcsStrings []string

		for _, Vpc = range vpcsUpdate {
			vpcsStrings = append(vpcsStrings, Vpc.Name)
		}
		o.Logger.Debugf("listVPCInCloudConnections: vpcsUpdate = %v", vpcsStrings)
		o.Logger.Debugf("listVPCInCloudConnections: endpointUpdateVpc = %+v", endpointUpdateVpc)

		cloudConnectionUpdateNew, jobReference, err = o.cloudConnectionClient.Update(*cloudConnection.CloudConnectionID, &cloudConnectionUpdate)
		if err != nil {
			log.Fatalf("Failed to update cloud connection %v", err)
		}

		o.Logger.Debugf("listVPCInCloudConnections: cloudConnectionUpdateNew = %+v", cloudConnectionUpdateNew)
		o.Logger.Debugf("listVPCInCloudConnections: jobReference = %+v", jobReference)

		result = append(result, cloudResource{
			key:      *jobReference.ID,
			name:     *jobReference.ID,
			status:   "",
			typeName: jobTypeName,
			id:       *jobReference.ID,
		})
	}

	if !foundOne {
		o.Logger.Debugf("listVPCInCloudConnections: NO matching cloud connections")
		for _, cloudConnection = range cloudConnections.CloudConnections {
			o.Logger.Debugf("listVPCInCloudConnections: only found cloud connection: %s", *cloudConnection.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

// destroyVPCInCloudConnections removes all VPCs in cloud connections that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyVPCInCloudConnections() error {
	var (
		found cloudResources
		err   error
		ctx   context.Context
		items []cloudResource
	)

	found, err = o.listVPCInCloudConnections()
	if err != nil {
		return err
	}

	items = o.insertPendingItems(jobTypeName, found.list())

	ctx, _ = o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyVPCInCloudConnections: case <-o.Context.Done()")
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
		return errors.Errorf("destroyVPCInCloudConnections: %d undeleted items pending", len(items))
	}
	return nil
}
