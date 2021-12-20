package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type ServiceAliasQueryFilter struct {
	AccountID         string
	ServiceInstanceID string
	Name              string // TODO: RC API currently not support name filtering
}

type CreateServiceAliasParams struct {
	Name              string                 `json:"name"`
	ServiceInstanceID string                 `json:"resource_instance_id"`
	ScopeCRN          crn.CRN                `json:"scope_crn"`
	Tags              []string               `json:"tags,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
}

type UpdateServiceAliasParams struct {
	Name       string                 `json:"name,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

//ErrCodeResourceServiceAliasDoesnotExist ...
const ErrCodeResourceServiceAliasDoesnotExist = "ResourceServiceAliasDoesnotExist"

type ResourceServiceAliasRepository interface {
	Alias(aliasID string) (models.ServiceAlias, error)
	Aliases(*ServiceAliasQueryFilter) ([]models.ServiceAlias, error)
	AliasesWithCallback(*ServiceAliasQueryFilter, func(models.ServiceAlias) bool) error

	InstanceAliases(serviceInstanceID string) ([]models.ServiceAlias, error)
	InstanceAliasByName(serviceInstanceID string, name string) ([]models.ServiceAlias, error)

	CreateAlias(params CreateServiceAliasParams) (models.ServiceAlias, error)
	UpdateAlias(aliasID string, params UpdateServiceAliasParams) (models.ServiceAlias, error)
	DeleteAlias(aliasID string) error
}

type serviceAliasRepository struct {
	client *client.Client
}

func newResourceServiceAliasRepository(c *client.Client) ResourceServiceAliasRepository {
	return &serviceAliasRepository{
		client: c,
	}
}

func (r *serviceAliasRepository) InstanceAliases(serviceInstanceID string) ([]models.ServiceAlias, error) {
	return r.Aliases(&ServiceAliasQueryFilter{ServiceInstanceID: serviceInstanceID})
}

func (r *serviceAliasRepository) InstanceAliasByName(serviceInstanceID string, name string) ([]models.ServiceAlias, error) {
	return r.Aliases(&ServiceAliasQueryFilter{ServiceInstanceID: serviceInstanceID, Name: name})
}

func (r *serviceAliasRepository) Alias(aliasID string) (models.ServiceAlias, error) {
	var alias models.ServiceAlias
	resp, err := r.client.Get("/v1/resource_aliases/"+url.PathEscape(aliasID), &alias)
	if resp.StatusCode == http.StatusNotFound {
		return alias, bmxerror.New(ErrCodeResourceServiceAliasDoesnotExist,
			fmt.Sprintf("Given service alias : %q doesn't exist", aliasID))
	}
	return alias, err
}
func (r *serviceAliasRepository) Aliases(filter *ServiceAliasQueryFilter) ([]models.ServiceAlias, error) {
	var aliases []models.ServiceAlias
	err := r.AliasesWithCallback(filter, func(a models.ServiceAlias) bool {
		aliases = append(aliases, a)
		return true
	})
	return aliases, err
}

func (r *serviceAliasRepository) AliasesWithCallback(filter *ServiceAliasQueryFilter, cb func(models.ServiceAlias) bool) error {
	listRequest := rest.GetRequest("/v1/resource_aliases")
	if filter != nil {
		if filter.AccountID != "" {
			listRequest.Query("account_id", filter.AccountID)
		}
		if filter.ServiceInstanceID != "" {
			listRequest.Query("resource_instance_id", url.PathEscape(filter.ServiceInstanceID))
		}
	}

	req, err := listRequest.Build()
	if err != nil {
		return err
	}

	// TODO: GetPaginated's first argument should be a request instead if url path
	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewRCPaginatedResources(models.ServiceAlias{}),
		func(resource interface{}) bool {
			// if alias, ok := resource.(models.ServiceAlias); ok {
			// 	return cb(alias)
			// }
			// TODO: once RC API support name filtering, remove name check in cb
			if alias, ok := resource.(models.ServiceAlias); ok {
				if filter.Name == "" || strings.EqualFold(filter.Name, alias.Name) {
					return cb(alias)
				}
				return true
			}
			return false
		})

	return err
}

func (r *serviceAliasRepository) CreateAlias(params CreateServiceAliasParams) (models.ServiceAlias, error) {
	alias := models.ServiceAlias{}
	_, err := r.client.Post("/v1/resource_aliases", params, &alias)
	return alias, err
}

func (r *serviceAliasRepository) UpdateAlias(aliasID string, params UpdateServiceAliasParams) (models.ServiceAlias, error) {
	alias := models.ServiceAlias{}
	resp, err := r.client.Patch("/v1/resource_aliases/"+url.PathEscape(aliasID), params, &alias)
	if resp.StatusCode == http.StatusNotFound {
		return alias, bmxerror.New(ErrCodeResourceServiceAliasDoesnotExist,
			fmt.Sprintf("Given service alias : %q doesn't exist", aliasID))
	}
	return alias, err
}

func (r *serviceAliasRepository) DeleteAlias(aliasID string) error {
	resp, err := r.client.Delete("/v1/resource_aliases/" + url.PathEscape(aliasID))
	if resp.StatusCode == http.StatusNotFound {
		return bmxerror.New(ErrCodeResourceServiceAliasDoesnotExist,
			fmt.Sprintf("Given service alias : %q doesn't exist", aliasID))
	}
	return err
}
