package powervs

import (
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
	gohttp "net/http"
	"strings"
)

const subnetTypeName = "subnet"

// listSubnets lists subnets in the cloud.
func (o *ClusterUninstaller) listSubnets() (cloudResources, error) {
	o.Logger.Debugf("Listing Subnets")

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listSubnets: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	options := o.vpcSvc.NewListSubnetsOptions()
	subnets, detailedResponse, err := o.vpcSvc.ListSubnets(options)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to list subnets and the response is: %s", detailedResponse)
	}

	var foundOne = false

	result := []cloudResource{}
	for _, subnet := range subnets.Subnets {
		if strings.Contains(*subnet.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listSubnets: FOUND: %s, %s", *subnet.ID, *subnet.Name)
			result = append(result, cloudResource{
				key:      *subnet.ID,
				name:     *subnet.Name,
				status:   "",
				typeName: subnetTypeName,
				id:       *subnet.ID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listSubnets: NO matching subnet against: %s", o.InfraID)
		for _, subnet := range subnets.Subnets {
			o.Logger.Debugf("listSubnets: subnet: %s", *subnet.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteSubnet(item cloudResource) error {
	var getOptions *vpcv1.GetSubnetOptions
	var response *core.DetailedResponse
	var err error

	getOptions = o.vpcSvc.NewGetSubnetOptions(item.id)
	_, response, err = o.vpcSvc.GetSubnet(getOptions)

	if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted subnet %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == gohttp.StatusInternalServerError {
		o.Logger.Infof("deleteSubnet: internal server error")
		return nil
	}

	o.Logger.Debugf("Deleting subnet %q", item.name)

	ctx, _ := o.contextWithTimeout()

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("deleteSubnet: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	deleteOptions := o.vpcSvc.NewDeleteSubnetOptions(item.id)
	_, err = o.vpcSvc.DeleteSubnetWithContext(ctx, deleteOptions)

	if err != nil {
		return errors.Wrapf(err, "failed to delete subnet %s", item.name)
	}

	return nil
}

// destroySubnets removes all subnet resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroySubnets() error {
	found, err := o.listSubnets()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(subnetTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroySubnets: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted subnet %q", item.name)
				continue
			}
			err = o.deleteSubnet(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(subnetTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(subnetTypeName); len(items) > 0 {
		return errors.Errorf("destroySubnets: %d undeleted items pending", len(items))
	}
	return nil
}
