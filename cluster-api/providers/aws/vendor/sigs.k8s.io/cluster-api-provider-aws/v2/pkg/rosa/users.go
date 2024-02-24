package rosa

import (
	"fmt"
	"net/http"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// CreateUserIfNotExist creates a new user with `username` and adds it to the group if it doesn't already exist.
func (c *RosaClient) CreateUserIfNotExist(clusterID string, group, username string) (*cmv1.User, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		Groups().Group(group).
		Users().User(username).
		Get().
		Send()
	if err == nil {
		return response.Body(), nil
	} else if response.Error().Status() != http.StatusNotFound {
		return nil, handleErr(response.Error(), err)
	}

	user, err := cmv1.NewUser().ID(username).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create user '%s' for cluster '%s'", username, clusterID)
	}

	return c.CreateUser(clusterID, group, user)
}

// CreateUser adds a new user to the group.
func (c *RosaClient) CreateUser(clusterID string, group string, user *cmv1.User) (*cmv1.User, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		Groups().Group(group).
		Users().
		Add().Body(user).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

// DeleteUser deletes the user from the cluster.
func (c *RosaClient) DeleteUser(clusterID string, group string, username string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		Groups().Group(group).
		Users().User(username).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}
