package management

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/models"
)

type quotaDefinitionQueryResult struct {
	QuotaDefinitions []models.QuotaDefinition `json:"resources"`
}

//ErrCodeResourceQuotaDoesnotExist ...
const ErrCodeResourceQuotaDoesnotExist = "ResourceQuotaDoesnotExist"

type ResourceQuotaRepository interface {
	// List all quota definitions
	List() ([]models.QuotaDefinition, error)
	// Query quota definitions having specific name
	FindByName(name string) ([]models.QuotaDefinition, error)
	// Get quota definition by ID
	Get(id string) (*models.QuotaDefinition, error)
}

type resourceQuota struct {
	client *client.Client
}

func newResourceQuotaAPI(c *client.Client) ResourceQuotaRepository {
	return &resourceQuota{
		client: c,
	}
}

func (r *resourceQuota) List() ([]models.QuotaDefinition, error) {
	resp := quotaDefinitionQueryResult{}
	// TODO: change to use pagination if it's available on backend
	_, err := r.client.Get("/v1/quota_definitions", &resp)
	if err != nil {
		return []models.QuotaDefinition{}, err
	}
	return resp.QuotaDefinitions, nil
}

func (r *resourceQuota) FindByName(name string) ([]models.QuotaDefinition, error) {
	allQuotas, err := r.List()
	if err != nil {
		return []models.QuotaDefinition{}, err
	}

	quotas := []models.QuotaDefinition{}
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

func (r *resourceQuota) Get(id string) (*models.QuotaDefinition, error) {
	quota := models.QuotaDefinition{}
	_, err := r.client.Get("/v1/quota_definitions/"+id, &quota)
	if err != nil {
		return nil, err
	}
	return &quota, nil
}
