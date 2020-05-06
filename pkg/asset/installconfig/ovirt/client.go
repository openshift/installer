package ovirt

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
)

// getConnection is a convenience method to get a connection to ovirt api
// form a Config Object.
func getConnection(ovirtConfig Config) (*ovirtsdk.Connection, error) {
	con, err := ovirtsdk.NewConnectionBuilder().
		URL(ovirtConfig.URL).
		Username(ovirtConfig.Username).
		Password(ovirtConfig.Password).
		CAFile(ovirtConfig.CAFile).
		Insecure(ovirtConfig.Insecure).
		Build()
	if err != nil {
		return nil, err
	}
	return con, nil
}

// NewConnection returns a new client connection to oVirt's API endpoint.
// It is the responsibility of the caller to close the connection.
func NewConnection() (*ovirtsdk.Connection, error) {
	ovirtConfig, err := NewConfig()
	if err != nil {
		return nil, errors.Wrap(err, "getting ovirt configuration")
	}
	con, err := getConnection(ovirtConfig)
	if err != nil {
		return nil, errors.Wrap(err, "establishing ovirt connection")
	}
	return con, nil
}

// FetchVNICProfileByClusterNetwork returns a list of profiles for the given cluster and network name.
func FetchVNICProfileByClusterNetwork(con *ovirtsdk.Connection, clusterID string, networkName string) ([]*ovirtsdk.VnicProfile, error) {
	clusterResponse, err := con.SystemService().ClustersService().ClusterService(clusterID).Get().Follow("networks").Send()
	if err != nil {
		return nil, err
	}

	cluster, ok := clusterResponse.Cluster()
	if !ok {
		return nil, fmt.Errorf("failed to find cluster with id %s", clusterID)
	}

	networks, ok := cluster.Networks()
	if !ok {
		return nil, fmt.Errorf("no cluster networks for cluster %s [%s]", cluster.MustName(), clusterID)
	}

	for _, n := range networks.Slice() {
		if n.MustName() != networkName {
			continue
		}

		profilesGet, err := con.SystemService().NetworksService().NetworkService(n.MustId()).VnicProfilesService().List().Send()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch vNic profiles")
		}

		return profilesGet.MustProfiles().Slice(), nil
	}
	return nil, fmt.Errorf("there are no vNic profiles for the given cluster ID %s and network name %s", clusterID, networkName)
}
