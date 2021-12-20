package iampapv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type CreateRoleRequest struct {
	Name        string   `json:"name"`
	ServiceName string   `json:"service_name"`
	AccountID   string   `json:"account_id"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	Actions     []string `json:"actions,omitempty"`
}
type UpdateRoleRequest struct {
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	Actions     []string `json:"actions,omitempty"`
}

type Role struct {
	CreateRoleRequest
	ID               string `json:"id"`
	Crn              string `json:"crn"`
	CreatedAt        string `json:"created_at"`
	CreatedByID      string `json:"created_by_id"`
	LastModifiedAt   string `json:"last_modified_at"`
	LastModifiedByID string `json:"last_modified_by_id"`
}

type ListResponse struct {
	CustomRoles  []Role `json:"custom_roles"`
	ServiceRoles []Role `json:"service_roles"`
	SystemRoles  []Role `json:"system_roles"`
}

type RoleRepository interface {
	Get(roleID string) (Role, string, error)
	Create(request CreateRoleRequest) (Role, error)
	Update(request UpdateRoleRequest, roleID, etag string) (Role, error)
	Delete(roleID string) error
	ListCustomRoles(accountID, serviceName string) ([]Role, error)
	ListSystemDefinedRoles() ([]Role, error)
	ListServiceRoles(serviceName string) ([]Role, error)
	ListAll(query RoleQuery) ([]Role, error)
}

type roleRepository struct {
	client *client.Client
}

func NewRoleRepository(c *client.Client) RoleRepository {
	return &roleRepository{
		client: c,
	}
}

type RoleQueryFormatParameter string

type RoleQuery struct {
	AccountID   string
	ServiceName string
	Format      RoleQueryFormatParameter
}

// SetQuery will set query parameter to the passed-in request
func (q RoleQuery) SetQuery(req *rest.Request) {
	if q.AccountID != "" {
		req.Query("account_id", q.AccountID)
	}
	if q.ServiceName != "" {
		req.Query("service_name", q.ServiceName)
	}
	if q.Format != "" {
		req.Query("format", string(q.Format))
	}
}

func (r *roleRepository) Create(request CreateRoleRequest) (Role, error) {
	res := Role{}
	_, err := r.client.Post(fmt.Sprintf("/v2/roles"), &request, &res)
	if err != nil {
		return Role{}, err
	}
	return res, nil
}

func (r *roleRepository) Get(roleID string) (Role, string, error) {
	res := Role{}
	response, err := r.client.Get(fmt.Sprintf("/v2/roles/%s", roleID), &res)
	if err != nil {
		return Role{}, "", err
	}
	return res, response.Header.Get("Etag"), nil
}

func (r *roleRepository) Update(request UpdateRoleRequest, roleID, etag string) (Role, error) {
	res := Role{}
	header := make(map[string]string)

	header["IF-Match"] = etag
	_, err := r.client.Put(fmt.Sprintf("/v2/roles/%s", roleID), &request, &res, header)
	if err != nil {
		return Role{}, err
	}
	return res, nil
}

//Delete Function
func (r *roleRepository) Delete(roleID string) error {
	_, err := r.client.Delete(fmt.Sprintf("/v2/roles/%s", roleID))
	return err
}

func (r *roleRepository) ListCustomRoles(accountID, serviceName string) ([]Role, error) {
	res := ListResponse{}
	var requestpath string

	requestpath = fmt.Sprintf("/v2/roles?account_id=%s", accountID)

	_, err := r.client.Get(requestpath, &res)
	if err != nil {
		return []Role{}, err
	}
	if serviceName == "" {
		return res.CustomRoles, nil
	} else {
		var matchingRoles []Role
		for _, role := range res.CustomRoles {
			if role.ServiceName == serviceName {
				matchingRoles = append(matchingRoles, role)
			}
		}
		return matchingRoles, nil
	}

}

func (r *roleRepository) ListSystemDefinedRoles() ([]Role, error) {
	res := ListResponse{}
	var requestpath string
	requestpath = fmt.Sprintf("/v2/roles")
	_, err := r.client.Get(requestpath, &res)
	if err != nil {
		return []Role{}, err
	}
	return res.SystemRoles, nil
}

func (r *roleRepository) ListServiceRoles(serviceName string) ([]Role, error) {
	res := ListResponse{}
	var requestpath string
	requestpath = fmt.Sprintf("/v2/roles?service_name=%s", serviceName)
	_, err := r.client.Get(requestpath, &res)
	if err != nil {
		return []Role{}, err
	}
	return res.ServiceRoles, nil
}

func (r *roleRepository) ListAll(query RoleQuery) ([]Role, error) {
	response, err := r.query(query)
	if err != nil {
		return []Role{}, err
	}
	return append(response.CustomRoles, append(response.ServiceRoles, response.SystemRoles...)...), nil
}

func (r *roleRepository) query(query RoleQuery) (ListResponse, error) {
	req := rest.GetRequest(*r.client.Config.Endpoint + "/v2/roles")
	query.SetQuery(req)

	var response ListResponse
	_, err := r.client.SendRequest(req, &response)

	return response, err
}
