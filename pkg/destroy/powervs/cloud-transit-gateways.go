package powervs

import (
	"context"
	"fmt"
	"math"
	gohttp "net/http"
	"strings"
	"time"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	transitGatewayTypeName           = "transitGateway"
	transitGatewayConnectionTypeName = "transitGatewayConnection"
)

// listTransitGateways lists Transit Gateways matching either name or tag in the IBM Cloud.
func (o *ClusterUninstaller) listTransitGateways() (cloudResources, error) {
	var (
		tgIDs    []string
		tgID     string
		ctx      context.Context
		cancel   context.CancelFunc
		result   = make([]cloudResource, 0, 1)
		options  *transitgatewayapisv1.GetTransitGatewayOptions
		tg       *transitgatewayapisv1.TransitGateway
		response *core.DetailedResponse
		err      error
	)

	if o.searchByTag {
		// Should we list by tag matching?
		tgIDs, err = o.listByTag(TagTypeTransitGateway)
	} else {
		// Otherwise list will list by name matching.
		tgIDs, err = o.listTransitGatewaysByName()
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, tgID = range tgIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listLoadBalancers: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		options = o.tgClient.NewGetTransitGatewayOptions(tgID)

		tg, response, err = o.tgClient.GetTransitGatewayWithContext(ctx, options)
		if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
			// The transit gateway could have been deleted just after a list was created.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get transit gateway (%s): err = %w, response = %v", tgID, err, response)
		}

		result = append(result, cloudResource{
			key:      *tg.ID,
			name:     *tg.Name,
			status:   "",
			typeName: transitGatewayTypeName,
			id:       *tg.ID,
		})
	}

	return cloudResources{}.insert(result...), nil
}

// listTransitGatewaysByName lists Transit Gateways matching by name in the IBM Cloud.
func (o *ClusterUninstaller) listTransitGatewaysByName() ([]string, error) {
	var (
		ctx                        context.Context
		cancel                     context.CancelFunc
		listTransitGatewaysOptions *transitgatewayapisv1.ListTransitGatewaysOptions
		gatewayCollection          *transitgatewayapisv1.TransitGatewayCollection
		gateway                    transitgatewayapisv1.TransitGateway
		response                   *core.DetailedResponse
		err                        error
		foundOne                         = false
		perPage                    int64 = 32
		moreData                         = true
		result                           = make([]string, 0, 1)
	)

	o.Logger.Debugf("Listing Transit Gateways (%s) by NAME", o.InfraID)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listTransitGatewaysByName: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	listTransitGatewaysOptions = o.tgClient.NewListTransitGatewaysOptions()
	listTransitGatewaysOptions.Limit = &perPage

	for moreData {
		// https://github.com/IBM/networking-go-sdk/blob/master/transitgatewayapisv1/transit_gateway_apis_v1.go#L184
		gatewayCollection, response, err = o.tgClient.ListTransitGatewaysWithContext(ctx, listTransitGatewaysOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list transit gateways: %w and the respose is: %s", err, response)
		}

		for _, gateway = range gatewayCollection.TransitGateways {
			if strings.Contains(*gateway.Name, o.InfraID) {
				foundOne = true
				o.Logger.Debugf("listTransitGatewaysByName: FOUND: %s, %s", *gateway.ID, *gateway.Name)
				result = append(result, *gateway.ID)
			}
		}

		if gatewayCollection.First != nil {
			o.Logger.Debugf("listTransitGatewaysByName: First = %+v", *gatewayCollection.First.Href)
		} else {
			o.Logger.Debugf("listTransitGatewaysByName: First = nil")
		}
		if gatewayCollection.Limit != nil {
			o.Logger.Debugf("listTransitGatewaysByName: Limit = %v", *gatewayCollection.Limit)
		}
		if gatewayCollection.Next != nil {
			start, err := gatewayCollection.GetNextStart()
			if err != nil {
				o.Logger.Debugf("listTransitGatewaysByName: err = %v", err)
				return nil, fmt.Errorf("listTransitGatewaysByName: failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("listTransitGatewaysByName: start = %v", *start)
				listTransitGatewaysOptions.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("listTransitGatewaysByName: Next = nil")
			moreData = false
		}
	}
	if !foundOne {
		o.Logger.Debugf("listTransitGatewaysByName: NO matching transit gateway against: %s", o.InfraID)

		listTransitGatewaysOptions = o.tgClient.NewListTransitGatewaysOptions()
		listTransitGatewaysOptions.Limit = &perPage
		moreData = true

		for moreData {
			gatewayCollection, response, err = o.tgClient.ListTransitGatewaysWithContext(ctx, listTransitGatewaysOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to list transit gateways: %w and the respose is: %s", err, response)
			}
			for _, gateway = range gatewayCollection.TransitGateways {
				o.Logger.Debugf("listTransitGatewaysByName: FOUND: %s, %s", *gateway.ID, *gateway.Name)
			}
			if gatewayCollection.First != nil {
				o.Logger.Debugf("listTransitGatewaysByName: First = %+v", *gatewayCollection.First.Href)
			} else {
				o.Logger.Debugf("listTransitGatewaysByName: First = nil")
			}
			if gatewayCollection.Limit != nil {
				o.Logger.Debugf("listTransitGatewaysByName: Limit = %v", *gatewayCollection.Limit)
			}
			if gatewayCollection.Next != nil {
				start, err := gatewayCollection.GetNextStart()
				if err != nil {
					o.Logger.Debugf("listTransitGatewaysByName: err = %v", err)
					return nil, fmt.Errorf("listTransitGatewaysByName: failed to GetNextStart: %w", err)
				}
				if start != nil {
					o.Logger.Debugf("listTransitGatewaysByName: start = %v", *start)
					listTransitGatewaysOptions.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("listTransitGatewaysByName: Next = nil")
				moreData = false
			}
		}
	}

	return result, nil
}

// removeTGSIConnection finds and removes the specified Service Instance connection in any Transit Gateway.
func (o *ClusterUninstaller) removeTGSIConnection(siID string) error {
	var (
		ctx                          context.Context
		cancel                       context.CancelFunc
		listConnectionsOptions       *transitgatewayapisv1.ListConnectionsOptions
		transitConnectionCollections *transitgatewayapisv1.TransitConnectionCollection
		transitConnection            transitgatewayapisv1.TransitConnection
		crnStruct                    crn.CRN
		item                         cloudResource
		response                     *core.DetailedResponse
		err                          error
		perPage                      int64 = 32
		moreData                           = true
	)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	listConnectionsOptions = o.tgClient.NewListConnectionsOptions()
	listConnectionsOptions.SetLimit(perPage)
	listConnectionsOptions.SetNetworkID("")

	for moreData {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("removeTGSIConnection: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		transitConnectionCollections, response, err = o.tgClient.ListConnectionsWithContext(ctx, listConnectionsOptions)
		if err != nil {
			o.Logger.Debugf("removeTGSIConnection: ListConnections returns %v and the response is: %s", err, response)
			return err
		}

		for _, transitConnection = range transitConnectionCollections.Connections {
			select {
			case <-ctx.Done():
				o.Logger.Debugf("removeTGSIConnection: case <-ctx.Done()")
				return ctx.Err() // we're cancelled, abort
			default:
			}

			if *transitConnection.NetworkType != transitgatewayapisv1.TransitConnection_NetworkType_PowerVirtualServer {
				continue
			}

			o.Logger.Debugf("removeTGSIConnection: transitConnection.NetworkID = %s", *transitConnection.NetworkID)
			o.Logger.Debugf("removeTGSIConnection: transitConnection.Name = %s", *transitConnection.Name)
			o.Logger.Debugf("removeTGSIConnection: transitConnection.TransitGateway.Name = %s", *transitConnection.TransitGateway.Name)

			crnStruct, err = crn.Parse(*transitConnection.NetworkID)
			if err != nil {
				o.Logger.Debugf("removeTGSIConnection: failed to crn.Parse: %v", err)
				return fmt.Errorf("removeTGSIConnection: failed to crn.Parse: %w", err)
			}
			o.Logger.Debugf("removeTGSIConnection: crnStruct.ServiceInstance = %v", crnStruct.ServiceInstance)

			if crnStruct.ServiceInstance != siID {
				continue
			}

			item = cloudResource{
				key:      *transitConnection.ID,
				name:     *transitConnection.Name,
				status:   *transitConnection.TransitGateway.ID,
				typeName: transitGatewayConnectionTypeName,
				id:       *transitConnection.ID,
			}
			o.Logger.Debugf("removeTGSIConnection: item = %+v", item)

			err = o.destroyAndWaitTransitConnection(siID, item)
			return err
		}

		if transitConnectionCollections.First != nil {
			o.Logger.Debugf("removeTGSIConnection: First = %+v", *transitConnectionCollections.First)
		} else {
			o.Logger.Debugf("removeTGSIConnection: First = nil")
		}
		if transitConnectionCollections.Limit != nil {
			o.Logger.Debugf("removeTGSIConnection: Limit = %v", *transitConnectionCollections.Limit)
		}
		if transitConnectionCollections.Next != nil {
			start, err := transitConnectionCollections.GetNextStart()
			if err != nil {
				o.Logger.Debugf("removeTGSIConnection: err = %v", err)
				return fmt.Errorf("removeTGSIConnection: failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("removeTGSIConnection: start = %v", *start)
				listConnectionsOptions.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("removeTGSIConnection: Next = nil")
			moreData = false
		}
	}

	return nil
}

// destroyTransitGateway destroy a specified transit gateway.
func (o *ClusterUninstaller) destroyTransitGateway(item cloudResource) error {
	var (
		deleteTransitGatewayOptions *transitgatewayapisv1.DeleteTransitGatewayOptions
		response                    *core.DetailedResponse
		err                         error

		ctx    context.Context
		cancel context.CancelFunc
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

// destroyTransitGatewayConnections destroy the connections for a specified transit gateway.
func (o *ClusterUninstaller) destroyTransitGatewayConnections(item cloudResource) error {
	var (
		firstPassList cloudResources

		err error

		items []cloudResource

		ctx    context.Context
		cancel context.CancelFunc

		backoff = wait.Backoff{Duration: 15 * time.Second,
			Factor: 1.5,
			Cap:    10 * time.Minute,
			Steps:  math.MaxInt32}
	)

	firstPassList, err = o.listTransitConnectionsByName(item)
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

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyTransitGateway: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
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

		secondPassList, err2 = o.listTransitConnectionsByName(item)
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

// destroyTransitConnection destroy a specified transit gateway connection.
func (o *ClusterUninstaller) destroyTransitConnection(item cloudResource) error {
	var (
		ctx    context.Context
		cancel context.CancelFunc

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

// destroyAndWaitTransitConnection destroy a specified transit gateway connection and wait for it to complete.
func (o *ClusterUninstaller) destroyAndWaitTransitConnection(siID string, item cloudResource) error {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		err    error
	)

	err = o.destroyTransitConnection(item)
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	backoff := wait.Backoff{Duration: 15 * time.Second,
		Factor: 1.5,
		Cap:    10 * time.Minute,
		Steps:  math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		var (
			listConnectionsOptions       *transitgatewayapisv1.ListConnectionsOptions
			transitConnectionCollections *transitgatewayapisv1.TransitConnectionCollection
			transitConnection            transitgatewayapisv1.TransitConnection
			crnStruct                    crn.CRN
			response                     *core.DetailedResponse
			perPage                      int64 = 32
			moreData                           = true
			err2                         error
		)

		listConnectionsOptions = o.tgClient.NewListConnectionsOptions()
		listConnectionsOptions.SetLimit(perPage)
		listConnectionsOptions.SetNetworkID("")

		for moreData {
			select {
			case <-ctx.Done():
				o.Logger.Debugf("destroyAndWaitTransitConnection: case <-ctx.Done()")
				return false, ctx.Err() // we're cancelled, abort
			default:
			}

			transitConnectionCollections, response, err2 = o.tgClient.ListConnectionsWithContext(ctx, listConnectionsOptions)
			if err2 != nil {
				o.Logger.Debugf("destroyAndWaitTransitConnection: ListConnections returns %v and the response is: %s", err2, response)
				return false, err2
			}

			for _, transitConnection = range transitConnectionCollections.Connections {
				select {
				case <-ctx.Done():
					o.Logger.Debugf("destroyAndWaitTransitConnection: case <-ctx.Done()")
					return false, ctx.Err() // we're cancelled, abort
				default:
				}

				if *transitConnection.NetworkType != transitgatewayapisv1.TransitConnection_NetworkType_PowerVirtualServer {
					continue
				}

				crnStruct, err2 = crn.Parse(*transitConnection.NetworkID)
				if err2 != nil {
					o.Logger.Debugf("destroyAndWaitTransitConnection: failed to crn.Parse: %v", err2)
					return false, fmt.Errorf("destroyAndWaitTransitConnection: failed to crn.Parse: %w", err2)
				}

				if crnStruct.ServiceInstance != siID {
					o.Logger.Debugf("destroyAndWaitTransitConnection: SKIP  %s", crnStruct.ServiceInstance)
					continue
				}
				o.Logger.Debugf("destroyAndWaitTransitConnection: FOUND %s", crnStruct.ServiceInstance)

				// We have found a connection!
				return false, nil
			}

			if transitConnectionCollections.First != nil {
				o.Logger.Debugf("destroyAndWaitTransitConnection: First = %+v", *transitConnectionCollections.First)
			} else {
				o.Logger.Debugf("destroyAndWaitTransitConnection: First = nil")
			}
			if transitConnectionCollections.Limit != nil {
				o.Logger.Debugf("destroyAndWaitTransitConnection: Limit = %v", *transitConnectionCollections.Limit)
			}
			if transitConnectionCollections.Next != nil {
				start, err2 := transitConnectionCollections.GetNextStart()
				if err2 != nil {
					o.Logger.Debugf("destroyAndWaitTransitConnection: err2 = %v", err2)
					return false, fmt.Errorf("destroyAndWaitTransitConnection: failed to GetNextStart: %w", err2)
				}
				if start != nil {
					o.Logger.Debugf("destroyAndWaitTransitConnection: start = %v", *start)
					listConnectionsOptions.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("destroyAndWaitTransitConnection: Next = nil")
				moreData = false
			}
		}

		// We haven't found any connection!
		return true, nil
	})
	if err != nil {
		o.Logger.Fatalf("destroyAndWaitTransitConnection: ExponentialBackoffWithContext returns %v", err)
	}

	return nil
}

// listTransitConnectionsByName lists Transit Connections for a Transit Gateway in the IBM Cloud.
func (o *ClusterUninstaller) listTransitConnectionsByName(item cloudResource) (cloudResources, error) {
	o.Logger.Debugf("Listing Transit Gateways Connections (%s)", item.name)

	var (
		ctx                          context.Context
		cancel                       context.CancelFunc
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

	o.Logger.Debugf("listTransitConnectionsByName: searching for ID %s", item.id)

	listConnectionsOptions = o.tgClient.NewListConnectionsOptions()
	listConnectionsOptions.SetLimit(perPage)
	listConnectionsOptions.SetNetworkID("")

	result := []cloudResource{}

	for moreData {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listTransitConnectionsByName: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		transitConnectionCollections, response, err = o.tgClient.ListConnectionsWithContext(ctx, listConnectionsOptions)
		if err != nil {
			o.Logger.Debugf("listTransitConnectionsByName: ListConnections returns %v and the response is: %s", err, response)
			return nil, err
		}
		for _, transitConnection = range transitConnectionCollections.Connections {
			if *transitConnection.TransitGateway.ID != item.id {
				o.Logger.Debugf("listTransitConnectionsByName: SKIP: %s, %s, %s", *transitConnection.ID, *transitConnection.Name, *transitConnection.TransitGateway.Name)
				continue
			}

			foundOne = true
			o.Logger.Debugf("listTransitConnectionsByName: FOUND: %s, %s, %s", *transitConnection.ID, *transitConnection.Name, *transitConnection.TransitGateway.Name)
			result = append(result, cloudResource{
				key:      *transitConnection.ID,
				name:     *transitConnection.Name,
				status:   *transitConnection.TransitGateway.ID,
				typeName: transitGatewayConnectionTypeName,
				id:       *transitConnection.ID,
			})
		}

		if transitConnectionCollections.First != nil {
			o.Logger.Debugf("listTransitConnectionsByName: First = %+v", *transitConnectionCollections.First)
		} else {
			o.Logger.Debugf("listTransitConnectionsByName: First = nil")
		}
		if transitConnectionCollections.Limit != nil {
			o.Logger.Debugf("listTransitConnectionsByName: Limit = %v", *transitConnectionCollections.Limit)
		}
		if transitConnectionCollections.Next != nil {
			start, err := transitConnectionCollections.GetNextStart()
			if err != nil {
				o.Logger.Debugf("listTransitConnectionsByName: err = %v", err)
				return nil, fmt.Errorf("listTransitConnectionsByName: failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("listTransitConnectionsByName: start = %v", *start)
				listConnectionsOptions.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("listTransitConnectionsByName: Next = nil")
			moreData = false
		}
	}
	if !foundOne {
		o.Logger.Debugf("listTransitConnectionsByName: NO matching transit connections against: %s", o.InfraID)

		listConnectionsOptions = o.tgClient.NewListConnectionsOptions()
		listConnectionsOptions.SetLimit(perPage)
		listConnectionsOptions.SetNetworkID("")
		moreData = true

		for moreData {
			select {
			case <-ctx.Done():
				o.Logger.Debugf("listTransitConnectionsByName: case <-ctx.Done()")
				return nil, ctx.Err() // we're cancelled, abort
			default:
			}

			transitConnectionCollections, response, err = o.tgClient.ListConnectionsWithContext(ctx, listConnectionsOptions)
			if err != nil {
				o.Logger.Debugf("listTransitConnectionsByName: ListConnections returns %v and the response is: %s", err, response)
				return nil, err
			}
			for _, transitConnection = range transitConnectionCollections.Connections {
				o.Logger.Debugf("listTransitConnectionsByName: FOUND: %s, %s, %s", *transitConnection.ID, *transitConnection.Name, *transitConnection.TransitGateway.Name)
			}
			if transitConnectionCollections.First != nil {
				o.Logger.Debugf("listTransitConnectionsByName: First = %+v", *transitConnectionCollections.First)
			} else {
				o.Logger.Debugf("listTransitConnectionsByName: First = nil")
			}
			if transitConnectionCollections.Limit != nil {
				o.Logger.Debugf("listTransitConnectionsByName: Limit = %v", *transitConnectionCollections.Limit)
			}
			if transitConnectionCollections.Next != nil {
				start, err := transitConnectionCollections.GetNextStart()
				if err != nil {
					o.Logger.Debugf("listTransitConnectionsByName: err = %v", err)
					return nil, fmt.Errorf("listTransitConnectionsByName: failed to GetNextStart: %w", err)
				}
				if start != nil {
					o.Logger.Debugf("listTransitConnectionsByName: start = %v", *start)
					listConnectionsOptions.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("listTransitConnectionsByName: Next = nil")
				moreData = false
			}
		}
	}

	return cloudResources{}.insert(result...), nil
}

// destroyTransitGateways we either deal with an existing TG or destroy TGs matching a name.
func (o *ClusterUninstaller) destroyTransitGateways() error {
	var (
		err error
	)

	// Old style: delete all TGs matching by name
	if o.TransitGatewayName == "" {
		return o.innerDestroyTransitGateways()
	}

	// New style: before we can delete the created Service Instance, we need to remove the
	// Transit Gateway connections!
	if !o.siPreconfigured {
		err = o.removeTGSIConnection(o.ServiceGUID)
		if err != nil {
			return err
		}
	}

	// New style: leave the TG and its existing connections alone
	o.Logger.Infof("Not cleaning up persistent Transit Gateway since tgName was specified")
	return nil
}

// innerDestroyTransitGateways searches for transit gateways that have a name that starts with
// the cluster's infra ID.
func (o *ClusterUninstaller) innerDestroyTransitGateways() error {
	var (
		firstPassList cloudResources

		err error

		items []cloudResource

		ctx    context.Context
		cancel context.CancelFunc

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
			o.Logger.Debugf("innerDestroyTransitGateways: case <-ctx.Done()")
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
			o.Logger.Fatalf("innerDestroyTransitGateways: ExponentialBackoffWithContext (destroy) returns %v", err)
		}
	}

	if items = o.getPendingItems(transitGatewayTypeName); len(items) > 0 {
		return fmt.Errorf("innerDestroyTransitGateways: %d undeleted items pending", len(items))
	}

	select {
	case <-ctx.Done():
		o.Logger.Debugf("innerDestroyTransitGateways: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
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
			o.Logger.Debugf("innerDestroyTransitGateways: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatalf("innerDestroyTransitGateways: ExponentialBackoffWithContext (list) returns %v", err)
	}

	return nil
}
