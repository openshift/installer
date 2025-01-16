package powervs

import (
	"context"
	"fmt"
	"math"
	"net/http"
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

// listCOSInstances lists COS service instances.
// ibmcloud resource service-instances --output JSON --service-name cloud-object-storage | jq -r '.[] | select(.name|test("rdr-hamzy.*")) | "\(.name) - \(.id)"'
func (o *ClusterUninstaller) listCOSInstances() (cloudResources, error) {
	o.Logger.Debugf("Listing COS instances")

	ctx, cancel := contextWithTimeout()
	defer cancel()

	var (
		// https://github.com/IBM/platform-services-go-sdk/blob/main/resourcecontrollerv2/resource_controller_v2.go#L3086
		options *resourcecontrollerv2.ListResourceInstancesOptions

		perPage int64 = 64

		// https://github.com/IBM/platform-services-go-sdk/blob/main/resourcecontrollerv2/resource_controller_v2.go#L4525-L4534
		resources *resourcecontrollerv2.ResourceInstancesList

		err error

		foundOne = false
		moreData = true
	)

	options = o.controllerSvc.NewListResourceInstancesOptions()
	options.Limit = &perPage
	options.SetResourceID(cosResourceID)
	options.SetResourceGroupID(o.resourceGroupID)
	options.SetType("service_instance")

	result := []cloudResource{}

	for moreData {
		// https://github.com/IBM/platform-services-go-sdk/blob/main/resourcecontrollerv2/resource_controller_v2.go#L173
		resources, _, err = o.controllerSvc.ListResourceInstancesWithContext(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to list COS instances: %w", err)
		}
		o.Logger.Debugf("listCOSInstances: RowsCount %v", *resources.RowsCount)

		for _, instance := range resources.Resources {
			// Match the COS instances created by both the installer and the
			// cluster-image-registry-operator.
			if strings.Contains(*instance.Name, o.InfraID) {
				if !(strings.HasSuffix(*instance.Name, "-cos") ||
					strings.HasSuffix(*instance.Name, "-image-registry")) {
					continue
				}
				foundOne = true
				o.Logger.Debugf("listCOSInstances: FOUND %s %s", *instance.Name, *instance.GUID)
				result = append(result, cloudResource{
					key:      *instance.ID,
					name:     *instance.Name,
					status:   *instance.State,
					typeName: cosTypeName,
					id:       *instance.GUID,
				})
			}
		}

		if resources.NextURL != nil {
			start, err := resources.GetNextStart()
			if err != nil {
				o.Logger.Debugf("listCOSInstances: err = %v", err)
				return nil, fmt.Errorf("failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("listCOSInstances: start = %v", *start)
				options.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("listCOSInstances: NextURL = nil")
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
			// https://github.com/IBM/platform-services-go-sdk/blob/main/resourcecontrollerv2/resource_controller_v2.go#L173
			resources, _, err = o.controllerSvc.ListResourceInstancesWithContext(ctx, options)
			if err != nil {
				return nil, fmt.Errorf("failed to list COS instances: %w", err)
			}
			o.Logger.Debugf("listCOSInstances: RowsCount %v", *resources.RowsCount)
			if resources.NextURL != nil {
				o.Logger.Debugf("listCOSInstances: NextURL   %v", *resources.NextURL)
			}

			o.Logger.Debugf("listCOSInstances: NO matching COS instance against: %s", o.InfraID)
			for _, instance := range resources.Resources {
				o.Logger.Debugf("listCOSInstances: only found COS instance: %s", *instance.Name)
			}

			if resources.NextURL != nil {
				start, err := resources.GetNextStart()
				if err != nil {
					o.Logger.Debugf("listCOSInstances: err = %v", err)
					return nil, fmt.Errorf("failed to GetNextStart: %w", err)
				}
				if start != nil {
					o.Logger.Debugf("listCOSInstances: start = %v", *start)
					options.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("listCOSInstances: NextURL = nil")
				moreData = false
			}
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) findReclaimedCOSInstance(item cloudResource) (*resourcecontrollerv2.ResourceInstance, *resourcecontrollerv2.Reclamation) {
	var getReclamationOptions *resourcecontrollerv2.ListReclamationsOptions
	var reclamations *resourcecontrollerv2.ReclamationsList
	var response *core.DetailedResponse
	var err error
	var reclamation resourcecontrollerv2.Reclamation
	var getInstanceOptions *resourcecontrollerv2.GetResourceInstanceOptions
	var cosInstance *resourcecontrollerv2.ResourceInstance

	getReclamationOptions = o.controllerSvc.NewListReclamationsOptions()

	ctx, cancel := contextWithTimeout()
	defer cancel()

	reclamations, response, err = o.controllerSvc.ListReclamationsWithContext(ctx, getReclamationOptions)
	if err != nil {
		o.Logger.Debugf("Error: ListReclamationsWithContext: %v, response is %v", err, response)
		return nil, nil
	}

	// ibmcloud resource reclamations --output json
	for _, reclamation = range reclamations.Resources {
		getInstanceOptions = o.controllerSvc.NewGetResourceInstanceOptions(*reclamation.ResourceInstanceID)

		cosInstance, response, err = o.controllerSvc.GetResourceInstanceWithContext(ctx, getInstanceOptions)
		if err != nil {
			o.Logger.Debugf("Error: GetResourceInstanceWithContext: %v", err)
			return nil, nil
		}

		if *cosInstance.Name == item.name {
			return cosInstance, &reclamation
		}
	}

	return nil, nil
}

func (o *ClusterUninstaller) destroyCOSInstance(item cloudResource) error {
	var cosInstance *resourcecontrollerv2.ResourceInstance

	cosInstance, _ = o.findReclaimedCOSInstance(item)
	if cosInstance != nil {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted COS Instance %q", item.name)
		return nil
	}

	var getOptions *resourcecontrollerv2.GetResourceInstanceOptions
	var response *core.DetailedResponse
	var err error

	getOptions = o.controllerSvc.NewGetResourceInstanceOptions(item.id)

	ctx, cancel := contextWithTimeout()
	defer cancel()

	_, response, err = o.controllerSvc.GetResourceInstanceWithContext(ctx, getOptions)

	if err != nil && response != nil && response.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted COS Instance %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == http.StatusInternalServerError {
		o.Logger.Infof("destroyCOSInstance: internal server error")
		return nil
	}

	options := o.controllerSvc.NewDeleteResourceInstanceOptions(item.id)
	options.SetRecursive(true)

	response, err = o.controllerSvc.DeleteResourceInstanceWithContext(ctx, options)

	if err != nil && response != nil && response.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to delete COS instance %s: %w", item.name, err)
	}

	var reclamation *resourcecontrollerv2.Reclamation

	cosInstance, reclamation = o.findReclaimedCOSInstance(item)
	if cosInstance != nil {
		var reclamationActionOptions *resourcecontrollerv2.RunReclamationActionOptions

		reclamationActionOptions = o.controllerSvc.NewRunReclamationActionOptions(*reclamation.ID, "reclaim")

		_, response, err = o.controllerSvc.RunReclamationActionWithContext(ctx, reclamationActionOptions)
		if err != nil {
			return fmt.Errorf("failed RunReclamationActionWithContext: %w", err)
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

// COSInstanceID returns the ID of the Cloud Object Storage service instance
// created by the installer during installation.
func (o *ClusterUninstaller) COSInstanceID() (string, error) {
	if o.cosInstanceID != "" {
		return o.cosInstanceID, nil
	}
	cosInstances, err := o.listCOSInstances()
	if err != nil {
		return "", err
	}
	instanceList := cosInstances.list()
	if len(instanceList) == 0 {
		return "", fmt.Errorf("COS instance not found")
	}

	// Locate the installer's COS instance by name.
	for _, instance := range instanceList {
		if instance.name == fmt.Sprintf("%s-cos", o.InfraID) {
			o.cosInstanceID = instance.id
			return instance.id, nil
		}
	}
	return "", fmt.Errorf("COS instance not found")
}
