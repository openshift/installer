package mccpv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//ErrCodeServiceDoesnotExist ...
const ErrCodeServiceDoesnotExist = "ServiceDoesnotExist"

//ServiceOffering model
type ServiceOffering struct {
	GUID              string
	Label             string   `json:"label"`
	Provider          string   `json:"provider"`
	Description       string   `json:"description"`
	LongDescription   string   `json:"long_description"`
	Version           string   `json:"version"`
	URL               string   `json:"url"`
	InfoURL           string   `json:"info_url"`
	DocumentURL       string   `json:"documentation_url"`
	Timeout           string   `json:"timeout"`
	UniqueID          string   `json:"unique_id"`
	ServiceBrokerGUID string   `json:"service_broker_guid"`
	ServicePlansURL   string   `json:"service_plans_url"`
	Tags              []string `json:"tags"`
	Requires          []string `json:"requires"`
	IsActive          bool     `json:"active"`
	IsBindable        bool     `json:"bindable"`
	IsPlanUpdateable  bool     `json:"plan_updateable"`
}

//ServiceOfferingResource ...
type ServiceOfferingResource struct {
	Resource
	Entity ServiceOfferingEntity
}

//ServiceOfferingEntity ...
type ServiceOfferingEntity struct {
	Label             string   `json:"label"`
	Provider          string   `json:"provider"`
	Description       string   `json:"description"`
	LongDescription   string   `json:"long_description"`
	Version           string   `json:"version"`
	URL               string   `json:"url"`
	InfoURL           string   `json:"info_url"`
	DocumentURL       string   `json:"documentation_url"`
	Timeout           string   `json:"timeout"`
	UniqueID          string   `json:"unique_id"`
	ServiceBrokerGUID string   `json:"service_broker_guid"`
	ServicePlansURL   string   `json:"service_plans_url"`
	Tags              []string `json:"tags"`
	Requires          []string `json:"requires"`
	IsActive          bool     `json:"active"`
	IsBindable        bool     `json:"bindable"`
	IsPlanUpdateable  bool     `json:"plan_updateable"`
}

//ToFields ...
func (resource ServiceOfferingResource) ToFields() ServiceOffering {
	entity := resource.Entity

	return ServiceOffering{
		GUID:              resource.Metadata.GUID,
		Label:             entity.Label,
		Provider:          entity.Provider,
		Description:       entity.Description,
		LongDescription:   entity.LongDescription,
		Version:           entity.Version,
		URL:               entity.URL,
		InfoURL:           entity.InfoURL,
		DocumentURL:       entity.DocumentURL,
		Timeout:           entity.Timeout,
		UniqueID:          entity.UniqueID,
		ServiceBrokerGUID: entity.ServiceBrokerGUID,
		ServicePlansURL:   entity.ServicePlansURL,
		Tags:              entity.Tags,
		Requires:          entity.Requires,
		IsActive:          entity.IsActive,
		IsBindable:        entity.IsBindable,
		IsPlanUpdateable:  entity.IsPlanUpdateable,
	}
}

//ServiceOfferingFields ...
type ServiceOfferingFields struct {
	Metadata ServiceOfferingMetadata
	Entity   ServiceOffering
}

//ServiceOfferingMetadata ...
type ServiceOfferingMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//ServiceOfferings ...
type ServiceOfferings interface {
	FindByLabel(serviceName string) (*ServiceOffering, error)
	Get(svcOfferingGUID string) (*ServiceOfferingFields, error)
}

type serviceOfferrings struct {
	client *client.Client
}

func newServiceOfferingAPI(c *client.Client) ServiceOfferings {
	return &serviceOfferrings{
		client: c,
	}
}

func (s *serviceOfferrings) Get(svcGUID string) (*ServiceOfferingFields, error) {
	rawURL := fmt.Sprintf("/v2/services/%s", svcGUID)
	svcFields := ServiceOfferingFields{}
	_, err := s.client.Get(rawURL, &svcFields)
	if err != nil {
		return nil, err
	}
	return &svcFields, err
}

func (s *serviceOfferrings) FindByLabel(serviceName string) (*ServiceOffering, error) {
	req := rest.GetRequest("v2/services")
	if serviceName != "" {
		req.Query("q", "label:"+serviceName)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	var services ServiceOffering
	var found bool
	err = s.listServicesOfferingWithPath(path, func(serviceOfferingResource ServiceOfferingResource) bool {
		services = serviceOfferingResource.ToFields()
		found = true
		return false
	})

	if err != nil {
		return nil, err
	}

	if found {
		return &services, err
	}
	//May not be found and no error

	return nil, bmxerror.New(ErrCodeServiceDoesnotExist,
		fmt.Sprintf("Given service %q doesn't exist", serviceName))

}

func (s *serviceOfferrings) listServicesOfferingWithPath(path string, cb func(ServiceOfferingResource) bool) error {
	_, err := s.client.GetPaginated(path, NewCCPaginatedResources(ServiceOfferingResource{}), func(resource interface{}) bool {
		if serviceOfferingResource, ok := resource.(ServiceOfferingResource); ok {
			return cb(serviceOfferingResource)
		}
		return false
	})
	return err
}
