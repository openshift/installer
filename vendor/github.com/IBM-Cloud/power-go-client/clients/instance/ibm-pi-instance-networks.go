package instance

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	pvmclient "github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// Instance-scoped networks under a PVM:
type IBMPIInstanceNetworksClient struct {
	IBMPIClient
}

func NewIBMPIInstanceNetworksClient(ctx context.Context, sess *ibmpisession.IBMPISession, cloudInstanceID string) *IBMPIInstanceNetworksClient {
	return &IBMPIInstanceNetworksClient{*NewIBMPIClient(ctx, sess, cloudInstanceID)}
}

// GetAll returns the wrapper that contains all networks on a PVM instance.
func (c *IBMPIInstanceNetworksClient) GetAll(pvmInstanceID string) (*models.PVMInstanceNetworks, error) {
	params := pvmclient.NewPcloudPvminstancesNetworksGetallParams().
		WithContext(c.ctx).
		WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(c.cloudInstanceID).
		WithPvmInstanceID(pvmInstanceID)

	resp, err := c.session.Power.PCloudpVMInstances.PcloudPvminstancesNetworksGetall(params, c.session.AuthInfo(c.cloudInstanceID))
	if err != nil {
		return nil, ibmpisession.SDKFailWithAPIError(err, fmt.Errorf("failed to get networks for pvm instance %s: %w", pvmInstanceID, err))
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to get networks for pvm instance %s", pvmInstanceID)
	}
	return resp.Payload, nil
}

// Get returns the wrapper for a specific network on a PVM instance.
func (c *IBMPIInstanceNetworksClient) Get(pvmInstanceID, networkID string) (*models.PVMInstanceNetworks, error) {
	params := pvmclient.NewPcloudPvminstancesNetworksGetParams().
		WithContext(c.ctx).
		WithTimeout(helpers.PIGetTimeOut).
		WithCloudInstanceID(c.cloudInstanceID).
		WithPvmInstanceID(pvmInstanceID).
		WithNetworkID(networkID)

	resp, err := c.session.Power.PCloudpVMInstances.PcloudPvminstancesNetworksGet(params, c.session.AuthInfo(c.cloudInstanceID))
	if err != nil {
		return nil, ibmpisession.SDKFailWithAPIError(err, fmt.Errorf("failed to get network %s for pvm instance %s: %w", networkID, pvmInstanceID, err))
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to get network %s for pvm instance %s", networkID, pvmInstanceID)
	}
	return resp.Payload, nil
}
