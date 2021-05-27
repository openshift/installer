package iamv1

import (
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type serviceRoleQueryResponse struct {
	ServiceSpecificRoles []models.PolicyRole `json:"supportedRoles"`
	PlatformExtensions   struct {
		Roles []models.PolicyRole `json:"supportedRoles"`
	} `json:"platformExtensions"`
}

type ServiceRoleRepository interface {
	// List all roles of a given service, including those supported system defined roles
	ListServiceRoles(serviceName string) ([]models.PolicyRole, error)
	// List all system defined roles
	ListSystemDefinedRoles() ([]models.PolicyRole, error)
	// List servie specific roles
	ListServiceSpecificRoles(serviceName string) ([]models.PolicyRole, error)
	// List authorization roles
	ListAuthorizationRoles(sourceServiceName string, targetServiceName string) ([]models.PolicyRole, error)
}

type serviceRoleRepository struct {
	client *client.Client
}

func NewServiceRoleRepository(c *client.Client) ServiceRoleRepository {
	return &serviceRoleRepository{
		client: c,
	}
}

func (r *serviceRoleRepository) ListServiceRoles(serviceName string) ([]models.PolicyRole, error) {
	response := struct {
		ServiceSpecificRoles []models.PolicyRole `json:"supportedRoles"`
		PlatformExtensions   struct {
			Roles []models.PolicyRole `json:"supportedRoles"`
		} `json:"platformExtensions"`
	}{}

	_, err := r.client.Get("/acms/v1/roles?serviceName="+url.QueryEscape(serviceName), &response)
	if err != nil {
		return []models.PolicyRole{}, err
	}

	roles := append(response.ServiceSpecificRoles, response.PlatformExtensions.Roles...)

	return roles, nil
}

func (r *serviceRoleRepository) ListSystemDefinedRoles() ([]models.PolicyRole, error) {
	response := struct {
		Roles []models.Role `json:"systemDefinedRoles"`
	}{}

	_, err := r.client.Get("/acms/v1/roles", &response)
	if err != nil {
		return []models.PolicyRole{}, err
	}

	// system defined roles uses `crn` instead of `id`, need to conversion
	// TODO: remove this if IAM PAP unify the data model
	roles := []models.PolicyRole{}
	for _, role := range response.Roles {
		roles = append(roles, role.ToPolicyRole())
	}
	return roles, nil
}
func (r *serviceRoleRepository) ListServiceSpecificRoles(serviceName string) ([]models.PolicyRole, error) {
	var response serviceRoleQueryResponse
	var err error
	if response, err = r.queryServiceRoles(serviceName); err != nil {
		return []models.PolicyRole{}, err
	}
	return response.ServiceSpecificRoles, nil

}

func (r *serviceRoleRepository) queryServiceRoles(name string) (serviceRoleQueryResponse, error) {
	response := serviceRoleQueryResponse{}

	_, err := r.client.Get("/acms/v1/roles?serviceName="+url.QueryEscape(name), &response)
	if err != nil {
		return serviceRoleQueryResponse{}, err
	}

	return response, nil
}

func (r *serviceRoleRepository) ListAuthorizationRoles(sourceServiceName string, targetServiceName string) ([]models.PolicyRole, error) {
	req := rest.GetRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/acms/v1/roles"))
	req.Query("sourceServiceName", sourceServiceName)
	req.Query("serviceName", targetServiceName)
	req.Query("policyType", "authorization")

	var response serviceRoleQueryResponse
	_, err := r.client.SendRequest(req, &response)
	if err != nil {
		return []models.PolicyRole{}, err
	}

	roles := append(response.ServiceSpecificRoles, response.PlatformExtensions.Roles...)
	return roles, nil
}
