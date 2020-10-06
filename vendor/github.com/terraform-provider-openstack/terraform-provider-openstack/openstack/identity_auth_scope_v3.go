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

func GetTokenDetails(sc *gophercloud.ServiceClient) (*tokens3.User, *tokens3.Domain, *tokens3.Project, []tokens3.Role, error) {
	var (
		user    *tokens3.User
		domain  *tokens3.Domain
		project *tokens3.Project
		roles   []tokens3.Role
		err     error
	)

	r := sc.ProviderClient.GetAuthResult()
	switch result := r.(type) {
	case tokens3.CreateResult:
		user, err = result.ExtractUser()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		domain, err = result.ExtractDomain()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		project, err = result.ExtractProject()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		roles, err = result.ExtractRoles()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	case tokens3.GetResult:
		user, err = result.ExtractUser()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		domain, err = result.ExtractDomain()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		project, err = result.ExtractProject()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		roles, err = result.ExtractRoles()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	default:
		res := tokens3.Get(sc, sc.ProviderClient.TokenID)
		if res.Err != nil {
			return nil, nil, nil, nil, res.Err
		}
		user, err = res.ExtractUser()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		domain, err = res.ExtractDomain()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		project, err = res.ExtractProject()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		roles, err = res.ExtractRoles()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	return user, domain, project, roles, nil
}

func GetTokenInfo(sc *gophercloud.ServiceClient) (string, string, error) {
	r := sc.ProviderClient.GetAuthResult()
	switch r := r.(type) {
	case tokens2.CreateResult:
		return GetTokenInfoV2(r)
	case tokens3.CreateResult, tokens3.GetResult:
		return GetTokenInfoV3(r)
	default:
		token := tokens3.Get(sc, sc.ProviderClient.TokenID)
		if token.Err != nil {
			return "", "", token.Err
		}
		return GetTokenInfoV3(token)
	}
}

func GetTokenInfoV3(t interface{}) (string, string, error) {
	switch r := t.(type) {
	case tokens3.CreateResult:
		user, err := r.ExtractUser()
		if err != nil {
			return "", "", err
		}
		project, err := r.ExtractProject()
		if err != nil {
			return "", "", err
		}
		return user.ID, project.ID, nil
	case tokens3.GetResult:
		user, err := r.ExtractUser()
		if err != nil {
			return "", "", err
		}
		project, err := r.ExtractProject()
		if err != nil {
			return "", "", err
		}
		return user.ID, project.ID, nil
	default:
		return "", "", fmt.Errorf("got unexpected AuthResult type %t", r)
	}
}

func GetTokenInfoV2(t tokens2.CreateResult) (string, string, error) {
	var tokeninfo struct {
		Access struct {
			Token struct {
				Expires string         `json:"expires"`
				ID      string         `json:"id"`
				Tenant  tenants.Tenant `json:"tenant"`
			} `json:"token"`
			User tokens2.User `json:"user"`
		} `json:"access"`
	}

	err := t.ExtractInto(&tokeninfo)
	if err != nil {
		return "", "", err
	}
	return tokeninfo.Access.User.ID, tokeninfo.Access.Token.ID, nil
}
