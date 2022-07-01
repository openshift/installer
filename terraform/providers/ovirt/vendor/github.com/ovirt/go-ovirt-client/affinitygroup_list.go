package ovirtclient

import "fmt"

// nolint:dupl
func (o *oVirtClient) ListAffinityGroups(
	clusterID ClusterID,
	retries ...RetryStrategy,
) (result []AffinityGroup, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []AffinityGroup{}
	err = retry(
		fmt.Sprintf("listing affinity groups in cluster %s", clusterID),
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().ClustersService().ClusterService(string(clusterID)).AffinityGroupsService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.Groups()
			if !ok {
				return nil
			}
			result = make([]AffinityGroup, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKAffinityGroup(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert affinity group during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListAffinityGroups(
	clusterID ClusterID,
	_ ...RetryStrategy,
) ([]AffinityGroup, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	result := make([]AffinityGroup, len(m.affinityGroups[clusterID]))
	i := 0
	for _, affinityGroup := range m.affinityGroups[clusterID] {
		result[i] = affinityGroup
		i++
	}
	return result, nil
}
