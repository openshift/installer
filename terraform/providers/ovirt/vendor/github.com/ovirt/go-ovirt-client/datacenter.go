package ovirtclient

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

//go:generate go run scripts/rest/rest.go -i "DataCenter" -n "datacenter" -o "Datacenter" -T DatacenterID

// DatacenterID is the UUID of a datacenter.
type DatacenterID string

// DatacenterClient contains the functions related to handling datacenter objects in oVirt. Datacenters bind together
// resources of an environment (clusters, storage domains).
// See https://www.ovirt.org/documentation/administration_guide/#chap-Data_Centers for details.
type DatacenterClient interface {
	// GetDatacenter returns a single datacenter by its ID.
	GetDatacenter(id DatacenterID, retries ...RetryStrategy) (Datacenter, error)
	// ListDatacenters lists all datacenters in the oVirt engine.
	ListDatacenters(retries ...RetryStrategy) ([]Datacenter, error)
	// ListDatacenterClusters lists all clusters in the specified datacenter.
	ListDatacenterClusters(id DatacenterID, retries ...RetryStrategy) ([]Cluster, error)
}

// DatacenterData is the core of a Datacenter when client functions are not required.
type DatacenterData interface {
	ID() DatacenterID
	Name() string
}

// Datacenter is a logical entity that defines the set of resources used in a specific environment.
// See https://www.ovirt.org/documentation/administration_guide/#chap-Data_Centers for details.
type Datacenter interface {
	DatacenterData

	// Clusters lists the clusters for this datacenter. This is a network call and may be slow.
	Clusters(retries ...RetryStrategy) ([]Cluster, error)
	// HasCluster returns true if the cluster is in the datacenter. This is a network call and may be slow.
	HasCluster(clusterID ClusterID, retries ...RetryStrategy) (bool, error)
}

func convertSDKDatacenter(sdkObject *ovirtsdk4.DataCenter, client *oVirtClient) (Datacenter, error) {
	id, ok := sdkObject.Id()
	if !ok {
		return nil, newFieldNotFound("datacenter", "id")
	}
	name, ok := sdkObject.Name()
	if !ok {
		return nil, newFieldNotFound("datacenter", "name")
	}

	return &datacenter{
		client: client,
		id:     DatacenterID(id),
		name:   name,
	}, nil
}

type datacenter struct {
	client Client

	id   DatacenterID
	name string
}

func (d datacenter) Clusters(retries ...RetryStrategy) ([]Cluster, error) {
	return d.client.ListDatacenterClusters(d.id, retries...)
}

func (d datacenter) HasCluster(clusterID ClusterID, retries ...RetryStrategy) (bool, error) {
	clusters, err := d.client.ListDatacenterClusters(d.id, retries...)
	if err != nil {
		return false, err
	}
	for _, cluster := range clusters {
		if cluster.ID() == clusterID {
			return true, nil
		}
	}
	return false, nil
}

func (d datacenter) ID() DatacenterID {
	return d.id
}

func (d datacenter) Name() string {
	return d.name
}
