package powervs

import (
	"context"
	"fmt"
	"math"
	gohttp "net/http"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	serviceInstanceTypeName = "service instance"

	// resource ID for Power Systems Virtual Server in the Global catalog.
	virtualServerResourceID = "abd259f0-9990-11e8-acc8-b9f54a8f1661"
)

// convertResourceGroupNameToID converts a resource group name/id to an id.
func (o *ClusterUninstaller) convertResourceGroupNameToID(resourceGroupID string) (string, error) {
	listResourceGroupsOptions := o.managementSvc.NewListResourceGroupsOptions()

	resourceGroups, _, err := o.managementSvc.ListResourceGroups(listResourceGroupsOptions)
	if err != nil {
		return "", err
	}

	for _, resourceGroup := range resourceGroups.Resources {
		if *resourceGroup.Name == resourceGroupID {
			return *resourceGroup.ID, nil
		} else if *resourceGroup.ID == resourceGroupID {
			return resourceGroupID, nil
		}
	}

	return "", fmt.Errorf("failed to find resource group %v", resourceGroupID)
}

// listServiceInstances list service instances for the cluster.
func (o *ClusterUninstaller) listServiceInstances() (cloudResources, error) {
	o.Logger.Debugf("Listing service instances")

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listServiceInstances: case <-ctx.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	var (
		resourceGroupID string
		options         *resourcecontrollerv2.ListResourceInstancesOptions
		resources       *resourcecontrollerv2.ResourceInstancesList
		err             error
		perPage         int64 = 10
		moreData              = true
		nextURL         *string
	)

	resourceGroupID, err = o.convertResourceGroupNameToID(o.resourceGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert resourceGroupID: %w", err)
	}
	o.Logger.Debugf("listServiceInstances: converted %v to %v", o.resourceGroupID, resourceGroupID)

	options = o.controllerSvc.NewListResourceInstancesOptions()
	// options.SetType("resource_instance")
	options.SetResourceGroupID(resourceGroupID)
	options.SetResourceID(virtualServerResourceID)
	options.SetLimit(perPage)

	result := []cloudResource{}

	for moreData {
		if options.Start != nil {
			o.Logger.Debugf("listServiceInstances: options = %+v, options.Limit = %v, options.Start = %v, options.ResourceGroupID = %v", options, *options.Limit, *options.Start, *options.ResourceGroupID)
		} else {
			o.Logger.Debugf("listServiceInstances: options = %+v, options.Limit = %v, options.ResourceGroupID = %v", options, *options.Limit, *options.ResourceGroupID)
		}

		resources, _, err = o.controllerSvc.ListResourceInstancesWithContext(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to list resource instances: %w", err)
		}

		o.Logger.Debugf("listServiceInstances: resources.RowsCount = %v", *resources.RowsCount)

		for _, resource := range resources.Resources {
			var (
				getResourceOptions *resourcecontrollerv2.GetResourceInstanceOptions
				resourceInstance   *resourcecontrollerv2.ResourceInstance
				response           *core.DetailedResponse
			)

			o.Logger.Debugf("listServiceInstances: resource.Name = %s", *resource.Name)

			getResourceOptions = o.controllerSvc.NewGetResourceInstanceOptions(*resource.ID)

			resourceInstance, response, err = o.controllerSvc.GetResourceInstance(getResourceOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to get instance: %s: %w", response, err)
			}
			if response != nil && response.StatusCode == gohttp.StatusNotFound {
				o.Logger.Debugf("listServiceInstances: gohttp.StatusNotFound")
				continue
			} else if response != nil && response.StatusCode == gohttp.StatusInternalServerError {
				o.Logger.Debugf("listServiceInstances: gohttp.StatusInternalServerError")
				continue
			}

			if resourceInstance.Type == nil {
				o.Logger.Debugf("listServiceInstances: type: nil")
			} else {
				o.Logger.Debugf("listServiceInstances: type: %v", *resourceInstance.Type)
			}

			if resourceInstance.Type == nil || resourceInstance.GUID == nil {
				continue
			}
			if *resourceInstance.Type != "service_instance" && *resourceInstance.Type != "composite_instance" {
				continue
			}
			if !strings.Contains(*resource.Name, o.InfraID) {
				continue
			}

			if strings.Contains(*resource.Name, o.InfraID) {
				result = append(result, cloudResource{
					key:      *resource.ID,
					name:     *resource.Name,
					status:   *resource.GUID,
					typeName: serviceInstanceTypeName,
					id:       *resource.ID,
				})
			}
		}

		// Based on: https://cloud.ibm.com/apidocs/resource-controller/resource-controller?code=go#list-resource-instances
		nextURL, err = core.GetQueryParam(resources.NextURL, "start")
		if err != nil {
			return nil, fmt.Errorf("failed to GetQueryParam on start: %w", err)
		}
		if nextURL == nil {
			o.Logger.Debugf("nextURL = nil")
			options.SetStart("")
		} else {
			o.Logger.Debugf("nextURL = %v", *nextURL)
			options.SetStart(*nextURL)
		}

		moreData = *resources.RowsCount == perPage
	}

	return cloudResources{}.insert(result...), nil
}

// destroyServiceInstance destroys a service instance.
func (o *ClusterUninstaller) destroyServiceInstance(item cloudResource) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyServiceInstance: case <-ctx.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	o.Logger.Debugf("destroyServiceInstance: Preparing to delete, item.name = %v", item.name)

	var (
		getOptions *resourcecontrollerv2.GetResourceInstanceOptions
		response   *core.DetailedResponse
		err        error
	)

	getOptions = o.controllerSvc.NewGetResourceInstanceOptions(item.id)

	_, response, err = o.controllerSvc.GetResourceInstanceWithContext(ctx, getOptions)

	if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Service Instance %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == gohttp.StatusInternalServerError {
		o.Logger.Infof("destroyServiceInstance: internal server error")
		return nil
	}

	options := o.controllerSvc.NewDeleteResourceInstanceOptions(item.id)
	options.SetRecursive(true)

	response, err = o.controllerSvc.DeleteResourceInstanceWithContext(ctx, options)

	if err != nil && response != nil && response.StatusCode != gohttp.StatusNotFound {
		return fmt.Errorf("failed to delete service instance %s: %w", item.name, err)
	}

	o.Logger.Infof("Deleted Service Instance %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroyServiceInstances removes all service instances have a name containing
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyServiceInstances() error {
	firstPassList, err := o.listServiceInstances()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(serviceInstanceTypeName, firstPassList.list())

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyServiceInstances: case <-ctx.Done()")
			return o.Context.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.destroyServiceInstance(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyServiceInstances: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(serviceInstanceTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyServiceInstances: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyServiceInstances: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listServiceInstances()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyServiceInstances: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyServiceInstances: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
