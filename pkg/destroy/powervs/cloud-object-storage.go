package powervs

import (
	"context"
	"fmt"
	"math"
	gohttp "net/http"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	// https://github.com/IBM/platform-services-go-sdk/blob/v0.18.16/resourcecontrollerv2/resource_controller_v2.go
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"k8s.io/apimachinery/pkg/util/wait"
)

const cosTypeName = "cos instance"

// $ ibmcloud catalog service cloud-object-storage --output json | jq -r '.[].id'
// dff97f5c-bc5e-4455-b470-411c3edbe49c.
const cosResourceID = "dff97f5c-bc5e-4455-b470-411c3edbe49c"

// listCOSInstances list the Cloud Object Storage instances by name or tag in the IBM Cloud.
// ibmcloud resource service-instances --output JSON --service-name cloud-object-storage | jq -r '.[] | select(.name|test("rdr-hamzy.*")) | "\(.name) - \(.id)"' .
func (o *ClusterUninstaller) listCOSInstances() (cloudResources, error) {
	var (
		cosIDs   []string
		cosID    string
		ctx      context.Context
		cancel   context.CancelFunc
		result   = make([]cloudResource, 0, 1)
		options  *resourcecontrollerv2.GetResourceInstanceOptions
		instance *resourcecontrollerv2.ResourceInstance
		response *core.DetailedResponse
		err      error
	)

	if o.searchByTag {
		// Should we list by tag matching?
		cosIDs, err = o.listByTag(TagTypeCloudObjectStorage)
	} else {
		// Otherwise list will list by name matching.
		cosIDs, err = o.listCOSInstancesByName()
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, cosID = range cosIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listCOSInstances: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		options = o.controllerSvc.NewGetResourceInstanceOptions(cosID)

		instance, response, err = o.controllerSvc.GetResourceInstanceWithContext(ctx, options)
		if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
			// The COS could have been deleted just after a list was created.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get COS instance (%s): err = %w, response = %v", cosID, err, response)
		}

		result = append(result, cloudResource{
			key:      *instance.ID,
			name:     *instance.Name,
			status:   *instance.State,
			typeName: cosTypeName,
			id:       *instance.GUID,
		})
	}

	return cloudResources{}.insert(result...), nil
}

// listCOSInstancesByName list the Cloud Object Storage instances by name in the IBM Cloud.
func (o *ClusterUninstaller) listCOSInstancesByName() ([]string, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc

		// https://github.com/IBM/platform-services-go-sdk/blob/main/resourcecontrollerv2/resource_controller_v2.go#L3425-L3466
		options *resourcecontrollerv2.ListResourceInstancesOptions

		perPage int64 = 64

		// https://github.com/IBM/platform-services-go-sdk/blob/main/resourcecontrollerv2/resource_controller_v2.go#L5008-L5017
		resources *resourcecontrollerv2.ResourceInstancesList

		// https://github.com/IBM/platform-services-go-sdk/blob/main/resourcecontrollerv2/resource_controller_v2.go#L4485-L4608
		instance resourcecontrollerv2.ResourceInstance

		foundOne = false
		moreData = true

		result = make([]string, 0, 1)

		err error
	)

	o.Logger.Debugf("Listing COS instances by NAME")

	ctx, cancel = contextWithTimeout()
	defer cancel()

	options = o.controllerSvc.NewListResourceInstancesOptions()
	options.Limit = &perPage
	options.SetResourceID(cosResourceID)
	options.SetResourceGroupID(o.resourceGroupID)
	options.SetType("service_instance")

	for moreData {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listCOSInstancesByName: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		// https://github.com/IBM/platform-services-go-sdk/blob/main/resourcecontrollerv2/resource_controller_v2.go#L173
		resources, _, err = o.controllerSvc.ListResourceInstancesWithContext(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to list COS instances: %w", err)
		}
		o.Logger.Debugf("listCOSInstancesByName: RowsCount %v", *resources.RowsCount)

		for _, instance = range resources.Resources {
			// Match the COS instances created by both the installer and the
			// cluster-image-registry-operator.
			if strings.Contains(*instance.Name, o.InfraID) {
				if !(strings.HasSuffix(*instance.Name, "-cos") ||
					strings.HasSuffix(*instance.Name, "-image-registry")) {
					continue
				}
				foundOne = true
				o.Logger.Debugf("listCOSInstancesByName: FOUND %s %s", *instance.Name, *instance.GUID)
				result = append(result, *instance.ID)
			}
		}

		if resources.NextURL != nil {
			start, err := resources.GetNextStart()
			if err != nil {
				o.Logger.Debugf("listCOSInstancesByName: err = %v", err)
				return nil, fmt.Errorf("failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("listCOSInstancesByName: start = %v", *start)
				options.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("listCOSInstancesByName: NextURL = nil")
			moreData = false
		}
	}
	if !foundOne {
		options = o.controllerSvc.NewListResourceInstancesOptions()
		options.Limit = &perPage
		options.SetResourceID(cosResourceID)
		options.SetResourceGroupID(o.resourceGroupID)
		options.SetType("service_instance")

		moreData = true
		for moreData {
			select {
			case <-ctx.Done():
				o.Logger.Debugf("listCOSInstancesByName: case <-ctx.Done()")
				return nil, ctx.Err() // we're cancelled, abort
			default:
			}

			// https://github.com/IBM/platform-services-go-sdk/blob/main/resourcecontrollerv2/resource_controller_v2.go#L173
			resources, _, err = o.controllerSvc.ListResourceInstancesWithContext(ctx, options)
			if err != nil {
				return nil, fmt.Errorf("failed to list COS instances: %w", err)
			}
			o.Logger.Debugf("listCOSInstancesByName: RowsCount %v", *resources.RowsCount)
			if resources.NextURL != nil {
				o.Logger.Debugf("listCOSInstancesByName: NextURL   %v", *resources.NextURL)
			}

			o.Logger.Debugf("listCOSInstancesByName: NO matching COS instance against: %s", o.InfraID)
			for _, instance = range resources.Resources {
				o.Logger.Debugf("listCOSInstancesByName: only found COS instance: %s", *instance.Name)
			}

			if resources.NextURL != nil {
				start, err := resources.GetNextStart()
				if err != nil {
					o.Logger.Debugf("listCOSInstancesByName: err = %v", err)
					return nil, fmt.Errorf("failed to GetNextStart: %w", err)
				}
				if start != nil {
					o.Logger.Debugf("listCOSInstancesByName: start = %v", *start)
					options.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("listCOSInstancesByName: NextURL = nil")
				moreData = false
			}
		}
	}

	return result, nil
}

func (o *ClusterUninstaller) findReclaimedCOSInstance(item cloudResource) (*resourcecontrollerv2.ResourceInstance, *resourcecontrollerv2.Reclamation) {
	var (
		getReclamationOptions *resourcecontrollerv2.ListReclamationsOptions
		ctx                   context.Context
		cancel                context.CancelFunc
		reclamations          *resourcecontrollerv2.ReclamationsList
		response              *core.DetailedResponse
		err                   error
		reclamation           resourcecontrollerv2.Reclamation
		getInstanceOptions    *resourcecontrollerv2.GetResourceInstanceOptions
		cosInstance           *resourcecontrollerv2.ResourceInstance
	)

	getReclamationOptions = o.controllerSvc.NewListReclamationsOptions()

	ctx, cancel = contextWithTimeout()
	defer cancel()

	reclamations, response, err = o.controllerSvc.ListReclamationsWithContext(ctx, getReclamationOptions)
	if err != nil {
		o.Logger.Debugf("Error: ListReclamationsWithContext: %v, response is %v", err, response)
		return nil, nil
	}

	// ibmcloud resource reclamations --output json
	for _, reclamation = range reclamations.Resources {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("findReclaimedCOSInstance: case <-ctx.Done()")
			return nil, nil // we're cancelled, abort
		default:
		}

		getInstanceOptions = o.controllerSvc.NewGetResourceInstanceOptions(*reclamation.ResourceInstanceID)

		cosInstance, response, err = o.controllerSvc.GetResourceInstanceWithContext(ctx, getInstanceOptions)
		if err != nil {
			o.Logger.Debugf("Error: GetResourceInstanceWithContext: %v, response is %v", err, response)
			return nil, nil
		}

		if *cosInstance.Name == item.name {
			return cosInstance, &reclamation
		}
	}

	return nil, nil
}

func (o *ClusterUninstaller) destroyCOSInstance(item cloudResource) error {
	var (
		cosInstance              *resourcecontrollerv2.ResourceInstance
		getOptions               *resourcecontrollerv2.GetResourceInstanceOptions
		ctx                      context.Context
		cancel                   context.CancelFunc
		response                 *core.DetailedResponse
		reclamation              *resourcecontrollerv2.Reclamation
		reclamationActionOptions *resourcecontrollerv2.RunReclamationActionOptions
		deleteOptions            *resourcecontrollerv2.DeleteResourceInstanceOptions
		err                      error
	)

	cosInstance, _ = o.findReclaimedCOSInstance(item)
	if cosInstance != nil {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted COS Instance %q", item.name)
		return nil
	}

	getOptions = o.controllerSvc.NewGetResourceInstanceOptions(item.id)

	ctx, cancel = contextWithTimeout()
	defer cancel()

	_, response, err = o.controllerSvc.GetResourceInstanceWithContext(ctx, getOptions)

	if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted COS Instance %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == gohttp.StatusInternalServerError {
		o.Logger.Infof("destroyCOSInstance: internal server error")
		return nil
	}

	deleteOptions = o.controllerSvc.NewDeleteResourceInstanceOptions(item.id)
	deleteOptions.SetRecursive(true)

	response, err = o.controllerSvc.DeleteResourceInstanceWithContext(ctx, deleteOptions)
	if err != nil && response != nil && response.StatusCode != gohttp.StatusNotFound {
		return fmt.Errorf("failed to delete COS instance %s: %w", item.name, err)
	}

	cosInstance, reclamation = o.findReclaimedCOSInstance(item)
	if cosInstance != nil {
		reclamationActionOptions = o.controllerSvc.NewRunReclamationActionOptions(*reclamation.ID, "reclaim")

		_, response, err = o.controllerSvc.RunReclamationActionWithContext(ctx, reclamationActionOptions)
		if err != nil {
			return fmt.Errorf("failed RunReclamationActionWithContext: %w, response = %v", err, response)
		}
	}

	o.Logger.Infof("Deleted COS Instance %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroyCOSInstances removes the COS service instance resources that have a
// name prefixed with the cluster's infra ID.
func (o *ClusterUninstaller) destroyCOSInstances() error {
	firstPassList, err := o.listCOSInstances()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(cosTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyCOSInstances: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.destroyCOSInstance(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyCOSInstances: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(cosTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyCOSInstances: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyCOSInstances: %d undeleted items pending", len(items))
	}

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyCOSInstances: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listCOSInstances()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyCOSInstances: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyCOSInstances: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
