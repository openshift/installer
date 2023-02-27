// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

func (o *oVirtClient) ListClusters(retries ...RetryStrategy) (result []Cluster, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []Cluster{}
	err = retry(
		"listing clusters",
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().ClustersService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.Clusters()
			if !ok {
				return nil
			}
			result = make([]Cluster, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKCluster(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert cluster during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListClusters(_ ...RetryStrategy) ([]Cluster, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]Cluster, len(m.clusters))
	i := 0
	for _, item := range m.clusters {
		result[i] = item
		i++
	}
	return result, nil
}
