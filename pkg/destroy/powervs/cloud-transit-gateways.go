package powervs

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"k8s.io/apimachinery/pkg/util/wait"

	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
)

const (
	transitGatewayTypeName           = "transitGateway"
	transitGatewayConnectionTypeName = "transitGatewayConnection"
)

// listTransitGateways lists Transit Gateways in the IBM Cloud.
func (o *ClusterUninstaller) listTransitGateways() (cloudResources, error) {
	o.Logger.Debugf("Listing Transit Gateways (%s)", o.InfraID)

	var (
		ctx                        context.Context
		cancel                     func()
		listTransitGatewaysOptions *transitgatewayapisv1.ListTransitGatewaysOptions
		gatewayCollection          *transitgatewayapisv1.TransitGatewayCollection
		gateway                    transitgatewayapisv1.TransitGateway
		response                   *core.DetailedResponse
		err                        error
		foundOne                         = false
		perPage                    int64 = 32
		moreData                         = true
	)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	listTransitGatewaysOptions = o.tgClient.NewListTransitGatewaysOptions()
	listTransitGatewaysOptions.Limit = &perPage

	result := []cloudResource{}

	for moreData {
		// https://github.com/IBM/networking-go-sdk/blob/master/transitgatewayapisv1/transit_gateway_apis_v1.go#L184
		gatewayCollection, response, err = o.tgClient.ListTransitGatewaysWithContext(ctx, listTransitGatewaysOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list transit gateways: %w and the respose is: %s", err, response)
		}

		for _, gateway = range gatewayCollection.TransitGateways {
			if strings.Contains(*gateway.Name, o.InfraID) {
				foundOne = true
				o.Logger.Debugf("listTransitGateways: FOUND: %s, %s", *gateway.ID, *gateway.Name)
				result = append(result, cloudResource{
					key:      *gateway.ID,
					name:     *gateway.Name,
					status:   "",
					typeName: transitGatewayTypeName,
					id:       *gateway.ID,
				})
			}
		}

		if gatewayCollection.First != nil {
			o.Logger.Debugf("listTransitGateways: First = %+v", *gatewayCollection.First.Href)
		} else {
			o.Logger.Debugf("listTransitGateways: First = nil")
		}
		if gatewayCollection.Limit != nil {
			o.Logger.Debugf("listTransitGateways: Limit = %v", *gatewayCollection.Limit)
		}
		if gatewayCollection.Next != nil {
			start, err := gatewayCollection.GetNextStart()
			if err != nil {
				o.Logger.Debugf("listTransitGateways: err = %v", err)
				return nil, fmt.Errorf("listTransitGateways: failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("listTransitGateways: start = %v", *start)
				listTransitGatewaysOptions.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("listTransitGateways: Next = nil")
			moreData = false
		}
	}
	if !foundOne {
		o.Logger.Debugf("listTransitGateways: NO matching transit gateway against: %s", o.InfraID)

		listTransitGatewaysOptions = o.tgClient.NewListTransitGatewaysOptions()
		listTransitGatewaysOptions.Limit = &perPage
		moreData = true

		for moreData {
			gatewayCollection, response, err = o.tgClient.ListTransitGatewaysWithContext(ctx, listTransitGatewaysOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to list transit gateways: %w and the respose is: %s", err, response)
			}
			for _, gateway = range gatewayCollection.TransitGateways {
				o.Logger.Debugf("listTransitGateways: FOUND: %s, %s", *gateway.ID, *gateway.Name)
			}
			if gatewayCollection.First != nil {
				o.Logger.Debugf("listTransitGateways: First = %+v", *gatewayCollection.First.Href)
			} else {
				o.Logger.Debugf("listTransitGateways: First = nil")
			}
			if gatewayCollection.Limit != nil {
				o.Logger.Debugf("listTransitGateways: Limit = %v", *gatewayCollection.Limit)
			}
			if gatewayCollection.Next != nil {
				start, err := gatewayCollection.GetNextStart()
				if err != nil {
					o.Logger.Debugf("listTransitGateways: err = %v", err)
					return nil, fmt.Errorf("listTransitGateways: failed to GetNextStart: %w", err)
				}
				if start != nil {
					o.Logger.Debugf("listTransitGateways: start = %v", *start)
					listTransitGatewaysOptions.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("listTransitGateways: Next = nil")
				moreData = false
			}
		}
	}

	return cloudResources{}.insert(result...), nil
}

// Destroy a specified transit gateway.
func (o *ClusterUninstaller) destroyTransitGateway(item cloudResource) error {
	var (
		deleteTransitGatewayOptions *transitgatewayapisv1.DeleteTransitGatewayOptions
		response                    *core.DetailedResponse
		err                         error

		ctx    context.Context
		cancel func()
	)

	ctx, cancel = contextWithTimeout()
	defer cancel()

	err = o.destroyTransitGatewayConnections(item)
	if err != nil {
		return err
	}

	// We can delete the transit gateway now!
	deleteTransitGatewayOptions = o.tgClient.NewDeleteTransitGatewayOptions(item.id)

	response, err = o.tgClient.DeleteTransitGatewayWithContext(ctx, deleteTransitGatewayOptions)
	if err != nil {
		o.Logger.Fatalf("destroyTransitGateway: DeleteTransitGatewayWithContext returns %v with response %v", err, response)
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted Transit Gateway %q", item.name)

	return nil
}

// Destroy the connections for a specified transit gateway.
func (o *ClusterUninstaller) destroyTransitGatewayConnections(item cloudResource) error {
	var (
		firstPassList cloudResources

		err error

		items []cloudResource

		ctx    context.Context
		cancel func()

		backoff = wait.Backoff{Duration: 15 * time.Second,
			Factor: 1.5,
			Cap:    10 * time.Minute,
			Steps:  math.MaxInt32}
	)

	firstPassList, err = o.listTransitConnections(item)
	if err != nil {
		return err
	}

	items = o.insertPendingItems(transitGatewayConnectionTypeName, firstPassList.list())

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyTransitGateway: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.destroyTransitConnection(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatalf("destroyTransitGateway: ExponentialBackoffWithContext (destroy) returns %v", err)
		}
	}

	if items = o.getPendingItems(transitGatewayConnectionTypeName); len(items) > 0 {
		return fmt.Errorf("destroyTransitGateway: %d undeleted items pending", len(items))
	}

	backoff = wait.Backoff{Duration: 15 * time.Second,
		Factor: 1.5,
		Cap:    10 * time.Minute,
		Steps:  math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		var (
			secondPassList cloudResources

			err2 error
		)

		secondPassList, err2 = o.listTransitConnections(item)
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyTransitGateway: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatalf("destroyTransitGateway: ExponentialBackoffWithContext (list) returns %v", err)
	}

	return err
}

// Destroy a specified transit gateway connection.
func (o *ClusterUninstaller) destroyTransitConnection(item cloudResource) error {
	var (
		ctx    context.Context
		cancel func()

		deleteTransitGatewayConnectionOptions *transitgatewayapisv1.DeleteTransitGatewayConnectionOptions
		response                              *core.DetailedResponse
		err                                   error
	)

	ctx, cancel = contextWithTimeout()
	defer cancel()

	// ...Options(transitGatewayID string, id string)
	// NOTE: item.status is reused as the parent transit gateway id!
	deleteTransitGatewayConnectionOptions = o.tgClient.NewDeleteTransitGatewayConnectionOptions(item.status, item.id)

	response, err = o.tgClient.DeleteTransitGatewayConnectionWithContext(ctx, deleteTransitGatewayConnectionOptions)
	if err != nil {
		o.Logger.Fatalf("destroyTransitConnection: DeleteTransitGatewayConnectionWithContext returns %v with response %v", err, response)
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted Transit Gateway Connection %q", item.name)

	return nil
}

// listTransitConnections lists Transit Connections for a Transit Gateway in the IBM Cloud.
func (o *ClusterUninstaller) listTransitConnections(item cloudResource) (cloudResources, error) {
	o.Logger.Debugf("Listing Transit Gateways Connections (%s)", item.name)

	var (
		ctx                          context.Context
		cancel                       func()
		listConnectionsOptions       *transitgatewayapisv1.ListConnectionsOptions
		transitConnectionCollections *transitgatewayapisv1.TransitConnectionCollection
		transitConnection            transitgatewayapisv1.TransitConnection
		response                     *core.DetailedResponse
		err                          error
		foundOne                           = false
		perPage                      int64 = 32
		moreData                           = true
	)

	ctx, cancel = contextWithTimeout()
	defer cancel()

	o.Logger.Debugf("listTransitConnections: searching for ID %s", item.id)

	listConnectionsOptions = o.tgClient.NewListConnectionsOptions()
	listConnectionsOptions.SetLimit(perPage)
	listConnectionsOptions.SetNetworkID("")

	result := []cloudResource{}

	for moreData {
		transitConnectionCollections, response, err = o.tgClient.ListConnectionsWithContext(ctx, listConnectionsOptions)
		if err != nil {
			o.Logger.Debugf("listTransitConnections: ListConnections returns %v and the response is: %s", err, response)
			return nil, err
		}
		for _, transitConnection = range transitConnectionCollections.Connections {
			if *transitConnection.TransitGateway.ID != item.id {
				o.Logger.Debugf("listTransitConnections: SKIP: %s, %s, %s", *transitConnection.ID, *transitConnection.Name, *transitConnection.TransitGateway.Name)
				continue
			}

			foundOne = true
			o.Logger.Debugf("listTransitConnections: FOUND: %s, %s, %s", *transitConnection.ID, *transitConnection.Name, *transitConnection.TransitGateway.Name)
			result = append(result, cloudResource{
				key:      *transitConnection.ID,
				name:     *transitConnection.Name,
				status:   *transitConnection.TransitGateway.ID,
				typeName: transitGatewayConnectionTypeName,
				id:       *transitConnection.ID,
			})
		}

		if transitConnectionCollections.First != nil {
			o.Logger.Debugf("listTransitConnections: First = %+v", *transitConnectionCollections.First)
		} else {
			o.Logger.Debugf("listTransitConnections: First = nil")
		}
		if transitConnectionCollections.Limit != nil {
			o.Logger.Debugf("listTransitConnections: Limit = %v", *transitConnectionCollections.Limit)
		}
		if transitConnectionCollections.Next != nil {
			start, err := transitConnectionCollections.GetNextStart()
			if err != nil {
				o.Logger.Debugf("listTransitConnections: err = %v", err)
				return nil, fmt.Errorf("listTransitConnections: failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("listTransitConnections: start = %v", *start)
				listConnectionsOptions.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("listTransitConnections: Next = nil")
			moreData = false
		}
	}
	if !foundOne {
		o.Logger.Debugf("listTransitConnections: NO matching transit connections against: %s", o.InfraID)

		listConnectionsOptions = o.tgClient.NewListConnectionsOptions()
		listConnectionsOptions.SetLimit(perPage)
		listConnectionsOptions.SetNetworkID("")
		moreData = true

		for moreData {
			transitConnectionCollections, response, err = o.tgClient.ListConnectionsWithContext(ctx, listConnectionsOptions)
			if err != nil {
				o.Logger.Debugf("listTransitConnections: ListConnections returns %v and the response is: %s", err, response)
				return nil, err
			}
			for _, transitConnection = range transitConnectionCollections.Connections {
				o.Logger.Debugf("listTransitConnections: FOUND: %s, %s, %s", *transitConnection.ID, *transitConnection.Name, *transitConnection.TransitGateway.Name)
			}
			if transitConnectionCollections.First != nil {
				o.Logger.Debugf("listTransitConnections: First = %+v", *transitConnectionCollections.First)
			} else {
				o.Logger.Debugf("listTransitConnections: First = nil")
			}
			if transitConnectionCollections.Limit != nil {
				o.Logger.Debugf("listTransitConnections: Limit = %v", *transitConnectionCollections.Limit)
			}
			if transitConnectionCollections.Next != nil {
				start, err := transitConnectionCollections.GetNextStart()
				if err != nil {
					o.Logger.Debugf("listTransitConnections: err = %v", err)
					return nil, fmt.Errorf("listTransitConnections: failed to GetNextStart: %w", err)
				}
				if start != nil {
					o.Logger.Debugf("listTransitConnections: start = %v", *start)
					listConnectionsOptions.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("listTransitConnections: Next = nil")
				moreData = false
			}
		}
	}

	return cloudResources{}.insert(result...), nil
}

// We either deal with an existing TG or destroy TGs matching a name.
func (o *ClusterUninstaller) destroyTransitGateways() error {
	var (
		client *powervsconfig.Client
		tgID   string
		item   cloudResource
		err    error

		ctx    context.Context
		cancel func()
	)

	ctx, cancel = contextWithTimeout()
	defer cancel()

	// Old style: delete all TGs matching by name
	if o.TransitGatewayName == "" {
		return o.innerDestroyTransitGateways()
	}

	// New style: delete just TG connections for existing TG
	client, err = powervsconfig.NewClient()
	if err != nil {
		return err
	}

	tgID, err = client.TransitGatewayID(ctx, o.TransitGatewayName)
	if err != nil {
		return err
	}

	item = cloudResource{
		key:      tgID,
		name:     o.TransitGatewayName,
		status:   "",
		typeName: transitGatewayConnectionTypeName,
		id:       tgID,
	}

	err = o.destroyTransitGatewayConnections(item)

	return err
}

// innerDestroyTransitGateways searches for transit gateways that have a name that starts with
// the cluster's infra ID.
func (o *ClusterUninstaller) innerDestroyTransitGateways() error {
	var (
		firstPassList cloudResources

		err error

		items []cloudResource

		ctx    context.Context
		cancel func()

		backoff = wait.Backoff{Duration: 15 * time.Second,
			Factor: 1.5,
			Cap:    10 * time.Minute,
			Steps:  math.MaxInt32}
	)

	firstPassList, err = o.listTransitGateways()
	if err != nil {
		return err
	}

	items = o.insertPendingItems(transitGatewayTypeName, firstPassList.list())

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyTransitGateways: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.destroyTransitGateway(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatalf("destroyTransitGateways: ExponentialBackoffWithContext (destroy) returns %v", err)
		}
	}

	if items = o.getPendingItems(transitGatewayTypeName); len(items) > 0 {
		return fmt.Errorf("destroyTransitGateways: %d undeleted items pending", len(items))
	}

	backoff = wait.Backoff{Duration: 15 * time.Second,
		Factor: 1.5,
		Cap:    10 * time.Minute,
		Steps:  math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		var (
			secondPassList cloudResources

			err2 error
		)

		secondPassList, err2 = o.listTransitGateways()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyTransitGateways: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatalf("destroyTransitGateways: ExponentialBackoffWithContext (list) returns %v", err)
	}

	return nil
}
