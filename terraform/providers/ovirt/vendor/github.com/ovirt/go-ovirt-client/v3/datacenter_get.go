// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetDatacenter(id DatacenterID, retries ...RetryStrategy) (result Datacenter, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting datacenter %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().DataCentersService().DataCenterService(string(id)).Get().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.DataCenter()
			if !ok {
				return newError(
					ENotFound,
					"no datacenter returned when getting datacenter ID %s",
					id,
				)
			}
			result, err = convertSDKDatacenter(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert datacenter %s",
					id,
				)
			}
			return nil
		})
	return
}

func (m *mockClient) GetDatacenter(id DatacenterID, _ ...RetryStrategy) (Datacenter, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if item, ok := m.dataCenters[id]; ok {
		return item, nil
	}
	return nil, newError(ENotFound, "datacenter with ID %s not found", id)
}
