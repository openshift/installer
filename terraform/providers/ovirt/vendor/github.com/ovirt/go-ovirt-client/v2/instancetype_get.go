package ovirtclient //nolint:dupl

import (
	"fmt"
)

func (o *oVirtClient) GetInstanceType(id InstanceTypeID, retries ...RetryStrategy) (result InstanceType, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting instance type %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().InstanceTypesService().InstanceTypeService(string(id)).Get().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.InstanceType()
			if !ok {
				return newError(
					ENotFound,
					"no instance type returned when getting instance type ID %s",
					id,
				)
			}
			result, err = convertSDKInstanceType(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert instance type %s",
					id,
				)
			}
			return nil
		})
	return result, err
}

func (m *mockClient) GetInstanceType(id InstanceTypeID, _ ...RetryStrategy) (InstanceType, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if item, ok := m.instanceTypes[id]; ok {
		return item, nil
	}
	return nil, newError(ENotFound, "instance type with ID %s not found", id)
}
