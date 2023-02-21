package ovirtclient

import "fmt"

func (o *oVirtClient) RemoveAffinityGroup(clusterID ClusterID, id AffinityGroupID, retries ...RetryStrategy) error {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	return retry(
		fmt.Sprintf("removing affinity group %s from cluster %s", id, clusterID),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.
				SystemService().
				ClustersService().
				ClusterService(string(clusterID)).
				AffinityGroupsService().
				GroupService(string(id)).
				Remove().
				Send()
			return err
		},
	)
}

func (m *mockClient) RemoveAffinityGroup(clusterID ClusterID, id AffinityGroupID, retries ...RetryStrategy) error {

	retries = defaultRetries(retries, defaultWriteTimeouts(m))

	return retry(
		fmt.Sprintf("removing affinity group %s from cluster %s", id, clusterID),
		m.logger,
		retries,
		func() error {
			m.lock.Lock()
			defer m.lock.Unlock()

			clusterAffinityGroups, ok := m.affinityGroups[clusterID]
			if !ok {
				return newError(ENotFound, "Cluster with ID %s not found", clusterID)
			}
			if _, ok := clusterAffinityGroups[id]; !ok {
				return newError(ENotFound, "Affinity group with ID %s not found in cluster %s", id, clusterID)
			}

			delete(m.affinityGroups[clusterID], id)

			return nil
		})
}
