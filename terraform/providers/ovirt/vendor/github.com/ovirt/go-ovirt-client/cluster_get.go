// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetCluster(id ClusterID, retries ...RetryStrategy) (result Cluster, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting cluster %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().ClustersService().ClusterService(string(id)).Get().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Cluster()
			if !ok {
				return newError(
					ENotFound,
					"no cluster returned when getting cluster ID %s",
					id,
				)
			}
			result, err = convertSDKCluster(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert cluster %s",
					id,
				)
			}
			return nil
		})
	return
}

func (m *mockClient) GetCluster(id ClusterID, _ ...RetryStrategy) (Cluster, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if item, ok := m.clusters[id]; ok {
		return item, nil
	}
	return nil, newError(ENotFound, "cluster with ID %s not found", id)
}
