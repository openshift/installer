package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) GetAffinityGroupByName(clusterID ClusterID, name string, retries ...RetryStrategy) (result AffinityGroup, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting affinity group %s", name),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().ClustersService().ClusterService(string(clusterID)).AffinityGroupsService().List().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Groups()
			if !ok {
				return newError(
					ENotFound,
					"no affinity group returned when listing affinity groups in cluster ID %s",
					clusterID,
				)
			}
			slice := sdkObject.Slice()
			if len(slice) == 0 {
				return newError(ENotFound, "no affinity group named %s found in cluster %s", name, clusterID)
			}
			var results []*ovirtsdk.AffinityGroup
			for _, item := range slice {
				n, ok := item.Name()
				if ok && n == name {
					results = append(results, item)
				}
			}
			if len(results) == 0 {
				return newError(ENotFound, "no affinity group named %s found in cluster %s", name, clusterID)
			}
			if len(results) > 1 {
				return newError(EMultipleResults, "multiple affinity groups with the name %s found in cluster %s", name, clusterID)
			}
			result, err = convertSDKAffinityGroup(results[0], o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert affinity group %s",
					name,
				)
			}
			return nil
		})
	return result, err
}

func (m *mockClient) GetAffinityGroupByName(clusterID ClusterID, name string, retries ...RetryStrategy) (result AffinityGroup, err error) {

	retries = defaultRetries(retries, defaultWriteTimeouts(m))

	err = retry(
		fmt.Sprintf("getting affinity group %s from cluster %s", name, clusterID),
		m.logger,
		retries,
		func() error {
			m.lock.Lock()
			defer m.lock.Unlock()

			clusterAffinityGroups, ok := m.affinityGroups[clusterID]
			if !ok {
				return newError(ENotFound, "Cluster with ID %s not found", clusterID)
			}
			for _, ag := range clusterAffinityGroups {
				if ag.name == name {
					result = ag
					return nil
				}
			}
			return newError(ENotFound, "Affinity group with name %s not found in cluster %s", name, clusterID)
		})
	return
}
