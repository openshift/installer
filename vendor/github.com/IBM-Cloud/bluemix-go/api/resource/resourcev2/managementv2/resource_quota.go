package managementv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
)

type quotaDefinitionQueryResult struct {
	QuotaDefinitions []QuotaDefinition `json:"resources"`
}

type QuotaDefinition struct {
	ID                        string          `json:"_id,omitempty"`
	Revision                  string          `json:"_rev,omitempty"`
	Name                      string          `json:"name,omitmempty"`
	Type                      string          `json:"type,omitempty"`
	ServiceInstanceCountLimit int             `json:"number_of_service_instances,omitempty"`
	AppCountLimit             int             `json:"number_of_apps,omitempty"`
	AppInstanceCountLimit     int             `json:"instances_per_app,omitempty"`
	AppInstanceMemoryLimit    string          `json:"instance_memory,omitempty"`
	TotalAppMemoryLimit       string          `json:"total_app_memory,omitempty"`
	VSICountLimit             int             `json:"vsi_limit,omitempty"`
	ResourceQuotas            []ResourceQuota `json:"resource_quotas,omitempty"`
	CreatedAt                 string          `json:"created_at,omitempty"`
	UpdatedAt                 string          `json:"updated_at,omitempty"`
}

type ResourceQuota struct {
	ID         string `json:"_id,omitempty"`
	ResourceID string `json:"resource_id,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

//ErrCodeResourceQuotaDoesnotExist ...
const ErrCodeResourceQuotaDoesnotExist = "ResourceQuotaDoesnotExist"

type ResourceQuotaRepository interface {
	// List all quota definitions
	List() ([]QuotaDefinition, error)
	// Query quota definitions having specific name
	FindByName(name string) ([]QuotaDefinition, error)
	// Get quota definition by ID
	Get(id string) (*QuotaDefinition, error)
}

type resourceQuota struct {
	client *client.Client
}

func newResourceQuotaAPI(c *client.Client) ResourceQuotaRepository {
	return &resourceQuota{
		client: c,
	}
}

func (r *resourceQuota) List() ([]QuotaDefinition, error) {
	resp := quotaDefinitionQueryResult{}
	// TODO: change to use pagination if it's available on backend
	_, err := r.client.Get("/v2/quota_definitions", &resp)
	if err != nil {
		return []QuotaDefinition{}, err
	}
	return resp.QuotaDefinitions, nil
}

func (r *resourceQuota) FindByName(name string) ([]QuotaDefinition, error) {
	allQuotas, err := r.List()
	if err != nil {
		return []QuotaDefinition{}, err
	}

	quotas := []QuotaDefinition{}
	for _, quota := range allQuotas {
		if quota.Name == name {
			quotas = append(quotas, quota)
		}
	}

	if len(quotas) == 0 {
		return quotas, bmxerror.New(ErrCodeResourceQuotaDoesnotExist,
			fmt.Sprintf("Given quota : %q doesn't exist", name))
	}

	return quotas, nil
}

func (r *resourceQuota) Get(id string) (*QuotaDefinition, error) {
	quota := QuotaDefinition{}
	_, err := r.client.Get("/v2/quota_definitions/"+id, &quota)
	if err != nil {
		return nil, err
	}
	return &quota, nil
}
