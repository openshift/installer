package controller

import (
	"net/url"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

const (
	_Role_Crn       = "role_crn"
	_Service_ID_Crn = "serviceid_crn"
)

type CreateServiceKeyRequest struct {
	Name       string                 `json:"name"`
	SourceCRN  crn.CRN                `json:"source_crn"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

//ErrCodeResourceServiceInstanceDoesnotExist ...
const ErrCodeResourceServiceKeyDoesnotExist = "ResourceServiceInstanceDoesnotExist"

//ResourceServiceInstanceQuery ...
type ResourceServiceKeyRepository interface {
	GetKey(keyID string) (models.ServiceKey, error)
	GetKeys(keyName string) ([]models.ServiceKey, error)
	CreateKey(CreateServiceKeyRequest) (models.ServiceKey, error)
	DeleteKey(keyID string) error
}

type resourceServiceKey struct {
	client *client.Client
}

func newResourceServiceKeyAPI(c *client.Client) ResourceServiceKeyRepository {
	return &resourceServiceKey{
		client: c,
	}
}

func (r *resourceServiceKey) GetKeys(keyName string) ([]models.ServiceKey, error) {
	var keys []models.ServiceKey
	_, err := r.client.GetPaginated(
		"/v1/resource_keys",
		NewRCPaginatedResources(models.ServiceKey{}),
		func(resource interface{}) bool {
			if key, ok := resource.(models.ServiceKey); ok {
				keys = append(keys, key)
				return true
			}
			return false
		},
	)
	if err != nil {
		return []models.ServiceKey{}, err
	}

	if keyName != "" {
		keys = filterKeysByName(keys, keyName)
	}

	return keys, nil
}

func filterKeysByName(keys []models.ServiceKey, name string) []models.ServiceKey {
	ret := []models.ServiceKey{}
	for _, k := range keys {
		if strings.EqualFold(k.Name, name) {
			ret = append(ret, k)
		}
	}
	return ret
}

func (r *resourceServiceKey) GetKey(keyID string) (models.ServiceKey, error) {
	resp := models.ServiceKey{}
	request := rest.GetRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/resource_keys/"+url.PathEscape(keyID)))
	_, err := r.client.SendRequest(request, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (r *resourceServiceKey) CreateKey(request CreateServiceKeyRequest) (models.ServiceKey, error) {
	resp := models.ServiceKey{}
	_, err := r.client.Post("/v1/resource_keys", request, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (r *resourceServiceKey) DeleteKey(keyID string) error {
	request := rest.DeleteRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/resource_keys/"+url.PathEscape(keyID))).Query("id", url.PathEscape(keyID))
	_, err := r.client.SendRequest(request, nil)
	return err
}
