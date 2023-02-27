// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetTag(id TagID, retries ...RetryStrategy) (result Tag, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting tag %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().TagsService().TagService(string(id)).Get().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Tag()
			if !ok {
				return newError(
					ENotFound,
					"no tag returned when getting tag ID %s",
					id,
				)
			}
			result, err = convertSDKTag(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert tag %s",
					id,
				)
			}
			return nil
		})
	return
}

func (m *mockClient) GetTag(id TagID, _ ...RetryStrategy) (Tag, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if item, ok := m.tags[id]; ok {
		return item, nil
	}
	return nil, newError(ENotFound, "tag with ID %s not found", id)
}
