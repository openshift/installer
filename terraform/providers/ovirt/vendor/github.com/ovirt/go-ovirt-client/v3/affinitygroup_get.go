package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetAffinityGroup(clusterID ClusterID, id AffinityGroupID, retries ...RetryStrategy) (result AffinityGroup, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting affinity group %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().ClustersService().ClusterService(string(clusterID)).AffinityGroupsService().GroupService(string(id)).Get().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Group()
			if !ok {
				return newError(
					ENotFound,
					"no affinity group returned when getting affinity group ID %s in cluster ID %s",
					id,
					clusterID,
				)
			}
			result, err = convertSDKAffinityGroup(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert affinity group %s",
					id,
				)
			}
			return nil
		})
	return result, err
}

func (m *mockClient) GetAffinityGroup(clusterID ClusterID, id AffinityGroupID, retries ...RetryStrategy) (result AffinityGroup, err error) {

	retries = defaultRetries(retries, defaultWriteTimeouts(m))

	err = retry(
		fmt.Sprintf("getting affinity group %s from cluster %s", id, clusterID),
		m.logger,
		retries,
		func() error {
			m.lock.Lock()
			defer m.lock.Unlock()

			clusterAffinityGroups, ok := m.affinityGroups[clusterID]
			if !ok {
				return newError(ENotFound, "Cluster with ID %s not found", clusterID)
			}
			result, ok = clusterAffinityGroups[id]
			if !ok {
				return newError(ENotFound, "Affinity group with ID %s not found in cluster %s", id, clusterID)
			}

			return nil
		})
	return
}
