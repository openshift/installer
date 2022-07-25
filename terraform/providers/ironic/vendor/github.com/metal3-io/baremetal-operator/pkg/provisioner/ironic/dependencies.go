package ironic

import (
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/drivers"
	"github.com/gophercloud/gophercloud/pagination"
)

const (
	checkRequeueDelay = time.Second * 10
)

type ironicDependenciesChecker struct {
	client    *gophercloud.ServiceClient
	inspector *gophercloud.ServiceClient
	log       logr.Logger
}

func newIronicDependenciesChecker(client *gophercloud.ServiceClient, inspector *gophercloud.ServiceClient, log logr.Logger) *ironicDependenciesChecker {
	return &ironicDependenciesChecker{
		client:    client,
		inspector: inspector,
		log:       log,
	}
}

func (i *ironicDependenciesChecker) IsReady() (result bool, err error) {

	ready, err := i.checkIronic()
	if ready && err == nil {
		ready = i.checkIronicInspector()
	}

	return ready, err
}

func (i *ironicDependenciesChecker) checkEndpoint(client *gophercloud.ServiceClient) (ready bool) {

	// NOTE: Some versions of Ironic inspector returns 404 for /v1/ but 200 for /v1,
	// which seems to be the default behavior for Flask. Remove the trailing slash
	// from the client endpoint.
	endpoint := strings.TrimSuffix(client.Endpoint, "/")

	_, err := client.Get(endpoint, nil, nil)
	if err != nil {
		i.log.Info("error caught while checking endpoint", "endpoint", client.Endpoint, "error", err)
	}

	return err == nil
}

func (i *ironicDependenciesChecker) checkIronic() (ready bool, err error) {
	ready = i.checkEndpoint(i.client)
	if ready {
		ready, err = i.checkIronicConductor()
	}
	return ready, err
}

func (i *ironicDependenciesChecker) checkIronicConductor() (ready bool, err error) {

	pager := drivers.ListDrivers(i.client, drivers.ListDriversOpts{
		Detail: false,
	})
	err = pager.Err

	if err != nil {
		return ready, err
	}

	driverCount := 0
	pager.EachPage(func(page pagination.Page) (bool, error) {
		actual, driverErr := drivers.ExtractDrivers(page)
		if driverErr != nil {
			return false, driverErr
		}
		driverCount += len(actual)
		return true, nil
	})
	// If we have any drivers, conductor is up.
	ready = driverCount > 0

	return ready, err
}

func (i *ironicDependenciesChecker) checkIronicInspector() (ready bool) {
	return i.checkEndpoint(i.inspector)
}
