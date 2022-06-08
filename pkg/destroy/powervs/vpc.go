package powervs

import (
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
	gohttp "net/http"
	"strings"
)

const vpcTypeName = "vpc"

// listVPCs lists VPCs in the cloud.
func (o *ClusterUninstaller) listVPCs() (cloudResources, error) {
	o.Logger.Debugf("Listing VPCs")

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listVPCs: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	options := o.vpcSvc.NewListVpcsOptions()
	vpcs, _, err := o.vpcSvc.ListVpcs(options)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to list vps")
	}

	var foundOne = false

	result := []cloudResource{}
	for _, vpc := range vpcs.Vpcs {
		if strings.Contains(*vpc.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listVPCs: FOUND: %s, %s", *vpc.ID, *vpc.Name)
			result = append(result, cloudResource{
				key:      *vpc.ID,
				name:     *vpc.Name,
				status:   "",
				typeName: vpcTypeName,
				id:       *vpc.ID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listVPCs: NO matching vpc against: %s", o.InfraID)
		for _, vpc := range vpcs.Vpcs {
			o.Logger.Debugf("listVPCs: vpc: %s", *vpc.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteVPC(item cloudResource) error {
	var getOptions *vpcv1.GetVPCOptions
	var getResponse *core.DetailedResponse
	var deleteResponse *core.DetailedResponse
	var err error

	getOptions = o.vpcSvc.NewGetVPCOptions(item.id)
	_, getResponse, err = o.vpcSvc.GetVPC(getOptions)

	// Sadly, there is no way to get the status of this VPC to check on the results of the
	// delete call.

	if err != nil && getResponse != nil && getResponse.StatusCode == gohttp.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted vpc %q", item.name)
		return nil
	}
	if err != nil && getResponse != nil && getResponse.StatusCode == gohttp.StatusInternalServerError {
		o.Logger.Infof("deleteVPC: internal server error")
		return nil
	}

	o.Logger.Debugf("Deleting vpc %q", item.name)

	ctx, _ := o.contextWithTimeout()

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("deleteVPC: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	deleteOptions := o.vpcSvc.NewDeleteVPCOptions(item.id)
	deleteResponse, err = o.vpcSvc.DeleteVPCWithContext(ctx, deleteOptions)
	o.Logger.Debugf("deleteVPC: DeleteVPCWithContext returns %+v", deleteResponse)

	if err != nil {
		return errors.Wrapf(err, "failed to delete vpc %s", item.name)
	}

	return nil
}

// destroyVPCs removes all vpc resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyVPCs() error {
	found, err := o.listVPCs()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(vpcTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyVPCs: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted vpc %q", item.name)
				continue
			}
			err = o.deleteVPC(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(vpcTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(vpcTypeName); len(items) > 0 {
		return errors.Errorf("destroyVPCs: %d undeleted items pending", len(items))
	}
	return nil
}
