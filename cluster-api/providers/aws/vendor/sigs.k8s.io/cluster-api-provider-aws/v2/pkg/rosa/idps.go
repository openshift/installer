package rosa

import (
	"fmt"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

const (
	clusterAdminUserGroup = "cluster-admins"
	clusterAdminIDPname   = "cluster-admin"
)

// CreateAdminUserIfNotExist creates a new admin user withe username/password in the cluster if username doesn't already exist.
// the user is granted admin privileges by being added to a special IDP called `cluster-admin` which will be created if it doesn't already exist.
func CreateAdminUserIfNotExist(client OCMClient, clusterID, username, password string) error {
	existingClusterAdminIDP, userList, err := findExistingClusterAdminIDP(client, clusterID)
	if err != nil {
		return fmt.Errorf("failed to find existing cluster admin IDP: %w", err)
	}
	if existingClusterAdminIDP != nil {
		if hasUser(username, userList) {
			// user already exist in the cluster
			return nil
		}
	}

	// Add admin user to the cluster-admins group:
	user, err := CreateUserIfNotExist(client, clusterID, clusterAdminUserGroup, username)
	if err != nil {
		return fmt.Errorf("failed to add user '%s' to cluster '%s': %s",
			username, clusterID, err)
	}

	if existingClusterAdminIDP != nil {
		// add htpasswd user to existing idp
		err := client.AddHTPasswdUser(username, password, clusterID, existingClusterAdminIDP.ID())
		if err != nil {
			return fmt.Errorf("failed to add htpassawoed user cluster-admin to existing idp: %s", existingClusterAdminIDP.ID())
		}

		return nil
	}

	// No ClusterAdmin IDP exists, create an Htpasswd IDP
	htpasswdIDP := cmv1.NewHTPasswdIdentityProvider().Users(cmv1.NewHTPasswdUserList().Items(
		cmv1.NewHTPasswdUser().Username(username).Password(password),
	))
	clusterAdminIDP, err := cmv1.NewIdentityProvider().
		Type(cmv1.IdentityProviderTypeHtpasswd).
		Name(clusterAdminIDPname).
		Htpasswd(htpasswdIDP).
		Build()
	if err != nil {
		return fmt.Errorf(
			"failed to create '%s' identity provider for cluster '%s'",
			clusterAdminIDPname,
			clusterID,
		)
	}

	// Add HTPasswd IDP to cluster
	_, err = client.CreateIdentityProvider(clusterID, clusterAdminIDP)
	if err != nil {
		// since we could not add the HTPasswd IDP to the cluster, roll back and remove the cluster admin
		if err := client.DeleteUser(clusterID, clusterAdminUserGroup, user.ID()); err != nil {
			return fmt.Errorf("failed to revert the admin user for cluster '%s': %w",
				clusterID, err)
		}
		return fmt.Errorf("failed to create identity cluster-admin idp: %w", err)
	}

	return nil
}

// CreateUserIfNotExist creates a new user with `username` and adds it to the group if it doesn't already exist.
func CreateUserIfNotExist(client OCMClient, clusterID string, group, username string) (*cmv1.User, error) {
	user, err := client.GetUser(clusterID, group, username)
	if user != nil || err != nil {
		return user, err
	}

	userCfg, err := cmv1.NewUser().ID(username).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create user '%s' for cluster '%s': %w", username, clusterID, err)
	}
	return client.CreateUser(clusterID, group, userCfg)
}

func findExistingClusterAdminIDP(client OCMClient, clusterID string) (
	htpasswdIDP *cmv1.IdentityProvider, userList *cmv1.HTPasswdUserList, reterr error) {
	idps, err := client.GetIdentityProviders(clusterID)
	if err != nil {
		reterr = fmt.Errorf("failed to get identity providers for cluster '%s': %v", clusterID, err)
		return
	}

	for _, idp := range idps {
		if idp.Name() != clusterAdminIDPname {
			continue
		}

		itemUserList, err := client.GetHTPasswdUserList(clusterID, idp.ID())
		if err != nil {
			reterr = fmt.Errorf("failed to get user list of the HTPasswd IDP of '%s: %s': %v", idp.Name(), clusterID, err)
			return
		}

		htpasswdIDP = idp
		userList = itemUserList
		return
	}

	return
}

func hasUser(username string, userList *cmv1.HTPasswdUserList) bool {
	hasUser := false
	userList.Each(func(user *cmv1.HTPasswdUser) bool {
		if user.Username() == username {
			hasUser = true
			return false
		}
		return true
	})
	return hasUser
}
