package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v2/tenants"
	tokens2 "github.com/gophercloud/gophercloud/openstack/identity/v2/tokens"
	tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

func flattenIdentityAuthScopeV3Roles(roles []tokens3.Role) []map[string]string {
	allRoles := make([]map[string]string, len(roles))

	for i, r := range roles {
		allRoles[i] = map[string]string{
			"role_name": r.Name,
			"role_id":   r.ID,
		}
	}

	return allRoles
}

func flattenIdentityAuthScopeV3ServiceCatalog(catalog *tokens3.ServiceCatalog) []map[string]interface{} {
	ret := make([]map[string]interface{}, len(catalog.Entries))

	for iEntry, entry := range catalog.Entries {
		endpoints := make([]map[string]string, len(entry.Endpoints))
		for iEndpoint, endpoint := range entry.Endpoints {
			endpoints[iEndpoint] = map[string]string{
				"id":        endpoint.ID,
				"region":    endpoint.Region,
				"region_id": endpoint.RegionID,
				"interface": endpoint.Interface,
				"url":       endpoint.URL,
			}
		}
		ret[iEntry] = map[string]interface{}{
			"id":        entry.ID,
			"name":      entry.Name,
			"type":      entry.Type,
			"endpoints": endpoints,
		}
	}

	return ret
}

type authScopeTokenDetails struct {
	user    *tokens3.User
	domain  *tokens3.Domain
	project *tokens3.Project
	catalog *tokens3.ServiceCatalog
	roles   []tokens3.Role
}

func getTokenDetails(sc *gophercloud.ServiceClient) (authScopeTokenDetails, error) {
	var (
		details authScopeTokenDetails
		err     error
	)

	r := sc.ProviderClient.GetAuthResult()
	switch result := r.(type) {
	case tokens3.CreateResult:
		details.user, err = result.ExtractUser()
		if err != nil {
			return details, err
		}
		details.domain, err = result.ExtractDomain()
		if err != nil {
			return details, err
		}
		details.project, err = result.ExtractProject()
		if err != nil {
			return details, err
		}
		details.roles, err = result.ExtractRoles()
		if err != nil {
			return details, err
		}
		details.catalog, err = result.ExtractServiceCatalog()
		if err != nil {
			return details, err
		}
	case tokens3.GetResult:
		details.user, err = result.ExtractUser()
		if err != nil {
			return details, err
		}
		details.domain, err = result.ExtractDomain()
		if err != nil {
			return details, err
		}
		details.project, err = result.ExtractProject()
		if err != nil {
			return details, err
		}
		details.roles, err = result.ExtractRoles()
		if err != nil {
			return details, err
		}
		details.catalog, err = result.ExtractServiceCatalog()
		if err != nil {
			return details, err
		}
	default:
		res := tokens3.Get(sc, sc.ProviderClient.TokenID)
		if res.Err != nil {
			return details, res.Err
		}
		details.user, err = res.ExtractUser()
		if err != nil {
			return details, err
		}
		details.domain, err = res.ExtractDomain()
		if err != nil {
			return details, err
		}
		details.project, err = res.ExtractProject()
		if err != nil {
			return details, err
		}
		details.roles, err = res.ExtractRoles()
		if err != nil {
			return details, err
		}
		// AuthResult has no method ExtractServiceCatalog
	}

	return details, nil
}

type authScopeTokenInfo struct {
	userID    string
	projectID string
	tokenID   string
}

func getTokenInfo(sc *gophercloud.ServiceClient) (authScopeTokenInfo, error) {
	r := sc.ProviderClient.GetAuthResult()
	switch r := r.(type) {
	case tokens2.CreateResult:
		return getTokenInfoV2(r)
	case tokens3.CreateResult, tokens3.GetResult:
		return getTokenInfoV3(r)
	default:
		token := tokens3.Get(sc, sc.ProviderClient.TokenID)
		if token.Err != nil {
			return authScopeTokenInfo{}, token.Err
		}
		return getTokenInfoV3(token)
	}
}

func getTokenInfoV3(t interface{}) (authScopeTokenInfo, error) {
	var info authScopeTokenInfo
	switch r := t.(type) {
	case tokens3.CreateResult:
		user, err := r.ExtractUser()
		if err != nil {
			return info, err
		}
		project, err := r.ExtractProject()
		if err != nil {
			return info, err
		}
		info.userID = user.ID
		if project != nil {
			info.projectID = project.ID
		}
		return info, nil
	case tokens3.GetResult:
		user, err := r.ExtractUser()
		if err != nil {
			return info, err
		}
		project, err := r.ExtractProject()
		if err != nil {
			return info, err
		}
		info.userID = user.ID
		if project != nil {
			info.projectID = project.ID
		}
		return info, nil
	default:
		return info, fmt.Errorf("got unexpected AuthResult type %t", r)
	}
}

func getTokenInfoV2(t tokens2.CreateResult) (authScopeTokenInfo, error) {
	var info authScopeTokenInfo
	var s struct {
		Access struct {
			Token struct {
				Expires string         `json:"expires"`
				ID      string         `json:"id"`
				Tenant  tenants.Tenant `json:"tenant"`
			} `json:"token"`
			User tokens2.User `json:"user"`
		} `json:"access"`
	}

	err := t.ExtractInto(&s)
	if err != nil {
		return info, err
	}
	info.userID = s.Access.User.ID
	info.tokenID = s.Access.Token.ID
	return info, nil
}
