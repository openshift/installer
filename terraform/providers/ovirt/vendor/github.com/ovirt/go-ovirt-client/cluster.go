package ovirtclient

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

//go:generate go run scripts/rest/rest.go -i "Cluster" -n "cluster" -T ClusterID

// ClusterClient is a part of the Client that deals with clusters in the oVirt Engine. A cluster is a logical grouping
// of hosts that share the same storage domains and have the same type of CPU (either Intel or AMD). If the hosts have
// different generations of CPU models, they use only the features present in all models.
//
// See https://www.ovirt.org/documentation/administration_guide/#chap-Clusters for details.
type ClusterClient interface {
	// ListClusters returns a list of all clusters in the oVirt engine.
	ListClusters(retries ...RetryStrategy) ([]Cluster, error)
	// GetCluster returns a specific cluster based on the cluster ID. An error is returned if the cluster doesn't exist.
	GetCluster(id ClusterID, retries ...RetryStrategy) (Cluster, error)
}

// ClusterID is an identifier for a cluster.
type ClusterID string

// Cluster represents a cluster returned from a ListClusters or GetCluster call.
type Cluster interface {
	// ID returns the UUID of the cluster.
	ID() ClusterID
	// Name returns the textual name of the cluster.
	Name() string
}

func convertSDKCluster(sdkCluster *ovirtsdk4.Cluster, client Client) (Cluster, error) {
	id, ok := sdkCluster.Id()
	if !ok {
		return nil, newError(EFieldMissing, "failed to fetch ID for cluster")
	}

	name, ok := sdkCluster.Name()
	if !ok {
		return nil, newError(EFieldMissing, "failed to fetch name for cluster %s", id)
	}
	return &cluster{
		client: client,
		id:     ClusterID(id),
		name:   name,
	}, nil
}

type cluster struct {
	client Client

	id   ClusterID
	name string
}

func (c cluster) ID() ClusterID {
	return c.id
}

func (c cluster) Name() string {
	return c.name
}
