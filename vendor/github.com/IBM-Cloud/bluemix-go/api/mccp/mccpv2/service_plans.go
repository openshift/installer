package mccpv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//ErrCodeServicePlanDoesNotExist ...
const ErrCodeServicePlanDoesNotExist = "ServicePlanDoesNotExist"

//ServicePlan ...
type ServicePlan struct {
	GUID                string
	Name                string `json:"name"`
	Description         string `json:"description"`
	IsFree              bool   `json:"free"`
	IsPublic            bool   `json:"public"`
	IsActive            bool   `json:"active"`
	ServiceGUID         string `json:"service_guid"`
	UniqueID            string `json:"unique_id"`
	ServiceInstancesURL string `json:"service_instances_url"`
}

//ServicePlanResource ...
type ServicePlanResource struct {
	Resource
	Entity ServicePlanEntity
}

//ServicePlanEntity ...
type ServicePlanEntity struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	IsFree              bool   `json:"free"`
	IsPublic            bool   `json:"public"`
	IsActive            bool   `json:"active"`
	ServiceGUID         string `json:"service_guid"`
	UniqueID            string `json:"unique_id"`
	ServiceInstancesURL string `json:"service_instances_url"`
}

//ToFields ...
func (resource ServicePlanResource) ToFields() ServicePlan {
	entity := resource.Entity

	return ServicePlan{
		GUID:                resource.Metadata.GUID,
		Name:                entity.Name,
		Description:         entity.Description,
		IsFree:              entity.IsFree,
		IsPublic:            entity.IsPublic,
		IsActive:            entity.IsActive,
		ServiceGUID:         entity.ServiceGUID,
		UniqueID:            entity.UniqueID,
		ServiceInstancesURL: entity.ServiceInstancesURL,
	}
}

//ServicePlanFields ...
type ServicePlanFields struct {
	Metadata ServicePlanMetadata
	Entity   ServicePlan
}

//ServicePlanMetadata ...
type ServicePlanMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//ServicePlans ...
type ServicePlans interface {
	FindPlanInServiceOffering(serviceOfferingGUID string, planType string) (*ServicePlan, error)
	Get(planGUID string) (*ServicePlanFields, error)
}

type servicePlan struct {
	client *client.Client
}

func newServicePlanAPI(c *client.Client) ServicePlans {
	return &servicePlan{
		client: c,
	}
}

func (s *servicePlan) Get(planGUID string) (*ServicePlanFields, error) {
	rawURL := fmt.Sprintf("/v2/service_plans/%s", planGUID)
	planFields := ServicePlanFields{}
	_, err := s.client.Get(rawURL, &planFields)
	if err != nil {
		return nil, err
	}
	return &planFields, err
}

func (s *servicePlan) FindPlanInServiceOffering(serviceOfferingGUID string, planType string) (*ServicePlan, error) {
	req := rest.GetRequest("/v2/service_plans")
	if serviceOfferingGUID != "" {
		req.Query("q", "service_guid:"+serviceOfferingGUID)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	plans, err := s.listServicesPlanWithPath(path)
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, bmxerror.New(ErrCodeServicePlanDoesNotExist,
			fmt.Sprintf("Given plan %q doesn't  exist for the service %q", planType, serviceOfferingGUID))
	}
	for _, p := range plans {
		if p.Name == planType {
			return &p, nil
		}

	}
	return nil, bmxerror.New(ErrCodeServicePlanDoesNotExist,
		fmt.Sprintf("Given plan %q doesn't  exist for the service %q", planType, serviceOfferingGUID))

}

func (s *servicePlan) listServicesPlanWithPath(path string) ([]ServicePlan, error) {
	var servicePlans []ServicePlan
	_, err := s.client.GetPaginated(path, NewCCPaginatedResources(ServicePlanResource{}), func(resource interface{}) bool {
		if servicePlanResource, ok := resource.(ServicePlanResource); ok {
			servicePlans = append(servicePlans, servicePlanResource.ToFields())
			return true
		}
		return false
	})
	return servicePlans, err
}
