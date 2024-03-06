/*
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ocm

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	idputils "github.com/openshift-online/ocm-common/pkg/idp/utils"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	"github.com/openshift/rosa/pkg/helper"
)

const (
	HTPasswdIDPType = "HTPasswd"
	GithubIDPType   = "GitHub"
	GitlabIDPType   = "GitLab"
	GoogleIDPType   = "Google"
	LDAPIDPType     = "LDAP"
	OpenIDIDPType   = "OpenID"
)

func (c *Client) GetIdentityProviders(clusterID string) ([]*cmv1.IdentityProvider, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		IdentityProviders().
		List().Page(1).Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Items().Slice(), nil
}

func (c *Client) CreateIdentityProvider(clusterID string, idp *cmv1.IdentityProvider) (*cmv1.IdentityProvider, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		IdentityProviders().
		Add().Body(idp).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) GetHTPasswdUserList(clusterID, htpasswdIDPId string) (*cmv1.HTPasswdUserList, error) {
	listResponse, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).
		IdentityProviders().IdentityProvider(htpasswdIDPId).HtpasswdUsers().List().Send()
	if err != nil {
		if listResponse.Error().Status() == http.StatusNotFound {
			return nil, nil
		}
		return nil, handleErr(listResponse.Error(), err)
	}
	return listResponse.Items(), nil
}

func (c *Client) AddHTPasswdUser(username, password, clusterID, idpID string) error {
	hashedPwd, err := idputils.GenerateHTPasswdCompatibleHash(password)
	if err != nil {
		return fmt.Errorf("Failed to hash the password: %s", err)
	}
	htpasswdUser, _ := cmv1.NewHTPasswdUser().Username(username).HashedPassword(hashedPwd).Build()
	response, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).
		IdentityProviders().IdentityProvider(idpID).HtpasswdUsers().Add().Body(htpasswdUser).Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

func (c *Client) AddHTPasswdUsers(userList *cmv1.HTPasswdUserList, clusterID, idpID string) error {
	response, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).
		IdentityProviders().IdentityProvider(idpID).HtpasswdUsers().Import().Items(userList.Slice()).Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

func (c *Client) DeleteHTPasswdUser(username, clusterID string, htpasswdIDP *cmv1.IdentityProvider) error {
	var userID string

	listResponse, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).
		IdentityProviders().IdentityProvider(htpasswdIDP.ID()).HtpasswdUsers().List().Send()
	if err != nil {
		if listResponse.Error().Status() == http.StatusNotFound {
			return nil
		}
		return handleErr(listResponse.Error(), err)
	}
	listResponse.Items().Each(func(user *cmv1.HTPasswdUser) bool {
		if user.Username() == username {
			userID = user.ID()
		}
		return true
	})
	if userID == "" {
		return fmt.Errorf("HTPasswd user named '%s' on cluster '%s' does not exist", username, clusterID)
	}
	deleteResponse, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).
		IdentityProviders().IdentityProvider(htpasswdIDP.ID()).HtpasswdUsers().
		HtpasswdUser(userID).Delete().Send()
	if err != nil {
		return handleErr(deleteResponse.Error(), err)
	}
	return nil
}

func (c *Client) DeleteIdentityProvider(clusterID string, idpID string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		IdentityProviders().IdentityProvider(idpID).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

func IdentityProviderType(idp *cmv1.IdentityProvider) string {
	switch idp.Type() {
	case cmv1.IdentityProviderTypeGithub:
		return GithubIDPType
	case cmv1.IdentityProviderTypeGitlab:
		return GitlabIDPType
	case cmv1.IdentityProviderTypeGoogle:
		return GoogleIDPType
	case cmv1.IdentityProviderTypeHtpasswd:
		return HTPasswdIDPType
	case cmv1.IdentityProviderTypeLDAP:
		return LDAPIDPType
	case cmv1.IdentityProviderTypeOpenID:
		return OpenIDIDPType
	}

	return ""
}

func HasAuthURLSupport(idp *cmv1.IdentityProvider) bool {
	return !helper.Contains(getIDPListWithoutAuthURLSupport(), IdentityProviderType(idp))
}

// OAuthURLNeedsPort defines if an IDP needs a port for the callback URL
func OAuthURLNeedsPort(idpType cmv1.IdentityProviderType) bool {
	return idpType == cmv1.IdentityProviderTypeOpenID
}

func getIDPListWithoutAuthURLSupport() []string {
	return []string{HTPasswdIDPType, LDAPIDPType}
}

// BuildOAuthURL builds the correct OAuthURL depending on the cluster type
func BuildOAuthURL(cluster *cmv1.Cluster, idpType cmv1.IdentityProviderType) (string, error) {
	var oauthURL string
	if cluster.Hypershift().Enabled() {
		// https://api.example.com:443 -> https://oauth.example.com
		apiURL := cluster.API().URL()
		if OAuthURLNeedsPort(idpType) {
			// Some IDPs require the port, so just replace what is needed
			oauthURL = strings.Replace(apiURL, "api", "oauth", 1)
		} else {
			// Otherwise, remove the port and replace what is needed
			u, err := url.ParseRequestURI(apiURL)
			if err != nil {
				return oauthURL, err
			}
			host, _, err := net.SplitHostPort(u.Host)
			if err != nil {
				return oauthURL, err
			}
			u.Host = host
			oauthURL = strings.Replace(u.String(), "api", "oauth", 1)
		}
	} else {
		oauthURL = strings.Replace(cluster.Console().URL(), "console-openshift-console", "oauth-openshift", 1)
	}
	return oauthURL, nil
}

// GetOAuthURL builds the full OAuthURL depending on the cluster type and the idp name
func GetOAuthURL(cluster *cmv1.Cluster, idp *cmv1.IdentityProvider) (string, error) {
	if !HasAuthURLSupport(idp) {
		return "", nil
	}
	oauthURL, err := BuildOAuthURL(cluster, idp.Type())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/oauth2callback/%s", oauthURL, idp.Name()), nil
}
