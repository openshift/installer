package catalog

import (
	"fmt"
	"net/http"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/utils"
)

//ErrCodeServiceDoesnotExist ...
const ErrCodeServicePlanDoesnotExist = "ServicePlanDoesnotExist"
const ErrCodeServiceDoesnotExist = "ServiceDoesnotExist"
const ErrCodeServiceDeploymentNotFound = "ServiceDeploymentNotFound"

func newResourceCatalogAPI(c *client.Client) ResourceCatalogRepository {
	return &resourceCatalog{
		client: c,
	}
}

type resourceCatalog struct {
	client *client.Client
}

type ResourceCatalogRepository interface {
	Get(serviceID string, indepth bool) (models.Service, error)
	FindByName(name string, indepth bool) ([]models.Service, error)
	ListServices(cb func(service models.Service) bool) error
	ListServicePlans(cb func(servicePlan models.ServicePlan) bool, service models.Service) error
	GetServiceID(serviceName string) (string, error)
	GetServicePlanID(service models.Service, planName string) (string, error)
	GetServiceName(serviceID string) (string, error)
	GetServicePlanName(servicePlanID string) (string, error)
	ListDeployments(servicePlanID string) ([]models.ServiceDeployment, error)
	GetServicePlan(servicePlanID string) (models.ServicePlan, error)
	ListDeploymentAliases(servicePlanID string) ([]models.ServiceDeploymentAlias, error)
	GetDeploymentAlias(servicePlanID string, instanceTarget string, regionID string) (*models.ServiceDeploymentAlias, error)
	GetServices() ([]models.Service, error)
	GetServicePlans(service models.Service) ([]models.ServicePlan, error)
}

func (r *resourceCatalog) GetServicePlanID(service models.Service, planName string) (string, error) {
	var servicePlanID string
	err := r.ListServicePlans(func(servicePlan models.ServicePlan) bool {
		if servicePlan.Name == planName {
			servicePlanID = servicePlan.ID
			return false
		}
		return true
	}, service)
	if err != nil {
		return "", err
	}
	return servicePlanID, nil
}

func (r *resourceCatalog) GetServiceID(serviceName string) (string, error) {
	var serviceID string
	err := r.ListServices(func(service models.Service) bool {
		if service.Name == serviceName {
			serviceID = service.ID
			return false
		}
		return true
	})
	if err != nil {
		return "", err
	}
	return serviceID, nil
}

func (r *resourceCatalog) ListServices(cb func(service models.Service) bool) error {
	listRequest := rest.GetRequest("/api/v1/")
	req, err := listRequest.Build()
	if err != nil {
		return err
	}

	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewResourceCatalogPaginatedResources(models.Service{}, *r.client.Config.Endpoint),
		func(resource interface{}) bool {
			if catalogResource, ok := resource.(models.Service); ok {
				return cb(catalogResource)
			}
			return false
		})

	return err
}
func (r *resourceCatalog) ListServicePlans(cb func(service models.ServicePlan) bool, service models.Service) error {
	var urlSuffix string
	if service.Kind == "iaas" {
		urlSuffix = "/flavor"
	} else {
		urlSuffix = "/plan"
	}
	listRequest := rest.GetRequest("/api/v1/" + service.ID + urlSuffix)
	req, err := listRequest.Build()
	if err != nil {
		return err
	}
	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewResourceCatalogPaginatedResources(models.ServicePlan{}, *r.client.Config.Endpoint),
		func(resource interface{}) bool {
			if resourcePlan, ok := resource.(models.ServicePlan); ok {
				return cb(resourcePlan)
			}
			return false
		})

	return err
}

func (r *resourceCatalog) GetServiceName(serviceID string) (string, error) {
	request := rest.GetRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/api/v1/"+serviceID))
	resp := map[string]interface{}{}
	_, err := r.client.SendRequest(request, &resp)
	if err != nil {
		return "", err
	}
	if resp["kind"] == "runtime" || resp["kind"] == "service" || resp["kind"] == "iaas" || resp["kind"] == "platform_service" || resp["kind"] == "template" {
		if name, ok := resp["name"].(string); ok {
			return name, nil
		}
		return "", nil
	}
	return "", bmxerror.New(ErrCodeServiceDoesnotExist,
		fmt.Sprintf("Given service : %q doesn't exist", serviceID))
}

func (r *resourceCatalog) GetServicePlanName(servicePlanID string) (string, error) {
	request := rest.GetRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/api/v1/"+servicePlanID))
	resp := map[string]interface{}{}
	_, err := r.client.SendRequest(request, &resp)
	if err != nil {
		return "", err
	}
	if resp["kind"] == "flavor" || resp["kind"] == "plan" {
		if name, ok := resp["name"].(string); ok {
			return name, nil
		}
		return "", nil
	}
	return "", bmxerror.New(ErrCodeServicePlanDoesnotExist,
		fmt.Sprintf("Given service plan : %q doesn't exist", servicePlanID))
}

func (r *resourceCatalog) ListDeployments(servicePlanID string) ([]models.ServiceDeployment, error) {
	deployments := []models.ServiceDeployment{}
	listRequest := rest.GetRequest("/api/v1/" + servicePlanID + "/deployment?include=*")
	req, err := listRequest.Build()
	if err != nil {
		return deployments, err
	}
	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewResourceCatalogPaginatedResources(models.ServiceDeployment{}, *r.client.Config.Endpoint),
		func(resource interface{}) bool {
			if catalogDeployment, ok := resource.(models.ServiceDeployment); ok {
				deployments = append(deployments, catalogDeployment)
			}
			return true
		})
	return deployments, err
}

func (r *resourceCatalog) Get(serviceID string, indepth bool) (models.Service, error) {
	request := rest.GetRequest(helpers.GetFullURL(*r.client.Config.Endpoint, fmt.Sprintf("/api/v1/%s", serviceID)))
	if indepth {
		request = request.Query("include", "*")
	}
	service := models.Service{}
	resp, err := r.client.SendRequest(request, &service)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return models.Service{}, bmxerror.New(ErrCodeServicePlanDoesnotExist,
				fmt.Sprintf("Given service : %q doesn't exist", serviceID))
		}
		return models.Service{}, err
	}
	return service, nil
}

func (r *resourceCatalog) FindByName(name string, indepth bool) ([]models.Service, error) {
	services := []models.Service{}
	request := rest.GetRequest("/api/v1/").Query("q", name)
	if indepth {
		request = request.Query("include", "*")
	}
	req, err := request.Build()
	if err != nil {
		return services, err
	}
	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewResourceCatalogPaginatedResources(models.Service{}, *r.client.Config.Endpoint),
		func(rb interface{}) bool {
			if r, ok := rb.(models.Service); ok {
				services = append(services, visitServiceTree(r, name)...)
			}
			return true
		})
	if err != nil {
		return []models.Service{}, err
	}
	if len(services) == 0 {
		return services, bmxerror.New(ErrCodeServiceDoesnotExist,
			fmt.Sprintf("Given service : %q doesn't exist", name))
	}
	return services, err
}

func (r *resourceCatalog) GetServicePlan(servicePlanID string) (models.ServicePlan, error) {
	request := rest.GetRequest(helpers.GetFullURL(*r.client.Config.Endpoint, fmt.Sprintf("/api/v1/%s", servicePlanID)))
	servicePlan := models.ServicePlan{}
	resp, err := r.client.SendRequest(request, &servicePlan)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return models.ServicePlan{}, bmxerror.New(ErrCodeServicePlanDoesnotExist,
				fmt.Sprintf("Given service plan : %q doesn't exist", servicePlanID))
		}
		return models.ServicePlan{}, err
	}
	return servicePlan, nil
}

func (r *resourceCatalog) ListDeploymentAliases(serviceDeploymentID string) ([]models.ServiceDeploymentAlias, error) {
	aliases := []models.ServiceDeploymentAlias{}
	listRequest := rest.GetRequest("/api/v1/" + serviceDeploymentID + "/alias?include=*")
	req, err := listRequest.Build()
	if err != nil {
		return aliases, err
	}
	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewResourceCatalogPaginatedResources(models.ServiceDeploymentAlias{}, *r.client.Config.Endpoint),
		func(resource interface{}) bool {
			if deploymentAlias, ok := resource.(models.ServiceDeploymentAlias); ok {
				aliases = append(aliases, deploymentAlias)
			}
			return true
		})
	return aliases, err
}

func (r *resourceCatalog) GetDeploymentAlias(servicePlanID string, instanceTarget string, currentRegion string) (*models.ServiceDeploymentAlias, error) {
	deployments, err := r.ListDeployments(servicePlanID)
	if err != nil {
		return nil, err
	}
	found := false
	var deploymentID string
	for _, deployment := range deployments {
		deploymentLocation := utils.GetLocationFromTargetCRN(deployment.Metadata.Deployment.TargetCrn.Resource)
		if deploymentLocation == instanceTarget {
			deploymentID = deployment.ID
			found = true
			break
		}
	}
	if !found {
		//Should not go here since instanceTarget is get from deployments when create service instance
		return nil, bmxerror.New(ErrCodeServiceDeploymentNotFound,
			fmt.Sprintf("Service alias Deployment doesn't exist for %q", instanceTarget))
	}
	aliases, err := r.ListDeploymentAliases(deploymentID)
	if err != nil {
		return nil, err
	}
	for _, alias := range aliases {
		if alias.Metadata.Deployment.Location == currentRegion {
			return &alias, nil
		}
	}
	return nil, nil
}

func visitServiceTree(rootService models.Service, name string) []models.Service {
	services := []models.Service{}
	if rootService.Name == name && isService(rootService) {
		services = append(services, rootService)
	}
	for _, child := range rootService.Children {
		services = append(services, visitServiceTree(child, name)...)
	}
	return services
}

func (r *resourceCatalog) GetServices() ([]models.Service, error) {
	listRequest := rest.GetRequest("/api/v1/")
	var services []models.Service
	req, err := listRequest.Build()
	if err != nil {
		return nil, err
	}

	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewResourceCatalogPaginatedResources(models.Service{}, *r.client.Config.Endpoint),
		func(resource interface{}) bool {
			if catalogResource, ok := resource.(models.Service); ok {
				services = append(services, catalogResource)
				return true
			}
			return false
		})

	if err != nil {
		return []models.Service{}, err
	}

	return services, nil
}

func (r *resourceCatalog) GetServicePlans(service models.Service) ([]models.ServicePlan, error) {
	var urlSuffix string
	if service.Kind == "iaas" {
		urlSuffix = "/flavor"
	} else {
		urlSuffix = "/plan"
	}
	listRequest := rest.GetRequest("/api/v1/" + service.ID + urlSuffix)
	req, err := listRequest.Build()
	if err != nil {
		return nil, err
	}
	var servicePlans []models.ServicePlan
	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewResourceCatalogPaginatedResources(models.ServicePlan{}, *r.client.Config.Endpoint),
		func(resource interface{}) bool {
			if resourcePlan, ok := resource.(models.ServicePlan); ok {
				servicePlans = append(servicePlans, resourcePlan)
			}
			return false
		})

	if err != nil {
		return []models.ServicePlan{}, err
	}

	return servicePlans, nil
}

func isService(e models.Service) bool {
	// TODO: COS is 'iaas' kind, but considered to be a service
	if e.Kind == "service" || e.Kind == "iaas" {
		return true
	}
	return false
}
