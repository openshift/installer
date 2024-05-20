// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetNetwork(id NetworkID, retries ...RetryStrategy) (result Network, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting network %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().NetworksService().NetworkService(string(id)).Get().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Network()
			if !ok {
				return newError(
					ENotFound,
					"no network returned when getting network ID %s",
					id,
				)
			}
			result, err = convertSDKNetwork(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert network %s",
					id,
				)
			}
			return nil
		})
	return
}

func (m *mockClient) GetNetwork(id NetworkID, _ ...RetryStrategy) (Network, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if item, ok := m.networks[id]; ok {
		return item, nil
	}
	return nil, newError(ENotFound, "network with ID %s not found", id)
}
