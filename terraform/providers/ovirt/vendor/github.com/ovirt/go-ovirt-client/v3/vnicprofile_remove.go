package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) RemoveVNICProfile(id VNICProfileID, retries ...RetryStrategy) (err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("removing VNIC profile %s", id),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().VnicProfilesService().ProfileService(string(id)).Remove().Send()
			if err != nil {
				return err
			}
			return nil
		})
	return
}

func (m *mockClient) RemoveVNICProfile(id VNICProfileID, _ ...RetryStrategy) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.vnicProfiles[id]; !ok {
		return newError(ENotFound, "VNIC profile not found")
	}

	delete(m.vnicProfiles, id)

	return nil
}
